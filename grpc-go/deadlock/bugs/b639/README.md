### How to run [bug639](https://github.com/grpc/grpc-go/issues/639)?

```bash
$ cd ${GOPATH}
$ go get -u google.golang.org/grpc  
$ cd ${GOPATH}/src/google.golang.org/grpc
$ git checkout -b b639 306a1ee0fe2c012a074592da8ffe4e33b5204f2a
$ go run grpc-go-b639.go
```

```
2018/06/04 14:28:01 Dialing to server over a synchronous pipe...
2018/06/04 14:28:01 Pipe created: &{0xc42000e058 0xc42000e070} &{0xc42000e068 0xc42000e060}
2018/06/04 14:28:01 Pipe accepted: &{0xc42000e058 0xc42000e070} &{0xc42000e068 0xc42000e060}

// BUG: never reached
log.Printf("SUCCESS! Connected to server: %v", serverConn)

```
