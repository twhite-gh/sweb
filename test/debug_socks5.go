package main

import (
	"fmt"
	"net"
	"time"
)

// debugSOCKS5Connection 调试SOCKS5连接
func debugSOCKS5Connection() {
	fmt.Println("🔍 SOCKS5连接调试工具")
	fmt.Println("====================")
	
	// 连接到SOCKS5代理
	proxyAddr := "127.0.0.1:1080"
	fmt.Printf("连接到SOCKS5代理: %s\n", proxyAddr)
	
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		fmt.Printf("❌ 连接失败: %v\n", err)
		fmt.Println("请确保服务器已启动: ./sweb.exe -socks5")
		return
	}
	defer conn.Close()
	
	fmt.Println("✅ 连接成功")
	
	// 发送SOCKS5认证请求
	fmt.Println("\n1. 发送SOCKS5认证请求...")
	authReq := []byte{0x05, 0x01, 0x00} // VER=5, NMETHODS=1, METHOD=0(无认证)
	fmt.Printf("发送数据: %v\n", authReq)
	
	_, err = conn.Write(authReq)
	if err != nil {
		fmt.Printf("❌ 发送认证请求失败: %v\n", err)
		return
	}
	
	// 读取认证响应
	fmt.Println("等待认证响应...")
	authResp := make([]byte, 2)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(authResp)
	if err != nil {
		fmt.Printf("❌ 读取认证响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("接收到 %d 字节: %v\n", n, authResp[:n])
	
	if n >= 2 && authResp[0] == 0x05 && authResp[1] == 0x00 {
		fmt.Println("✅ 认证成功")
	} else {
		fmt.Printf("❌ 认证失败: VER=%d, METHOD=%d\n", authResp[0], authResp[1])
		return
	}
	
	// 发送连接请求
	fmt.Println("\n2. 发送连接请求...")
	// CONNECT to www.baidu.com:80
	connectReq := []byte{
		0x05,       // VER
		0x01,       // CMD (CONNECT)
		0x00,       // RSV
		0x03,       // ATYP (DOMAINNAME)
		0x0d,       // Domain length (13)
	}
	connectReq = append(connectReq, []byte("www.baidu.com")...) // Domain
	connectReq = append(connectReq, 0x00, 0x50)                 // Port 80
	
	fmt.Printf("发送连接请求: %v\n", connectReq)
	
	_, err = conn.Write(connectReq)
	if err != nil {
		fmt.Printf("❌ 发送连接请求失败: %v\n", err)
		return
	}
	
	// 读取连接响应
	fmt.Println("等待连接响应...")
	connectResp := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err = conn.Read(connectResp)
	if err != nil {
		fmt.Printf("❌ 读取连接响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("接收到 %d 字节: %v\n", n, connectResp[:n])
	
	if n >= 4 && connectResp[0] == 0x05 && connectResp[1] == 0x00 {
		fmt.Println("✅ 连接建立成功")
		
		// 发送简单的HTTP请求
		fmt.Println("\n3. 发送HTTP请求...")
		httpReq := "GET / HTTP/1.1\r\nHost: www.baidu.com\r\nConnection: close\r\n\r\n"
		_, err = conn.Write([]byte(httpReq))
		if err != nil {
			fmt.Printf("❌ 发送HTTP请求失败: %v\n", err)
			return
		}
		
		// 读取HTTP响应
		fmt.Println("等待HTTP响应...")
		httpResp := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err = conn.Read(httpResp)
		if err != nil {
			fmt.Printf("❌ 读取HTTP响应失败: %v\n", err)
			return
		}
		
		fmt.Printf("HTTP响应 (%d 字节):\n%s\n", n, string(httpResp[:n]))
		fmt.Println("✅ SOCKS5代理工作正常！")
		
	} else {
		fmt.Printf("❌ 连接失败: VER=%d, REP=%d\n", connectResp[0], connectResp[1])
	}
}

func main() {
	debugSOCKS5Connection()
}
