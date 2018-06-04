### How to run [bug1459](https://github.com/grpc/grpc-go/issues/1459)?

```bash
$ cd ${GOPATH}
$ go get -u google.golang.org/grpc  
$ cd ${GOPATH}/src/google.golang.org/grpc
$ git checkout -b b1459 7db1564ba1229bc42919bb1f6d9c4186f3aa8678
$ cp -rf b1459/ ${GOPATH}src/google.golang.org/grpc/examples/
$ go run greeter_server/main.go -logtostderr &
$ go run greeter_client/main.go -logtostderr
```
