package etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

// ServiceInfo 微服务注册信息
type ServiceInfo struct {
	Name    string
	Address string
	TTL     int64 // 秒
}

// Registry 服务注册器
type Registry struct {
	cli     *clientv3.Client
	manager endpoints.Manager
	leaseID clientv3.LeaseID
	info    ServiceInfo
}

// NewRegistry 创建注册器并建立连接
func NewRegistry(etcdEndpoints []string, info ServiceInfo) (*Registry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd connect failed: %w", err)
	}

	manager, err := endpoints.NewManager(cli, info.Name)
	if err != nil {
		cli.Close()
		return nil, fmt.Errorf("endpoints manager failed: %w", err)
	}

	return &Registry{cli: cli, manager: manager, info: info}, nil
}

// Register 注册服务到 ETCD，带上 TTL 租约
func (r *Registry) Register(ctx context.Context) error {
	lease := clientv3.NewLease(r.cli)
	grantResp, err := lease.Grant(ctx, r.info.TTL)
	if err != nil {
		return fmt.Errorf("lease grant failed: %w", err)
	}
	r.leaseID = grantResp.ID

	key := r.info.Name + "/" + r.info.Address
	err = r.manager.AddEndpoint(ctx, key,
		endpoints.Endpoint{Addr: r.info.Address},
		clientv3.WithLease(r.leaseID),
	)
	if err != nil {
		return fmt.Errorf("add endpoint failed: %w", err)
	}

	// 自动续约
	ch, err := lease.KeepAlive(ctx, r.leaseID)
	if err != nil {
		return fmt.Errorf("keepalive failed: %w", err)
	}
	go func() {
		for range ch {
		}
		log.Printf("[Etcd] lease keepalive stopped for %s/%s", r.info.Name, r.info.Address)
	}()

	log.Printf("[Etcd] registered %s -> %s (TTL=%ds)", r.info.Name, r.info.Address, r.info.TTL)
	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister(ctx context.Context) error {
	if r.manager != nil {
		key := r.info.Name + "/" + r.info.Address
		return r.manager.DeleteEndpoint(ctx, key)
	}
	return nil
}

// Close 关闭连接
func (r *Registry) Close() error {
	if r.leaseID != 0 {
		r.cli.Revoke(context.Background(), r.leaseID)
	}
	return r.cli.Close()
}

// Resolver 服务发现：从 ETCD 获取指定服务的地址列表
type Resolver struct {
	cli  *clientv3.Client
	name string
}

func NewResolver(etcdEndpoints []string, serviceName string) (*Resolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Resolver{cli: cli, name: serviceName}, nil
}

func (r *Resolver) ListEndpoints(ctx context.Context) ([]string, error) {
	manager, err := endpoints.NewManager(r.cli, r.name)
	if err != nil {
		return nil, err
	}
	eps, err := manager.List(ctx)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(eps))
	for _, ep := range eps {
		addrs = append(addrs, ep.Addr)
	}
	return addrs, nil
}

func (r *Resolver) Close() error {
	return r.cli.Close()
}
