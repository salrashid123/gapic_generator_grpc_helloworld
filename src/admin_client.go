package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	lropb "google.golang.org/genproto/googleapis/longrunning"
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
	serverName := flag.String("servername", "grpc.domain.com", "SNI for server")
	//operationid := flag.String("operationid", "", "OperationID for cancel/delte operations")
	flag.Parse()

	var err error
	var conn *grpc.ClientConn
	log.Printf("Running Admin Client")

	tok := &oauth2.Token{}
	rpcCreds := oauth.NewOauthAccess(tok)

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

	c := lropb.NewOperationsClient(conn)
	ctx := context.Background()

	log.Printf("ListOperations:")

	operationsList, err := c.ListOperations(ctx, &lropb.ListOperationsRequest{})
	if err != nil {
		log.Fatalf("did not list operations: %v", err)
	}
	for _, o := range operationsList.Operations {
		log.Printf("   %s", o.Name)
	}
}
