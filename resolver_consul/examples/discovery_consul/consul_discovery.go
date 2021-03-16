package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"log"

	_ "easy-grpc/resolver_consul"

	pbdata "myprotobuf/pb" //自己生成的pb文件包
	pbfunc "myprotobuf/pbheader"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"

	"net/http"
	_ "net/http/pprof"
)

func init() {
	go http.ListenAndServe(":9090", nil)
}

func main() {
	client, err := grpc.Dial(
		"consul://token/192.168.1.67:8500?fabio-demo",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	data, err := proto.Marshal(&pbdata.HelloRequest{
		Name: "helloRequest",
	})

	if err != nil {
		log.Println(err)
	}

	in := &pbfunc.xxxxxbpb{
		Cmd:    2200,
		Userid: 100011,
		Seq:    1111111,
		Key:    "xxxxxxx",
		Buf:    data,
	}

	defer func(t time.Time) {
		fmt.Println("Total CostTime:", time.Since(t))
	}(time.Now())

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1; j++ {
				//t := time.Now()

				r, err := pbfunc.NewGrpcfuncClient(client).GrpcDataFunc(context.Background(), in)
				if err != nil {
					log.Println("Call:", err)
					continue
				}

				//log.Println(r.String())

				helloReply := pbdata.HelloReply{}
				if err := proto.Unmarshal(r.Buf, &helloReply); err != nil {
					log.Println(err)
				}
				//log.Println(helloReply)
				//fmt.Println("CostTime:", time.Since(t))
				//time.Sleep(time.Second)
			}
		}()

		//log.Println(helloReply.String())
		//time.Sleep(1 * time.Second)
	}
	wg.Wait()
	log.Println("done")

}
