package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// SOCKS5 协议常量
const (
	SOCKS5_VERSION = 0x05

	// 认证方法
	AUTH_NO_AUTH       = 0x00
	AUTH_GSSAPI        = 0x01
	AUTH_PASSWORD      = 0x02
	AUTH_NO_ACCEPTABLE = 0xFF

	// 命令类型
	CMD_CONNECT = 0x01
	CMD_BIND    = 0x02
	CMD_UDP     = 0x03

	// 地址类型
	ADDR_IPV4   = 0x01
	ADDR_DOMAIN = 0x03
	ADDR_IPV6   = 0x04

	// 响应状态
	REP_SUCCESS                    = 0x00
	REP_GENERAL_FAILURE            = 0x01
	REP_CONNECTION_NOT_ALLOWED     = 0x02
	REP_NETWORK_UNREACHABLE        = 0x03
	REP_HOST_UNREACHABLE           = 0x04
	REP_CONNECTION_REFUSED         = 0x05
	REP_TTL_EXPIRED                = 0x06
	REP_COMMAND_NOT_SUPPORTED      = 0x07
	REP_ADDRESS_TYPE_NOT_SUPPORTED = 0x08
)

// SOCKS5服务器结构
type SOCKS5Server struct {
	port     int
	listener net.Listener
	running  bool
	mutex    sync.RWMutex
	stats    SOCKS5Stats
	auth     bool
	username string
	password string
}

// SOCKS5统计信息
type SOCKS5Stats struct {
	TotalConnections      int64
	ActiveConnections     int64
	SuccessfulConnections int64
	FailedConnections     int64
	BytesTransferred      int64
	mutex                 sync.RWMutex
}

// 全局SOCKS5服务器实例
var socks5Server *SOCKS5Server

// startSOCKS5Server 启动SOCKS5代理服务器
func startSOCKS5Server(port int, auth bool, username, password string) error {
	if socks5Server != nil && socks5Server.IsRunning() {
		return fmt.Errorf("SOCKS5服务器已在运行")
	}

	server := &SOCKS5Server{
		port:     port,
		auth:     auth,
		username: username,
		password: password,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("无法启动SOCKS5服务器: %v", err)
	}

	server.listener = listener
	server.running = true
	socks5Server = server

	fmt.Printf("🔌 SOCKS5代理服务器启动在端口 %d\n", port)

	go server.acceptConnections()

	return nil
}

// stopSOCKS5Server 停止SOCKS5代理服务器
func stopSOCKS5Server() error {
	if socks5Server == nil {
		return fmt.Errorf("SOCKS5服务器未运行")
	}

	socks5Server.mutex.Lock()
	defer socks5Server.mutex.Unlock()

	if !socks5Server.running {
		return fmt.Errorf("SOCKS5服务器未运行")
	}

	socks5Server.running = false
	if socks5Server.listener != nil {
		socks5Server.listener.Close()
	}

	fmt.Println("🔌 SOCKS5代理服务器已停止")
	return nil
}

// IsRunning 检查SOCKS5服务器是否正在运行
func (s *SOCKS5Server) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.running
}

// GetStats 获取SOCKS5服务器统计信息
func (s *SOCKS5Server) GetStats() SOCKS5Stats {
	s.stats.mutex.RLock()
	defer s.stats.mutex.RUnlock()
	return s.stats
}

