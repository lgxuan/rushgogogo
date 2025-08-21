# RushGoGoGo

一个基于Golang的HTTP敏感信息被动代理工具(小玩具)，可做为中间件与其它扫描器对接，用来记录扫描过程中的敏感信息。

## 功能特性

- HTTP/HTTPS代理服务器
- 敏感信息实时过滤（手机号、身份证、邮箱、银行卡等）
- 自动HTTPS证书管理
- 可配置的过滤规则
- 多线程消息处理

## 快速开始

### 安装

```bash
git clone <repository-url>
cd rushgogogo
go build -o rushgogogo .
```

### 安装证书
可以修改`./pkgs/proxy/cert.go`来修改证书防止被识别
首次使用需要安装HTTPS拦截证书：

```bash
# macOS/Linux (需要sudo权限)
sudo ./rushgogogo install-cert

# Windows (需要管理员权限)
rushgogogo install-cert
```

### 启动代理

```bash
# 默认端口8081，默认线程数10
./rushgogogo listen

# 指定端口
./rushgogogo listen :8080

# 设置线程数
./rushgogogo listen -t 20

# 指定端口和线程数
./rushgogogo listen :8080 -t 25
```

## 配置说明

配置文件 `config.yaml` 包含：

- 服务器配置（线程数等）
- 过滤规则配置
- 敏感信息检测模式

## 过滤规则

支持多种敏感信息检测：

- 个人信息：手机号、身份证、邮箱
- 金融信息：银行卡、密钥、API密钥
- 技术信息：JWT、数据库连接、内网IP
- 云服务：AWS、Azure、阿里云等AK/SK

## 注意事项

- 安装证书需要系统管理员权限
- 仅用于合法的安全测试和开发目的
- 请遵守相关法律法规和道德准则

