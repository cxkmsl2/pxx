# PXX 校园二手交易平台

一站式校园 O2O + C2C 二手交易平台，涵盖电商、论坛、拼团、租赁、以物换物等 6 大核心业务域。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + Vant 4 + Pinia |
| 后端 | Golang + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 + 进程内存缓存 (L2) + 布隆过滤器 (L1) |
| 消息队列 | Kafka (削峰 / 异步 / 延迟) |
| 反向代理 | Nginx |
| 部署 | Docker Compose |

## 三层缓存架构

```
请求 → L1 布隆过滤器 → L2 本地缓存 → L3 Redis → MySQL
         ↓ 穿透拦截      ↓ 零网络开销    ↓ singleflight 防击穿
```

- **防击穿**: `golang.org/x/sync/singleflight` —— 同一 key 只放行一个请求查 DB
- **防穿透**: 布隆过滤器 + Redis 缓存空值 (Null Marker, TTL 60s)
- **防雪崩**: 所有 Redis 缓存 TTL 附加随机偏移量 (`baseTTL + rand(0, baseTTL/4)`)

## Kafka 消息流

| Topic | 用途 |
|-------|------|
| `pxx.post.created` | 帖子发布后异步图片压缩 / 内容审核 |
| `pxx.order.created` | 订单创建削峰，异步写入 |
| `pxx.order.expired` | 15 分钟未支付取消订单 (延迟消费) |
| `pxx.groupbuy.flash` | 拼团抢购削峰，保护 MySQL |

## 快速启动

```bash
# 1. 克隆项目
cd pxx

# 2. 一键启动所有服务 (MySQL + Redis + Zookeeper + Kafka + Server + Nginx)
docker-compose up -d

# 3. 访问
# 前端: http://localhost:80
# API:  http://localhost:8080/api/v1

# 4. 停止
docker-compose down
```

## 业务模块

1. **社区论坛** — 求购悬赏 + 好物种草，帖子可关联商品链接，内容到交易闭环
2. **精细化电商** — 按分类 / 校区 / 楼栋筛选，最小化线下交易距离
3. **以物换物** — 上传闲置物品，平台撮合意向，线下当面交割
4. **校园本地服务** — 校内商家后台，毕业季上门收书
5. **宿舍拼团** — 满人成团 / 阶梯降价，团长可获免单或积分
6. **物品租赁** — 正装 / 相机 / 游戏机日租周租，校园卡认证替代高额押金

## 后续 AI 演进

- **物品智能识别 (CV)**: 拍照自动识别商品类别，自动填充标签
- **个性化推荐 (RecSys)**: 基于专业 / 年级 / 浏览行为的协同过滤推荐

## 项目结构

```
pxx/
├── server/                        # Go 后端
│   ├── cmd/main.go                # 入口
│   ├── internal/
│   │   ├── config/config.go       # 环境变量配置
│   │   ├── model/                 # GORM 数据模型 (11 张表)
│   │   ├── handler/               # HTTP 处理器 (7 个模块)
│   │   ├── service/               # 业务逻辑层
│   │   ├── cache/                 # 三层缓存核心
│   │   │   ├── redis.go           #   Redis 连接池
│   │   │   ├── local.go           #   进程内 LRU 缓存
│   │   │   ├── bloom.go           #   布隆过滤器
│   │   │   └── cache_manager.go   #   缓存门面 + singleflight
│   │   ├── middleware/            # JWT Auth / CORS / 限流
│   │   ├── mq/kafka.go           # Kafka Producer + Consumer
│   │   └── router/router.go      # 路由注册
│   ├── pkg/
│   │   ├── response/response.go   # 统一 API 响应
│   │   └── utils/utils.go        # 工具函数
│   └── Dockerfile
├── frontend/                      # Vue3 前端
│   ├── src/
│   │   ├── pages/                 # 8 个页面组件
│   │   ├── router/index.ts        # Hash 路由
│   │   ├── store/index.ts         # Pinia 状态管理
│   │   └── api/index.ts           # Axios 封装
│   └── Dockerfile
├── nginx/
│   ├── nginx.conf                 # 反向代理 + CDN 缓存配置
│   └── Dockerfile
├── sql/
│   └── init.sql                   # 建表 + 测试数据
├── docker-compose.yml             # 一键部署编排
└── README.md
```
