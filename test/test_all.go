package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ServerStatus 服务器状态结构
type ServerStatus struct {
	Upload struct {
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	} `json:"upload"`
	Files struct {
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	} `json:"files"`
	WebDAV struct {
		Enabled   bool   `json:"enabled"`
		Readonly  bool   `json:"readonly"`
		Directory string `json:"directory"`
		Status    string `json:"status"`
	} `json:"webdav"`
	HTTPS struct {
		Enabled    bool   `json:"enabled"`
		HTTPPort   int    `json:"httpPort"`
		HTTPSPort  int    `json:"httpsPort"`
		CertDir    string `json:"certDir"`
		CertStatus string `json:"certStatus"`
	} `json:"https"`
	SOCKS5 struct {
		Enabled               bool   `json:"enabled"`
		Status                string `json:"status"`
		Port                  int    `json:"port"`
		TotalConnections      int64  `json:"totalConnections"`
		ActiveConnections     int64  `json:"activeConnections"`
		SuccessfulConnections int64  `json:"successfulConnections"`
		FailedConnections     int64  `json:"failedConnections"`
		BytesTransferred      int64  `json:"bytesTransferred"`
	} `json:"socks5"`
	Proxy struct {
		Enabled               bool   `json:"enabled"`
		Status                string `json:"status"`
		Port                  int    `json:"port"`
		TotalConnections      int64  `json:"totalConnections"`
		ActiveConnections     int64  `json:"activeConnections"`
		SuccessfulConnections int64  `json:"successfulConnections"`
		FailedConnections     int64  `json:"failedConnections"`
		HTTPRequests          int64  `json:"httpRequests"`
		HTTPSRequests         int64  `json:"httpsRequests"`
		BytesTransferred      int64  `json:"bytesTransferred"`
	} `json:"proxy"`
}

// testAllFeatures 测试所有功能
func testAllFeatures() {
	fmt.Println("🧪 SWeb 综合功能测试")
	fmt.Println("====================")

	// 1. 检查服务器状态
	fmt.Println("1. 检查服务器状态...")
	status, err := getServerStatus()
	if err != nil {
		fmt.Printf("❌ 无法获取服务器状态: %v\n", err)
		fmt.Println("请确保服务器正在运行: './sweb.exe -upload -files -webdav -https -socks5 -proxy'")
		return
	}

	printServerStatus(status)

	// 2. 测试基础HTTP服务
	fmt.Println("\n2. 测试基础HTTP服务...")
	if err := testBasicHTTPService(); err != nil {
		fmt.Printf("❌ 基础HTTP服务测试失败: %v\n", err)
	} else {
		fmt.Println("✅ 基础HTTP服务测试成功")
	}

	// 3. 测试文件上传功能
	if status.Upload.Enabled {
		fmt.Println("\n3. 测试文件上传功能...")
		if err := testFileUpload(); err != nil {
			fmt.Printf("❌ 文件上传测试失败: %v\n", err)
		} else {
			fmt.Println("✅ 文件上传测试成功")
		}
	} else {
		fmt.Println("\n3. ⏭️  文件上传功能未启用，跳过测试")
	}

	// 4. 测试文件浏览功能
	if status.Files.Enabled {
		fmt.Println("\n4. 测试文件浏览功能...")
		if err := testFileBrowsing(); err != nil {
			fmt.Printf("❌ 文件浏览测试失败: %v\n", err)
		} else {
			fmt.Println("✅ 文件浏览测试成功")
		}
	} else {
		fmt.Println("\n4. ⏭️  文件浏览功能未启用，跳过测试")
	}

	// 5. 测试WebDAV功能
	if status.WebDAV.Enabled {
		fmt.Println("\n5. 测试WebDAV功能...")
		if err := testWebDAVService(); err != nil {
			fmt.Printf("❌ WebDAV测试失败: %v\n", err)
		} else {
			fmt.Println("✅ WebDAV测试成功")
		}
	} else {
		fmt.Println("\n5. ⏭️  WebDAV功能未启用，跳过测试")
	}

	// 6. 测试SOCKS5代理
	if status.SOCKS5.Enabled {
		fmt.Println("\n6. 测试SOCKS5代理功能...")
		if err := testSOCKS5ProxyBasic(); err != nil {
			fmt.Printf("❌ SOCKS5代理基础测试失败: %v\n", err)
		} else {
			fmt.Println("✅ SOCKS5代理基础测试成功")
			fmt.Println("   运行详细测试: go run test_socks5.go")
		}
	} else {
		fmt.Println("\n6. ⏭️  SOCKS5代理功能未启用，跳过测试")
	}

	// 7. 测试HTTP代理
	if status.Proxy.Enabled {
		fmt.Println("\n7. 测试HTTP代理功能...")
		if err := testHTTPProxyBasic(); err != nil {
			fmt.Printf("❌ HTTP代理基础测试失败: %v\n", err)
		} else {
			fmt.Println("✅ HTTP代理基础测试成功")
			fmt.Println("   运行详细测试: go run test_proxy.go")
		}
	} else {
		fmt.Println("\n7. ⏭️  HTTP代理功能未启用，跳过测试")
	}

	fmt.Println("\n🎉 综合功能测试完成！")
}

