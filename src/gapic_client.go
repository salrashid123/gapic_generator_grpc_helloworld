package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"

	"github.com/salrashid123/gapic_generator_grpc_helloworld/echo"

	"github.com/salrashid123/gapic_generator_grpc_helloworld/echoclient"

	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const ()

var ()

func main() {

	cacert := flag.String("cacert", "CA_crt.pem", "CACert for server")
	serverName := flag.String("servername", "grpc.domain.com", "SNI for server")
	flag.Parse()

	var err error
	tok := &oauth2.Token{}
	//rpcCreds := oauth.NewOauthAccess(tok)
	tokSrc := oauth2.StaticTokenSource(tok)

	var tlsCfg tls.Config
	rootCAs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(*cacert)
	if err != nil {
		log.Fatalf("failed to load root CA certificates  error=%v", err)
	}
	if !rootCAs.AppendCertsFromPEM(pem) {
		log.Fatalf("no root CA certs parsed from file ")
	}
	tlsCfg.RootCAs = rootCAs
	tlsCfg.ServerName = *serverName

	ce := credentials.NewTLS(&tlsCfg)

	log.Printf("Usign GAPIC")
	ctx := context.Background()

	gclient, err := echoclient.NewEchoServerClient(ctx,
		option.WithGRPCDialOption(grpc.WithTransportCredentials(ce)),
		option.WithTokenSource(tokSrc))
	//	option.WithoutAuthentication())
	//  option.WithTokenSource(tokSrc))
	if err != nil {
		log.Fatalf("could not get gapic client: %v", err)
	}

	resp, err := gclient.SayHello(ctx, &echo.EchoRequest{Name: "SayHello gapic msg"})
	if err != nil {
		log.Fatalf("could not get say hello gapic: %v", err)
	}
	log.Printf("%v", resp)

	resplro, err := gclient.SayHelloLRO(ctx, &echo.EchoRequest{Name: "SayHelloLRO gapic "})
	if err != nil {
		log.Fatalf("could not get say hello gapic: %v", err)
	}

	log.Printf("Starting operationID %v", resplro.Name())
	log.Printf("Done: %v", resplro.Done())

	echoReply, err := resplro.Wait(ctx)
	if err != nil {
		log.Fatalf("could not wait synchronously: %v", err)
	}

	log.Printf("%v", echoReply)
	log.Printf("Done: %v", resplro.Done())
}
