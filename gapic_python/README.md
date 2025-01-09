
## Python gapic with LRO

 -[Instructions for create a gRPC client for google cloud services](https://github.com/GoogleCloudPlatform/grpc-gcp-python/blob/master/doc/gRPC-client-user-guide.md#generate-client-api-from-proto-files)

### First some housekeeping (not necessary for the log)

- Requires `python 3.9` (really 3+)


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
virtualenv  env
source env/bin/activate
$ python --version
   Python 3.11.9
pip install grpcio-tools  gapic-generator   protobuf proto-plus google.api.core grpc-google-longrunning-v2
```

Get protos

```bash
git clone https://github.com/googleapis/googleapis.git
```

### Build gRPC server and gapic clients

```bash
  export PYTHONPATH=`pwd`/destserver:$PYTHONPATH

  python -m grpc_tools.protoc  -I googleapis/ --proto_path=. -I=/usr/local/include/google/protobuf/ -I . --python_out=destserver/  --grpc_python_out=destserver/ helloworld.proto

  python -m grpc_tools.protoc  -I googleapis/ --proto_path=. -I=/usr/local/include/google/protobuf/ -I . --python_out=./env/lib/python3.11/site-packages   --grpc_python_out=./env/lib/python3.11/site-packages ./googleapis/google/longrunning/operations.proto

  python -m grpc_tools.protoc -I=. -I=/usr/local/include/google/protobuf/  -I=googleapis/  --python_gapic_out=destclient/ --grpc_python_out=destclient/ --include_imports --include_source_info -o helloworld_descriptor.desc helloworld.proto
```

(note,, i'm manually copying the generating `operations.proto` and 'copying them into the python 311 folders)

### Start Server

```bash
python greeter_server.py
```


###  Setup client

in a *new* window

```bash
cd client/
virtualenv  env
source env/bin/activate
```

- for GAPIC

```bash
pip install ../server/destclient/
```

- for gRPC


Point the `PYTHONPATH` to point to where the `helloworld` compiled files sit

```bash
export PYTHONPATH=$PYTHONPATH:`pwd`/../server/destserver
```


### Run client

```bash
python greeter_client.py
```
