package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// createDefaultPageIfNeeded 检查并创建默认页面
func createDefaultPageIfNeeded(webDir string, uploadEnabled bool) {
	// 检查是否存在默认页面
	indexFiles := []string{"index.html", "index.htm", "about.html", "about.htm"}
	hasDefaultPage := false

	for _, indexFile := range indexFiles {
		indexPath := filepath.Join(webDir, indexFile)
		if _, err := os.Stat(indexPath); err == nil {
			hasDefaultPage = true
			break
		}
	}

	// 如果没有默认页面，创建一个
	if !hasDefaultPage {
		indexPath := filepath.Join(webDir, "about.html")
		indexContent := generateMultiLanguagePageContent()

		err := os.WriteFile(indexPath, []byte(indexContent), 0644)
		if err != nil {
			log.Printf("警告：无法创建默认页面: %v", err)
		} else {
			fmt.Println("已创建默认页面: about.html (支持多语言和状态检查)")
		}
	}
}

// statusHandler 处理状态查询请求
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"upload": map[string]interface{}{
			"enabled": uploadEnabled,
			"status": func() string {
				if uploadEnabled {
					return "enabled"
				}
				return "disabled"
			}(),
		},
		"files": map[string]interface{}{
			"enabled": filesEnabled,
			"status": func() string {
				if filesEnabled {
					return "enabled"
				}
				return "disabled"
			}(),
		},
		"webdav": map[string]interface{}{
			"enabled":   webdavEnabled,
			"readonly":  webdavReadonly,
			"directory": webdavDir,
			"port":      webdavPort,
			"status": func() string {
				if webdavEnabled {
					if webdavReadonly {
						return "enabled-readonly"
					}
					return "enabled-readwrite"
				}
				return "disabled"
			}(),
		},
		"https": map[string]interface{}{
			"enabled":   httpsEnabled,
			"httpPort":  httpPort,
			"httpsPort": httpsPort,
			"certDir":   certDir,
			"certStatus": func() string {
				if !httpsEnabled {
					return "disabled"
				}
				if err := checkCertificates(); err != nil {
					return "cert-error"
				}
				return "enabled"
			}(),
		},
		"socks5": getSOCKS5Status(),
		"proxy":  getHTTPProxyStatus(),
	}

	json.NewEncoder(w).Encode(response)
}

// showHelpInfo 显示帮助信息
func showHelpInfo(lang string) {
	if lang == "en" {
		showHelpInfoEnglish()
	} else {
		showHelpInfoChinese()
	}
}

