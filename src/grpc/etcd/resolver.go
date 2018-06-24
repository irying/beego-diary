package etcd

import (
	"google.golang.org/grpc/naming"
	"errors"
	etcdv3 "github.com/coreos/etcd/clientv3"
	"strings"
	"fmt"
)

type resolver struct {
	serviceName string
}

func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName:serviceName}
}

func (r *resolver) Resolve(target string) (naming.Watcher, error) {
	if r.serviceName == "" {
		return nil, errors.New("no service name provided")
	}

	client, err := etcdv3.New(etcdv3.Config{
		Endpoints:strings.Split(target, ","),
	})

	if err != nil {
		return nil, fmt.Errorf("grpc create etcd client failed: %s", err.Error())
	}

	return &watcher{re:r, client: *client}, nil
}