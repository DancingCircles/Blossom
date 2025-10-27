# Blossom 论坛前端

## 📦 技术栈

- HTML5
- CSS3（响应式设计、网格布局、3D效果）
- JavaScript（原生JS，无框架）
- API调用（Fetch）

## 🚀 快速开始

### 1. 安装 Node.js

确保已安装 Node.js（推荐 v16 或更高版本）

```bash
# 检查是否已安装
node -v
npm -v
```

### 2. 运行项目

```bash
# 进入前端目录
cd frontend

# 方式1：使用 live-server（推荐，支持热更新）
npm run dev

# 方式2：使用 http-server
npm run serve
```

浏览器会自动打开 http://localhost:3000

### 3. 无需安装依赖

本项目使用原生JavaScript，不需要安装任何npm依赖包。
`npm run dev` 会使用 npx 自动运行临时的开发服务器。

## 📁 项目结构

```
frontend/
├── index.html          # 主页
├── login.html          # 登录/注册页
├── post.html           # 发帖页
├── css/
│   ├── style.css       # 主样式
│   ├── auth.css        # 登录/注册样式
│   ├── post.css        # 发帖页样式
│   └── smooth-scroll.css # 滚动优化
├── js/
│   ├── main.js         # 主页逻辑
│   ├── auth.js         # 登录/注册逻辑
│   ├── post.js         # 发帖逻辑
│   └── api.js          # API调用封装
└── package.json        # npm配置
```

## 🔧 API 配置

在 `js/api.js` 中配置后端API地址：

```javascript
const API_BASE_URL = 'http://localhost:8082/api/v1';
```

根据你的后端地址修改这个配置。

## ✨ 功能特性

### 已实现功能

- ✅ 首页展示（话题列表、排序、搜索）
- ✅ 用户登录/注册
- ✅ 发布话题
- ✅ 点赞/点踩
- ✅ 用户菜单（登录状态切换）
- ✅ 响应式设计（支持移动端）
- ✅ 平滑滚动和动画效果
- ✅ 3D视觉效果
- ✅ 自定义滚动条

### 待实现功能

- ⏳ 话题详情页
- ⏳ 评论功能
- ⏳ 用户个人主页
- ⏳ 图片上传
- ⏳ 实时通知

## 🎨 设计特色

- 网格背景纹理
- 流畅的曲线装饰
- 气泡风格logo
- 现代卡片设计
- 平滑过渡动画
- 3D悬停效果

## 🌐 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge
- 移动端浏览器

## 📝 开发说明

### 本地存储

使用 localStorage 存储：
- `token`: JWT认证令牌
- `username`: 当前用户名

### API调用示例

```javascript
// 登录
const result = await login(username, password);

// 获取话题列表
const topics = await getTopics({ page: 1, sort: 'hot' });

// 创建话题
const result = await createTopic(title, content, category);

// 点赞
await likeTopic(topicId);
```

## 🐛 常见问题

### 1. 跨域问题
确保后端已配置CORS，允许前端域名访问。

### 2. API调用失败
- 检查后端是否启动
- 检查API地址配置是否正确
- 打开浏览器控制台查看错误信息

### 3. 登录状态丢失
- 检查localStorage中是否有token
- 检查token是否过期

## 📄 License

MIT

## 👨‍💻 作者

Your Name

