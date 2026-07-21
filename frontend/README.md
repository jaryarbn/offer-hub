# Offer Hub Frontend

Offer Hub 的 React 前端项目。

## 本地开发

```bash
npm install
npm run dev
```

开发服务器默认通过 Vite 代理访问 `http://127.0.0.1:8180` 的后端接口。
如需指定其他 API 地址，可复制 `.env.example` 为 `.env.local`，并设置
`VITE_API_BASE_URL`。

## 常用命令

```bash
npm run lint
npm run build
npm run preview
```

## 目录

```text
src/
├── components/     # 可复用组件和 Provider
├── pages/          # 路由页面
├── services/       # API 请求封装
├── hooks/          # 自定义 Hook
├── types/          # TypeScript 类型
├── lib/            # axios 实例、shadcn 工具
└── utils/          # 通用工具函数
```
