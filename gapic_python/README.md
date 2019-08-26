
## Python gapic with LRO

 -[Instructions for create a gRPC client for google cloud services](https://github.com/GoogleCloudPlatform/grpc-gcp-python/blob/master/doc/gRPC-client-user-guide.md#generate-client-api-from-proto-files)

### First some housekeeping (not necessary for the log)

- Requires `python 3.7`


```bash
cd server/
export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=`pwd`/CA_crt.pem
```

- Add to `/etc/hosts`
```
127.0.0.1 grpc.domain.com
```


### Clean up the previous iteration, if any

```bash
rm -rf destserver/ destclient/
mkdir destserver/ destclient/
rm -rf env
```

### Generate gRPC Server
```bash
virtualenv   --no-site-packages -p python3.7 env
source env/bin/activate
pip install grpcio-tools  gapic-generator   protobuf proto-plus google.api.core grpc-google-longrunning-v2
```

Get protos

```
git clone https://github.com/googleapis/googleapis.git
```

### Build gRPC server and gapic clients

```bash
  export PYTHONPATH=`pwd`/destserver:$PYTHONPATH
  export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=`pwd`/CA_crt.pem

  python -m grpc_tools.protoc  -I googleapis/ --proto_path=. -I=/usr/local/include/google/protobuf/ -I . --python_out=destserver/  --grpc_python_out=destserver/ helloworld.proto

  python -m grpc_tools.protoc  -I googleapis/ --proto_path=. -I=/usr/local/include/google/protobuf/ -I . --python_out=env/lib/python3.7/site-packages   --grpc_python_out=env/lib/python3.7/site-packages ./googleapis/google/longrunning/operations.proto

  python -m grpc_tools.protoc -I=. -I=/usr/local/include/google/protobuf/  -I=googleapis/  --python_gapic_out=destclient/ --grpc_python_out=destclient/ --include_imports --include_source_info -o helloworld_descriptor.desc helloworld.proto
```
(note,, i'm manually copying the generating `operations.proto` and 'copying them into the python 37 folders)

### Start Server

```
python greeter_server.py
```


###  Setup client

in a *new* window
```

cd client/
virtualenv   --no-site-packages -p python3.7 env
source env/bin/activate
export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=`pwd`/CA_crt.pem

```

- for GAPIC

```
pip install ../server/destclient/

```

- for gRPC


Point the `PYTHONPATH` to point to where the `helloworld` compiled files sit

```
export PYTHONPATH=$PYTHONPATH:`pwd`/../server/destserver
```


### Run client

```
python greeter_client.py
```
