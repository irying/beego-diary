package etcd

import (
	"time"
	"net"
	"fmt"
	etcdv3 "github.com/coreos/etcd/clientv3"
	"strings"
	"context"
)

var Prefix = "etcd_naming"
var Deregister = make(chan struct{})

func Register(name, host, port string, target string, interval time.Duration, ttl int) error {
	serviceValue := net.JoinHostPort(host, port)
	serviceKey := fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)

	var err error
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints:strings.Split(target,","),
	})

	if err != nil {
		return fmt.Errorf("grpc create etcd client failed: %s", err.Error())
	}

	// 第1步
	resp, err := client.Grant(context.TODO(), int64(ttl))

	if err != nil {
		return fmt.Errorf("grpc create etcd lease failed: %s", err.Error())
	}

	// 第2步
	if _, err := client.Put(context.TODO(), serviceKey, serviceValue, etcdv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("set service '%s' with ttl to etcd3 failed: %s", name, err.Error())
	}

	//  第3步
	if _, err := client.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf("refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
	}

	go func() {
		<-Deregister
		client.Delete(context.Background(), serviceKey)
		// todo 存疑
		Deregister <- struct{}{}
	}()

	return nil
}

func UnRegister()  {
	Deregister <- struct{}{}
	<-Deregister
}


