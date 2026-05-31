package etcd

import (
	"context"
	"fmt"
	"log"
	"sync"
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

// ============================ 连接池化 ServiceResolver ============================

// ServiceResolver 持久的服务发现器：复用 ETCD 长连接，round-robin 负载均衡
type ServiceResolver struct {
	cli        *clientv3.Client
	serviceName string
	endpoints  []string
	nextIdx    int
	mu         sync.Mutex
	refreshInterval time.Duration
}

// NewServiceResolver 创建持久的服务发现器
func NewServiceResolver(etcdEndpoints []string, serviceName string) (*ServiceResolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd connect: %w", err)
	}

	sr := &ServiceResolver{
		cli:             cli,
		serviceName:     serviceName,
		refreshInterval: 15 * time.Second,
	}

	// 首次加载
	if err := sr.refresh(context.Background()); err != nil {
		cli.Close()
		return nil, fmt.Errorf("initial refresh: %w", err)
	}

	// 后台定期刷新
	go sr.refreshLoop()

	log.Printf("[Etcd] ServiceResolver ready for %s (%d endpoints)", serviceName, len(sr.endpoints))
	return sr, nil
}

// refresh 从 ETCD 获取最新端点列表
func (sr *ServiceResolver) refresh(ctx context.Context) error {
	manager, err := endpoints.NewManager(sr.cli, sr.serviceName)
	if err != nil {
		return err
	}
	eps, err := manager.List(ctx)
	if err != nil {
		return err
	}

	addrs := make([]string, 0, len(eps))
	for _, ep := range eps {
		addrs = append(addrs, ep.Addr)
	}

	sr.mu.Lock()
	sr.endpoints = addrs
	// 重置索引避免越界
	if sr.nextIdx >= len(addrs) {
		sr.nextIdx = 0
	}
	sr.mu.Unlock()
	return nil
}

// refreshLoop 定期刷新端点列表
func (sr *ServiceResolver) refreshLoop() {
	ticker := time.NewTicker(sr.refreshInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := sr.refresh(context.Background()); err != nil {
			log.Printf("[Etcd] refresh %s failed: %v", sr.serviceName, err)
		}
	}
}

// GetEndpoint 获取一个端点地址（round-robin）
func (sr *ServiceResolver) GetEndpoint() (string, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if len(sr.endpoints) == 0 {
		return "", fmt.Errorf("no endpoints for %s", sr.serviceName)
	}

	addr := sr.endpoints[sr.nextIdx]
	sr.nextIdx = (sr.nextIdx + 1) % len(sr.endpoints)
	return addr, nil
}

// GetAllEndpoints 获取所有端点（用于首次连接之前的轮询）
func (sr *ServiceResolver) GetAllEndpoints() []string {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	cp := make([]string, len(sr.endpoints))
	copy(cp, sr.endpoints)
	return cp
}

// Close 关闭连接
func (sr *ServiceResolver) Close() error {
	return sr.cli.Close()
}
