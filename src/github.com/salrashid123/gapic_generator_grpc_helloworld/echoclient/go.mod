module echoclient

go 1.17

require (

	github.com/salrashid123/gapic_generator_grpc_helloworld/echo v0.0.0
)



replace (
	github.com/salrashid123/gapic_generator_grpc_helloworld/echo  => ../echo
)