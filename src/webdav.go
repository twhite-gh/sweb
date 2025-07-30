package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/webdav"
)

// checkWebDAVAuthentication 检查WebDAV认证
func checkWebDAVAuthentication(r *http.Request, username, password string) bool {
	// 获取Authorization头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	// 检查是否是Basic认证
	if !strings.HasPrefix(authHeader, "Basic ") {
		return false
	}

	// 解码Base64编码的认证信息
	encoded := authHeader[6:] // 去掉"Basic "前缀
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return false
	}

	// 解析用户名和密码
	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return false
	}

	user := parts[0]
	pass := parts[1]

	// 验证用户名和密码
	return user == username && pass == password
}

// webdavAuthWrapper 包装WebDAV处理器以添加认证
func webdavAuthWrapper(handler http.Handler, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 对于OPTIONS请求，Windows WebDAV客户端需要特殊处理
		if r.Method == "OPTIONS" {
			// 设置WebDAV相关的响应头
			w.Header().Set("DAV", "1, 2")
			w.Header().Set("MS-Author-Via", "DAV")
			w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, COPY, MOVE, MKCOL, PROPFIND, PROPPATCH, LOCK, UNLOCK")
			w.Header().Set("Public", "OPTIONS, GET, HEAD, POST, PUT, DELETE, TRACE, COPY, MOVE, MKCOL, PROPFIND, PROPPATCH, LOCK, UNLOCK")
			w.Header().Set("Cache-Control", "private")
			w.Header().Set("Server", "Microsoft-IIS/10.0")

			// 如果没有认证信息，要求认证
			if !checkWebDAVAuthentication(r, username, password) {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"WebDAV\", charset=\"UTF-8\"")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if !checkWebDAVAuthentication(r, username, password) {
			// Windows WebDAV客户端兼容性改进
			// 设置必要的响应头以确保Windows客户端正确处理认证
			w.Header().Set("WWW-Authenticate", "Basic realm=\"WebDAV\", charset=\"UTF-8\"")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			w.Header().Set("Connection", "close")
			w.Header().Set("Server", "Microsoft-IIS/10.0") // 伪装成IIS服务器以提高Windows兼容性

			// 发送401 Unauthorized响应
			w.WriteHeader(http.StatusUnauthorized)

			// 提供更详细的错误信息
			errorHTML := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>WebDAV Authentication Required</title>
</head>
<body>
    <h1>401 Unauthorized</h1>
    <p>This WebDAV resource requires authentication.</p>
    <p>Please provide valid credentials.</p>
</body>
</html>`
			w.Write([]byte(errorHTML))
			return
		}

		// 为所有成功的请求添加WebDAV相关头部
		w.Header().Set("DAV", "1, 2")
		w.Header().Set("MS-Author-Via", "DAV")
		w.Header().Set("Server", "Microsoft-IIS/10.0")

		handler.ServeHTTP(w, r)
	}
}

// setupWebDAVHandler 设置WebDAV处理器
func setupWebDAVHandler() {
	// 确保WebDAV目录存在
	if _, err := os.Stat(webdavDir); os.IsNotExist(err) {
		err := os.MkdirAll(webdavDir, 0755)
		if err != nil {
			log.Fatalf("无法创建WebDAV目录: %v", err)
		}
	}

	// 创建WebDAV处理器
	handler := &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(webdavDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				// 过滤掉一些常见的非关键错误
				errStr := err.Error()
				// 忽略文件不存在的PROPFIND错误（这在文件创建过程中是正常的）
				if r.Method == "PROPFIND" && (strings.Contains(errStr, "cannot find the file specified") ||
					strings.Contains(errStr, "no such file or directory") ||
					strings.Contains(errStr, "file does not exist")) {
					// 这些是正常的操作流程，不记录错误
					return
				}
				// 记录其他重要错误
				log.Printf("WebDAV操作: %s %s - %v", r.Method, r.URL.Path, err)
			}
		},
	}

	// 创建最终的处理器
	var finalHandler http.Handler = handler

	// 如果启用认证，包装认证处理器
	if webdavAuth {
		finalHandler = webdavAuthWrapper(finalHandler, webdavUsername, webdavPassword)
	}

	// 如果是只读模式，包装只读处理器
	if webdavReadonly {
		readonlyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				finalHandler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
		http.Handle("/webdav/", readonlyHandler)
		http.Handle("/webdav", readonlyHandler)
	} else {
		http.Handle("/webdav/", finalHandler)
		http.Handle("/webdav", finalHandler)
	}
}

// webdavDisabledHandler 处理WebDAV功能被禁用时的请求
func webdavDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>WebDAV服务已禁用</title>
            <style>
                body {
                    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                    max-width: 600px;
                    margin: 50px auto;
                    padding: 20px;
                    text-align: center;
                    background-color: #f8f9fa;
                }
                .container {
                    background: white;
                    padding: 40px;
                    border-radius: 10px;
                    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
                    border-left: 5px solid #dc3545;
                }
                h1 {
                    color: #dc3545;
                    margin-bottom: 20px;
                }
                .icon {
                    font-size: 64px;
                    margin-bottom: 20px;
                }
                .command {
                    background: #f8f9fa;
                    padding: 10px;
                    border-radius: 5px;
                    font-family: 'Courier New', monospace;
                    margin: 10px 0;
                    border: 1px solid #dee2e6;
                }
                .back-link {
                    display: inline-block;
                    background: #007acc;
                    color: white;
                    padding: 10px 20px;
                    text-decoration: none;
                    border-radius: 5px;
                    margin-top: 20px;
                }
                .back-link:hover {
                    background: #005a9e;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="icon">🔒</div>
                <h1>WebDAV服务已禁用</h1>
                <p>出于安全考虑，WebDAV服务默认处于禁用状态。</p>
                <p>如需启用WebDAV服务，请使用以下命令重新启动服务器：</p>

                <div class="command">sweb.exe -webdav</div>
                <p>或</p>
                <div class="command">sweb.exe --enable-webdav</div>

                <p><strong>可选参数：</strong></p>
                <div class="command">sweb.exe -webdav -webdav-dir /path/to/directory</div>
                <div class="command">sweb.exe -webdav -webdav-readonly</div>

                <p>您也可以使用 <code>sweb.exe -help</code> 查看所有可用选项。</p>

                <a href="/" class="back-link">← 返回首页</a>
            </div>
        </body>
        </html>
    `))
}

