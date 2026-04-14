# 自动化运维部署平台

自动化运维部署平台，包含 Go 后端和 Vue 前端。

## 项目结构

```
auto-deploy/
├── backend/          # Go 后端 (Gin + GORM + MySQL)
├── frontend/         # Vue 3 前端 (Element Plus)
├── docker-compose.yml
└── .github/workflows/
    ├── backend.yml   # 后端 CI/CD
    └── frontend.yml  # 前端 CI/CD
```

## 快速开始

### 后端

```bash
cd backend
cp .env.example .env
# 编辑 .env 配置数据库等信息
go mod tidy
go run main.go
```

### 前端

```bash
cd frontend
npm install
npm run serve
```

### Docker 部署

```bash
docker-compose pull
docker-compose up -d
```

## CI/CD

- 推送到 `master` 分支时自动构建 Docker 镜像并推送到阿里云 ACR
- 后端和前端通过路径过滤独立触发
- 需要在 GitHub Secrets 中配置：
  - `ALIYUN_USERNAME`：阿里云 ACR 用户名
  - `ALIYUN_PASSWORD`：阿里云 ACR 密码

## 镜像地址

- 后端：`crpi-4z8h6zhc9e28drvr.ap-southeast-1.personal.cr.aliyuncs.com/images_d/auto-deploy-backend`
- 前端：`crpi-4z8h6zhc9e28drvr.ap-southeast-1.personal.cr.aliyuncs.com/images_d/auto-deploy-frontend`
