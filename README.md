# Offer Hub

Offer Hub 是一个面向技术面试准备的前后端分离项目，提供题库浏览、题目检索、热门榜单、用户认证、评论互动和学习状态标记等功能。

## 已实现功能

- 题库分类、题目筛选与分页、热门题目、题目详情
- 用户注册、登录、退出和 JWT 鉴权
- 游客内容预览，登录后查看完整题目与解析
- 评论、回复、编辑、删除和点赞
- 题目点赞及掌握状态标记
- Redis 登录状态校验和认证接口限流
- 桌面端与移动端响应式页面

## 技术栈

### 前端

- React 18、TypeScript、Vite
- Tailwind CSS、Radix UI、Lucide Icons
- React Router、TanStack Query、Axios
- React Hook Form、Zod、React Markdown

### 后端

- Go 1.22、Gin
- GORM、MySQL 8
- MongoDB 官方 Go 驱动
- Redis、JWT、bcrypt

数据存储职责：MySQL 保存用户信息，MongoDB 保存题库、题目、评论和互动数据，Redis 保存登录令牌状态及限流计数。

## 环境要求

- Go 1.22+
- Node.js 20.19+ 或 22.12+
- Docker、Docker Compose
- Bash、curl、python3（运行接口回归脚本时需要）
- pre-commit（可选，用于提交前检查）

## 快速开始

### 1. 启动本地数据库

```bash
cd backend
docker compose up -d
```

默认端口如下：

| 服务 | 本地端口 |
| --- | --- |
| MySQL | `3307` |
| MongoDB | `27017` |
| Redis | `6379` |

### 2. 初始化测试数据

首次启动时创建 MySQL 用户表：

```bash
docker exec -i mysql-dev mysql -uroot -p1234 < testdata/init_mysql_schema.sql
```

写入 MongoDB 题库、系列和题目数据：

```bash
docker exec -i mongodb-dev mongosh < testdata/insert_bank_and_series.js
docker exec -i mongodb-dev mongosh < testdata/insert_question.js
```

MongoDB 脚本会先清空对应的本地测试集合，再重新写入数据。请勿在生产数据库上执行。

### 3. 启动后端

在 `backend` 目录执行：

```bash
go mod download
go run ./src
```

后端默认监听 `http://127.0.0.1:8180`。可通过健康检查确认 MySQL 和 MongoDB 的连接状态：

```bash
curl http://127.0.0.1:8180/health
```

默认读取 `backend/config/config-test.toml`。如需切换环境，可设置 `APP_ENV`，后端将读取 `config/config-<环境>.toml`。包含真实凭据的配置文件不要提交到仓库。

### 4. 启动前端

新开终端，在 `frontend` 目录执行：

```bash
npm ci
npm run dev
```

访问 `http://127.0.0.1:3001`。开发环境通过 Vite 代理访问本地后端。

## 常用命令

### 后端

```bash
cd backend
go test ./...
make fmt
make vet
./testsh/question_guest.sh
```

`question_guest.sh` 要求后端已经启动，并且 MongoDB 测试数据已完成初始化。可通过 `BASE_URL` 覆盖默认后端地址。

### 前端

```bash
cd frontend
npm run build
npm run lint
npm run format:check
```

### 提交前检查

```bash
pre-commit install
pre-commit run --all-files
```

提交信息遵循 Conventional Commits，例如：

```text
feat(frontend): 完成首页与热门题目页面
fix(auth): 修复登录状态恢复异常
```

## 目录结构

```text
offer-hub/
├── backend/
│   ├── config/          # TOML 配置
│   ├── src/
│   │   ├── config/      # 配置读取
│   │   ├── model/       # 请求与响应结构
│   │   ├── data/        # MySQL、MongoDB、Redis 访问
│   │   ├── service/     # 业务逻辑
│   │   ├── ctrl/        # HTTP 控制器与中间件
│   │   └── router/      # 路由注册
│   ├── testdata/        # 本地表结构与 MongoDB 种子数据
│   └── testsh/          # HTTP 回归脚本
├── frontend/
│   └── src/
│       ├── components/  # 通用组件与业务组件
│       ├── pages/       # 页面
│       ├── services/    # API 请求封装
│       ├── hooks/       # React Query Hooks
│       ├── types/       # TypeScript 类型
│       └── lib/         # Axios 与通用工具
└── .pre-commit-config.yaml
```

## 主要接口

- `GET /health`：服务与数据库健康检查
- `/auth/*`：注册、登录、退出
- `/api/v1/question/*`：题库与题目查询
- `/api/v1/open/list_comments`：公开评论列表
- `/api/v1/comment/*`：评论写操作
- `/api/v1/interaction/*`：题目和评论点赞
- `/api/v1/safe/tag_question`：题目掌握状态

需要登录的接口统一使用 `Authorization: Bearer <token>`。

## License

[MIT](LICENSE)
