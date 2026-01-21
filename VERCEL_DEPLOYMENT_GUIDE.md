# 🚀 Vercel 部署指南

## 📋 部署步骤

### 1. Fork GitHub 仓库

1. 访问 https://github.com/NoccW/auto-grad
2. 点击右上角的 "Fork" 按钮
3. 选择你的 GitHub 账户进行 Fork

### 2. 连接 Vercel

1. 访问 [vercel.com](https://vercel.com)
2. 使用 GitHub 账号登录
3. 点击 "New Project"
4. 在 "Import Git Repository" 中选择你 fork 的 `auto-grad` 仓库
5. 点击 "Import"

### 3. 配置项目设置

Vercel 会自动检测到这是一个 Vue.js 项目，并进行以下配置：

#### Build Settings

- **Framework Preset**: Vue.js
- **Root Directory**: `frontend`
- **Build Command**: `npm run build`
- **Output Directory**: `dist`
- **Install Command**: `npm install`

#### Environment Variables (可选)

如果你有后端API地址，可以添加：

```
VITE_API_BASE_URL=https://your-backend-url.com
```

### 4. 部署项目

1. 点击 "Deploy" 按钮
2. 等待构建完成（通常需要1-2分钟）
3. 部署成功后会获得一个 `.vercel.app` 的域名

## 🎯 部署后配置

### 1. 访问应用

部署成功后，你可以通过以下地址访问：

- `https://your-project-name.vercel.app`

### 2. 测试功能

- 使用测试账号登录：用户名 `123123`，密码 `123123`
- 选择角色（家长端或教师端）
- 测试各项功能

### 3. 配置后端API（如果需要）

如果你有独立的后端服务，需要：

1. 在 Vercel 项目设置中添加环境变量
2. 更新 `VITE_API_BASE_URL` 为你的后端地址
3. 重新部署项目

## 🔧 故障排除

### 常见问题

#### 1. 构建失败

- 检查 `frontend/package.json` 中的依赖是否正确
- 确保所有依赖都已安装
- 查看 Vercel 的构建日志

#### 2. 路由问题

- 确保前端路由配置正确
- 检查 `vercel.json` 中的路由设置

#### 3. API 调用失败

- 检查环境变量配置
- 确保后端服务可访问
- 检查 CORS 设置

#### 4. 静态资源加载失败

- 检查构建输出目录
- 确保资源路径正确

### 调试方法

1. **查看构建日志**
   - 在 Vercel 控制台查看详细的构建日志
   - 检查是否有错误或警告

2. **本地测试**
   - 在本地运行 `npm run build` 确保构建成功
   - 使用 `npm run preview` 预览构建结果

3. **检查网络请求**
   - 使用浏览器开发者工具检查 API 调用
   - 确认请求地址和参数正确

## 🌱 自定义域名

### 配置自定义域名

1. 在 Vercel 项目设置中点击 "Domains"
2. 添加你的自定义域名
3. 按照提示配置 DNS 记录

### SSL 证书

Vercel 会自动为你的域名提供免费的 SSL 证书

## 📊 监控和分析

### Vercel Analytics

- 访问量统计
- 性能监控
- 错误追踪

### 自定义分析

你可以集成 Google Analytics 或其他分析工具：

1. 在 `frontend/index.html` 中添加分析代码
2. 重新部署项目

## 🔄 自动部署

Vercel 支持自动部署：

- 每次推送到 `main` 分支会自动触发部署
- 可以配置不同分支对应不同环境

## 📱 移动端优化

本项目已经配置了响应式设计：

- 在移动设备上会自动适配
- 触摸友好的界面设计
- 优化的加载性能

## 🎉 部署成功！

恭喜！你的智能改卷统一系统现在已经成功部署到 Vercel！

### 下一步

1. 分享你的应用链接
2. 收集用户反馈
3. 根据需要进行功能迭代

### 技术支持

如果遇到问题，可以：

- 查看 [Vercel 文档](https://vercel.com/docs)
- 提交 [GitHub Issue](https://github.com/NoccW/auto-grad/issues)
- 联系技术支持

---

🎊 享受你的智能改卷系统吧！
