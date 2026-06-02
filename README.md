# PXX 分布式校园生活服务平台 (Distributed Campus Multi-Service Platform)

PXX 是一个基于 **微服务架构** 的高性能、高可用校园一站式生活服务平台。项目涵盖了电商、社交、租赁、拼团、订阅、以物换物、任务众包及校园仲裁等多个核心业务域。

本项目已从早期的单体架构演进为基于 **gRPC** 的分布式架构，深度实践了分布式事务 (TCC)、服务发现 (Etcd)、分布式共识 (Raft) 及多级缓存等核心技术。

## 🏗️ 系统架构

系统采用微服务分层设计，各服务通过 gRPC 进行高性能通信：

-   **API Gateway (BFF)**: 统一入口，负责路由转发、鉴权、限流、协议转换 (HTTP/JSON to gRPC) 及跨服务数据聚合。
-   **Account Service**: 处理用户账号、身份认证、权限管理及余额变动（TCC 事务支持）。
-   **Item Service**: 核心业务域，管理商品 (Product)、租赁 (Rental)、拼团 (GroupBuy)、订阅 (Subscription) 及以物换物 (Barter)。
-   **Trade Service**: 负责订单生命周期、支付逻辑及分布式事务协调。
-   **Feed Service**: 社交互动核心，包括动态墙 (Posts)、即时消息 (Websocket)、任务众包 (Tasks) 及争议仲裁 (Tribunal/Dispute)。
-   **Raft-Scheduler**: 基于 Hashicorp Raft 实现的高可用分布式调度器，用于处理定时任务、超时订单清理等。

## 🛠️ 技术栈

### 核心框架
-   **后端**: [Go](https://go.dev/) (1.25+) + [gRPC](https://grpc.io/) + [Protobuf](https://protobuf.dev/)
-   **Web 框架**: [Gin](https://gin-gonic.com/) (Gateway 层)
-   **前端**: [Vue 3](https://vuejs.org/) + [Vite](https://vitejs.dev/) + [Vant 4](https://vant-ui.github.io/vant/) + [Pinia](https://pinia.vuejs.org/)
-   **ORM**: [GORM](https://gorm.io/)

### 基础设施 & 中间件
-   **服务发现**: [Etcd](https://etcd.io/) (Service Registry & Discovery)
-   **消息队列**: [Kafka](https://kafka.apache.org/) (解耦、削峰、延迟任务)
-   **缓存**: [Redis](https://redis.io/) (L2) + [Local Cache](https://github.com/hashicorp/golang-lru) (L1) + [Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter) (防止穿透)
-   **分布式共识**: [Hashicorp Raft](https://github.com/hashicorp/raft) (用于 Scheduler 高可用)
-   **存储**: [MySQL 8.0](https://www.mysql.com/) (主库) + [MinIO](https://min.io/) (对象存储)
-   **网关/代理**: [Nginx](https://www.nginx.com/) + [Docker Compose](https://docs.docker.com/compose/)

## 💎 核心特性

1.  **分布式事务 (TCC)**: 在 `Trade` 与 `Account` 服务间采用 Try-Confirm-Cancel 模式，确保下单扣费操作的最终一致性，并实现防悬挂与幂等控制。
2.  **三层缓存架构**:
    -   `L1`: 进程内 LRU 缓存，极致响应速度。
    -   `L2`: Redis 分布式缓存，支持数据共享。
    -   `L3`: 布隆过滤器 + `singleflight` 防击穿，保护底层 MySQL。
3.  **高性能 Feed 流**: 采用异步推送模式，结合 Kafka 实现动态发布后的图片处理与多端分发。
4.  **高可用调度**: `Raft-Scheduler` 确保在集群环境下，定时任务（如自动确认收货）仅由一个 Leader 节点触发，具备自动选主与故障转移能力。
5.  **即时通讯**: 基于 Websocket 建立的双工通信频道，支持校园私信与实时系统通知。

## 📂 项目结构

```bash
pxx/
├── server/                     # 后端微服务集群
│   ├── gateway/                # API 网关 (Gin + gRPC Client)
│   ├── account/                # 用户服务 (gRPC Server)
│   ├── item/                   # 商品与业务逻辑服务 (gRPC Server)
│   ├── trade/                  # 交易与订单服务 (gRPC Server)
│   ├── feed/                   # 社交与 Feed 流服务 (gRPC Server)
│   ├── raft-scheduler/         # 分布式调度器 (Raft Consensus)
│   ├── idl/                    # gRPC 接口定义 (Protobuf)
│   └── internal/               # 共享业务逻辑与核心库
├── frontend/                   # Vue 3 移动端应用
├── nginx/                      # 反向代理与前端静态资源托管
├── sql/                        # 各服务的数据库初始化脚本
└── docker-compose.yml          # 全容器化编排 (支持一键拉起全栈环境)
```

## 🚀 快速启动

### 1. 环境准备
确保已安装 [Docker](https://www.docker.com/) 和 [Docker Compose](https://docs.docker.com/compose/)。

### 2. 一键拉起
```bash
docker-compose up -d --build
```
该命令会自动构建所有微服务镜像，并拉起 MySQL, Redis, Kafka, Etcd, MinIO 等中间件。

### 3. 访问地址
-   **前端界面**: [http://localhost](http://localhost)
-   **API 入口**: `http://localhost/api/v1`
-   **MinIO 控制台**: `http://localhost:9001`

## ⚖️ 业务域

-   **论坛社区**: 校园动态、求购悬赏、种草分享。
-   **精细化交易**: 支持商品、租赁、订阅、以物换物多种交易模式。
-   **拼团秒杀**: 宿舍组团，阶梯降价，Kafka 削峰处理并发请求。
-   **众包任务**: 校园跑腿、技能互助。
-   **校园仲裁**: 针对交易纠纷的陪审团制度（Tribunal）。

---
*PXX - 让校园生活更简单、更高效。*