// WebDAV服务器结构
type WebDAVServer struct {
	port     int
	auth     bool
	username string
	password string
	dir      string
	readonly bool
	listener net.Listener
	running  bool
}

// 全局WebDAV服务器实例
var webdavServer *WebDAVServer

// startWebDAVServer 启动独立的WebDAV服务器
func startWebDAVServer(port int, auth bool, username, password, dir string, readonly bool) error {
	if webdavServer != nil && webdavServer.running {
		return fmt.Errorf("WebDAV服务器已在运行")
	}

	// 确保WebDAV目录存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("无法创建WebDAV目录: %v", err)
		}
	}

	// 创建WebDAV处理器
	handler := &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				// 过滤掉一些常见的非关键错误
				errStr := err.Error()
				// 忽略文件不存在的PROPFIND错误（这在文件创建过程中是正常的）
				if r.Method == "PROPFIND" && (strings.Contains(errStr, "cannot find the file specified") ||
					strings.Contains(errStr, "no such file or directory") ||
					strings.Contains(errStr, "file does not exist")) {
					// 这些是正常的操作流程，不记录错误
					return
				}
				// 记录其他重要错误
				log.Printf("WebDAV操作: %s %s - %v", r.Method, r.URL.Path, err)
			}
		},
	}

	// 创建最终的处理器
	var finalHandler http.Handler = handler

	// 如果启用认证，包装认证处理器
	if auth {
		finalHandler = webdavAuthWrapper(finalHandler, username, password)
	}

	// 如果是只读模式，包装只读处理器
	if readonly {
		readonlyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				finalHandler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
		finalHandler = readonlyHandler
	}

	// 创建HTTP服务器
	mux := http.NewServeMux()
	mux.Handle("/webdav/", finalHandler)
	mux.Handle("/webdav", finalHandler)

	// 添加根路径处理，显示WebDAV信息页面
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			webdavInfoHandler(w, r, port, auth, username, dir, readonly)
		} else {
			http.NotFound(w, r)
		}
	})

	// 启动服务器
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("无法启动WebDAV服务器: %v", err)
	}

	server := &WebDAVServer{
		port:     port,
		auth:     auth,
		username: username,
		password: password,
		dir:      dir,
		readonly: readonly,
		listener: listener,
		running:  true,
	}

	webdavServer = server

	// 显示启动信息
	mode := "读写模式"
	if readonly {
		mode = "只读模式"
	}
	authInfo := "无认证"
	if auth {
		authInfo = fmt.Sprintf("认证: %s", username)
	}
	fmt.Printf("📁 WebDAV服务器启动在 http://localhost:%d/webdav (%s, %s) - 目录: %s\n", port, mode, authInfo, dir)

	// 在独立的goroutine中运行服务器
	go func() {
		err := http.Serve(listener, mux)
		if err != nil && server.running {
			log.Printf("WebDAV服务器错误: %v", err)
		}
	}()

	return nil
}

