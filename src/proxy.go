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

// HTTP代理服务器结构
type HTTPProxyServer struct {
	port     int
	listener net.Listener
	running  bool
	mutex    sync.RWMutex
	stats    HTTPProxyStats
	auth     bool
	username string
	password string
}

// HTTP代理统计信息
type HTTPProxyStats struct {
	TotalConnections      int64
	ActiveConnections     int64
	SuccessfulConnections int64
	FailedConnections     int64
	HTTPRequests          int64
	HTTPSRequests         int64
	BytesTransferred      int64
	mutex                 sync.RWMutex
}

// 全局HTTP代理服务器实例
var httpProxyServer *HTTPProxyServer

// startHTTPProxyServer 启动HTTP代理服务器
func startHTTPProxyServer(port int, auth bool, username, password string) error {
	if httpProxyServer != nil && httpProxyServer.IsRunning() {
		return fmt.Errorf("HTTP代理服务器已在运行")
	}

	server := &HTTPProxyServer{
		port:     port,
		auth:     auth,
		username: username,
		password: password,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("无法启动HTTP代理服务器: %v", err)
	}

	server.listener = listener
	server.running = true
	httpProxyServer = server

	if auth {
		fmt.Printf("🌐 HTTP代理服务器启动在端口 %d (认证: %s)\n", port, username)
	} else {
		fmt.Printf("🌐 HTTP代理服务器启动在端口 %d (无认证)\n", port)
	}

	go server.acceptConnections()

	return nil
}

// stopHTTPProxyServer 停止HTTP代理服务器
func stopHTTPProxyServer() error {
	if httpProxyServer == nil {
		return fmt.Errorf("HTTP代理服务器未运行")
	}

	httpProxyServer.mutex.Lock()
	defer httpProxyServer.mutex.Unlock()

	if !httpProxyServer.running {
		return fmt.Errorf("HTTP代理服务器未运行")
	}

	httpProxyServer.running = false
	if httpProxyServer.listener != nil {
		httpProxyServer.listener.Close()
	}

	fmt.Println("🌐 HTTP代理服务器已停止")
	return nil
}

// IsRunning 检查HTTP代理服务器是否正在运行
func (s *HTTPProxyServer) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetStats 获取HTTP代理服务器统计信息
func (s *HTTPProxyServer) GetStats() HTTPProxyStats {
	s.stats.mutex.RLock()
	defer s.stats.mutex.RUnlock()
	return s.stats
}

// acceptConnections 接受客户端连接
func (s *HTTPProxyServer) acceptConnections() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("HTTP代理服务器panic: %v", r)
		}
	}()

	for s.IsRunning() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.IsRunning() {
				log.Printf("HTTP代理接受连接错误: %v", err)
			}
			continue
		}

		// 更新统计信息
		s.stats.mutex.Lock()
		s.stats.TotalConnections++
		s.stats.ActiveConnections++
		s.stats.mutex.Unlock()

		// 并发处理连接
		go s.handleConnection(conn)
	}
}

// checkAuthentication 检查HTTP代理认证
func (s *HTTPProxyServer) checkAuthentication(request *http.Request) bool {
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

	username := parts[0]
	password := parts[1]

	// 验证用户名和密码
	return username == s.username && password == s.password
}

// handleConnection 处理单个客户端连接
func (s *HTTPProxyServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		// 更新活跃连接数
		s.stats.mutex.Lock()
		s.stats.ActiveConnections--
		s.stats.mutex.Unlock()

		if r := recover(); r != nil {
			log.Printf("HTTP代理连接处理panic: %v", r)
		}
	}()

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// 读取HTTP请求
	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("HTTP代理读取请求失败: %v", err)
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
				"Proxy-Authenticate: Basic realm=\"HTTP Proxy\"\r\n" +
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
		s.handleHTTPRequest(conn, request)
	}

	s.stats.mutex.Lock()
	s.stats.SuccessfulConnections++
	s.stats.mutex.Unlock()
}

// handleHTTPSConnect 处理HTTPS CONNECT请求
func (s *HTTPProxyServer) handleHTTPSConnect(clientConn net.Conn, request *http.Request) {
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
		log.Printf("HTTP代理连接目标服务器失败: %v", err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		s.stats.mutex.Lock()
		s.stats.FailedConnections++
		s.stats.mutex.Unlock()
		return
	}
	defer targetConn.Close()

	// 发送200 Connection Established响应
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		log.Printf("HTTP代理发送CONNECT响应失败: %v", err)
		return
	}

	// 开始数据转发
	s.relay(clientConn, targetConn)
}

