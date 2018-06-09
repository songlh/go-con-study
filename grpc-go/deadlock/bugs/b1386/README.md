### How to run [bug1386](https://github.com/grpc/grpc-go/issues/1386)?

```bash
$ cd ${GOPATH}
$ go get -u google.golang.org/grpc  
$ cd ${GOPATH}/src/google.golang.org/grpc
$ git checkout -b b1386 833680729394fcca4904ec569758c78b78411ee8
$ cp -rf eggs/ ${GOPATH}/src/bugs/
$ cd ${GOPATH}/src/bugs/eggs
$ go main.go
```

```
2018/06/09 15:30:51 Calling
2018/06/09 15:30:51 Server sleeping, will timeout
2018/06/09 15:30:56 client call failed: rpc error: code = DeadlineExceeded desc = context deadline exceeded
```
