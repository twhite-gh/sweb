package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// testHTTPProxy 测试HTTP代理功能
func testHTTPProxy() {
	fmt.Println("🌐 开始测试HTTP代理功能...")

	// 测试代理服务器地址
	proxyAddr := "127.0.0.1:10808"

	// 1. 测试代理服务器连接
	fmt.Printf("1. 测试代理服务器连接 (%s)...\n", proxyAddr)
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		fmt.Printf("❌ 无法连接到HTTP代理服务器: %v\n", err)
		fmt.Println("请确保使用 './sweb.exe -proxy' 启动服务器")
		return
	}
	conn.Close()
	fmt.Println("✅ 代理服务器连接成功")

	// 2. 测试HTTP GET请求
	fmt.Println("2. 测试HTTP GET请求...")
	if err := testHTTPGetRequest(proxyAddr); err != nil {
		fmt.Printf("❌ HTTP GET请求失败: %v\n", err)
		return
	}
	fmt.Println("✅ HTTP GET请求成功")

	// 3. 测试HTTPS CONNECT请求
	fmt.Println("3. 测试HTTPS CONNECT请求...")
	if err := testHTTPSConnectRequest(proxyAddr); err != nil {
		fmt.Printf("❌ HTTPS CONNECT请求失败: %v\n", err)
		return
	}
	fmt.Println("✅ HTTPS CONNECT请求成功")

	// 4. 测试使用Go HTTP客户端
	fmt.Println("4. 测试使用Go HTTP客户端...")
	if err := testHTTPClientWithProxy(proxyAddr); err != nil {
		fmt.Printf("❌ HTTP客户端测试失败: %v\n", err)
		return
	}
	fmt.Println("✅ HTTP客户端测试成功")

	// 5. 测试HTTPS代理功能
	fmt.Println("5. 测试HTTPS代理功能...")
	if err := testHTTPSProxyWithClient(proxyAddr); err != nil {
		fmt.Printf("❌ HTTPS代理测试失败: %v\n", err)
		return
	}
	fmt.Println("✅ HTTPS代理测试成功")

	// 6. 测试并发请求
	fmt.Println("6. 测试并发请求...")
	if err := testHTTPProxyConcurrent(proxyAddr); err != nil {
		fmt.Printf("❌ 并发测试失败: %v\n", err)
		return
	}
	fmt.Println("✅ 并发测试成功")

	fmt.Println("🎉 HTTP代理功能测试全部通过！")
}

// testHTTPGetRequest 测试HTTP GET请求
func testHTTPGetRequest(proxyAddr string) error {
	// 直接连接到代理服务器
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("连接代理服务器失败: %v", err)
	}
	defer conn.Close()

	// 发送HTTP GET请求
	request := "GET http://httpbin.org/ip HTTP/1.1\r\n" +
		"Host: httpbin.org\r\n" +
		"User-Agent: sweb-test/1.0\r\n" +
		"Connection: close\r\n\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}

	// 读取响应
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, nil)
	if err != nil {
		// 如果httpbin.org不可访问，尝试百度
		conn.Close()
		conn, err = net.DialTimeout("tcp", proxyAddr, 5*time.Second)
		if err != nil {
			return fmt.Errorf("重新连接失败: %v", err)
		}
		defer conn.Close()

		request = "GET http://www.baidu.com HTTP/1.1\r\n" +
			"Host: www.baidu.com\r\n" +
			"User-Agent: sweb-test/1.0\r\n" +
			"Connection: close\r\n\r\n"

		_, err = conn.Write([]byte(request))
		if err != nil {
			return fmt.Errorf("发送备用请求失败: %v", err)
		}

		reader = bufio.NewReader(conn)
		response, err = http.ReadResponse(reader, nil)
		if err != nil {
			return fmt.Errorf("读取响应失败: %v", err)
		}
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("HTTP响应状态码错误: %d", response.StatusCode)
	}

	// 读取响应内容（限制大小）
	body := make([]byte, 1024)
	n, _ := response.Body.Read(body)
	fmt.Printf("   响应状态: %s\n", response.Status)
	fmt.Printf("   响应内容预览: %s\n", string(body[:min(n, 100)]))

	return nil
}

