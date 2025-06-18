package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

// testSOCKS5Proxy 测试SOCKS5代理功能
func testSOCKS5Proxy() {
	fmt.Println("🔌 开始测试SOCKS5代理功能...")

	// 测试代理服务器地址
	proxyAddr := "127.0.0.1:1080"

	// 1. 测试代理服务器连接
	fmt.Printf("1. 测试代理服务器连接 (%s)...\n", proxyAddr)
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		fmt.Printf("❌ 无法连接到SOCKS5代理服务器: %v\n", err)
		fmt.Println("请确保使用 './sweb.exe -socks5' 启动服务器")
		return
	}
	conn.Close()
	fmt.Println("✅ 代理服务器连接成功")

	// 2. 测试SOCKS5协议握手
	fmt.Println("2. 测试SOCKS5协议握手...")
	if err := testSOCKS5Handshake(proxyAddr); err != nil {
		fmt.Printf("❌ SOCKS5握手失败: %v\n", err)
		return
	}
	fmt.Println("✅ SOCKS5握手成功")

	// 3. 测试HTTP请求通过SOCKS5代理
	fmt.Println("3. 测试HTTP请求通过SOCKS5代理...")
	if err := testHTTPThroughSOCKS5(proxyAddr); err != nil {
		fmt.Printf("❌ HTTP请求失败: %v\n", err)
		return
	}
	fmt.Println("✅ HTTP请求成功")

	// 4. 测试并发连接
	fmt.Println("4. 测试并发连接...")
	if err := testSOCKS5Concurrent(proxyAddr); err != nil {
		fmt.Printf("❌ 并发测试失败: %v\n", err)
		return
	}
	fmt.Println("✅ 并发测试成功")

	fmt.Println("🎉 SOCKS5代理功能测试全部通过！")
}

// testSOCKS5Handshake 测试SOCKS5握手过程
func testSOCKS5Handshake(proxyAddr string) error {
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 发送认证方法协商
	authReq := []byte{0x05, 0x01, 0x00} // VER=5, NMETHODS=1, METHOD=0(无认证)
	_, err = conn.Write(authReq)
	if err != nil {
		return fmt.Errorf("发送认证请求失败: %v", err)
	}

	// 读取认证响应
	authResp := make([]byte, 2)
	_, err = conn.Read(authResp)
	if err != nil {
		return fmt.Errorf("读取认证响应失败: %v", err)
	}

	if authResp[0] != 0x05 || authResp[1] != 0x00 {
		return fmt.Errorf("认证响应无效: %v", authResp)
	}

	return nil
}

// testHTTPThroughSOCKS5 测试通过SOCKS5代理发送HTTP请求
func testHTTPThroughSOCKS5(proxyAddr string) error {
	// 创建SOCKS5代理拨号器
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return fmt.Errorf("创建SOCKS5拨号器失败: %v", err)
	}

	// 创建HTTP客户端使用SOCKS5代理
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}

	// 发送HTTP请求
	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		// 如果httpbin.org不可访问，尝试其他测试站点
		resp, err = client.Get("http://www.baidu.com")
		if err != nil {
			return fmt.Errorf("HTTP请求失败: %v", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP响应状态码错误: %d", resp.StatusCode)
	}

	// 读取响应内容（限制大小）
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	fmt.Printf("   响应内容预览: %s\n", string(body[:min(n, 100)]))

	return nil
}

// testSOCKS5Concurrent 测试SOCKS5代理的并发处理能力
func testSOCKS5Concurrent(proxyAddr string) error {
	const numConnections = 5
	errChan := make(chan error, numConnections)

	for i := 0; i < numConnections; i++ {
		go func(id int) {
			// 创建SOCKS5代理拨号器
			dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
				errChan <- fmt.Errorf("连接%d: 创建拨号器失败: %v", id, err)
				return
			}

			// 通过代理连接到测试服务器
			conn, err := dialer.Dial("tcp", "www.baidu.com:80")
			if err != nil {
				errChan <- fmt.Errorf("连接%d: 代理连接失败: %v", id, err)
				return
			}
			defer conn.Close()

			// 发送简单的HTTP请求
			request := "GET / HTTP/1.1\r\nHost: www.baidu.com\r\nConnection: close\r\n\r\n"
			_, err = conn.Write([]byte(request))
			if err != nil {
				errChan <- fmt.Errorf("连接%d: 发送请求失败: %v", id, err)
				return
			}

			// 读取响应
			response := make([]byte, 1024)
			_, err = conn.Read(response)
			if err != nil && err != io.EOF {
				errChan <- fmt.Errorf("连接%d: 读取响应失败: %v", id, err)
				return
			}

			errChan <- nil // 成功
		}(i)
	}

	// 等待所有连接完成
	successCount := 0
	for i := 0; i < numConnections; i++ {
		err := <-errChan
		if err != nil {
			fmt.Printf("   并发连接错误: %v\n", err)
		} else {
			successCount++
		}
	}

	fmt.Printf("   并发连接成功: %d/%d\n", successCount, numConnections)

	if successCount < numConnections/2 {
		return fmt.Errorf("并发连接成功率过低: %d/%d", successCount, numConnections)
	}

	return nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("🧪 SOCKS5代理测试程序")
	fmt.Println("===================")
	testSOCKS5Proxy()
}
