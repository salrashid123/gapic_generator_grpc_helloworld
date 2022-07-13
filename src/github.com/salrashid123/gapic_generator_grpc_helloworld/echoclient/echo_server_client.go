// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go_gapic. DO NOT EDIT.

package echoclient

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"cloud.google.com/go/longrunning"
	lroauto "cloud.google.com/go/longrunning/autogen"
	gax "github.com/googleapis/gax-go/v2"
	echopb "github.com/salrashid123/gapic_generator_grpc_helloworld/echo"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var newEchoServerClientHook clientHook

// EchoServerCallOptions contains the retry settings for each method of EchoServerClient.
type EchoServerCallOptions struct {
	SayHello []gax.CallOption
	SayHelloLRO []gax.CallOption
}

func defaultEchoServerGRPCClientOptions() []option.ClientOption {
	return []option.ClientOption{
		internaloption.WithDefaultEndpoint("grpc.domain.com:50051"),
		internaloption.WithDefaultMTLSEndpoint("grpc.domain.com:50051"),
		internaloption.WithDefaultAudience("https://grpc.domain.com/"),
		internaloption.WithDefaultScopes(DefaultAuthScopes()...),
		internaloption.EnableJwtWithScope(),
		option.WithGRPCDialOption(grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

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
		SayHelloLRO: []gax.CallOption{
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

// internalEchoServerClient is an interface that defines the methods available from .
type internalEchoServerClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	SayHello(context.Context, *echopb.EchoRequest, ...gax.CallOption) (*echopb.EchoReply, error)
	SayHelloLRO(context.Context, *echopb.EchoRequest, ...gax.CallOption) (*SayHelloLROOperation, error)
	SayHelloLROOperation(name string) *SayHelloLROOperation
}

// EchoServerClient is a client for interacting with .
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type EchoServerClient struct {
	// The internal transport-dependent client.
	internalClient internalEchoServerClient

	// The call options for this service.
	CallOptions *EchoServerCallOptions

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient *lroauto.OperationsClient

}

// Wrapper methods routed to the internal client.

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *EchoServerClient) Close() error {
	return c.internalClient.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *EchoServerClient) setGoogleClientInfo(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *EchoServerClient) Connection() *grpc.ClientConn {
	return c.internalClient.Connection()
}

func (c *EchoServerClient) SayHello(ctx context.Context, req *echopb.EchoRequest, opts ...gax.CallOption) (*echopb.EchoReply, error) {
	return c.internalClient.SayHello(ctx, req, opts...)
}

func (c *EchoServerClient) SayHelloLRO(ctx context.Context, req *echopb.EchoRequest, opts ...gax.CallOption) (*SayHelloLROOperation, error) {
	return c.internalClient.SayHelloLRO(ctx, req, opts...)
}

// SayHelloLROOperation returns a new SayHelloLROOperation from a given name.
// The name must be that of a previously created SayHelloLROOperation, possibly from a different process.
func (c *EchoServerClient) SayHelloLROOperation(name string) *SayHelloLROOperation {
	return c.internalClient.SayHelloLROOperation(name)
}

// echoServerGRPCClient is a client for interacting with  over gRPC transport.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type echoServerGRPCClient struct {
	// Connection pool of gRPC connections to the service.
	connPool gtransport.ConnPool

	// flag to opt out of default deadlines via GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE
	disableDeadlines bool

	// Points back to the CallOptions field of the containing EchoServerClient
	CallOptions **EchoServerCallOptions

	// The gRPC API client.
	echoServerClient echopb.EchoServerClient

	// LROClient is used internally to handle long-running operations.
	// It is exposed so that its CallOptions can be modified if required.
	// Users should not Close this client.
	LROClient **lroauto.OperationsClient

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewEchoServerClient creates a new echo server client based on gRPC.
// The returned client must be Closed when it is done being used to clean up its underlying connections.
func NewEchoServerClient(ctx context.Context, opts ...option.ClientOption) (*EchoServerClient, error) {
	clientOpts := defaultEchoServerGRPCClientOptions()
	if newEchoServerClientHook != nil {
		hookOpts, err := newEchoServerClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}

	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}

	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	client := EchoServerClient{CallOptions: defaultEchoServerCallOptions()}

	c := &echoServerGRPCClient{
		connPool:    connPool,
		disableDeadlines: disableDeadlines,
		echoServerClient: echopb.NewEchoServerClient(connPool),
		CallOptions: &client.CallOptions,

	}
	c.setGoogleClientInfo()

	client.internalClient = c

	client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))
	if err != nil {
		// This error "should not happen", since we are just reusing old connection pool
		// and never actually need to dial.
		// If this does happen, we could leak connp. However, we cannot close conn:
		// If the user invoked the constructor with option.WithGRPCConn,
		// we would close a connection that's still in use.
		// TODO: investigate error conditions.
		return nil, err
	}
	c.LROClient = &client.LROClient
	return &client, nil
}

// Connection returns a connection to the API service.
//
// Deprecated.
func (c *echoServerGRPCClient) Connection() *grpc.ClientConn {
	return c.connPool.Conn()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *echoServerGRPCClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *echoServerGRPCClient) Close() error {
	return c.connPool.Close()
}

func (c *echoServerGRPCClient) SayHello(ctx context.Context, req *echopb.EchoRequest, opts ...gax.CallOption) (*echopb.EchoReply, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000 * time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).SayHello[0:len((*c.CallOptions).SayHello):len((*c.CallOptions).SayHello)], opts...)
	var resp *echopb.EchoReply
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoServerClient.SayHello(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *echoServerGRPCClient) SayHelloLRO(ctx context.Context, req *echopb.EchoRequest, opts ...gax.CallOption) (*SayHelloLROOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 600000 * time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))

	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).SayHelloLRO[0:len((*c.CallOptions).SayHelloLRO):len((*c.CallOptions).SayHelloLRO)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.echoServerClient.SayHelloLRO(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &SayHelloLROOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, resp),
	}, nil
}

// SayHelloLROOperation manages a long-running operation from SayHelloLRO.
type SayHelloLROOperation struct {
	lro *longrunning.Operation
}

// SayHelloLROOperation returns a new SayHelloLROOperation from a given name.
// The name must be that of a previously created SayHelloLROOperation, possibly from a different process.
func (c *echoServerGRPCClient) SayHelloLROOperation(name string) *SayHelloLROOperation {
	return &SayHelloLROOperation{
		lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),
	}
}

// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.
//
// See documentation of Poll for error-handling information.
func (op *SayHelloLROOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*echopb.EchoReply, error) {
	var resp echopb.EchoReply
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Poll fetches the latest state of the long-running operation.
//
// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and
// the operation has completed with failure, the error is returned and op.Done will return true.
// If Poll succeeds and the operation has completed successfully,
// op.Done will return true, and the response of the operation is returned.
// If Poll succeeds and the operation has not completed, the returned response and error are both nil.
func (op *SayHelloLROOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*echopb.EchoReply, error) {
	var resp echopb.EchoReply
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}

// Done reports whether the long-running operation has completed.
func (op *SayHelloLROOperation) Done() bool {
	return op.lro.Done()
}

// Name returns the name of the long-running operation.
// The name is assigned by the server and is unique within the service from which the operation is created.
func (op *SayHelloLROOperation) Name() string {
	return op.lro.Name()
}