// acceptConnections 接受客户端连接
func (s *SOCKS5Server) acceptConnections() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("SOCKS5服务器panic: %v", r)
		}
	}()

	for s.IsRunning() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.IsRunning() {
				log.Printf("SOCKS5接受连接错误: %v", err)
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

// handleConnection 处理单个客户端连接
func (s *SOCKS5Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		// 更新活跃连接数
		s.stats.mutex.Lock()
		s.stats.ActiveConnections--
		s.stats.mutex.Unlock()

		if r := recover(); r != nil {
			log.Printf("SOCKS5连接处理panic: %v", r)
		}
	}()

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// 1. 处理认证协商
	if err := s.handleAuthentication(conn); err != nil {
		log.Printf("SOCKS5认证失败: %v", err)
		s.stats.mutex.Lock()
		s.stats.FailedConnections++
		s.stats.mutex.Unlock()
		return
	}

	// 2. 处理连接请求
	if err := s.handleRequest(conn); err != nil {
		log.Printf("SOCKS5请求处理失败: %v", err)
		s.stats.mutex.Lock()
		s.stats.FailedConnections++
		s.stats.mutex.Unlock()
		return
	}

	s.stats.mutex.Lock()
	s.stats.SuccessfulConnections++
	s.stats.mutex.Unlock()
}

// handleAuthentication 处理SOCKS5认证协商
func (s *SOCKS5Server) handleAuthentication(conn net.Conn) error {
	// 读取客户端认证请求
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("读取认证请求失败: %v", err)
	}

	// 添加详细的调试信息
	log.Printf("SOCKS5认证请求: 接收到 %d 字节数据: %v", n, buf[:n])

	// 检查最小数据长度
	if n < 2 {
		return fmt.Errorf("认证请求数据太短: %d 字节", n)
	}

	// 检查SOCKS版本
	version := buf[0]
	if version != SOCKS5_VERSION {
		return fmt.Errorf("无效的SOCKS5版本: 期望 %d, 收到 %d", SOCKS5_VERSION, version)
	}

	// 检查认证方法数量
	nmethods := int(buf[1])
	if nmethods == 0 {
		return fmt.Errorf("认证方法数量为0")
	}

	if n < 2+nmethods {
		return fmt.Errorf("认证方法数据不完整: 期望 %d 字节, 收到 %d 字节", 2+nmethods, n)
	}

	log.Printf("SOCKS5认证: 版本=%d, 方法数量=%d", version, nmethods)

	// 检查支持的认证方法
	supportedMethod := AUTH_NO_ACCEPTABLE
	for i := 0; i < nmethods; i++ {
		method := buf[2+i]
		log.Printf("SOCKS5认证: 客户端支持方法 %d", method)
		if s.auth && method == AUTH_PASSWORD {
			supportedMethod = AUTH_PASSWORD
			break
		} else if !s.auth && method == AUTH_NO_AUTH {
			supportedMethod = AUTH_NO_AUTH
			break
		}
	}

	log.Printf("SOCKS5认证: 选择的认证方法 %d", supportedMethod)

	// 发送认证方法选择响应
	response := []byte{SOCKS5_VERSION, byte(supportedMethod)}
	_, err = conn.Write(response)
	if err != nil {
		return fmt.Errorf("发送认证响应失败: %v", err)
	}

	if supportedMethod == AUTH_NO_ACCEPTABLE {
		return fmt.Errorf("不支持的认证方法")
	}

	// 如果选择了用户名密码认证，进行认证验证
	if supportedMethod == AUTH_PASSWORD {
		if err := s.handlePasswordAuthentication(conn); err != nil {
			return fmt.Errorf("用户名密码认证失败: %v", err)
		}
	}

	log.Printf("SOCKS5认证: 认证协商成功")
	return nil
}

