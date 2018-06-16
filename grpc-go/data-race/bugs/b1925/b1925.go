package main

import (
	"context"
	"fmt"
	"reflect"
	"google.golang.org/grpc/metadata"
)

func TestAppendToOutgoingContext() {
	var ctx = context.Background()

	var ctx1 context.Context
	var ctx2 context.Context

	for i := 0; i < 100; i = i + 2 {
		ctx1 = metadata.AppendToOutgoingContext(ctx, fmt.Sprintf("k1-%d", i), fmt.Sprintf("v1-%d", i+1))
		ctx2 = metadata.AppendToOutgoingContext(ctx, fmt.Sprintf("k2-%d", i), fmt.Sprintf("v2-%d", i+1))

		md1, _ := metadata.FromOutgoingContext(ctx1)
		md2, _ := metadata.FromOutgoingContext(ctx2)

		if reflect.DeepEqual(md1, md2) {
			panic(fmt.Sprintf("md1:(%v) and md2:(%v) should not be equal", md1, md2))
		}

		ctx = ctx1
	}
}

func main()  {
	TestAppendToOutgoingContext()
}
