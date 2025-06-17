# HTTPS功能使用说明

## 概述

简单Web文件服务器现已支持HTTPS协议，可以同时运行HTTP和HTTPS服务，为用户提供加密的安全连接。

## 功能特性

### ✅ 已实现的功能
- 🔒 **同时支持HTTP和HTTPS** - 两个服务器并行运行
- 📁 **自动证书加载** - 从指定目录加载SSL证书
- 🔧 **灵活配置** - 可配置端口和证书目录
- 📊 **实时状态监控** - API和页面显示HTTPS状态
- 🌐 **完整功能支持** - 所有功能（上传、WebDAV）都支持HTTPS

### 🔒 安全特性
- SSL/TLS加密传输
- 自签名证书支持（测试用）
- 证书状态实时检查
- 默认禁用（需明确启用）

## 启用HTTPS服务

### 基本命令
```bash
# 启用HTTPS服务（使用默认端口和证书目录）
sweb.exe -https

# 指定HTTPS端口
sweb.exe -https -https-port 9443

# 指定证书目录
sweb.exe -https -cert-dir /path/to/certs

# 同时启用所有功能
sweb.exe -upload -webdav -https
```

### 命令行参数
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-https` 或 `--enable-https` | 启用HTTPS服务 | 禁用 |
| `-https-port` | HTTPS服务器端口 | 8443 |
| `-cert-dir` | SSL证书目录 | ./cert |
| `-port` 或 `-p` | HTTP服务器端口 | 8080 |

## SSL证书配置

### 证书文件要求
HTTPS服务需要以下证书文件：
- `server.crt` - SSL证书文件
- `server.key` - 私钥文件

### 证书目录结构
```
cert/
├── server.crt    # SSL证书
└── server.key    # 私钥文件
```

### 生成自签名证书（测试用）

#### 方法1：使用提供的Go程序
```bash
# 运行证书生成程序
go run generate-cert.go
```

#### 方法2：使用OpenSSL
```bash
# 创建证书目录
mkdir cert

# 生成自签名证书
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes -subj "/C=CN/ST=Beijing/L=Beijing/O=Test/OU=Test/CN=localhost"
```

#### 方法3：使用PowerShell（Windows）
```powershell
# 生成自签名证书
$cert = New-SelfSignedCertificate -DnsName "localhost" -CertStoreLocation "cert:\LocalMachine\My"
$pwd = ConvertTo-SecureString -String "password" -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath "cert\server.pfx" -Password $pwd
```

## 访问地址

启用HTTPS后，服务器将同时监听两个端口：

### HTTP服务（默认）
- **主页**: http://localhost:8080
- **上传**: http://localhost:8080/upload
- **WebDAV**: http://localhost:8080/webdav
- **状态API**: http://localhost:8080/api/upload-status

### HTTPS服务（启用后）
- **主页**: https://localhost:8443
- **上传**: https://localhost:8443/upload
- **WebDAV**: https://localhost:8443/webdav
- **状态API**: https://localhost:8443/api/upload-status

## 实时状态监控

### Web界面状态显示
主页会实时显示HTTPS服务状态：
- 🔒 **已禁用** - HTTPS服务未启用
- ✅ **已启用** - HTTPS服务正常运行，证书有效
- ⚠️ **证书错误** - HTTPS服务已启用，但证书文件有问题

### API状态查询
```bash
# 查询服务状态
curl http://localhost:8080/api/upload-status
```

API响应示例：
```json
{
  "https": {
    "enabled": true,
    "httpPort": 8080,
    "httpsPort": 8443,
    "certDir": "./cert",
    "certStatus": "enabled"
  },
  "upload": {
    "enabled": true,
    "status": "enabled"
  },
  "webdav": {
    "enabled": true,
    "readonly": false,
    "directory": ".",
    "status": "enabled-readwrite"
  }
}
```

## 使用场景

### 开发测试
```bash
# 快速启动HTTPS测试环境
sweb.exe -https
```

### 文件共享（安全）
```bash
# 启用加密文件上传
sweb.exe -upload -https
```

### WebDAV over HTTPS
```bash
# 安全的WebDAV服务
sweb.exe -webdav -https
```

### 完整功能（生产环境）
```bash
# 启用所有功能，使用自定义端口
sweb.exe -upload -webdav -https -port 80 -https-port 443
```

## 故障排除

### 常见问题

#### 1. 证书文件不存在
**错误**: `证书文件不存在: ./cert/server.crt`
**解决**: 
- 确保证书目录存在
- 生成或复制正确的证书文件
- 检查文件权限

#### 2. 端口被占用
**错误**: `HTTPS服务器错误: listen tcp :8443: bind: address already in use`
**解决**:
```bash
# 检查端口占用
netstat -an | findstr :8443

# 使用其他端口
sweb.exe -https -https-port 9443
```

#### 3. 浏览器安全警告
**现象**: 浏览器显示"不安全连接"警告
**原因**: 使用自签名证书
**解决**:
- 测试环境：点击"高级" → "继续访问"
- 生产环境：使用受信任CA签发的证书

#### 4. 证书格式错误
**错误**: `x509: failed to parse certificate`
**解决**:
- 确保证书文件格式正确（PEM格式）
- 重新生成证书文件
- 检查证书文件完整性

### 日志调试
服务器启动时会显示详细信息：
```
🔒 HTTPS服务器启动在 https://localhost:8443
📁 证书目录: ./cert
🌐 HTTP服务器启动在 http://localhost:8080
```

## 安全建议

### 测试环境
- ✅ 使用自签名证书
- ✅ 接受浏览器安全警告
- ✅ 限制网络访问范围

### 生产环境
- 🔒 使用受信任CA签发的证书
- 🔒 定期更新证书
- 🔒 配置防火墙规则
- 🔒 使用强密码保护私钥
- 🔒 定期备份证书文件

### 最佳实践
1. **证书管理**: 使用Let's Encrypt等免费CA
2. **端口配置**: 生产环境使用标准端口（80/443）
3. **访问控制**: 结合防火墙限制访问
4. **监控**: 定期检查证书有效期
5. **备份**: 安全存储证书和私钥文件

---

**HTTPS功能** - 为您的文件服务器提供企业级安全保护！ 🔒
