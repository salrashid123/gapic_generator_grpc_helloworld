module echoclient

go 1.21

toolchain go1.22.4

require (

	github.com/salrashid123/gapic_generator_grpc_helloworld/echo v0.0.0
)



replace (
	github.com/salrashid123/gapic_generator_grpc_helloworld/echo  => ../echo
)