// testHTTPSConnectRequest 测试HTTPS CONNECT请求
func testHTTPSConnectRequest(proxyAddr string) error {
	// 连接到代理服务器
	conn, err := net.DialTimeout("tcp", proxyAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("连接代理服务器失败: %v", err)
	}
	defer conn.Close()

	// 发送CONNECT请求
	request := "CONNECT www.baidu.com:443 HTTP/1.1\r\n" +
		"Host: www.baidu.com:443\r\n" +
		"User-Agent: sweb-test/1.0\r\n\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("发送CONNECT请求失败: %v", err)
	}

	// 读取CONNECT响应
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, nil)
	if err != nil {
		return fmt.Errorf("读取CONNECT响应失败: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("CONNECT响应状态码错误: %d %s", response.StatusCode, response.Status)
	}

	fmt.Printf("   CONNECT响应: %s\n", response.Status)

	// 现在可以通过隧道发送数据（这里只是测试连接建立）
	return nil
}

// testHTTPClientWithProxy 测试使用Go HTTP客户端通过代理
func testHTTPClientWithProxy(proxyAddr string) error {
	// 创建代理URL
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return fmt.Errorf("解析代理URL失败: %v", err)
	}

	// 创建HTTP客户端使用代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	// 发送HTTP请求
	resp, err := client.Get("http://httpbin.org/user-agent")
	if err != nil {
		// 如果httpbin.org不可访问，尝试其他站点
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
	fmt.Printf("   客户端响应状态: %s\n", resp.Status)
	fmt.Printf("   客户端响应内容预览: %s\n", string(body[:min(n, 100)]))

	return nil
}

// testHTTPSProxyWithClient 测试HTTPS代理功能
func testHTTPSProxyWithClient(proxyAddr string) error {
	// 创建代理URL
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return fmt.Errorf("解析代理URL失败: %v", err)
	}

	// 创建HTTP客户端使用代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			// 跳过SSL验证以便测试
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 15 * time.Second,
	}

	// 测试HTTPS网站列表
	httpsUrls := []string{
		"https://www.baidu.com",
		"https://httpbin.org/ip",
		"https://www.google.com",
	}

	successCount := 0
	for i, testUrl := range httpsUrls {
		fmt.Printf("   测试HTTPS站点 %d: %s\n", i+1, testUrl)

		resp, err := client.Get(testUrl)
		if err != nil {
			fmt.Printf("   ⚠️  HTTPS请求失败: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			// 读取响应内容（限制大小）
			body := make([]byte, 512)
			n, _ := resp.Body.Read(body)
			fmt.Printf("   ✅ HTTPS响应成功 (状态: %s, 内容: %d 字节)\n", resp.Status, n)
			successCount++
		} else {
			fmt.Printf("   ⚠️  HTTPS响应状态码: %d\n", resp.StatusCode)
		}
	}

	if successCount == 0 {
		return fmt.Errorf("所有HTTPS测试都失败了")
	}

	fmt.Printf("   HTTPS代理测试成功: %d/%d 个站点可访问\n", successCount, len(httpsUrls))
	return nil
}

// testHTTPProxyConcurrent 测试HTTP代理的并发处理能力
func testHTTPProxyConcurrent(proxyAddr string) error {
	const numRequests = 5
	errChan := make(chan error, numRequests)

	// 创建代理URL
	proxyURL, err := url.Parse("http://" + proxyAddr)
	if err != nil {
		return fmt.Errorf("解析代理URL失败: %v", err)
	}

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			// 创建HTTP客户端
			client := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
				Timeout: 10 * time.Second,
			}

			// 发送HTTP请求
			resp, err := client.Get("http://www.baidu.com")
			if err != nil {
				errChan <- fmt.Errorf("请求%d: HTTP请求失败: %v", id, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				errChan <- fmt.Errorf("请求%d: HTTP响应状态码错误: %d", id, resp.StatusCode)
				return
			}

			// 读取部分响应内容
			body := make([]byte, 512)
			_, err = resp.Body.Read(body)
			if err != nil && err != io.EOF {
				errChan <- fmt.Errorf("请求%d: 读取响应失败: %v", id, err)
				return
			}

			errChan <- nil // 成功
		}(i)
	}

	// 等待所有请求完成
	successCount := 0
	for i := 0; i < numRequests; i++ {
		err := <-errChan
		if err != nil {
			fmt.Printf("   并发请求错误: %v\n", err)
		} else {
			successCount++
		}
	}

	fmt.Printf("   并发请求成功: %d/%d\n", successCount, numRequests)

	if successCount < numRequests/2 {
		return fmt.Errorf("并发请求成功率过低: %d/%d", successCount, numRequests)
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
	fmt.Println("🧪 HTTP代理测试程序")
	fmt.Println("==================")
	testHTTPProxy()
}
