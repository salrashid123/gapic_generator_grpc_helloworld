from __future__ import print_function

import logging
import time


def run_gapic_lro():
    from helloworld import GreeterClient
    from helloworld import HelloRequest
    from helloworld import HelloReply
    from google.auth.credentials import AnonymousCredentials
    import grpc
    from helloworld.services.greeter.transports import GreeterGrpcTransport

    creds = AnonymousCredentials()
    # normally use creds but since we have custom CA certs..
    #c = GreeterClient(credentials=cred)

    # load custom certs
    with open('CA_crt.pem', 'rb') as fh:
      root_ca = fh.read()
    channel_creds = grpc.ssl_channel_credentials(
        root_certificates=root_ca,
    )
    transport = GreeterGrpcTransport(
        credentials = creds,
        ssl_channel_credentials = channel_creds,
    )
    c = GreeterClient(transport=transport)


    req = HelloRequest()
    req.name = 'salLRO'

    def my_callback(future):
        result = future.result()
        logging.info(result)

    op = c.say_hello_lro(request=req)
    rr = op.result()
    print(rr)
    op.add_done_callback(my_callback)

    while True:
        print(time.strftime("%c"))
        time.sleep(1)


def run_gapic():
    from helloworld import GreeterClient
    from helloworld import HelloRequest
    from helloworld import HelloReply
    import grpc
    from helloworld.services.greeter.transports import GreeterGrpcTransport
    from google.auth.credentials import AnonymousCredentials

    creds = AnonymousCredentials()
    # normally use creds but since we have custom CA certs..
    #c = GreeterClient(credentials=cred)

    # load custom certs
    with open('CA_crt.pem', 'rb') as fh:
      root_ca = fh.read()
    channel_creds = grpc.ssl_channel_credentials(
        root_certificates=root_ca,
    )
    transport = GreeterGrpcTransport(
        credentials = creds,
        ssl_channel_credentials = channel_creds,
    )
    c = GreeterClient(transport=transport)



    req = HelloRequest()
    req.name = 'sal'
    resp = c.say_hello(request=req)
    logging.info(resp.message)


def run_standard():
    import grpc
    import helloworld_pb2
    import helloworld_pb2_grpc

    with open('CA_crt.pem', 'rb') as fh:
      root_ca = fh.read()

    channel_creds = grpc.ssl_channel_credentials(
        root_certificates=root_ca,
    )
    channel = grpc.secure_channel(target='grpc.domain.com:50051', credentials=channel_creds)
    stub = helloworld_pb2_grpc.GreeterStub(channel)
    response = stub.SayHello(helloworld_pb2.HelloRequest(name='sal'))

    logging.info("Greeter client received: " + response.message)


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    #run_gapic()
    run_gapic_lro()
    #run_standard()
