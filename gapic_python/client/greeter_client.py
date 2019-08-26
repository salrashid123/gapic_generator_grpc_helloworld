from __future__ import print_function

import logging
import time


def run_gapic_lro():
    from helloworld import Greeter
    from helloworld import HelloRequest
    from helloworld import HelloReply

    from google.auth.credentials import AnonymousCredentials

    creds = AnonymousCredentials()
    c = Greeter(credentials=creds)
    req = HelloRequest()
    req.name = 'salLRO'

    def my_callback(future):
        result = future.result()
        logging.info(result)

    op = c.say_hello_lro(request=req)
    #rr = op.result()
    # print(rr)
    op.add_done_callback(my_callback)

    while True:
        print(time.strftime("%c"))
        time.sleep(1)


def run_gapic():
    from helloworld import Greeter
    from helloworld import HelloRequest
    from helloworld import HelloReply

    from google.auth.credentials import AnonymousCredentials

    creds = AnonymousCredentials()
    c = Greeter(credentials=creds)
    req = HelloRequest()
    req.name = 'sal'
    resp = c.say_hello(request=req)
    logging.info(resp.message)


def run_standard():
    import grpc
    import helloworld_pb2
    import helloworld_pb2_grpc

    creds = grpc.ssl_channel_credentials()
    channel = grpc.secure_channel('grpc.domain.com:50051', creds)
    stub = helloworld_pb2_grpc.GreeterStub(channel)
    response = stub.SayHello(helloworld_pb2.HelloRequest(name='sal'))

    logging.info("Greeter client received: " + response.message)


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    run_gapic()
    run_gapic_lro()
    # run_standard()
