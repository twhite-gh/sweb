package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPSProxyStats HTTPS代理统计信息
type HTTPSProxyStats struct {
	mutex                 sync.RWMutex
	TotalConnections      int64
	ActiveConnections     int64
	SuccessfulConnections int64
	FailedConnections     int64
	HTTPSRequests         int64
	BytesTransferred      int64
	StartTime             time.Time
}

// HTTPSProxyServer HTTPS代理服务器
type HTTPSProxyServer struct {
	port     int
	auth     bool
	username string
	password string
	listener net.Listener
	running  bool
	stats    HTTPSProxyStats
	certFile string
	keyFile  string
}

// 全局HTTPS代理服务器实例
var httpsProxyServer *HTTPSProxyServer

// startHTTPSProxyServer 启动HTTPS代理服务器
func startHTTPSProxyServer(port int, auth bool, username, password, certFile, keyFile string) error {
	if httpsProxyServer != nil && httpsProxyServer.IsRunning() {
		return fmt.Errorf("HTTPS代理服务器已在运行")
	}

	server := &HTTPSProxyServer{
		port:     port,
		auth:     auth,
		username: username,
		password: password,
		certFile: certFile,
		keyFile:  keyFile,
	}

	// 加载TLS证书
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("无法加载HTTPS代理证书: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if err != nil {
		return fmt.Errorf("无法启动HTTPS代理服务器: %v", err)
	}

	server.listener = listener
	server.running = true
	server.stats.StartTime = time.Now()
	httpsProxyServer = server

	if auth {
		fmt.Printf("🔒 HTTPS代理服务器启动在端口 %d (认证: %s)\n", port, username)
	} else {
		fmt.Printf("🔒 HTTPS代理服务器启动在端口 %d (无认证)\n", port)
	}

	go server.acceptConnections()

	return nil
}

// stopHTTPSProxyServer 停止HTTPS代理服务器
func stopHTTPSProxyServer() error {
	if httpsProxyServer == nil {
		return fmt.Errorf("HTTPS代理服务器未运行")
	}
	return httpsProxyServer.Stop()
}

// IsRunning 检查服务器是否正在运行
func (s *HTTPSProxyServer) IsRunning() bool {
	return s.running
}

// Stop 停止服务器
func (s *HTTPSProxyServer) Stop() error {
	if !s.running {
		return fmt.Errorf("服务器未运行")
	}
	s.running = false
	return s.listener.Close()
}

// GetStats 获取统计信息
func (s *HTTPSProxyServer) GetStats() HTTPSProxyStats {
	s.stats.mutex.RLock()
	defer s.stats.mutex.RUnlock()
	statsCopy := s.stats
	return statsCopy
}

// acceptConnections 接受客户端连接
func (s *HTTPSProxyServer) acceptConnections() {
	defer func() {
		s.running = false
		if r := recover(); r != nil {
			log.Printf("HTTPS代理服务器panic: %v", r)
		}
	}()

	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running {
				log.Printf("HTTPS代理接受连接失败: %v", err)
			}
			continue
		}

		// 更新连接统计
		s.stats.mutex.Lock()
		s.stats.TotalConnections++
		s.stats.ActiveConnections++
		s.stats.mutex.Unlock()

		go s.handleConnection(conn)
	}
}

// handleConnection 处理单个客户端连接
func (s *HTTPSProxyServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.stats.mutex.Lock()
		s.stats.ActiveConnections--
		s.stats.mutex.Unlock()

		if r := recover(); r != nil {
			log.Printf("HTTPS代理连接处理panic: %v", r)
		}
	}()

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// 读取HTTP请求
	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("HTTPS代理读取请求失败: %v", err)
		s.stats.mutex.Lock()
		s.stats.FailedConnections++
		s.stats.mutex.Unlock()
		return
	}

	// 检查认证
	if s.auth {
		if !s.checkAuthentication(request) {
			// 发送407 Proxy Authentication Required响应
			response := "HTTP/1.1 407 Proxy Authentication Required\r\n" +
				"Proxy-Authenticate: Basic realm=\"HTTPS Proxy\"\r\n" +
				"Content-Length: 0\r\n\r\n"
			conn.Write([]byte(response))
			s.stats.mutex.Lock()
			s.stats.FailedConnections++
			s.stats.mutex.Unlock()
			return
		}
	}

	// 根据请求方法处理
	if request.Method == "CONNECT" {
		s.handleHTTPSConnect(conn, request)
	} else {
		s.handleHTTPSRequest(conn, request)
	}

	s.stats.mutex.Lock()
	s.stats.SuccessfulConnections++
	s.stats.mutex.Unlock()
}