// webdavInfoHandler 显示WebDAV服务信息页面
func webdavInfoHandler(w http.ResponseWriter, r *http.Request, port int, auth bool, username, dir string, readonly bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	mode := "读写模式"
	if readonly {
		mode = "只读模式"
	}

	authInfo := "无认证"
	authExample := ""
	if auth {
		authInfo = fmt.Sprintf("认证: %s", username)
		authExample = fmt.Sprintf(`
		<p><strong>认证信息:</strong></p>
		<ul>
			<li>用户名: <code>%s</code></li>
			<li>密码: <code>webdav</code></li>
		</ul>`, username)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebDAV 服务器</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .info-box { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .url { background: #e8f4fd; padding: 10px; border-radius: 4px; font-family: monospace; }
        .command { background: #2d3748; color: #e2e8f0; padding: 10px; border-radius: 4px; font-family: monospace; }
        .status { color: #38a169; font-weight: bold; }
        ul { margin: 10px 0; }
        li { margin: 5px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📁 WebDAV 服务器</h1>
        <p class="status">✅ 服务正在运行</p>
    </div>

    <div class="info-box">
        <h3>服务信息</h3>
        <ul>
            <li><strong>端口:</strong> %d</li>
            <li><strong>模式:</strong> %s</li>
            <li><strong>认证:</strong> %s</li>
            <li><strong>目录:</strong> %s</li>
        </ul>
    </div>

    <div class="info-box">
        <h3>WebDAV 访问地址</h3>
        <div class="url">http://localhost:%d/webdav</div>
        %s
    </div>

    <div class="info-box">
        <h3>Windows 挂载命令</h3>
        <div class="command">net use Z: http://localhost:%d/webdav /user:%s webdav</div>
    </div>

    <div class="info-box">
        <h3>使用说明</h3>
        <p>您可以通过以下方式访问WebDAV服务:</p>
        <ul>
            <li>Windows: 使用上面的命令挂载为网络驱动器</li>
            <li>文件资源管理器: 添加网络位置</li>
            <li>第三方客户端: WinSCP, CarotDAV 等</li>
        </ul>
    </div>
</body>
</html>`, port, mode, authInfo, dir, port, authExample, port, username)

	w.Write([]byte(html))
}
