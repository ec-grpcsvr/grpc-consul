package resolver_consul

import (
	"fmt"
	"time"

	"strings"

	consul "github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc/resolver"
)

type consulResolver struct {
	cc           resolver.ClientConn
	Addr         string
	Token        string
	ServiceName  string
	consulClient *consul.Client
	lastIndex    uint64
	addrs        []string
	quit         chan bool
	cancelFunc   context.CancelFunc
}

func (this *consulResolver) init() error {
	conf := &consul.Config{
		Scheme:  "http",
		Address: this.Addr,
	}

	client, err := consul.NewClient(conf)
	if err != nil {
		return fmt.Errorf("wonaming: creat consul error: %v", err)
	}

	this.consulClient = client

	if this.addrs == nil {
		serviceEntry, _ := this.queryConsul(nil)
		this.updateAddrs(serviceEntry)
		if len(this.addrs) > 0 {
			this.cc.NewAddress(this.covertResolverAddress())
		}
	}

	return nil
}

func (this *consulResolver) watcher() {
	ctx, cancel := context.WithCancel(context.Background())
	this.cancelFunc = cancel

	for {
		opt := &consul.QueryOptions{AllowStale: false, WaitIndex: this.lastIndex}
		serviceEntry, err := this.queryConsul(opt.WithContext(ctx))
		if err != nil {
			if strings.Contains(err.Error(), context.Canceled.Error()) {
				break
			}

			time.Sleep(50 * time.Microsecond)
			continue
		}
		this.updateAddrs(serviceEntry)
		this.cc.NewAddress(this.covertResolverAddress())
	}
	this.quit <- true
}

func (this *consulResolver) ResolveNow(opt resolver.ResolveNowOption) {
	//fmt.Println("ResolveNow")
}

func (this *consulResolver) Close() {
	if this.cancelFunc != nil {
		this.cancelFunc()
		<-this.quit
	}
}

func (this *consulResolver) queryConsul(q *consul.QueryOptions) ([]*consul.ServiceEntry, error) {
	serviceEntry, meta, err := this.consulClient.Health().Service(this.ServiceName, "", true, q)
	if err != nil {
		return nil, err
	}

	this.lastIndex = meta.LastIndex

	return serviceEntry, nil
}

func (this *consulResolver) updateAddrs(serviceEntry []*consul.ServiceEntry) {
	//data, _ := json.Marshal(serviceEntry)
	//fmt.Println("ServiceEntry:", string(data))
	addrs := []string{}
	for _, se := range serviceEntry {
		if se.Checks.AggregatedStatus() == consul.HealthPassing {
			addrs = append(addrs, se.Service.Address+":"+fmt.Sprint(se.Service.Port))
		}
	}

	//fmt.Println("addrs:", addrs)

	this.addrs = addrs
}

func (this *consulResolver) covertResolverAddress() []resolver.Address {
	addrs := []resolver.Address{}
	for _, addr := range this.addrs {
		addrs = append(addrs, resolver.Address{
			Addr: addr,
		})
	}
	return addrs
}
