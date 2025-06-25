package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/webdav"
)

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

	// 如果是只读模式，包装处理器
	if webdavReadonly {
		http.HandleFunc("/webdav/", func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				handler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
		// 处理根路径
		http.HandleFunc("/webdav", func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				handler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
	} else {
		http.Handle("/webdav/", handler)
		http.Handle("/webdav", handler)
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
