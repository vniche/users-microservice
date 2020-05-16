# users-microservice

## DGraph

```shell
docker run -it --rm \
    -p 6080:6080 -p 8080:8080 \
    -p 9080:9080 -p 8000:8000 \
    dgraph/standalone:v20.03.0
```

## Protobuff

```shell
protoc -I grpc/ proto/users.proto --go_out=plugins=grpc:grpc
```