// handlePasswordAuthentication 处理用户名密码认证
func (s *SOCKS5Server) handlePasswordAuthentication(conn net.Conn) error {
	// 读取用户名密码认证请求
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("读取认证数据失败: %v", err)
	}

	if n < 3 {
		return fmt.Errorf("认证数据太短")
	}

	// 检查认证版本（应该是1）
	if buf[0] != 0x01 {
		return fmt.Errorf("无效的认证版本: %d", buf[0])
	}

	// 解析用户名
	usernameLen := int(buf[1])
	if n < 2+usernameLen+1 {
		return fmt.Errorf("用户名数据不完整")
	}
	username := string(buf[2 : 2+usernameLen])

	// 解析密码
	passwordLen := int(buf[2+usernameLen])
	if n < 2+usernameLen+1+passwordLen {
		return fmt.Errorf("密码数据不完整")
	}
	password := string(buf[2+usernameLen+1 : 2+usernameLen+1+passwordLen])

	log.Printf("SOCKS5认证: 用户名=%s", username)

	// 验证用户名和密码
	success := (username == s.username && password == s.password)

	// 发送认证结果
	var response []byte
	if success {
		response = []byte{0x01, 0x00} // 认证成功
		log.Printf("SOCKS5认证: 用户名密码认证成功")
	} else {
		response = []byte{0x01, 0x01} // 认证失败
		log.Printf("SOCKS5认证: 用户名密码认证失败")
	}

	_, err = conn.Write(response)
	if err != nil {
		return fmt.Errorf("发送认证结果失败: %v", err)
	}

	if !success {
		return fmt.Errorf("用户名或密码错误")
	}

	return nil
}

// handleRequest 处理SOCKS5连接请求
func (s *SOCKS5Server) handleRequest(conn net.Conn) error {
	// 读取连接请求
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("读取连接请求失败: %v", err)
	}

	// 添加详细的调试信息
	debugLen := n
	if debugLen > 20 {
		debugLen = 20
	}
	log.Printf("SOCKS5连接请求: 接收到 %d 字节数据: %v", n, buf[:debugLen])

	if n < 4 {
		s.sendErrorResponse(conn, REP_GENERAL_FAILURE)
		return fmt.Errorf("连接请求数据太短: %d 字节", n)
	}

	if buf[0] != SOCKS5_VERSION {
		s.sendErrorResponse(conn, REP_GENERAL_FAILURE)
		return fmt.Errorf("无效的SOCKS5请求版本: 期望 %d, 收到 %d", SOCKS5_VERSION, buf[0])
	}

	cmd := buf[1]
	if cmd != CMD_CONNECT {
		s.sendErrorResponse(conn, REP_COMMAND_NOT_SUPPORTED)
		return fmt.Errorf("不支持的命令: %d", cmd)
	}

	// 解析目标地址
	targetAddr, err := s.parseAddress(buf[3:n])
	if err != nil {
		s.sendErrorResponse(conn, REP_ADDRESS_TYPE_NOT_SUPPORTED)
		return fmt.Errorf("解析地址失败: %v", err)
	}

	// 连接到目标服务器
	targetConn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		s.sendErrorResponse(conn, REP_HOST_UNREACHABLE)
		return fmt.Errorf("连接目标服务器失败: %v", err)
	}
	defer targetConn.Close()

	// 发送成功响应
	if err := s.sendSuccessResponse(conn, targetConn.LocalAddr()); err != nil {
		return fmt.Errorf("发送成功响应失败: %v", err)
	}

	// 开始数据转发
	s.relay(conn, targetConn)

	return nil
}

// parseAddress 解析SOCKS5地址
func (s *SOCKS5Server) parseAddress(data []byte) (string, error) {
	if len(data) < 1 {
		return "", fmt.Errorf("地址数据为空")
	}

	addrType := data[0]
	switch addrType {
	case ADDR_IPV4:
		if len(data) < 7 {
			return "", fmt.Errorf("IPv4地址数据不完整")
		}
		ip := net.IP(data[1:5])
		port := binary.BigEndian.Uint16(data[5:7])
		return fmt.Sprintf("%s:%d", ip.String(), port), nil

	case ADDR_DOMAIN:
		if len(data) < 2 {
			return "", fmt.Errorf("域名地址数据不完整")
		}
		domainLen := int(data[1])
		if len(data) < 2+domainLen+2 {
			return "", fmt.Errorf("域名地址数据不完整")
		}
		domain := string(data[2 : 2+domainLen])
		port := binary.BigEndian.Uint16(data[2+domainLen : 2+domainLen+2])
		return fmt.Sprintf("%s:%d", domain, port), nil

	case ADDR_IPV6:
		if len(data) < 19 {
			return "", fmt.Errorf("IPv6地址数据不完整")
		}
		ip := net.IP(data[1:17])
		port := binary.BigEndian.Uint16(data[17:19])
		return fmt.Sprintf("[%s]:%d", ip.String(), port), nil

	default:
		return "", fmt.Errorf("不支持的地址类型: %d", addrType)
	}
}

