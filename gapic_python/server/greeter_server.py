
import logging
import random
import time
import uuid
from concurrent import futures

import google.longrunning.operations_pb2
import google.protobuf.any_pb2
import grpc
import helloworld_pb2
import helloworld_pb2_grpc
from google.longrunning import operations_pb2
from google.longrunning.operations_pb2_grpc import OperationsServicer
from google.protobuf.any_pb2 import Any

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


workRequestMap = {}


class GreeterLRO(google.longrunning.operations_pb2_grpc.OperationsServicer):

    def GetOperation(self, request, context):
        logging.info("GetOperation " + request.name)
        if (request.name in workRequestMap):
            if (random.randint(1, 101) >= 70):
                logging.info("Operation complete " + request.name)
                rr = workRequestMap[request.name]
                respobj = helloworld_pb2.HelloReply(
                    message='HelloLRO ' + rr.name)
                some_any = google.protobuf.any_pb2.Any()
                some_any.Pack(respobj)
                workRequestMap.pop(request.name, None)
                return operations_pb2.Operation(name=request.name, done=True, response=some_any)
            else:
                return operations_pb2.Operation(name=request.name, done=False)
        context.set_details("operationID not found " + request.name)
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)


class Greeter(helloworld_pb2_grpc.GreeterServicer):

    def SayHello(self, request, context):
        logging.info(">>> got gRPC.. ")
        meta = dict(context.invocation_metadata())
        logging.info(meta)
        return helloworld_pb2.HelloReply(message='Hello, %s!' % request.name)

    def SayHelloLRO(self, request, context):
        logging.info(">>> got gRPC LRO.. ")
        meta = dict(context.invocation_metadata())
        id = str(uuid.uuid4())
        logging.info("Enqueue LRO with ID: " + id)
        workRequestMap[id] = request
        return operations_pb2.Operation(name=id, done=False)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))

    helloworld_pb2_grpc.add_GreeterServicer_to_server(Greeter(), server)
    google.longrunning.operations_pb2_grpc.add_OperationsServicer_to_server(
        GreeterLRO(), server)

    with open('server_key.pem', 'rb') as f:
        private_key = f.read()
    with open('server_crt.pem', 'rb') as f:
        server_crt = f.read()
    sc = grpc.ssl_server_credentials(((private_key, server_crt), ))

    server.add_secure_port('grpc.domain.com:50051', server_credentials=sc)
    logging.info(">>> Starting server..")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()
