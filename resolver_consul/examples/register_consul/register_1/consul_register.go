package main

import (
	"log"
	"math/rand"
	"syscall"
	"time"

	"git.workec.grpc/grpcframe/registry"
	"git.workec.grpc/grpcframe/registry/consul"
	"github.com/judwhite/go-svc/svc"
)

var reg registry.Registry

var s1 = &registry.Service{
	Name: "fabio-demo",
	Nodes: []*registry.Node{
		&registry.Node{
			Id:       "fabio-demo[10.0.108.91:50051]",
			Address:  "10.0.108.91",
			Port:     50051,
			Metadata: registry.NewMetaData("50051", "tcp"),
		},
	},
}

//var s2 = &registry.Service{
//	Name: "fabio-demo",
//	Nodes: []*registry.Node{
//		&registry.Node{
//			Id:       "fabio-demo[10.0.108.92:50051]",
//			Address:  "10.0.108.91",
//			Port:     50052,
//			Metadata: registry.NewMetaData("50052", "tcp"),
//		},
//	},
//}

func main() {

	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

type program struct {
}

func (p *program) Init(env svc.Environment) error {
	reg = consul.NewRegistry(
		registry.Token("xxxxxx"),
		registry.Addrs("192.168.1.67:8500"),
	)

	return nil
}

func (p *program) Start() error {
	go func() {
		for {
			n := rand.Intn(10)
			if err := reg.Register(s1, registry.RegisterTTL(10*time.Second)); err != nil {
				log.Println("s1 register error:", err)
			}
			//
			//if err := reg.Register(s2, registry.RegisterTTL(10*time.Second)); err != nil {
			//	log.Println("s1 register error:", err)
			//}

			//log.Println("sleep:", n)
			time.Sleep(time.Second * time.Duration(n))
		}
	}()
	return nil
}

func (p *program) Stop() error {
	reg.Deregister(s1)
	//reg.Deregister(s2)
	return nil
}
