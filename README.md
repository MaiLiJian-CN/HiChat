# HiChat - 实时即时通信系统

一个基于Go语言、Gin框架和Gorm ORM构建的即时通信系统，支持用户认证、实时消息传递等核心功能。

## 项目版本演进

### v1.0 - 项目基础构建
- 完成项目整体架构设计，确定目录结构
- 初始化Git仓库和Go模块配置
- 搭建基础开发环境，集成Gin框架和Gorm ORM
- 设计并创建数据库表结构（以用户表为核心）
- 实现基础的错误处理和日志记录机制

### v2.0 - DAO层代码设计
- 完成数据访问层（DAO）封装，实现与数据库的交互
- 针对`UserBasic`等核心模型编写CRUD操作
- 实现数据库连接池配置和优化
- 添加数据校验逻辑（基于govalidator）
- 完成数据库迁移脚本，支持环境初始化

### v3.0 - Server层实现与JWT鉴权
- 开发业务逻辑层（Server），封装核心业务处理
- 实现用户注册、登录功能
- 集成JWT认证机制，完成令牌生成与验证
- 编写认证中间件，实现接口权限控制
- 完善API错误码体系，统一响应格式

## 技术栈
- 后端框架：Go 1.24+、Gin v1.9.1
- ORM工具：Gorm v2.0
- 认证方式：JWT (JSON Web Token)
- 数据库：MySQL 8.0
- 数据校验：govalidator
- 密码加密：bcrypt（带盐值哈希）

## 快速开始

### 前置要求
- Go 1.24及以上版本
- MySQL 8.0数据库
- Git

### 安装步骤
1. 克隆仓库
   ```bash
   git clone https://github.com/MaiLiJian-CN/HiChat.git
   cd HiChat
   ```

2. 安装依赖
   ```bash
   go mod tidy
   ```

3. 配置数据库
   - 修改`db.go`文件中的数据库连接信息

4. 初始化数据库
   ```bash
   # 执行数据库迁移
   go run test/main.go
   ```

5. 启动服务
   ```bash
   go run main.go
   ```

## 目录结构HiChat/
├── config/         # 配置文件

├── dao/            # 数据访问层

├── middleware/     # 中间件（含JWT认证）

├── common/     # 中间件（MD5）

├── model/          # 数据模型定义

├── router/         # 路由配置

├── server/         # 业务逻辑层

├── main.go         # 程序入口
└── README.md       # 项目说明
## 后续规划
- v4.0：实现WebSocket实时通信功能
- v5.0：添加好友管理和消息历史记录功能
- v6.0：优化性能，支持集群部署
## 参考文档资料
网址：https://learnku.com/articles/74278
开发者：iceymoss
## 许可证
本项目采用MIT许可证 - 详见LICENSE文件
