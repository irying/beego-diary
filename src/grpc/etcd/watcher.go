package etcd

import (
	etcdv3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
	"fmt"
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)


type watcher struct {
	re *resolver
	client etcdv3.Client
	isInitialized bool
}

func (w *watcher) Close()  {

}

func (w *watcher) Next() ([]*naming.Update, error)  {
	prefix := fmt.Sprintf("/%s/%s/", Prefix, w.re.serviceName)
// todo 可以分开的
	if !w.isInitialized {
		resp, err := w.client.Get(context.Background(), prefix, etcdv3.WithPrefix())

		if err == nil {
			w.isInitialized = true
			addrs := extractAddrs(resp)
			if l := len(addrs); l != 0 {
				updates := make([]*naming.Update, l)
				for i := range addrs {
					updates[i] = &naming.Update{Op:naming.Add, Addr:addrs[i]}
				}

				return updates, nil
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rch := w.client.Watch(ctx, prefix, etcdv3.WithPrefix())
	defer cancel()

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				return []*naming.Update{{Op:naming.Add, Addr:string(ev.Kv.Value)}}, nil

			case mvccpb.DELETE:
				return []*naming.Update{{Op:naming.Delete, Addr:string(ev.Kv.Value)}}, nil

			}
		}
	}

	return nil, nil
}

func extractAddrs(resp *etcdv3.GetResponse) []string  {
	var addrs []string

	if resp == nil || resp.Kvs == nil {
		return addrs
	}
// todo 存疑
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil{
			addrs = append(addrs, string(v))
		}
	}

	return addrs

}