// sendErrorResponse 发送错误响应
func (s *SOCKS5Server) sendErrorResponse(conn net.Conn, rep byte) {
	response := []byte{
		SOCKS5_VERSION, rep, 0x00, // VER, REP, RSV
		ADDR_IPV4, 0, 0, 0, 0, 0, 0, // ATYP, BND.ADDR, BND.PORT
	}
	conn.Write(response)
}

// sendSuccessResponse 发送成功响应
func (s *SOCKS5Server) sendSuccessResponse(conn net.Conn, bindAddr net.Addr) error {
	response := []byte{SOCKS5_VERSION, REP_SUCCESS, 0x00}

	// 解析绑定地址
	tcpAddr, ok := bindAddr.(*net.TCPAddr)
	if !ok {
		// 如果无法解析，使用默认值
		response = append(response, ADDR_IPV4)
		response = append(response, 0, 0, 0, 0) // 0.0.0.0
		response = append(response, 0, 0)       // port 0
	} else {
		if tcpAddr.IP.To4() != nil {
			// IPv4
			response = append(response, ADDR_IPV4)
			response = append(response, tcpAddr.IP.To4()...)
		} else {
			// IPv6
			response = append(response, ADDR_IPV6)
			response = append(response, tcpAddr.IP.To16()...)
		}
		portBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(portBytes, uint16(tcpAddr.Port))
		response = append(response, portBytes...)
	}

	_, err := conn.Write(response)
	return err
}

// relay 在客户端和目标服务器之间转发数据
func (s *SOCKS5Server) relay(client, target net.Conn) {
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
			log.Printf("SOCKS5数据转发错误(客户端->目标): %v", err)
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
			log.Printf("SOCKS5数据转发错误(目标->客户端): %v", err)
		}
		s.stats.mutex.Lock()
		s.stats.BytesTransferred += bytes
		s.stats.mutex.Unlock()
	}()

	wg.Wait()
}

// copyData 复制数据并返回传输的字节数
func (s *SOCKS5Server) copyData(dst, src net.Conn) (int64, error) {
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

// getSOCKS5Status 获取SOCKS5服务器状态
func getSOCKS5Status() map[string]interface{} {
	if socks5Server == nil {
		return map[string]interface{}{
			"enabled": false,
			"status":  "disabled",
			"port":    0,
		}
	}

	stats := socks5Server.GetStats()
	return map[string]interface{}{
		"enabled": socks5Server.IsRunning(),
		"status": func() string {
			if socks5Server.IsRunning() {
				return "running"
			}
			return "stopped"
		}(),
		"port":                  socks5Server.port,
		"totalConnections":      stats.TotalConnections,
		"activeConnections":     stats.ActiveConnections,
		"successfulConnections": stats.SuccessfulConnections,
		"failedConnections":     stats.FailedConnections,
		"bytesTransferred":      stats.BytesTransferred,
	}
}

// socks5DisabledHandler 处理SOCKS5服务禁用时的请求
func socks5DisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>SOCKS5代理服务已禁用</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 50px; text-align: center; }
        .container { max-width: 600px; margin: 0 auto; }
        .error { color: #d32f2f; }
        .info { color: #1976d2; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="error">🔒 SOCKS5代理服务已禁用</h1>
        <p>SOCKS5代理功能当前未启用。</p>
        <div class="info">
            <p>要启用SOCKS5代理服务，请使用以下参数重新启动服务器：</p>
            <code>./sweb.exe -socks5</code>
            <p>或指定自定义端口：</p>
            <code>./sweb.exe -socks5 -socks5-port 1080</code>
        </div>
    </div>
</body>
</html>`

	w.Write([]byte(html))
}
