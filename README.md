# 智能改卷助手

一个基于Go后端和Vue前端的智能试卷批改系统，支持OCR识别和AI评分功能。

## 项目架构

- **后端**: Go + Fiber框架 + GORM + MySQL
- **前端**: Vue 3 + Element Plus + Vite
- **数据库**: MySQL 8.0
- **文件存储**: 本地文件系统
- **AI服务**: 百度OCR + DeepSeek API

## 功能特性

### 家长端功能

- ✅ 用户注册和登录
- ✅ 试卷图片上传
- ✅ OCR文字识别
- ✅ AI智能评分
- ✅ 批改结果展示
- ✅ 历史记录管理
- ✅ 错题分析

### 系统特性

- 🔐 JWT认证保护
- 📱 响应式设计
- 🚀 高性能异步处理
- 📊 详细的数据统计
- 🔍 搜索和筛选功能

## 项目结构

```
F:\Auto-grad-parents\
├── backend\              # Go后端
│   ├── main.go          # 主入口
│   ├── .env             # 环境配置
│   ├── internal\
│   │   ├── api\         # API路由和处理器
│   │   ├── config\      # 配置管理
│   │   ├── db\          # 数据库连接
│   │   ├── models\      # 数据模型
│   │   ├── services\    # 业务服务
│   │   └── storage\     # 文件存储
│   └── uploads\         # 文件上传目录
├── frontend\             # Vue前端
│   ├── src\
│   │   ├── views\       # 页面组件
│   │   ├── components\  # 通用组件
│   │   ├── router\      # 路由配置
│   │   └── api\         # API调用
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
├── database\             # 数据库相关
│   └── schema.sql      # 数据库表结构
└── docs\                # 项目文档
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 2. 数据库设置

```sql
CREATE DATABASE auto_grad_web CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 后端启动

```bash
cd backend
go mod tidy
go build -o app.exe
./app.exe
```

### 4. 前端启动

```bash
cd frontend
npm install
npm run dev
```

### 5. 访问系统

- 前端地址: http://localhost:5173
- 后端API: http://localhost:3000
- API文档: http://localhost:3000/health

## API接口

### 认证相关

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取用户信息

### 批改相关

- `POST /api/upload` - 文件上传
- `POST /api/grading` - 创建批改记录
- `GET /api/grading` - 获取批改列表
- `GET /api/grading/:id` - 获取批改详情
- `POST /api/grading/:id/process` - 处理批改
- `DELETE /api/grading/:id` - 删除记录

## 配置说明

### 后端配置 (.env)

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=auto_grad_web

# 服务器配置
SERVER_PORT=3000
JWT_SECRET=your-jwt-secret-key

# API密钥配置
BAIDU_API_KEY=your_baidu_api_key
BAIDU_SECRET_KEY=your_baidu_secret_key
DEEPSEEK_API_KEY=your_deepseek_api_key

# 文件存储配置
UPLOAD_PATH=./uploads
```

## 部署说明

### 生产环境部署

1. 编译前端: `npm run build`
2. 编译后端: `go build -o server`
3. 配置Nginx反向代理
4. 使用PM2或systemd管理后端服务

### Docker部署（可选）

项目支持Docker容器化部署，详见docker-compose.yml

## 开发说明

### 添加新功能

1. 后端: 在`internal/services`中添加业务逻辑
2. 前端: 在`src/views`中添加页面组件
3. API: 在`internal/api/handlers.go`中添加接口

### 数据库迁移

使用GORM的AutoMigrate功能自动创建数据表

## 注意事项

1. **API密钥**: 请确保正确配置百度OCR和DeepSeek API密钥
2. **文件权限**: 确保uploads目录有读写权限
3. **数据库**: 确保MySQL服务正在运行并且字符集为utf8mb4
4. **端口占用**: 确保3000和5173端口未被占用

## 故障排除

### 常见问题

1. **数据库连接失败**: 检查数据库配置和服务状态
2. **API调用失败**: 检查网络连接和API密钥
3. **文件上传失败**: 检查文件权限和磁盘空间

## 技术栈详情

### 后端技术

- **框架**: Fiber v1.14.6 (高性能Web框架)
- **ORM**: GORM (Go语言ORM库)
- **数据库**: MySQL 8.0
- **认证**: JWT (JSON Web Token)
- **文件处理**: Multipart Form Data

### 前端技术

- **框架**: Vue 3.5.26 (响应式前端框架)
- **UI库**: Element Plus 2.13.1 (Vue 3组件库)
- **构建工具**: Vite 7.3.1 (现代前端构建工具)
- **路由**: Vue Router 4.6.4 (官方路由管理)
- **HTTP客户端**: Axios 1.13.2 (Promise风格HTTP库)

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进项目。

## 联系方式

如有问题，请通过Issue联系我们。
