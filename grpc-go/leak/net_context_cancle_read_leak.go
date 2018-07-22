package main

import (
	"fmt"
	"time"
	"golang.org/x/net/context"
	"net"
	"runtime"
	"os"
)


func ReadConn(ctx context.Context, rawConn net.Conn) {
	tmp := make([]byte, 256)
	errCh := make(chan error, 1)
	go func() {
		_, err := rawConn.Read(tmp)
		errCh <- err
	}()
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("Failed to read the server name from the server %v", err)
		}
	case <-ctx.Done():
	}
	fmt.Printf("read length %d\n", len(tmp))

}

func ReadConnLeak(ctx context.Context, rawConn net.Conn) {
	tmp := make([]byte, 256)
	_, err := rawConn.Read(tmp)
	if err != nil {
		// handle error
	}
	fmt.Printf("read length %d\n", len(tmp))
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1)
	defer cancel()
	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		fmt.Println("dial error:", err)
	}
	defer conn.Close()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Number of goroutine: ", runtime.NumGoroutine())
		buf := make([]byte, 1<<11)
		runtime.Stack(buf, true)
		fmt.Printf("%s", buf)
		os.Exit(1)
	}()

	time.Sleep(2 * time.Second)
	//ReadConnLeak(ctx, conn)
	ReadConn(ctx, conn)
	fmt.Println("end!")
}