// checkAuthentication 检查HTTPS代理认证
func (s *HTTPSProxyServer) checkAuthentication(request *http.Request) bool {
	// 获取Proxy-Authorization头
	authHeader := request.Header.Get("Proxy-Authorization")
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
	return user == s.username && pass == s.password
}

// handleHTTPSConnect 处理HTTPS CONNECT请求
func (s *HTTPSProxyServer) handleHTTPSConnect(clientConn net.Conn, request *http.Request) {
	// 更新HTTPS请求统计
	s.stats.mutex.Lock()
	s.stats.HTTPSRequests++
	s.stats.mutex.Unlock()

	// 解析目标地址
	targetAddr := request.Host
	if !strings.Contains(targetAddr, ":") {
		targetAddr += ":443"
	}

	// 连接到目标服务器
	targetConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		log.Printf("HTTPS代理连接目标服务器失败: %v", err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		s.stats.mutex.Lock()
		s.stats.FailedConnections++
		s.stats.mutex.Unlock()
		return
	}
	defer targetConn.Close()

	// 发送200 Connection Established响应
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// 开始双向数据转发
	s.relayData(clientConn, targetConn)
}

// handleHTTPSRequest 处理普通HTTPS请求
func (s *HTTPSProxyServer) handleHTTPSRequest(clientConn net.Conn, request *http.Request) {
	// 更新HTTPS请求统计
	s.stats.mutex.Lock()
	s.stats.HTTPSRequests++
	s.stats.mutex.Unlock()

	// 对于HTTPS代理，所有请求都应该是CONNECT方法
	// 如果不是CONNECT，返回错误
	response := "HTTP/1.1 405 Method Not Allowed\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 23\r\n\r\n" +
		"Method Not Allowed\r\n"
	clientConn.Write([]byte(response))
}

// relayData 在两个连接之间转发数据
func (s *HTTPSProxyServer) relayData(conn1, conn2 net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	// 从conn1到conn2
	go func() {
		defer wg.Done()
		bytes, _ := io.Copy(conn2, conn1)
		s.stats.mutex.Lock()
		s.stats.BytesTransferred += bytes
		s.stats.mutex.Unlock()
	}()

	// 从conn2到conn1
	go func() {
		defer wg.Done()
		bytes, _ := io.Copy(conn1, conn2)
		s.stats.mutex.Lock()
		s.stats.BytesTransferred += bytes
		s.stats.mutex.Unlock()
	}()

	wg.Wait()
}

// httpsProxyDisabledHandler 处理HTTPS代理功能被禁用时的请求
func httpsProxyDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>HTTPS代理服务已禁用</title>
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
                    font-size: 48px;
                    margin-bottom: 20px;
                }
                .command {
                    background: #f8f9fa;
                    padding: 10px;
                    border-radius: 5px;
                    font-family: monospace;
                    margin: 10px 0;
                    border-left: 3px solid #007acc;
                }
                .back-link {
                    display: inline-block;
                    margin-top: 20px;
                    color: #007acc;
                    text-decoration: none;
                    padding: 10px 20px;
                    border: 1px solid #007acc;
                    border-radius: 5px;
                    transition: all 0.3s;
                }
                .back-link:hover {
                    background: #007acc;
                    color: white;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="icon">🔒</div>
                <h1>HTTPS代理服务已禁用</h1>
                <p>出于安全考虑，HTTPS代理服务默认处于禁用状态。</p>
                <p>如需启用HTTPS代理服务，请使用以下命令重新启动服务器：</p>

                <div class="command">sweb.exe -https-proxy</div>
                <p>或</p>
                <div class="command">sweb.exe --enable-https-proxy</div>

                <p><strong>可选参数：</strong></p>
                <div class="command">sweb.exe -https-proxy -https-proxy-port 10443</div>
                <div class="command">sweb.exe -https-proxy -https-proxy-auth</div>

                <p>您也可以使用 <code>sweb.exe -help</code> 查看所有可用选项。</p>

                <a href="/" class="back-link">← 返回首页</a>
            </div>
        </body>
        </html>
    `))
}