// handleHTTPRequest 处理普通HTTP请求
func (s *HTTPProxyServer) handleHTTPRequest(clientConn net.Conn, request *http.Request) {
	// 更新HTTP请求统计
	s.stats.mutex.Lock()
	s.stats.HTTPRequests++
	s.stats.mutex.Unlock()

	// 解析目标URL
	targetURL := request.URL
	if targetURL.Scheme == "" {
		targetURL.Scheme = "http"
	}
	if targetURL.Host == "" {
		targetURL.Host = request.Host
	}

	// 创建新的HTTP请求
	proxyReq, err := http.NewRequest(request.Method, targetURL.String(), request.Body)
	if err != nil {
		log.Printf("HTTP代理创建请求失败: %v", err)
		clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	// 复制请求头，但移除代理相关的头
	for name, values := range request.Header {
		if !isHopByHopHeader(name) {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
	}

	// 发送请求到目标服务器
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("HTTP代理请求目标服务器失败: %v", err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer resp.Body.Close()

	// 发送响应到客户端
	s.writeResponse(clientConn, resp)
}

// writeResponse 将HTTP响应写入客户端连接
func (s *HTTPProxyServer) writeResponse(conn net.Conn, resp *http.Response) {
	// 写入状态行
	statusLine := fmt.Sprintf("HTTP/%d.%d %d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.StatusCode, resp.Status)
	conn.Write([]byte(statusLine))

	// 写入响应头
	for name, values := range resp.Header {
		if !isHopByHopHeader(name) {
			for _, value := range values {
				headerLine := fmt.Sprintf("%s: %s\r\n", name, value)
				conn.Write([]byte(headerLine))
			}
		}
	}

	// 写入空行分隔符
	conn.Write([]byte("\r\n"))

	// 写入响应体
	bytes, err := io.Copy(conn, resp.Body)
	if err != nil {
		log.Printf("HTTP代理写入响应体失败: %v", err)
	}

	// 更新传输字节统计
	s.stats.mutex.Lock()
	s.stats.BytesTransferred += bytes
	s.stats.mutex.Unlock()
}

// isHopByHopHeader 检查是否为逐跳头部
func isHopByHopHeader(name string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	name = strings.ToLower(name)
	for _, header := range hopByHopHeaders {
		if strings.ToLower(header) == name {
			return true
		}
	}
	return false
}

// relay 在客户端和目标服务器之间转发数据
func (s *HTTPProxyServer) relay(client, target net.Conn) {
	// 移除连接超时，开始长连接数据转发
	client.SetDeadline(time.Time{})
	target.SetDeadline(time.Time{})

	// 使用两个goroutine进行双向数据转发
	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端到目标服务器
	go func() {
		defer wg.Done()
		bytes, err := s.copyData(target, client)
		if err != nil && err != io.EOF {
			log.Printf("HTTP代理数据转发错误(客户端->目标): %v", err)
		}
		s.stats.mutex.Lock()
		s.stats.BytesTransferred += bytes
		s.stats.mutex.Unlock()
	}()

	// 目标服务器到客户端
	go func() {
		defer wg.Done()
		bytes, err := s.copyData(client, target)
		if err != nil && err != io.EOF {
			log.Printf("HTTP代理数据转发错误(目标->客户端): %v", err)
		}
		s.stats.mutex.Lock()
		s.stats.BytesTransferred += bytes
		s.stats.mutex.Unlock()
	}()

	wg.Wait()
}

// copyData 复制数据并返回传输的字节数
func (s *HTTPProxyServer) copyData(dst, src net.Conn) (int64, error) {
	buf := make([]byte, 32*1024) // 32KB缓冲区
	var total int64

	for {
		n, err := src.Read(buf)
		if err != nil {
			return total, err
		}

		if n > 0 {
			written, err := dst.Write(buf[:n])
			if err != nil {
				return total, err
			}
			total += int64(written)
		}
	}
}

// getHTTPProxyStatus 获取HTTP代理服务器状态
func getHTTPProxyStatus() map[string]interface{} {
	if httpProxyServer == nil {
		return map[string]interface{}{
			"enabled": false,
			"status":  "disabled",
			"port":    0,
		}
	}

	stats := httpProxyServer.GetStats()
	return map[string]interface{}{
		"enabled": httpProxyServer.IsRunning(),
		"status": func() string {
			if httpProxyServer.IsRunning() {
				return "running"
			}
			return "stopped"
		}(),
		"port":                  httpProxyServer.port,
		"totalConnections":      stats.TotalConnections,
		"activeConnections":     stats.ActiveConnections,
		"successfulConnections": stats.SuccessfulConnections,
		"failedConnections":     stats.FailedConnections,
		"httpRequests":          stats.HTTPRequests,
		"httpsRequests":         stats.HTTPSRequests,
		"bytesTransferred":      stats.BytesTransferred,
	}
}

// httpProxyDisabledHandler 处理HTTP代理服务禁用时的请求
func httpProxyDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>HTTP代理服务已禁用</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 50px; text-align: center; }
        .container { max-width: 600px; margin: 0 auto; }
        .error { color: #d32f2f; }
        .info { color: #1976d2; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="error">🔒 HTTP代理服务已禁用</h1>
        <p>HTTP/HTTPS代理功能当前未启用。</p>
        <div class="info">
            <p>要启用HTTP代理服务，请使用以下参数重新启动服务器：</p>
            <code>./sweb.exe -proxy</code>
            <p>或指定自定义端口：</p>
            <code>./sweb.exe -proxy -proxy-port 10808</code>
        </div>
    </div>
</body>
</html>`

	w.Write([]byte(html))
}