// getServerStatus 获取服务器状态
func getServerStatus() (*ServerStatus, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8080/api/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("状态码错误: %d", resp.StatusCode)
	}

	var status ServerStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// printServerStatus 打印服务器状态
func printServerStatus(status *ServerStatus) {
	fmt.Println("   服务器功能状态:")
	fmt.Printf("   📤 文件上传: %s\n", getStatusText(status.Upload.Enabled))
	fmt.Printf("   📂 文件浏览: %s\n", getStatusText(status.Files.Enabled))
	fmt.Printf("   📁 WebDAV: %s", getStatusText(status.WebDAV.Enabled))
	if status.WebDAV.Enabled {
		if status.WebDAV.Readonly {
			fmt.Printf(" (只读)")
		} else {
			fmt.Printf(" (读写)")
		}
	}
	fmt.Println()
	fmt.Printf("   🔒 HTTPS: %s\n", getStatusText(status.HTTPS.Enabled))
	fmt.Printf("   🔌 SOCKS5代理: %s", getStatusText(status.SOCKS5.Enabled))
	if status.SOCKS5.Enabled {
		fmt.Printf(" (端口:%d)", status.SOCKS5.Port)
	}
	fmt.Println()
	fmt.Printf("   🌐 HTTP代理: %s", getStatusText(status.Proxy.Enabled))
	if status.Proxy.Enabled {
		fmt.Printf(" (端口:%d)", status.Proxy.Port)
	}
	fmt.Println()
}

// getStatusText 获取状态文本
func getStatusText(enabled bool) string {
	if enabled {
		return "✅ 已启用"
	}
	return "🔒 已禁用"
}

// testBasicHTTPService 测试基础HTTP服务
func testBasicHTTPService() error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("状态码错误: %d", resp.StatusCode)
	}

	return nil
}

// testFileUpload 测试文件上传功能
func testFileUpload() error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8080/upload")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("上传页面状态码错误: %d", resp.StatusCode)
	}

	return nil
}

// testFileBrowsing 测试文件浏览功能
func testFileBrowsing() error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8080/files")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("文件浏览页面状态码错误: %d", resp.StatusCode)
	}

	return nil
}

// testWebDAVService 测试WebDAV服务
func testWebDAVService() error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("PROPFIND", "http://localhost:8080/webdav/", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 207 && resp.StatusCode != 200 {
		return fmt.Errorf("WebDAV PROPFIND状态码错误: %d", resp.StatusCode)
	}

	return nil
}

// testSOCKS5ProxyBasic 基础SOCKS5代理测试
func testSOCKS5ProxyBasic() error {
	// 简单的连接测试
	conn, err := net.DialTimeout("tcp", "localhost:1080", 5*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接SOCKS5代理: %v", err)
	}
	conn.Close()
	return nil
}

// testHTTPProxyBasic 基础HTTP代理测试
func testHTTPProxyBasic() error {
	// 简单的连接测试
	conn, err := net.DialTimeout("tcp", "localhost:10808", 5*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接HTTP代理: %v", err)
	}
	conn.Close()
	return nil
}

func main() {
	testAllFeatures()
}