// showHelpInfoChinese 显示中文帮助信息
func showHelpInfoChinese() {
	fmt.Println("简单Web文件服务器 - 基于Go语言开发")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  sweb.exe [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -upload, --enable-upload    启用文件上传功能 (默认: 禁用)")
	fmt.Println("  -files, --enable-files      启用文件浏览功能 (默认: 禁用)")
	fmt.Println("  -webdav, --enable-webdav    启用WebDAV服务 (默认: 禁用)")
	fmt.Println("  -webdav-dir <目录>          WebDAV服务的根目录 (默认: 当前目录)")
	fmt.Println("  -webdav-port <端口>         指定WebDAV服务器端口 (默认: 8081)")
	fmt.Println("  -webdav-readonly            WebDAV服务只读模式 (默认: 读写)")
	fmt.Println("  -webdav-auth                启用WebDAV认证 (默认: 禁用)")
	fmt.Println("  -webdav-username <用户名>   WebDAV认证用户名 (默认: webdav)")
	fmt.Println("  -webdav-password <密码>     WebDAV认证密码 (默认: webdav)")
	fmt.Println("  -http, --enable-http        启用HTTP服务 (默认: 启用)")
	fmt.Println("  -https, --enable-https      启用HTTPS服务 (默认: 禁用)")
	fmt.Println("  -socks5, --enable-socks5    启用SOCKS5代理服务 (默认: 禁用)")
	fmt.Println("  -socks5-auth                启用SOCKS5代理认证 (默认: 禁用)")
	fmt.Println("  -socks5-username <用户名>   SOCKS5代理认证用户名 (默认: socks5)")
	fmt.Println("  -socks5-password <密码>     SOCKS5代理认证密码 (默认: socks5)")
	fmt.Println("  -proxy, --enable-proxy      启用HTTP代理服务 (默认: 禁用)")
	fmt.Println("  -proxy-auth                 启用HTTP代理认证 (默认: 禁用)")
	fmt.Println("  -proxy-username <用户名>    HTTP代理认证用户名 (默认: http)")
	fmt.Println("  -proxy-password <密码>      HTTP代理认证密码 (默认: http)")
	fmt.Println("  -https-proxy, --enable-https-proxy  启用HTTPS代理服务 (默认: 禁用)")
	fmt.Println("  -https-proxy-auth           启用HTTPS代理认证 (默认: 禁用)")
	fmt.Println("  -https-proxy-username <用户名>  HTTPS代理认证用户名 (默认: https)")
	fmt.Println("  -https-proxy-password <密码>    HTTPS代理认证密码 (默认: https)")
	fmt.Println("  -port, -p <端口>           指定HTTP服务器端口 (默认: 8080)")
	fmt.Println("  -https-port <端口>         指定HTTPS服务器端口 (默认: 8443)")
	fmt.Println("  -socks5-port <端口>        指定SOCKS5代理端口 (默认: 1080)")
	fmt.Println("  -proxy-port <端口>         指定HTTP代理端口 (默认: 10808)")
	fmt.Println("  -https-proxy-port <端口>   指定HTTPS代理端口 (默认: 10443)")
	fmt.Println("  -cert-dir <目录>           SSL证书目录 (默认: ./cert)")
	fmt.Println("  -help, -h                  显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  sweb.exe                           # 启动HTTP服务器，仅提供静态文件服务")
	fmt.Println("  sweb.exe -upload                   # 启动服务器并启用文件上传功能")
	fmt.Println("  sweb.exe -files                    # 启动服务器并启用文件浏览功能")
	fmt.Println("  sweb.exe -https                    # 启动HTTP和HTTPS服务器")
	fmt.Println("  sweb.exe -webdav                   # 启动WebDAV服务 (端口8081)")
	fmt.Println("  sweb.exe -webdav -webdav-auth      # 启动WebDAV服务并启用认证")
	fmt.Println("  sweb.exe -webdav -webdav-readonly  # 启动只读WebDAV服务")
	fmt.Println("  sweb.exe -webdav -webdav-dir /data # 指定WebDAV目录")
	fmt.Println("  sweb.exe -webdav -webdav-port 8082 # 指定WebDAV端口")
	fmt.Println("  sweb.exe -socks5                   # 启动SOCKS5代理服务")
	fmt.Println("  sweb.exe -socks5 -socks5-auth      # 启动SOCKS5代理服务并启用认证")
	fmt.Println("  sweb.exe -proxy                    # 启动HTTP代理服务")
	fmt.Println("  sweb.exe -proxy -proxy-auth        # 启动HTTP代理服务并启用认证")
	fmt.Println("  sweb.exe -https-proxy              # 启动HTTPS代理服务")
	fmt.Println("  sweb.exe -https-proxy -https-proxy-auth # 启动HTTPS代理服务并启用认证")
	fmt.Println("  sweb.exe -socks5 -socks5-port 1080 # 指定SOCKS5代理端口")
	fmt.Println("  sweb.exe -proxy -proxy-port 10808  # 指定HTTP代理端口")
	fmt.Println("  sweb.exe -https-proxy -https-proxy-port 10443 # 指定HTTPS代理端口")
	fmt.Println("  sweb.exe -upload -files -webdav -https -socks5 -proxy -https-proxy # 启用所有功能")
	fmt.Println("  sweb.exe -https -https-port 9443   # 指定HTTPS端口")
	fmt.Println("  sweb.exe -http=false -https        # 仅启用HTTPS服务")
	fmt.Println()
	fmt.Println("HTTPS证书:")
	fmt.Println("  证书文件: ./cert/server.crt")
	fmt.Println("  私钥文件: ./cert/server.key")
	fmt.Println("  可以使用openssl生成自签名证书用于测试")
	fmt.Println()
	fmt.Println("访问地址:")
	fmt.Println("  HTTP: http://localhost:8080")
	fmt.Println("  HTTPS: https://localhost:8443 (如果启用)")
	fmt.Println("  WebDAV: http://localhost:8080/webdav")
	fmt.Println("  SOCKS5代理: localhost:1080 (如果启用)")
	fmt.Println("  HTTP代理: localhost:10808 (如果启用)")
	fmt.Println("  HTTPS代理: localhost:10443 (如果启用)")
	fmt.Println()
	fmt.Println("安全说明:")
	fmt.Println("  文件上传、文件浏览、WebDAV、HTTPS、SOCKS5代理、HTTP代理和HTTPS代理功能默认禁用以确保服务器安全。")
	fmt.Println("  只有在明确需要时才使用相应参数启用。")
}

// showHelpInfoEnglish 显示英文帮助信息
func showHelpInfoEnglish() {
	fmt.Println("Simple Web File Server - Developed in Go")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sweb.exe [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -upload, --enable-upload    Enable file upload feature (default: disabled)")
	fmt.Println("  -files, --enable-files      Enable file browsing feature (default: disabled)")
	fmt.Println("  -webdav, --enable-webdav    Enable WebDAV service (default: disabled)")
	fmt.Println("  -webdav-dir <directory>     WebDAV service root directory (default: current directory)")
	fmt.Println("  -webdav-port <port>         Specify WebDAV server port (default: 8081)")
	fmt.Println("  -webdav-readonly            Enable WebDAV in read-only mode (default: read-write)")
	fmt.Println("  -webdav-auth                Enable WebDAV authentication (default: disabled)")
	fmt.Println("  -webdav-username <username> WebDAV authentication username (default: webdav)")
	fmt.Println("  -webdav-password <password> WebDAV authentication password (default: webdav)")
	fmt.Println("  -http, --enable-http        Enable HTTP service (default: enabled)")
	fmt.Println("  -https, --enable-https      Enable HTTPS service (default: disabled)")
	fmt.Println("  -socks5, --enable-socks5    Enable SOCKS5 proxy service (default: disabled)")
	fmt.Println("  -socks5-auth                Enable SOCKS5 proxy authentication (default: disabled)")
	fmt.Println("  -socks5-username <username> SOCKS5 proxy authentication username (default: socks5)")
	fmt.Println("  -socks5-password <password> SOCKS5 proxy authentication password (default: socks5)")
	fmt.Println("  -proxy, --enable-proxy      Enable HTTP proxy service (default: disabled)")
	fmt.Println("  -proxy-auth                 Enable HTTP proxy authentication (default: disabled)")
	fmt.Println("  -proxy-username <username>  HTTP proxy authentication username (default: http)")
	fmt.Println("  -proxy-password <password>  HTTP proxy authentication password (default: http)")
	fmt.Println("  -https-proxy, --enable-https-proxy  Enable HTTPS proxy service (default: disabled)")
	fmt.Println("  -https-proxy-auth           Enable HTTPS proxy authentication (default: disabled)")
	fmt.Println("  -https-proxy-username <username>  HTTPS proxy authentication username (default: https)")
	fmt.Println("  -https-proxy-password <password>  HTTPS proxy authentication password (default: https)")
	fmt.Println("  -port, -p <port>           Specify HTTP server port (default: 8080)")
	fmt.Println("  -https-port <port>         Specify HTTPS server port (default: 8443)")
	fmt.Println("  -socks5-port <port>        Specify SOCKS5 proxy port (default: 1080)")
	fmt.Println("  -proxy-port <port>         Specify HTTP proxy port (default: 10808)")
	fmt.Println("  -https-proxy-port <port>   Specify HTTPS proxy port (default: 10443)")
	fmt.Println("  -cert-dir <directory>      SSL certificate directory (default: ./cert)")
	fmt.Println("  -help-lang <language>      Help information language (zh/en) (default: zh)")
	fmt.Println("  -help, -h                  Show this help information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sweb.exe                   # Start basic HTTP server")
	fmt.Println("  sweb.exe -upload -files    # Enable upload and file browsing")
	fmt.Println("  sweb.exe -webdav           # Enable WebDAV service (port 8081)")
	fmt.Println("  sweb.exe -webdav-auth      # Enable WebDAV with authentication")
	fmt.Println("  sweb.exe -webdav -webdav-port 8082 # Specify WebDAV port")
	fmt.Println("  sweb.exe -https            # Enable HTTPS service")
	fmt.Println("  sweb.exe -socks5           # Enable SOCKS5 proxy")
	fmt.Println("  sweb.exe -socks5-auth      # Enable SOCKS5 proxy with authentication")
	fmt.Println("  sweb.exe -proxy            # Enable HTTP proxy")
	fmt.Println("  sweb.exe -proxy-auth       # Enable HTTP proxy with authentication")
	fmt.Println("  sweb.exe -https-proxy      # Enable HTTPS proxy")
	fmt.Println("  sweb.exe -https-proxy-auth # Enable HTTPS proxy with authentication")
	fmt.Println("  sweb.exe -upload -files -webdav -https -socks5 -proxy -https-proxy # Enable all features")
	fmt.Println("  sweb.exe -https -https-port 9443   # Specify HTTPS port")
	fmt.Println("  sweb.exe -http=false -https        # Enable HTTPS only")
	fmt.Println()
	fmt.Println("HTTPS Certificates:")
	fmt.Println("  Certificate file: ./cert/server.crt")
	fmt.Println("  Private key file: ./cert/server.key")
	fmt.Println("  You can use openssl to generate self-signed certificates for testing")
	fmt.Println()
	fmt.Println("Access URLs:")
	fmt.Println("  HTTP: http://localhost:8080")
	fmt.Println("  HTTPS: https://localhost:8443 (if enabled)")
	fmt.Println("  WebDAV: http://localhost:8080/webdav")
	fmt.Println("  SOCKS5 proxy: localhost:1080 (if enabled)")
	fmt.Println("  HTTP proxy: localhost:10808 (if enabled)")
	fmt.Println("  HTTPS proxy: localhost:10443 (if enabled)")
	fmt.Println()
	fmt.Println("Security Notice:")
	fmt.Println("  File upload, file browsing, WebDAV, HTTPS, SOCKS5 proxy, HTTP proxy and HTTPS proxy features")
	fmt.Println("  are disabled by default to ensure server security.")
	fmt.Println("  Only enable them when explicitly needed.")
}

// generateEnhancedDefaultPageContent 生成包含WebDAV功能的增强版默认页面
func generateEnhancedDefaultPageContent() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>简单Web文件服务器</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            border-bottom: 3px solid #007acc;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        .feature {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #007acc;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            background: #007acc;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 5px;
            transition: background 0.3s;
        }
        .button:hover {
            background: #005a9e;
        }
        .button.disabled {
            background: #6c757d;
            cursor: not-allowed;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        code {
            background: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .status-indicator {
            font-weight: bold;
            padding: 2px 6px;
            border-radius: 3px;
        }
        .status-enabled {
            color: #28a745;
        }
        .status-disabled {
            color: #dc3545;
        }
        .status-readonly {
            color: #ffc107;
        }
        .loading {
            color: #6c757d;
        }
        .hidden {
            display: none;
        }
        .webdav-info {
            background: #e7f3ff;
            border-left: 4px solid #007acc;
            padding: 10px;
            margin: 10px 0;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌐 简单Web文件服务器</h1>

        <p>欢迎使用这个简单而实用的Web文件服务器！这是一个用Go语言编写的轻量级文件管理工具。</p>

        <h2>📋 项目功能</h2>

        <div class="feature">
            <strong>📁 静态文件服务</strong><br>
            自动服务web目录下的所有文件，支持HTML、CSS、JavaScript、图片等各种文件类型。
        </div>

        <div class="feature">
            <strong>📤 文件上传</strong><br>
            <span id="upload-feature-description">文件上传功能可通过命令行参数启用，确保服务器安全。</span>
            <span id="upload-status" class="status-indicator loading">🔄 检查中...</span>
        </div>

        <div class="feature">
            <strong>🌐 WebDAV服务</strong><br>
            <span id="webdav-feature-description">WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。</span>
            <span id="webdav-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="webdav-info" class="webdav-info hidden">
                <strong>WebDAV地址:</strong> <code id="webdav-url">http://localhost:8081/webdav</code><br>
                <strong>挂载目录:</strong> <code id="webdav-directory">.</code><br>
                <strong>访问模式:</strong> <span id="webdav-mode">读写</span>
            </div>
        </div>

        <div class="feature">
            <strong>🔒 HTTPS服务</strong><br>
            <span id="https-feature-description">HTTPS服务可通过命令行参数启用，提供加密的安全连接。</span>
            <span id="https-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="https-info" class="webdav-info hidden">
                <strong>HTTP地址:</strong> <code id="http-url">http://localhost:8080</code><br>
                <strong>HTTPS地址:</strong> <code id="https-url">https://localhost:8443</code><br>
                <strong>证书目录:</strong> <code id="cert-directory">./cert</code><br>
                <strong>证书状态:</strong> <span id="cert-status">检查中</span>
            </div>
        </div>        

        <div class="feature">
            <strong>📂 文件浏览下载</strong><br>
            <span id="files-feature-description">专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。</span>
            <span id="files-status" class="status-indicator loading">🔄 检查中...</span>
        </div>

        <div class="feature">
            <strong>🔌 SOCKS5代理服务</strong><br>
            <span id="socks5-feature-description">高性能SOCKS5代理服务器，支持TCP连接代理。</span>
            <span id="socks5-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="socks5-info" class="webdav-info hidden">
                <strong>代理地址:</strong> <code id="socks5-address">localhost:1080</code><br>
                <strong>协议版本:</strong> <span>SOCKS5</span><br>
                <strong>认证方式:</strong> <span>无认证</span>
            </div>
        </div>

        <div class="feature">
            <strong>🌐 HTTP代理服务</strong><br>
            <span id="proxy-feature-description">高性能HTTP/HTTPS代理服务器，支持Web浏览器代理。</span>
            <span id="proxy-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="proxy-info" class="webdav-info hidden">
                <strong>代理地址:</strong> <code id="proxy-address">localhost:10808</code><br>
                <strong>支持协议:</strong> <span>HTTP, HTTPS</span><br>
                <strong>连接方式:</strong> <span>CONNECT, GET, POST</span>
            </div>
        </div>

        <div class="feature">
            <strong>📊 功能状态展示</strong><br>
            默认页面实时显示服务器各项功能的启用状态，包括文件上传、WebDAV、HTTPS和代理服务。
        </div>

        <div class="feature">
            <strong>🔧 自动配置</strong><br>
            自动创建必要的目录结构，无需手动配置即可使用。
        </div>

        <h2>🚀 快速开始</h2>

        <div id="upload-section">
            <p><strong>文件上传：</strong></p>
            <div id="upload-enabled-content" class="hidden">
                <a href="/upload" class="button" id="upload-button">📤 上传文件</a>
            </div>
            <div id="upload-disabled-content" class="hidden">
                <p>要启用文件上传功能，请使用以下命令启动服务器：</p>
                <code>sweb.exe -upload</code> 或 <code>sweb.exe --enable-upload</code>
                <br><br>
                <a href="/upload" class="button disabled" id="upload-button-disabled">📤 上传功能已禁用</a>
            </div>
        </div>

        <div id="webdav-section">
            <p><strong>WebDAV服务：</strong></p>
            <div id="webdav-enabled-content" class="hidden">
                <a href="/webdav" class="button" id="webdav-button">🌐 访问WebDAV</a>
                <p><small>可以在文件管理器中添加网络位置：<code id="webdav-mount-url">http://localhost:8081/webdav</code></small></p>
            </div>
            <div id="webdav-disabled-content" class="hidden">
                <p>要启用WebDAV服务，请使用以下命令启动服务器：</p>
                <code>sweb.exe -webdav</code> 或 <code>sweb.exe --enable-webdav</code>
                <br><br>
                <a href="/webdav" class="button disabled" id="webdav-button-disabled">🌐 WebDAV服务已禁用</a>
            </div>
        </div>

        <div id="files-section">
            <p><strong>文件浏览：</strong></p>
            <div id="files-enabled-content" class="hidden">
                <a href="/files" class="button" id="files-button">📂 浏览文件</a>
                <p><small>查看和下载web目录中的所有文件</small></p>
            </div>
            <div id="files-disabled-content" class="hidden">
                <p>要启用文件浏览功能，请使用以下命令启动服务器：</p>
                <code>sweb.exe -files</code> 或 <code>sweb.exe --enable-files</code>
                <br><br>
                <a href="/files" class="button disabled" id="files-button-disabled">📂 文件浏览已禁用</a>
            </div>
        </div>

        <p><strong>服务器信息：</strong></p>
        <ul>
            <li>服务端口: <code>8080</code></li>
            <li>文件目录: <code>./web</code></li>
            <li>上传地址: <code>/upload</code></li>
            <li>文件浏览: <code>/files</code></li>
            <li>WebDAV地址: <code>/webdav</code></li>
            <li>SOCKS5代理: <code>localhost:1080</code></li>
            <li>HTTP代理: <code>localhost:10808</code></li>
        </ul>

        <h2>💡 使用说明</h2>
        <div id="usage-instructions">
            <div id="usage-upload-enabled" class="hidden">
                <h3>文件上传</h3>
                <ol>
                    <li><strong>上传文件</strong>：点击上方的"上传文件"按钮，选择要上传的文件</li>
                    <li><strong>访问文件</strong>：上传成功后，文件将保存在web目录下，可以直接通过URL访问</li>
                    <li><strong>管理文件</strong>：所有上传的文件都会显示在主页的文件列表中</li>
                </ol>
            </div>
            <div id="usage-webdav-enabled" class="hidden">
                <h3>WebDAV服务</h3>
                <ol>
                    <li><strong>Windows</strong>：在文件资源管理器中，右键"此电脑" → "映射网络驱动器" → 输入WebDAV地址</li>
                    <li><strong>macOS</strong>：在Finder中，按Cmd+K → 输入WebDAV地址</li>
                    <li><strong>Linux</strong>：使用davfs2或其他WebDAV客户端挂载</li>
                    <li><strong>移动设备</strong>：使用支持WebDAV的文件管理应用</li>
                </ol>
            </div>
            <div id="usage-disabled" class="hidden">
                <ol>
                    <li><strong>启用功能</strong>：使用相应的命令行参数启动服务器</li>
                    <li><strong>浏览文件</strong>：当前可以浏览和下载web目录中的现有文件</li>
                    <li><strong>安全考虑</strong>：高级功能默认禁用，确保服务器安全</li>
                </ol>
            </div>
        </div>

        <h2>🛠️ 技术特性</h2>

        <ul>
            <li>使用Go语言标准库开发，轻量级无外部依赖</li>
            <li>支持多部分表单数据上传</li>
            <li>完整的WebDAV协议支持（RFC 4918）</li>
            <li>高性能SOCKS5代理服务器（RFC 1928）</li>
            <li>HTTP/HTTPS代理服务器支持</li>
            <li>可配置的读写权限控制</li>
            <li>自动MIME类型检测</li>
            <li>UTF-8编码支持，完美处理中文</li>
            <li>跨平台兼容（Windows、Linux、macOS）</li>
            <li>并发连接处理，高性能数据转发</li>
        </ul>

        <div class="footer">
            <p>🔗 <strong>简单Web文件服务器</strong> | 基于Go语言开发</p>
            <div id="footer-content">
                <div id="footer-enabled" class="hidden">
                    <p>开始使用：
                        <span id="footer-files-link" class="hidden"><a href="/files" class="button">浏览文件</a></span>
                        <span id="footer-upload-link" class="hidden"><a href="/upload" class="button">上传文件</a></span>
                        <span id="footer-webdav-link" class="hidden"><a href="/webdav" class="button">访问WebDAV</a></span>
                    </p>
                </div>
                <div id="footer-disabled" class="hidden">
                    <p>安全模式：高级功能已禁用</p>
                    <p>使用 <code>sweb.exe -help</code> 查看所有可用选项</p>
                </div>
            </div>
        </div>
    </div>

    <script>
        // 检查服务状态
        function checkServiceStatus() {
            fetch('/api/status')
                .then(response => response.json())
                .then(data => {
                    updateUploadStatus(data.upload);
                    updateFilesStatus(data.files);
                    updateWebDAVStatus(data.webdav);
                    updateHTTPSStatus(data.https);
                    updateSOCKS5Status(data.socks5);
                    updateProxyStatus(data.proxy);
                    updateUsageInstructions(data.upload, data.files, data.webdav, data.https, data.socks5, data.proxy);
                })
                .catch(error => {
                    console.error('检查服务状态失败:', error);
                    // 如果API调用失败，显示默认的禁用状态
                    updateUploadStatus({enabled: false, status: 'disabled'});
                    updateFilesStatus({enabled: false, status: 'disabled'});
                    updateWebDAVStatus({enabled: false, status: 'disabled'});
                    updateHTTPSStatus({enabled: false, status: 'disabled'});
                    updateSOCKS5Status({enabled: false, status: 'disabled'});
                    updateProxyStatus({enabled: false, status: 'disabled'});
                    updateUsageInstructions({enabled: false}, {enabled: false}, {enabled: false}, {enabled: false}, {enabled: false}, {enabled: false});
                });
        }

        // 更新页面上的上传状态显示
        function updateUploadStatus(uploadData) {
            const statusElement = document.getElementById('upload-status');
            const descriptionElement = document.getElementById('upload-feature-description');
            const uploadEnabledContent = document.getElementById('upload-enabled-content');
            const uploadDisabledContent = document.getElementById('upload-disabled-content');

            if (uploadData.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = '通过简单的Web界面上传文件到服务器，支持各种文件格式。';

                uploadEnabledContent.classList.remove('hidden');
                uploadDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = '文件上传功能可通过命令行参数启用，确保服务器安全。';

                uploadEnabledContent.classList.add('hidden');
                uploadDisabledContent.classList.remove('hidden');
            }
        }

        // 更新页面上的文件浏览状态显示
        function updateFilesStatus(filesData) {
            const statusElement = document.getElementById('files-status');
            const descriptionElement = document.getElementById('files-feature-description');
            const filesEnabledContent = document.getElementById('files-enabled-content');
            const filesDisabledContent = document.getElementById('files-disabled-content');

            if (filesData.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = '专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。';

                filesEnabledContent.classList.remove('hidden');
                filesDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = '文件浏览功能可通过命令行参数启用，确保服务器安全。';

                filesEnabledContent.classList.add('hidden');
                filesDisabledContent.classList.remove('hidden');
            }
        }

        // 更新页面上的WebDAV状态显示
        function updateWebDAVStatus(webdavData) {
            const statusElement = document.getElementById('webdav-status');
            const descriptionElement = document.getElementById('webdav-feature-description');
            const webdavEnabledContent = document.getElementById('webdav-enabled-content');
            const webdavDisabledContent = document.getElementById('webdav-disabled-content');
            const webdavInfo = document.getElementById('webdav-info');
            const webdavDirectory = document.getElementById('webdav-directory');
            const webdavMode = document.getElementById('webdav-mode');

            if (webdavData.enabled) {
                if (webdavData.readonly) {
                    statusElement.textContent = '📖 只读模式';
                    statusElement.className = 'status-indicator status-readonly';
                    descriptionElement.textContent = 'WebDAV服务已启用（只读模式），支持文件浏览和下载。';
                    webdavMode.textContent = '只读';
                } else {
                    statusElement.textContent = '✅ 读写模式';
                    statusElement.className = 'status-indicator status-enabled';
                    descriptionElement.textContent = 'WebDAV服务已启用（读写模式），支持完整的文件管理操作。';
                    webdavMode.textContent = '读写';
                }

                webdavDirectory.textContent = webdavData.directory || '.';

                // 动态更新WebDAV URL（使用端口信息）
                const webdavPort = webdavData.port || 8081;
                const webdavUrlText = ` + "`" + `http://localhost:${webdavPort}/webdav` + "`" + `;
                const webdavUrl = document.getElementById('webdav-url');
                const webdavMountUrl = document.getElementById('webdav-mount-url');
                if (webdavUrl) webdavUrl.textContent = webdavUrlText;
                if (webdavMountUrl) webdavMountUrl.textContent = webdavUrlText;

                webdavInfo.classList.remove('hidden');
                webdavEnabledContent.classList.remove('hidden');
                webdavDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。';

                webdavInfo.classList.add('hidden');
                webdavEnabledContent.classList.add('hidden');
                webdavDisabledContent.classList.remove('hidden');
            }
        }

        // 更新页面上的HTTPS状态显示
        function updateHTTPSStatus(httpsData) {
            const statusElement = document.getElementById('https-status');
            const descriptionElement = document.getElementById('https-feature-description');
            const httpsInfo = document.getElementById('https-info');
            const httpUrl = document.getElementById('http-url');
            const httpsUrl = document.getElementById('https-url');
            const certDirectory = document.getElementById('cert-directory');
            const certStatus = document.getElementById('cert-status');

            if (httpsData.enabled) {
                if (httpsData.certStatus === 'enabled') {
                    statusElement.textContent = '✅ 已启用';
                    statusElement.className = 'status-indicator status-enabled';
                    descriptionElement.textContent = 'HTTPS服务已启用，提供加密的安全连接。';
                    certStatus.textContent = '证书正常';
                    certStatus.className = 'status-enabled';
                } else if (httpsData.certStatus === 'cert-error') {
                    statusElement.textContent = '⚠️ 证书错误';
                    statusElement.className = 'status-indicator status-readonly';
                    descriptionElement.textContent = 'HTTPS服务已启用，但证书文件有问题。';
                    certStatus.textContent = '证书错误';
                    certStatus.className = 'status-disabled';
                }

                httpUrl.textContent = 'http://localhost:' + httpsData.httpPort;
                httpsUrl.textContent = 'https://localhost:' + httpsData.httpsPort;
                certDirectory.textContent = httpsData.certDir || './cert';
                httpsInfo.classList.remove('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'HTTPS服务可通过命令行参数启用，提供加密的安全连接。';

                httpsInfo.classList.add('hidden');
            }
        }

        // 更新页面上的SOCKS5代理状态显示
        function updateSOCKS5Status(socks5Data) {
            const statusElement = document.getElementById('socks5-status');
            const descriptionElement = document.getElementById('socks5-feature-description');
            const socks5Info = document.getElementById('socks5-info');
            const socks5Address = document.getElementById('socks5-address');

            if (socks5Data.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = 'SOCKS5代理服务已启用，支持高性能TCP连接代理。';

                socks5Address.textContent = 'localhost:' + socks5Data.port;
                socks5Info.classList.remove('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'SOCKS5代理服务可通过命令行参数启用。';

                socks5Info.classList.add('hidden');
            }
        }

        // 更新页面上的HTTP代理状态显示
        function updateProxyStatus(proxyData) {
            const statusElement = document.getElementById('proxy-status');
            const descriptionElement = document.getElementById('proxy-feature-description');
            const proxyInfo = document.getElementById('proxy-info');
            const proxyAddress = document.getElementById('proxy-address');

            if (proxyData.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = 'HTTP代理服务已启用，支持Web浏览器和HTTP客户端代理。';

                proxyAddress.textContent = 'localhost:' + proxyData.port;
                proxyInfo.classList.remove('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'HTTP代理服务可通过命令行参数启用。';

                proxyInfo.classList.add('hidden');
            }
        }

        // 更新使用说明和页脚
        function updateUsageInstructions(uploadData, filesData, webdavData, httpsData, socks5Data, proxyData) {
            const usageUploadEnabled = document.getElementById('usage-upload-enabled');
            const usageWebdavEnabled = document.getElementById('usage-webdav-enabled');
            const usageDisabled = document.getElementById('usage-disabled');
            const footerEnabled = document.getElementById('footer-enabled');
            const footerDisabled = document.getElementById('footer-disabled');
            const footerFilesLink = document.getElementById('footer-files-link');
            const footerUploadLink = document.getElementById('footer-upload-link');
            const footerWebdavLink = document.getElementById('footer-webdav-link');

            const anyEnabled = uploadData.enabled || filesData.enabled || webdavData.enabled || httpsData.enabled || socks5Data.enabled || proxyData.enabled;

            if (anyEnabled) {
                footerEnabled.classList.remove('hidden');
                footerDisabled.classList.add('hidden');
                usageDisabled.classList.add('hidden');

                if (filesData.enabled) {
                    footerFilesLink.classList.remove('hidden');
                } else {
                    footerFilesLink.classList.add('hidden');
                }

                if (uploadData.enabled) {
                    usageUploadEnabled.classList.remove('hidden');
                    footerUploadLink.classList.remove('hidden');
                } else {
                    usageUploadEnabled.classList.add('hidden');
                    footerUploadLink.classList.add('hidden');
                }

                if (webdavData.enabled) {
                    usageWebdavEnabled.classList.remove('hidden');
                    footerWebdavLink.classList.remove('hidden');
                } else {
                    usageWebdavEnabled.classList.add('hidden');
                    footerWebdavLink.classList.add('hidden');
                }
            } else {
                footerEnabled.classList.add('hidden');
                footerDisabled.classList.remove('hidden');
                usageDisabled.classList.remove('hidden');
                usageUploadEnabled.classList.add('hidden');
                usageWebdavEnabled.classList.add('hidden');
            }
        }

        // 页面加载时检查状态
        document.addEventListener('DOMContentLoaded', function() {
            checkServiceStatus();

            // 每30秒检查一次状态
            setInterval(checkServiceStatus, 30000);
        });
    </script>
</body>
</html>`
}

// generateMultiLanguagePageContent 生成支持多语言的增强版默认页面
func generateMultiLanguagePageContent() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title id="page-title">简单Web文件服务器</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .lang-switch {
            display: flex;
            gap: 10px;
        }
        .lang-btn {
            padding: 5px 10px;
            border: 1px solid #007acc;
            background: white;
            color: #007acc;
            border-radius: 3px;
            cursor: pointer;
            text-decoration: none;
            font-size: 12px;
        }
        .lang-btn.active {
            background: #007acc;
            color: white;
        }
        .lang-btn:hover {
            background: #005a9e;
            color: white;
        }
        h1 {
            color: #333;
            text-align: center;
            border-bottom: 3px solid #007acc;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        .feature {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #007acc;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            background: #007acc;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 5px;
            transition: background 0.3s;
        }
        .button:hover {
            background: #005a9e;
        }
        .button.disabled {
            background: #6c757d;
            cursor: not-allowed;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        code {
            background: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .status-indicator {
            font-weight: bold;
            padding: 2px 6px;
            border-radius: 3px;
        }
        .status-enabled {
            color: #28a745;
        }
        .status-disabled {
            color: #dc3545;
        }
        .status-readonly {
            color: #ffc107;
        }
        .loading {
            color: #6c757d;
        }
        .hidden {
            display: none;
        }
        .webdav-info {
            background: #e7f3ff;
            border-left: 4px solid #007acc;
            padding: 10px;
            margin: 10px 0;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div></div>
            <div class="lang-switch">
                <a href="#" class="lang-btn" id="lang-zh" onclick="switchLanguage('zh')">中文</a>
                <a href="#" class="lang-btn" id="lang-en" onclick="switchLanguage('en')">English</a>
            </div>
        </div>

        <!-- 中文内容 -->
        <div id="content-zh">
            <h1>🌐 简单Web文件服务器</h1>

            <p>欢迎使用这个简单而实用的Web文件服务器！这是一个用Go语言编写的轻量级文件管理工具。</p>

            <h2>📋 项目功能</h2>

            <div class="feature">
                <strong>📁 静态文件服务</strong><br>
                自动服务web目录下的所有文件，支持HTML、CSS、JavaScript、图片等各种文件类型。
            </div>

            <div class="feature">
                <strong>📤 文件上传</strong><br>
                <span id="upload-feature-description-zh">文件上传功能可通过命令行参数启用，确保服务器安全。</span>
                <span id="upload-status" class="status-indicator loading">🔄 检查中...</span>
            </div>

            <div class="feature">
                <strong>🌐 WebDAV服务</strong><br>
                <span id="webdav-feature-description-zh">WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。</span>
                <span id="webdav-status" class="status-indicator loading">🔄 检查中...</span>
                <div id="webdav-info" class="webdav-info hidden">
                    <strong>WebDAV地址:</strong> <code id="webdav-url">http://localhost:8081/webdav</code><br>
                    <strong>挂载目录:</strong> <code id="webdav-directory">.</code><br>
                    <strong>访问模式:</strong> <span id="webdav-mode">读写</span>
                </div>
            </div>

            <div class="feature">
                <strong>🔒 HTTPS服务</strong><br>
                <span id="https-feature-description-zh">HTTPS服务可通过命令行参数启用，提供加密的安全连接。</span>
                <span id="https-status" class="status-indicator loading">🔄 检查中...</span>
                <div id="https-info" class="webdav-info hidden">
                    <strong>HTTP地址:</strong> <code id="http-url">http://localhost:8080</code><br>
                    <strong>HTTPS地址:</strong> <code id="https-url">https://localhost:8443</code><br>
                    <strong>证书目录:</strong> <code id="cert-directory">./cert</code><br>
                    <strong>证书状态:</strong> <span id="cert-status">检查中</span>
                </div>
            </div>

            <div class="feature">
                <strong>📂 文件浏览下载</strong><br>
                <span id="files-feature-description-zh">专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。</span>
                <span id="files-status" class="status-indicator loading">🔄 检查中...</span>
            </div>

            <div class="feature">
                <strong>🔌 SOCKS5代理服务</strong><br>
                <span id="socks5-feature-description-zh">高性能SOCKS5代理服务器，支持TCP连接代理和Basic认证。</span>
                <span id="socks5-status" class="status-indicator loading">🔄 检查中...</span>
                <div id="socks5-info" class="webdav-info hidden">
                    <strong>代理地址:</strong> <code id="socks5-address">localhost:1080</code><br>
                    <strong>协议版本:</strong> <span>SOCKS5</span><br>
                    <strong>认证方式:</strong> <span id="socks5-auth-mode">无认证</span>
                </div>
            </div>

            <div class="feature">
                <strong>🌐 HTTP代理服务</strong><br>
                <span id="proxy-feature-description-zh">高性能HTTP/HTTPS代理服务器，支持Web浏览器代理和Basic认证。</span>
                <span id="proxy-status" class="status-indicator loading">🔄 检查中...</span>
                <div id="proxy-info" class="webdav-info hidden">
                    <strong>代理地址:</strong> <code id="proxy-address">localhost:10808</code><br>
                    <strong>支持协议:</strong> <span>HTTP, HTTPS</span><br>
                    <strong>连接方式:</strong> <span>CONNECT, GET, POST</span><br>
                    <strong>认证方式:</strong> <span id="proxy-auth-mode">无认证</span>
                </div>
            </div>

            <div class="feature">
                <strong>📊 功能状态展示</strong><br>
                默认页面实时显示服务器各项功能的启用状态，包括文件上传、WebDAV、HTTPS和代理服务。
            </div>

            <div class="feature">
                <strong>🔧 自动配置</strong><br>
                自动创建必要的目录结构，无需手动配置即可使用。
            </div>

            <h2>🚀 快速开始</h2>

            <div id="upload-section">
                <p><strong>文件上传：</strong></p>
                <div id="upload-enabled-content" class="hidden">
                    <a href="/upload" class="button" id="upload-button">📤 上传文件</a>
                </div>
                <div id="upload-disabled-content" class="hidden">
                    <p>要启用文件上传功能，请使用以下命令启动服务器：</p>
                    <code>sweb.exe -upload</code> 或 <code>sweb.exe --enable-upload</code>
                    <br><br>
                    <a href="/upload" class="button disabled" id="upload-button-disabled">📤 上传功能已禁用</a>
                </div>
            </div>

            <div id="webdav-section">
                <p><strong>WebDAV服务：</strong></p>
                <div id="webdav-enabled-content" class="hidden">
                    <a href="/webdav" class="button" id="webdav-button">🌐 访问WebDAV</a>
                    <p><small>可以在文件管理器中添加网络位置：<code id="webdav-mount-url">http://localhost:8081/webdav</code></small></p>
                </div>
                <div id="webdav-disabled-content" class="hidden">
                    <p>要启用WebDAV服务，请使用以下命令启动服务器：</p>
                    <code>sweb.exe -webdav</code> 或 <code>sweb.exe --enable-webdav</code>
                    <br><br>
                    <a href="/webdav" class="button disabled" id="webdav-button-disabled">🌐 WebDAV服务已禁用</a>
                </div>
            </div>

            <div id="files-section">
                <p><strong>文件浏览：</strong></p>
                <div id="files-enabled-content" class="hidden">
                    <a href="/files" class="button" id="files-button">📂 浏览文件</a>
                    <p><small>查看和下载web目录中的所有文件</small></p>
                </div>
                <div id="files-disabled-content" class="hidden">
                    <p>要启用文件浏览功能，请使用以下命令启动服务器：</p>
                    <code>sweb.exe -files</code> 或 <code>sweb.exe --enable-files</code>
                    <br><br>
                    <a href="/files" class="button disabled" id="files-button-disabled">📂 文件浏览已禁用</a>
                </div>
            </div>

            <p><strong>服务器信息：</strong></p>
            <ul>
                <li>服务端口: <code>8080</code></li>
                <li>文件目录: <code>./web</code></li>
                <li>上传地址: <code>/upload</code></li>
                <li>文件浏览: <code>/files</code></li>
                <li>WebDAV地址: <code>/webdav</code></li>
                <li>SOCKS5代理: <code>localhost:1080</code></li>
                <li>HTTP代理: <code>localhost:10808</code></li>
            </ul>

            <h2>💡 使用说明</h2>
            <div id="usage-instructions">
                <div id="usage-upload-enabled" class="hidden">
                    <h3>文件上传</h3>
                    <ol>
                        <li><strong>上传文件</strong>：点击上方的"上传文件"按钮，选择要上传的文件</li>
                        <li><strong>访问文件</strong>：上传成功后，文件将保存在web目录下，可以直接通过URL访问</li>
                        <li><strong>管理文件</strong>：所有上传的文件都会显示在主页的文件列表中</li>
                    </ol>
                </div>
                <div id="usage-webdav-enabled" class="hidden">
                    <h3>WebDAV服务</h3>
                    <ol>
                        <li><strong>Windows</strong>：在文件资源管理器中，右键"此电脑" → "映射网络驱动器" → 输入WebDAV地址</li>
                        <li><strong>macOS</strong>：在Finder中，按Cmd+K → 输入WebDAV地址</li>
                        <li><strong>Linux</strong>：使用davfs2或其他WebDAV客户端挂载</li>
                        <li><strong>移动设备</strong>：使用支持WebDAV的文件管理应用</li>
                    </ol>
                </div>
                <div id="usage-disabled" class="hidden">
                    <ol>
                        <li><strong>启用功能</strong>：使用相应的命令行参数启动服务器</li>
                        <li><strong>浏览文件</strong>：当前可以浏览和下载web目录中的现有文件</li>
                        <li><strong>安全考虑</strong>：高级功能默认禁用，确保服务器安全</li>
                    </ol>
                </div>
            </div>

            <h2>🛠️ 技术特性</h2>

            <ul>
                <li>使用Go语言标准库开发，轻量级无外部依赖</li>
                <li>支持多部分表单数据上传</li>
                <li>完整的WebDAV协议支持（RFC 4918）</li>
                <li>高性能SOCKS5代理服务器（RFC 1928）</li>
                <li>HTTP/HTTPS代理服务器支持</li>
                <li>Basic认证保护（WebDAV、SOCKS5、HTTP代理）</li>
                <li>可配置的读写权限控制</li>
                <li>自动MIME类型检测</li>
                <li>UTF-8编码支持，完美处理中文</li>
                <li>跨平台兼容（Windows、Linux、macOS）</li>
                <li>并发连接处理，高性能数据转发</li>
                <li>多语言界面支持（中文/英文）</li>
            </ul>

            <div class="footer">
                <p>🔗 <strong>简单Web文件服务器</strong> | 基于Go语言开发</p>
                <div id="footer-content">
                    <div id="footer-enabled" class="hidden">
                        <p>开始使用：
                            <span id="footer-files-link" class="hidden"><a href="/files" class="button">浏览文件</a></span>
                            <span id="footer-upload-link" class="hidden"><a href="/upload" class="button">上传文件</a></span>
                            <span id="footer-webdav-link" class="hidden"><a href="/webdav" class="button">访问WebDAV</a></span>
                        </p>
                    </div>
                    <div id="footer-disabled" class="hidden">
                        <p>安全模式：高级功能已禁用</p>
                        <p>使用 <code>sweb.exe -help</code> 查看所有可用选项</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 英文内容 -->
        <div id="content-en" class="hidden">
            <h1>🌐 Simple Web File Server</h1>

            <p>Welcome to this simple and practical Web file server! This is a lightweight file management tool written in Go.</p>

            <h2>📋 Project Features</h2>

            <div class="feature">
                <strong>📁 Static File Service</strong><br>
                Automatically serves all files in the web directory, supporting HTML, CSS, JavaScript, images and other file types.
            </div>

            <div class="feature">
                <strong>📤 File Upload</strong><br>
                <span id="upload-feature-description-en">File upload feature can be enabled via command line parameters to ensure server security.</span>
                <span id="upload-status-en" class="status-indicator loading">🔄 Checking...</span>
            </div>

            <div class="feature">
                <strong>🌐 WebDAV Service</strong><br>
                <span id="webdav-feature-description-en">WebDAV service can be enabled via command line parameters, supporting file management client connections.</span>
                <span id="webdav-status-en" class="status-indicator loading">🔄 Checking...</span>
                <div id="webdav-info-en" class="webdav-info hidden">
                    <strong>WebDAV URL:</strong> <code id="webdav-url-en">http://localhost:8081/webdav</code><br>
                    <strong>Mount Directory:</strong> <code id="webdav-directory-en">.</code><br>
                    <strong>Access Mode:</strong> <span id="webdav-mode-en">Read-Write</span>
                </div>
            </div>

            <div class="feature">
                <strong>🔒 HTTPS Service</strong><br>
                <span id="https-feature-description-en">HTTPS service can be enabled via command line parameters, providing encrypted secure connections.</span>
                <span id="https-status-en" class="status-indicator loading">🔄 Checking...</span>
                <div id="https-info-en" class="webdav-info hidden">
                    <strong>HTTP URL:</strong> <code id="http-url-en">http://localhost:8080</code><br>
                    <strong>HTTPS URL:</strong> <code id="https-url-en">https://localhost:8443</code><br>
                    <strong>Certificate Directory:</strong> <code id="cert-directory-en">./cert</code><br>
                    <strong>Certificate Status:</strong> <span id="cert-status-en">Checking</span>
                </div>
            </div>

            <div class="feature">
                <strong>📂 File Browsing & Download</strong><br>
                <span id="files-feature-description-en">Dedicated file browsing page showing all files in the web directory, supporting online browsing and downloading.</span>
                <span id="files-status-en" class="status-indicator loading">🔄 Checking...</span>
            </div>

            <div class="feature">
                <strong>🔌 SOCKS5 Proxy Service</strong><br>
                <span id="socks5-feature-description-en">High-performance SOCKS5 proxy server supporting TCP connection proxy and Basic authentication.</span>
                <span id="socks5-status-en" class="status-indicator loading">🔄 Checking...</span>
                <div id="socks5-info-en" class="webdav-info hidden">
                    <strong>Proxy Address:</strong> <code id="socks5-address-en">localhost:1080</code><br>
                    <strong>Protocol Version:</strong> <span>SOCKS5</span><br>
                    <strong>Authentication:</strong> <span id="socks5-auth-mode-en">No Authentication</span>
                </div>
            </div>

            <div class="feature">
                <strong>🌐 HTTP Proxy Service</strong><br>
                <span id="proxy-feature-description-en">High-performance HTTP/HTTPS proxy server supporting web browser proxy and Basic authentication.</span>
                <span id="proxy-status-en" class="status-indicator loading">🔄 Checking...</span>
                <div id="proxy-info-en" class="webdav-info hidden">
                    <strong>Proxy Address:</strong> <code id="proxy-address-en">localhost:10808</code><br>
                    <strong>Supported Protocols:</strong> <span>HTTP, HTTPS</span><br>
                    <strong>Connection Methods:</strong> <span>CONNECT, GET, POST</span><br>
                    <strong>Authentication:</strong> <span id="proxy-auth-mode-en">No Authentication</span>
                </div>
            </div>

            <div class="feature">
                <strong>📊 Feature Status Display</strong><br>
                The default page displays real-time status of server features including file upload, WebDAV, HTTPS and proxy services.
            </div>

            <div class="feature">
                <strong>🔧 Auto Configuration</strong><br>
                Automatically creates necessary directory structure, ready to use without manual configuration.
            </div>

            <h2>🚀 Quick Start</h2>

            <div id="upload-section-en">
                <p><strong>File Upload:</strong></p>
                <div id="upload-enabled-content-en" class="hidden">
                    <a href="/upload" class="button" id="upload-button-en">📤 Upload Files</a>
                </div>
                <div id="upload-disabled-content-en" class="hidden">
                    <p>To enable file upload feature, start the server with:</p>
                    <code>sweb.exe -upload</code> or <code>sweb.exe --enable-upload</code>
                    <br><br>
                    <a href="/upload" class="button disabled" id="upload-button-disabled-en">📤 Upload Feature Disabled</a>
                </div>
            </div>

            <div id="webdav-section-en">
                <p><strong>WebDAV Service:</strong></p>
                <div id="webdav-enabled-content-en" class="hidden">
                    <a href="/webdav" class="button" id="webdav-button-en">🌐 Access WebDAV</a>
                    <p><small>You can add network location in file manager: <code id="webdav-mount-url-en">http://localhost:8081/webdav</code></small></p>
                </div>
                <div id="webdav-disabled-content-en" class="hidden">
                    <p>To enable WebDAV service, start the server with:</p>
                    <code>sweb.exe -webdav</code> or <code>sweb.exe --enable-webdav</code>
                    <br><br>
                    <a href="/webdav" class="button disabled" id="webdav-button-disabled-en">🌐 WebDAV Service Disabled</a>
                </div>
            </div>

            <div id="files-section-en">
                <p><strong>File Browsing:</strong></p>
                <div id="files-enabled-content-en" class="hidden">
                    <a href="/files" class="button" id="files-button-en">📂 Browse Files</a>
                    <p><small>View and download all files in the web directory</small></p>
                </div>
                <div id="files-disabled-content-en" class="hidden">
                    <p>To enable file browsing feature, start the server with:</p>
                    <code>sweb.exe -files</code> or <code>sweb.exe --enable-files</code>
                    <br><br>
                    <a href="/files" class="button disabled" id="files-button-disabled-en">📂 File Browsing Disabled</a>
                </div>
            </div>

            <p><strong>Server Information:</strong></p>
            <ul>
                <li>Service Port: <code>8080</code></li>
                <li>File Directory: <code>./web</code></li>
                <li>Upload URL: <code>/upload</code></li>
                <li>File Browsing: <code>/files</code></li>
                <li>WebDAV URL: <code>/webdav</code></li>
                <li>SOCKS5 Proxy: <code>localhost:1080</code></li>
                <li>HTTP Proxy: <code>localhost:10808</code></li>
            </ul>

            <h2>💡 Usage Instructions</h2>
            <div id="usage-instructions-en">
                <div id="usage-upload-enabled-en" class="hidden">
                    <h3>File Upload</h3>
                    <ol>
                        <li><strong>Upload Files</strong>: Click the "Upload Files" button above and select files to upload</li>
                        <li><strong>Access Files</strong>: After successful upload, files are saved in the web directory and can be accessed directly via URL</li>
                        <li><strong>Manage Files</strong>: All uploaded files are displayed in the file list on the homepage</li>
                    </ol>
                </div>
                <div id="usage-webdav-enabled-en" class="hidden">
                    <h3>WebDAV Service</h3>
                    <ol>
                        <li><strong>Windows</strong>: In File Explorer, right-click "This PC" → "Map network drive" → Enter WebDAV address</li>
                        <li><strong>macOS</strong>: In Finder, press Cmd+K → Enter WebDAV address</li>
                        <li><strong>Linux</strong>: Use davfs2 or other WebDAV clients to mount</li>
                        <li><strong>Mobile Devices</strong>: Use file management apps that support WebDAV</li>
                    </ol>
                </div>
                <div id="usage-disabled-en" class="hidden">
                    <ol>
                        <li><strong>Enable Features</strong>: Start the server with corresponding command line parameters</li>
                        <li><strong>Browse Files</strong>: Currently you can browse and download existing files in the web directory</li>
                        <li><strong>Security Considerations</strong>: Advanced features are disabled by default to ensure server security</li>
                    </ol>
                </div>
            </div>

            <h2>🛠️ Technical Features</h2>

            <ul>
                <li>Developed using Go standard library, lightweight with no external dependencies</li>
                <li>Support for multipart form data upload</li>
                <li>Complete WebDAV protocol support (RFC 4918)</li>
                <li>High-performance SOCKS5 proxy server (RFC 1928)</li>
                <li>HTTP/HTTPS proxy server support</li>
                <li>Basic authentication protection (WebDAV, SOCKS5, HTTP proxy)</li>
                <li>Configurable read-write permission control</li>
                <li>Automatic MIME type detection</li>
                <li>UTF-8 encoding support, perfect Chinese character handling</li>
                <li>Cross-platform compatibility (Windows, Linux, macOS)</li>
                <li>Concurrent connection handling, high-performance data forwarding</li>
                <li>Multi-language interface support (Chinese/English)</li>
            </ul>

            <div class="footer">
                <p>🔗 <strong>Simple Web File Server</strong> | Developed in Go</p>
                <div id="footer-content-en">
                    <div id="footer-enabled-en" class="hidden">
                        <p>Get Started:
                            <span id="footer-files-link-en" class="hidden"><a href="/files" class="button">Browse Files</a></span>
                            <span id="footer-upload-link-en" class="hidden"><a href="/upload" class="button">Upload Files</a></span>
                            <span id="footer-webdav-link-en" class="hidden"><a href="/webdav" class="button">Access WebDAV</a></span>
                        </p>
                    </div>
                    <div id="footer-disabled-en" class="hidden">
                        <p>Security Mode: Advanced features are disabled</p>
                        <p>Use <code>sweb.exe -help</code> to see all available options</p>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        // 语言切换功能
        let currentLang = 'zh';

        // 检测浏览器语言
        function detectBrowserLanguage() {
            const browserLang = navigator.language || navigator.userLanguage;
            if (browserLang.startsWith('zh')) {
                return 'zh';
            } else {
                return 'en';
            }
        }

        // 切换语言
        function switchLanguage(lang) {
            currentLang = lang;
            localStorage.setItem('preferred-language', lang);

            // 更新页面标题
            const title = lang === 'zh' ? '简单Web文件服务器' : 'Simple Web File Server';
            document.getElementById('page-title').textContent = title;
            document.title = title;

            // 显示/隐藏对应语言内容
            document.getElementById('content-zh').style.display = lang === 'zh' ? 'block' : 'none';
            document.getElementById('content-en').style.display = lang === 'en' ? 'block' : 'none';

            // 更新语言按钮状态
            document.getElementById('lang-zh').classList.toggle('active', lang === 'zh');
            document.getElementById('lang-en').classList.toggle('active', lang === 'en');

            // 重新检查服务状态
            checkServiceStatus();
        }

        // 初始化语言
        function initLanguage() {
            const savedLang = localStorage.getItem('preferred-language');
            const detectedLang = detectBrowserLanguage();
            const initialLang = savedLang || detectedLang;
            switchLanguage(initialLang);
        }

        // 检查服务状态
        function checkServiceStatus() {
            fetch('/api/status')
                .then(response => response.json())
                .then(data => {
                    updateUploadStatus(data.upload);
                    updateFilesStatus(data.files);
                    updateWebDAVStatus(data.webdav);
                    updateHTTPSStatus(data.https);
                    updateSOCKS5Status(data.socks5);
                    updateProxyStatus(data.proxy);
                    updateUsageInstructions(data.upload, data.files, data.webdav, data.https, data.socks5, data.proxy);
                })
                .catch(error => {
                    console.error('检查服务状态失败:', error);
                    // 如果API调用失败，显示默认的禁用状态
                    updateUploadStatus({enabled: false, status: 'disabled'});
                    updateFilesStatus({enabled: false, status: 'disabled'});
                    updateWebDAVStatus({enabled: false, status: 'disabled'});
                    updateHTTPSStatus({enabled: false, status: 'disabled'});
                    updateSOCKS5Status({enabled: false, status: 'disabled'});
                    updateProxyStatus({enabled: false, status: 'disabled'});
                    updateUsageInstructions({enabled: false}, {enabled: false}, {enabled: false}, {enabled: false}, {enabled: false}, {enabled: false});
                });
        }

        // 更新页面上的上传状态显示
        function updateUploadStatus(uploadData) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('upload-status');
            const statusElementEn = document.getElementById('upload-status-en');
            const descriptionElement = document.getElementById('upload-feature-description-' + currentLang);
            const uploadEnabledContent = document.getElementById('upload-enabled-content');
            const uploadDisabledContent = document.getElementById('upload-disabled-content');
            const uploadEnabledContentEn = document.getElementById('upload-enabled-content-en');
            const uploadDisabledContentEn = document.getElementById('upload-disabled-content-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';
            const enabledDesc = isZh ? '通过简单的Web界面上传文件到服务器，支持各种文件格式。' : 'Upload files to server through simple web interface, supporting various file formats.';
            const disabledDesc = isZh ? '文件上传功能可通过命令行参数启用，确保服务器安全。' : 'File upload feature can be enabled via command line parameters to ensure server security.';

            if (uploadData.enabled) {
                if (statusElement) {
                    statusElement.textContent = enabledText;
                    statusElement.className = 'status-indicator status-enabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = enabledText;
                    statusElementEn.className = 'status-indicator status-enabled';
                }
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                if (uploadEnabledContent) uploadEnabledContent.classList.remove('hidden');
                if (uploadDisabledContent) uploadDisabledContent.classList.add('hidden');
                if (uploadEnabledContentEn) uploadEnabledContentEn.classList.remove('hidden');
                if (uploadDisabledContentEn) uploadDisabledContentEn.classList.add('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (uploadEnabledContent) uploadEnabledContent.classList.add('hidden');
                if (uploadDisabledContent) uploadDisabledContent.classList.remove('hidden');
                if (uploadEnabledContentEn) uploadEnabledContentEn.classList.add('hidden');
                if (uploadDisabledContentEn) uploadDisabledContentEn.classList.remove('hidden');
            }
        }

        // 更新页面上的文件浏览状态显示
        function updateFilesStatus(filesData) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('files-status');
            const statusElementEn = document.getElementById('files-status-en');
            const descriptionElement = document.getElementById('files-feature-description-' + currentLang);
            const filesEnabledContent = document.getElementById('files-enabled-content');
            const filesDisabledContent = document.getElementById('files-disabled-content');
            const filesEnabledContentEn = document.getElementById('files-enabled-content-en');
            const filesDisabledContentEn = document.getElementById('files-disabled-content-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';
            const enabledDesc = isZh ? '专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。' : 'Dedicated file browsing page showing all files in the web directory, supporting online browsing and downloading.';
            const disabledDesc = isZh ? '文件浏览功能可通过命令行参数启用，确保服务器安全。' : 'File browsing feature can be enabled via command line parameters to ensure server security.';

            if (filesData.enabled) {
                if (statusElement) {
                    statusElement.textContent = enabledText;
                    statusElement.className = 'status-indicator status-enabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = enabledText;
                    statusElementEn.className = 'status-indicator status-enabled';
                }
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                if (filesEnabledContent) filesEnabledContent.classList.remove('hidden');
                if (filesDisabledContent) filesDisabledContent.classList.add('hidden');
                if (filesEnabledContentEn) filesEnabledContentEn.classList.remove('hidden');
                if (filesDisabledContentEn) filesDisabledContentEn.classList.add('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (filesEnabledContent) filesEnabledContent.classList.add('hidden');
                if (filesDisabledContent) filesDisabledContent.classList.remove('hidden');
                if (filesEnabledContentEn) filesEnabledContentEn.classList.add('hidden');
                if (filesDisabledContentEn) filesDisabledContentEn.classList.remove('hidden');
            }
        }

        // 更新页面上的WebDAV状态显示
        function updateWebDAVStatus(webdavData) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('webdav-status');
            const statusElementEn = document.getElementById('webdav-status-en');
            const descriptionElement = document.getElementById('webdav-feature-description-' + currentLang);
            const webdavEnabledContent = document.getElementById('webdav-enabled-content');
            const webdavDisabledContent = document.getElementById('webdav-disabled-content');
            const webdavEnabledContentEn = document.getElementById('webdav-enabled-content-en');
            const webdavDisabledContentEn = document.getElementById('webdav-disabled-content-en');
            const webdavInfo = document.getElementById('webdav-info');
            const webdavInfoEn = document.getElementById('webdav-info-en');
            const webdavDirectory = document.getElementById('webdav-directory');
            const webdavDirectoryEn = document.getElementById('webdav-directory-en');
            const webdavMode = document.getElementById('webdav-mode');
            const webdavModeEn = document.getElementById('webdav-mode-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';
            const readonlyText = isZh ? '📖 只读模式' : '📖 Read-Only Mode';

            if (webdavData.enabled) {
                if (webdavData.readonly) {
                    if (statusElement) {
                        statusElement.textContent = readonlyText;
                        statusElement.className = 'status-indicator status-readonly';
                    }
                    if (statusElementEn) {
                        statusElementEn.textContent = readonlyText;
                        statusElementEn.className = 'status-indicator status-readonly';
                    }
                    if (webdavMode) webdavMode.textContent = isZh ? '只读' : 'Read-Only';
                    if (webdavModeEn) webdavModeEn.textContent = isZh ? '只读' : 'Read-Only';
                } else {
                    if (statusElement) {
                        statusElement.textContent = enabledText;
                        statusElement.className = 'status-indicator status-enabled';
                    }
                    if (statusElementEn) {
                        statusElementEn.textContent = enabledText;
                        statusElementEn.className = 'status-indicator status-enabled';
                    }
                    if (webdavMode) webdavMode.textContent = isZh ? '读写' : 'Read-Write';
                    if (webdavModeEn) webdavModeEn.textContent = isZh ? '读写' : 'Read-Write';
                }

                const enabledDesc = isZh ? 'WebDAV服务已启用，支持文件管理客户端连接和文件操作。' : 'WebDAV service is enabled, supporting file management client connections and file operations.';
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                if (webdavDirectory) webdavDirectory.textContent = webdavData.directory || '.';
                if (webdavDirectoryEn) webdavDirectoryEn.textContent = webdavData.directory || '.';

                // 动态更新WebDAV URL（使用端口信息）
                const webdavPort = webdavData.port || 8081;
                const webdavUrlText = ` + "`" + `http://localhost:${webdavPort}/webdav` + "`" + `;
                const webdavUrl = document.getElementById('webdav-url');
                const webdavUrlEn = document.getElementById('webdav-url-en');
                const webdavMountUrl = document.getElementById('webdav-mount-url');
                const webdavMountUrlEn = document.getElementById('webdav-mount-url-en');
                if (webdavUrl) webdavUrl.textContent = webdavUrlText;
                if (webdavUrlEn) webdavUrlEn.textContent = webdavUrlText;
                if (webdavMountUrl) webdavMountUrl.textContent = webdavUrlText;
                if (webdavMountUrlEn) webdavMountUrlEn.textContent = webdavUrlText;

                if (webdavEnabledContent) webdavEnabledContent.classList.remove('hidden');
                if (webdavDisabledContent) webdavDisabledContent.classList.add('hidden');
                if (webdavEnabledContentEn) webdavEnabledContentEn.classList.remove('hidden');
                if (webdavDisabledContentEn) webdavDisabledContentEn.classList.add('hidden');
                if (webdavInfo) webdavInfo.classList.remove('hidden');
                if (webdavInfoEn) webdavInfoEn.classList.remove('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }

                const disabledDesc = isZh ? 'WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。' : 'WebDAV service can be enabled via command line parameters, supporting file management client connections.';
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (webdavEnabledContent) webdavEnabledContent.classList.add('hidden');
                if (webdavDisabledContent) webdavDisabledContent.classList.remove('hidden');
                if (webdavEnabledContentEn) webdavEnabledContentEn.classList.add('hidden');
                if (webdavDisabledContentEn) webdavDisabledContentEn.classList.remove('hidden');
                if (webdavInfo) webdavInfo.classList.add('hidden');
                if (webdavInfoEn) webdavInfoEn.classList.add('hidden');
            }
        }

        // 更新HTTPS状态显示
        function updateHTTPSStatus(httpsData) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('https-status');
            const statusElementEn = document.getElementById('https-status-en');
            const descriptionElement = document.getElementById('https-feature-description-' + currentLang);
            const httpsInfo = document.getElementById('https-info');
            const httpsInfoEn = document.getElementById('https-info-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';

            if (httpsData.enabled) {
                if (statusElement) {
                    statusElement.textContent = enabledText;
                    statusElement.className = 'status-indicator status-enabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = enabledText;
                    statusElementEn.className = 'status-indicator status-enabled';
                }

                const enabledDesc = isZh ? 'HTTPS服务已启用，提供加密的安全连接。' : 'HTTPS service is enabled, providing encrypted secure connections.';
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                if (httpsInfo) httpsInfo.classList.remove('hidden');
                if (httpsInfoEn) httpsInfoEn.classList.remove('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }

                const disabledDesc = isZh ? 'HTTPS服务可通过命令行参数启用，提供加密的安全连接。' : 'HTTPS service can be enabled via command line parameters, providing encrypted secure connections.';
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (httpsInfo) httpsInfo.classList.add('hidden');
                if (httpsInfoEn) httpsInfoEn.classList.add('hidden');
            }
        }

        // 更新SOCKS5状态显示
        function updateSOCKS5Status(socks5Data) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('socks5-status');
            const statusElementEn = document.getElementById('socks5-status-en');
            const descriptionElement = document.getElementById('socks5-feature-description-' + currentLang);
            const socks5Info = document.getElementById('socks5-info');
            const socks5InfoEn = document.getElementById('socks5-info-en');
            const socks5AuthMode = document.getElementById('socks5-auth-mode');
            const socks5AuthModeEn = document.getElementById('socks5-auth-mode-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';

            if (socks5Data.enabled) {
                if (statusElement) {
                    statusElement.textContent = enabledText;
                    statusElement.className = 'status-indicator status-enabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = enabledText;
                    statusElementEn.className = 'status-indicator status-enabled';
                }

                const enabledDesc = isZh ? 'SOCKS5代理服务已启用，支持TCP连接代理和Basic认证。' : 'SOCKS5 proxy service is enabled, supporting TCP connection proxy and Basic authentication.';
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                const authText = socks5Data.auth ? (isZh ? 'Basic认证' : 'Basic Authentication') : (isZh ? '无认证' : 'No Authentication');
                if (socks5AuthMode) socks5AuthMode.textContent = authText;
                if (socks5AuthModeEn) socks5AuthModeEn.textContent = authText;

                if (socks5Info) socks5Info.classList.remove('hidden');
                if (socks5InfoEn) socks5InfoEn.classList.remove('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }

                const disabledDesc = isZh ? 'SOCKS5代理服务可通过命令行参数启用，支持TCP连接代理和Basic认证。' : 'SOCKS5 proxy service can be enabled via command line parameters, supporting TCP connection proxy and Basic authentication.';
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (socks5Info) socks5Info.classList.add('hidden');
                if (socks5InfoEn) socks5InfoEn.classList.add('hidden');
            }
        }

        // 更新HTTP代理状态显示
        function updateProxyStatus(proxyData) {
            const isZh = currentLang === 'zh';
            const statusElement = document.getElementById('proxy-status');
            const statusElementEn = document.getElementById('proxy-status-en');
            const descriptionElement = document.getElementById('proxy-feature-description-' + currentLang);
            const proxyInfo = document.getElementById('proxy-info');
            const proxyInfoEn = document.getElementById('proxy-info-en');
            const proxyAuthMode = document.getElementById('proxy-auth-mode');
            const proxyAuthModeEn = document.getElementById('proxy-auth-mode-en');

            const enabledText = isZh ? '✅ 已启用' : '✅ Enabled';
            const disabledText = isZh ? '🔒 已禁用' : '🔒 Disabled';

            if (proxyData.enabled) {
                if (statusElement) {
                    statusElement.textContent = enabledText;
                    statusElement.className = 'status-indicator status-enabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = enabledText;
                    statusElementEn.className = 'status-indicator status-enabled';
                }

                const enabledDesc = isZh ? 'HTTP代理服务已启用，支持Web浏览器代理和Basic认证。' : 'HTTP proxy service is enabled, supporting web browser proxy and Basic authentication.';
                if (descriptionElement) {
                    descriptionElement.textContent = enabledDesc;
                }

                const authText = proxyData.auth ? (isZh ? 'Basic认证' : 'Basic Authentication') : (isZh ? '无认证' : 'No Authentication');
                if (proxyAuthMode) proxyAuthMode.textContent = authText;
                if (proxyAuthModeEn) proxyAuthModeEn.textContent = authText;

                if (proxyInfo) proxyInfo.classList.remove('hidden');
                if (proxyInfoEn) proxyInfoEn.classList.remove('hidden');
            } else {
                if (statusElement) {
                    statusElement.textContent = disabledText;
                    statusElement.className = 'status-indicator status-disabled';
                }
                if (statusElementEn) {
                    statusElementEn.textContent = disabledText;
                    statusElementEn.className = 'status-indicator status-disabled';
                }

                const disabledDesc = isZh ? 'HTTP代理服务可通过命令行参数启用，支持Web浏览器代理和Basic认证。' : 'HTTP proxy service can be enabled via command line parameters, supporting web browser proxy and Basic authentication.';
                if (descriptionElement) {
                    descriptionElement.textContent = disabledDesc;
                }

                if (proxyInfo) proxyInfo.classList.add('hidden');
                if (proxyInfoEn) proxyInfoEn.classList.add('hidden');
            }
        }

        // 更新使用说明
        function updateUsageInstructions(uploadData, filesData, webdavData, httpsData, socks5Data, proxyData) {
            const usageUploadEnabled = document.getElementById('usage-upload-enabled');
            const usageWebdavEnabled = document.getElementById('usage-webdav-enabled');
            const usageDisabled = document.getElementById('usage-disabled');
            const usageUploadEnabledEn = document.getElementById('usage-upload-enabled-en');
            const usageWebdavEnabledEn = document.getElementById('usage-webdav-enabled-en');
            const usageDisabledEn = document.getElementById('usage-disabled-en');

            const footerEnabled = document.getElementById('footer-enabled');
            const footerDisabled = document.getElementById('footer-disabled');
            const footerEnabledEn = document.getElementById('footer-enabled-en');
            const footerDisabledEn = document.getElementById('footer-disabled-en');

            const footerFilesLink = document.getElementById('footer-files-link');
            const footerUploadLink = document.getElementById('footer-upload-link');
            const footerWebdavLink = document.getElementById('footer-webdav-link');
            const footerFilesLinkEn = document.getElementById('footer-files-link-en');
            const footerUploadLinkEn = document.getElementById('footer-upload-link-en');
            const footerWebdavLinkEn = document.getElementById('footer-webdav-link-en');

            if (uploadData.enabled || filesData.enabled || webdavData.enabled) {
                if (usageUploadEnabled) usageUploadEnabled.classList.toggle('hidden', !uploadData.enabled);
                if (usageWebdavEnabled) usageWebdavEnabled.classList.toggle('hidden', !webdavData.enabled);
                if (usageDisabled) usageDisabled.classList.add('hidden');
                if (usageUploadEnabledEn) usageUploadEnabledEn.classList.toggle('hidden', !uploadData.enabled);
                if (usageWebdavEnabledEn) usageWebdavEnabledEn.classList.toggle('hidden', !webdavData.enabled);
                if (usageDisabledEn) usageDisabledEn.classList.add('hidden');

                if (footerEnabled) footerEnabled.classList.remove('hidden');
                if (footerDisabled) footerDisabled.classList.add('hidden');
                if (footerEnabledEn) footerEnabledEn.classList.remove('hidden');
                if (footerDisabledEn) footerDisabledEn.classList.add('hidden');

                if (footerFilesLink) footerFilesLink.classList.toggle('hidden', !filesData.enabled);
                if (footerUploadLink) footerUploadLink.classList.toggle('hidden', !uploadData.enabled);
                if (footerWebdavLink) footerWebdavLink.classList.toggle('hidden', !webdavData.enabled);
                if (footerFilesLinkEn) footerFilesLinkEn.classList.toggle('hidden', !filesData.enabled);
                if (footerUploadLinkEn) footerUploadLinkEn.classList.toggle('hidden', !uploadData.enabled);
                if (footerWebdavLinkEn) footerWebdavLinkEn.classList.toggle('hidden', !webdavData.enabled);
            } else {
                if (usageUploadEnabled) usageUploadEnabled.classList.add('hidden');
                if (usageWebdavEnabled) usageWebdavEnabled.classList.add('hidden');
                if (usageDisabled) usageDisabled.classList.remove('hidden');
                if (usageUploadEnabledEn) usageUploadEnabledEn.classList.add('hidden');
                if (usageWebdavEnabledEn) usageWebdavEnabledEn.classList.add('hidden');
                if (usageDisabledEn) usageDisabledEn.classList.remove('hidden');

                if (footerEnabled) footerEnabled.classList.add('hidden');
                if (footerDisabled) footerDisabled.classList.remove('hidden');
                if (footerEnabledEn) footerEnabledEn.classList.add('hidden');
                if (footerDisabledEn) footerDisabledEn.classList.remove('hidden');
            }
        }

        // 页面加载时初始化
        document.addEventListener('DOMContentLoaded', function() {
            initLanguage();
            checkServiceStatus();

            // 每30秒检查一次状态
            setInterval(checkServiceStatus, 30000);
        });
    </script>
</body>
</html>`
}
