package main

import (
	"crypto/tls"
	"crypto/x509"
	"echo"
	"flag"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"golang.org/x/oauth2"
)

const ()

var ()

func main() {

	address := flag.String("host", "localhost:50051", "host:port of gRPC server")
	cacert := flag.String("cacert", "CA_crt.pem", "CACert for server")
	serverName := flag.String("servername", "server.domain.com", "SNI for server")
	flag.Parse()

	var err error
	var conn *grpc.ClientConn

	tok := &oauth2.Token{}
	rpcCreds := oauth.NewOauthAccess(tok)

	log.Printf("Usign gRPC")

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

	conn, err = grpc.Dial(*address, grpc.WithTransportCredentials(ce), grpc.WithPerRPCCredentials(rpcCreds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServerClient(conn)
	ctx := context.Background()

	r, err := c.SayHello(ctx, &echo.EchoRequest{Name: "SayHello grpc msg "})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("SayHello Response: %v", r)

}
