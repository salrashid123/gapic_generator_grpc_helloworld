# Using Google's Client Library Generation system

This tutorial on how to create a small "helloworld" client library using  Google's `GAPIC`  (Generated API Client) Generator.

Specifically, this repository sets up a simple gRPC client-server application where the GAPIC generator creates the client library based on the protocol buffer and service definition files.  These definition files allows you preset settings into the generated library that describe things like backoff-retry, the hostname to connecton etc.  GAPIC based clients can also offload handling various tasks such as managing [Long Running Operations](https://aip.dev/151).  Most google API libraries use GAPIC to help generate library language and automatically inject declarative configuration.  By declarative, I mean you define specifications in human-readable cofnig files and the generator creates the client side library framework accordingly.

The code provided here runs a golang gRPC server that implements handling long running requests by returning the `Operations` object as described in the links provided above.  

- Server `grpc_server.go`
  - `SayHello`:  This simply echo's back a message
  - `SayHelloLRO`:  Endpoint returns the `Operations` object that the GAPIC client side code can understand and manage

- Client `gapic_client.go`
  - GAPIC gRPC Client

- Client `grpc_client.go`
  - Plain gRPC Client

- Admin Client `admin_client.go`
  - gRPC client that manages the LRO objects on `grpc_server.go`

- Envoy gRPC Gateway
  - Envoy proxy that transfrom http to gRPC requests (eg, [gRPC Transcoding](https://www.envoyproxy.io/docs/envoy/v1.9.0/configuration/http_filters/grpc_json_transcoder_filter))


## Setup

- Install [protoc](https://github.com/protocolbuffers/protobuf/releases), `go 1.14+`

```
$ protoc --version
libprotoc 3.14.0
```

add the following entry to `/etc/hosts/`
(the SNI for the certificates included in this repo uses this name)
```
127.0.0.1 server.domain.com
```

Clone the repo and acqure prerequsites

```bash
git clone https://github.com/salrashid123/gapic_generator_grpc_helloworld.git
cd gapic_generator_grpc_helloworld

export GOPATH=$GOPATH:`pwd`

git clone https://github.com/googleapis/api-common-protos

go get golang.org/x/net/context \
        golang.org/x/oauth2/google \
        golang.org/x/net/http2 \
        google.golang.org/grpc \
        google.golang.org/grpc/credentials \
        google.golang.org/grpc/health \
        google.golang.org/grpc/health/grpc_health_v1 \
        google.golang.org/grpc/metadata \
        google.golang.org/api/option \
        google.golang.org/api/transport \
        github.com/google/uuid \
        github.com/googleapis/gax-go/v2 \
        github.com/golang/protobuf/protoc-gen-go \
        github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic
```

### Configure service config

Configure the backoff retry specifications at the method or service level

For more information on the service config files, see [service_config.proto](https://github.com/grpc/grpc-proto/blob/master/grpc/service_config/service_config.proto)

- echo_grpc_service_config.json

```json
{
    "methodConfig": [
      {
        "name": [
          {
            "service": "echo.EchoServer",
            "method": "SayHello"
          },
          {
            "service": "echo.EchoServer",
            "method": "SayHelloLRO"
          }
        ],
        "timeout": "600s",
        "retryPolicy": {
          "initialBackoff": "0.200s",
          "maxBackoff": "60s",
          "backoffMultiplier": 1.3,
          "retryableStatusCodes": [
            "UNKNOWN",
            "UNAVAILABLE"
          ]
        }
      }     
    ]     
  }
```

### Generate gRPC stubs and client library

The following directives compiles the proto files, generates the descriptor (used by envoy), then sets up the gapic clients (as `echoclient` package)

```bash
 protoc -I ./api-common-protos  -I src/echo  --descriptor_set_out=src/echo/echo.proto.pb  --include_imports   --go_out=plugins=grpc:src/echo/ --go_gapic_out src/     --go_gapic_opt="go-gapic-package=echoclient"';echoclient'       --go_gapic_opt="grpc-service-config=echo_grpc_service_config.json" src/echo/echo.proto
```

For more information on, see [gapic-generator](https://github.com/googleapis/gapic-generator) options.


Notice that the generated clients under `src/echoclient/echo_server_client.go` includes the retry scheme configured above:

```golang
func defaultEchoServerCallOptions() *EchoServerCallOptions {
	return &EchoServerCallOptions{
		SayHello: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.Unknown,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    200 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.30,
				})
			}),
		},
	}
}
```

### Start Server

```bash
go run src/grpc_server.go  -grpcport :50051
```


### Run gRPC Client

To verify things are working, run the plain gRPC client:

```
go run src/grpc_client.go  -cacert CA_crt.pem -host localhost:50051 -servername server.domain.com
```

You should see output on the  client and server similar to:

```bash
$ go run src/grpc_server.go  -grpcport :50051
2019/08/23 17:08:17 Starting gRPC sever on port :50051
2019/08/23 17:08:25 Got SayHello -->  SayHello grpc msg 


$ go run src/grpc_client.go  -cacert CA_crt.pem -host localhost:50051 -servername server.domain.com
2019/08/23 17:08:25 Usign gRPC
2019/08/23 17:08:25 SayHello Response: message:"Hello SayHello grpc msg   from hostname srashid1" 
```

Note that while using a plain gRPC client, you need to configure low-level gRPC specifications and the host settings on the `conn` object

```golang
import 	"echo"  // this is the geneated gRPC stubs

conn, err = grpc.Dial("grpc.domain.com:50051", grpc.WithTransportCredentials(ce), grpc.WithPerRPCCredentials(rpcCreds))
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
```

### Run GAPIC 

Run the GAPIC client

```
go run src/gapic_client.go  -cacert CA_crt.pem  -servername server.domain.com
```

You should see an output similar to:

```bash
$ go run src/gapic_client.go  -cacert CA_crt.pem  -servername server.domain.com
2019/08/23 17:13:31 Usign GAPIC
2019/08/23 17:13:31 message:"Hello SayHello gapic msg  from hostname srashid1" 
2019/08/23 17:13:31 Starting operationID 00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:31 Done: false
2019/08/23 17:13:33 message:"Hello Callback SayHelloLRO gapic " 
2019/08/23 17:13:33 Done: true


$ go run src/grpc_server.go  -grpcport :50051
2019/08/23 17:09:02 Starting gRPC sever on port :50051
2019/08/23 17:13:31 Got SayHello -->  SayHello gapic msg
2019/08/23 17:13:31 Got SayHelloLRO -->  SayHelloLRO gapic 
2019/08/23 17:13:31 Work request queued 00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:31 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:31 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:33 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:33 LRO Complete 00fdc6c2-c604-11e9-88d6-e86a641d5560
```

What the GAPIC client executes is a direct API call to `/SayHello` and another _long running_ api call to `/SayHelloLRO`

Notice the invocation and comapre to the gRPC example above.  The library involves less less low level configuratin for gRPC.  For example, the library has the service host/port configuration already built in so there is no reason to `Dial` the endpoint as with gRPC clients direct.

- `/SayHello`:

```golang
import 	"echoclient"  // this is the GAPIC generated client

gclient, err := echoclient.NewEchoServerClient(ctx,
option.WithGRPCDialOption(grpc.WithTransportCredentials(ce)),
option.WithoutAuthentication())
if err != nil {
	log.Fatalf("could not get gapic client: %v", err)
}

resp, err := gclient.SayHello(ctx, &echo.EchoRequest{Name: "SayHello gapic msg"})
if err != nil {
	log.Fatalf("could not get say hello gapic: %v", err)
}
log.Printf("%v", resp)
```

- `/SayHelloLRO`

Executes a long running operations. 

```golang
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
```


On the server side, all that `/SayHelloLRO` initially returns is an `Operations` object back to the GAPIC client.  THe GAPIC client takes the unique id embedded within that object to internally poll the service for the final outcome (you can see the polling as lines

```
2019/08/23 17:13:31 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:31 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
2019/08/23 17:13:33 GetOperation:  00fdc6c2-c604-11e9-88d6-e86a641d5560
```

On the server end, each `GetOperation` has a given probability of returning the final outcome.  In the case here, i've set the probability of success for any one call to `70%`:

```golang
func (s *operationsServer) GetOperation(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	log.Println("GetOperation: ", in.Name)

	var answer *lropb.Operation
	chances := rand.Intn(100)

	if wr, ok := workRequestMap[in.Name]; ok {
		if chances >= 70 {
			log.Printf("LRO Complete %v", in.Name)
			answer = &lropb.Operation{
				Name: in.Name,
				Done: true,
			}
			resp, _ := ptypes.MarshalAny(
				&echo.EchoReply{Message: "Hello Callback " + wr.Req.Name},
			)
			delete(workRequestMap, in.Name)
			answer.Result = &lropb.Operation_Response{Response: resp}
		} else {
			answer = &lropb.Operation{
				Name: in.Name,
				Done: false,
			}
		}
	} else {
		return &lropb.Operation{
			Name: in.Name,
			Done: false,
		}, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}
	return answer, nil
}
```

### Run LRO Admin Client

The admin client is basically just a gRPC client intended query and update the service for any inflight LRO or to cancel one in flight.  Its basically just an interface to manage LRO that are currently being handled.

All the admin client currently does is lists the operations that maybe in flight.

If one is, you may see its Name listed in the map, otherwise, it's an empty map

```bash
$ go run src/admin_client.go  -cacert CA_crt.pem -host localhost:50051 -servername server.domain.com 
2019/08/23 17:27:40 Running Admin Client
2019/08/23 17:27:40 ListOperations:
2019/08/23 17:27:40    {[name:"fb00f242-c605-11e9-88d6-e86a641d5560" ]  {} [] 0}
```

### Envoy and gRPC Transcoding

THis repo also contains an envoy proxy to do [gRPC Tanscoding](https://github.com/grpc-ecosystem/grpc-httpjson-transcoding).  Basically, it transforms an HTTP Rest call to gRPC.  The intent of this section is to just show the gRPC endpoints accessed via REST as well as the id field for an LRO for later use.  

> TODO: update this section with a LRO example w/ HTTP REST...i don't know if gapic generates that part though or if i have to do that manually

The flow in this section is

`curl`  --> (http) -->  `envoy`  --> (grpc) --> `gRPC Service`

First start gRPC server and envoy.  In the example below, i'm using the envoy binary directly...you can use docker or extract it from the official container as shown [here](https://github.com/salrashid123/envoy_discovery/issues/3#issuecomment-522073806).

```bash
docker cp `docker create envoyproxy/envoy-dev:latest`:/usr/local/bin/envoy .

./envoy -c grpc-transcoding.yaml
```

Invoke the endpoints via curl:

- `/v1/sayhello/{name}`:

```bash
curl -v  --cacert CA_crt.pem --resolve server.domain.com:8080:127.0.0.1 https://server.domain.com:8080/v1/sayhello/foo
< HTTP/1.1 200 OK
< content-type: application/json
< x-envoy-upstream-service-time: 0
< grpc-status: 0
< grpc-message: 
< content-length: 52
< date: Fri, 23 Aug 2019 21:08:33 GMT
< server: envoy
< 
{
 "message": "Hello foo  from hostname yourhost"
}
```


- `v1/sayhellolro/{name}`

```bash
curl -v  --cacert CA_crt.pem --resolve server.domain.com:8080:127.0.0.1 https://server.domain.com:8080/v1/sayhellolro/foo
< HTTP/1.1 200 OK
< content-type: application/json
< x-envoy-upstream-service-time: 0
< grpc-status: 0
< grpc-message: 
< content-length: 68
< date: Fri, 23 Aug 2019 21:07:50 GMT
< server: envoy
< 
{
 "name": "10e654f2-c5ea-11e9-9045-e86a641d5560",
 "done": false
}
```

The LRO endpoint returned the name of the LRO to track as well as its state.


### Caveats/Notes

In the course of working through the golang examples, I came across a couple of caveats to consider

#### Authentication

The GAPIC generated clients always tries to acquire your google [Application Default Credentials](https://cloud.google.com/docs/authentication/production)..._even if you don't declare it or need it_.  THis should get addresses since the target service maynot even know or care about google-centric credential objects.  For now, you can override the credentials by either setting null types or explictly declaring no Credentials:

- golang

```golang
gclient, err := echoclient.NewEchoServerClient(ctx,
	option.WithoutAuthentication())
```

or

```golang
tok := &oauth2.Token{}
tokSrc := oauth2.StaticTokenSource(tok)
gclient, err := echoclient.NewEchoServerClient(ctx,
	option.WithTokenSource(tokSrc))
```


The token source emitted via any client can be anything..not just google oauth2 tokens.  For example, you can configure gRPC to emit OpenID Connect tokens.

- (gRPC Authentication with Google OpenID Connect tokens)[https://github.com/salrashid123/grpc_google_id_tokens]
  Note, I'm using google's ID token there but you can use any provider as the token with some overrides.

See the examples there on how to acqurie and emit ID tokens as well as how to configure a gRPC interceptor to validate all handlers.



#### Links

- GAPIC, gRPC

- [API Improvement Proposals (AIPs)](https://aip.dev/client-libraries)
- [GAPIC go Generator](https://github.com/googleapis/gapic-generator-go)
- [gRPC Godocs](https://godoc.org/google.golang.org/grpc)
- [GAPIC Showcase](https://github.com/googleapis/gapic-showcase)
- [google.golang.org/genproto/googleapis/longrunning](https://godoc.org/google.golang.org/genproto/googleapis/longrunning)


- Misc

- [gRPC HealthCheck Proxy](https://github.com/salrashid123/grpc_health_proxy)
  Use this library to perform http-based gRPC Healthchecks
- [golang/oauth2](https://github.com/salrashid123/grpc_google_id_tokens)
   _unofficial_ `google/oauth2` library for OIDC.  Use this to acqurie OIDC tokens for authentication support.

#### Appendix

Trace logging environment variables for gRPC:

- [gRPC Environment variables](https://github.com/grpc/grpc/blob/master/doc/environment_variables.md)

```bash
export GRPC_TRACE=all
export GRPC_VERBOSITY=INFO
export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=`pwd`/CA_crt.pem
```

