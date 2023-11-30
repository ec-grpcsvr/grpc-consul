package resolver_consul

import (
	"strings"

	"fmt"
	"google.golang.org/grpc/resolver"
)

// consul://token/ip:port?serviceName
const name = "consul"

type consulBuilder struct {
}

func init() {
	resolver.Register(NewBuilder())
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func (consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	token := target.Endpoint()

	arr := strings.Split(token, "?")
	if len(arr) != 2 {
		return nil, fmt.Errorf("target error, consul://token/ip:port?servieName")
	}

	addr, serviceName := arr[0], arr[1]
	fmt.Println("build:", arr)
	cr := &consulResolver{
		Addr:        addr,
		Token:       token,
		ServiceName: serviceName,
		cc:          cc,
		quit:        make(chan bool, 1),
	}

	if err := cr.init(); err != nil {
		return nil, err
	}

	go cr.watcher()

	return cr, nil
}

func (consulBuilder) Scheme() string {
	return name
}
