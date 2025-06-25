package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// 全局变量存储功能状态
var (
	uploadEnabled  bool
	webdavEnabled  bool
	webdavDir      string
	webdavReadonly bool
	httpsEnabled   bool
	httpPort       int
	httpsPort      int
	certDir        string
	filesEnabled   bool
	socks5Enabled  bool
	socks5Port     int
	proxyEnabled   bool
	proxyPort      int
)

func main() {
	// 解析命令行参数
	var showHelp bool

	flag.BoolVar(&uploadEnabled, "upload", false, "启用文件上传功能")
	flag.BoolVar(&uploadEnabled, "enable-upload", false, "启用文件上传功能")
	flag.BoolVar(&filesEnabled, "files", false, "启用文件浏览功能")
	flag.BoolVar(&filesEnabled, "enable-files", false, "启用文件浏览功能")
	flag.BoolVar(&webdavEnabled, "webdav", false, "启用WebDAV服务")
	flag.BoolVar(&webdavEnabled, "enable-webdav", false, "启用WebDAV服务")
	flag.StringVar(&webdavDir, "webdav-dir", ".", "WebDAV服务的根目录")
	flag.BoolVar(&webdavReadonly, "webdav-readonly", false, "WebDAV服务只读模式")
	flag.BoolVar(&httpsEnabled, "https", false, "启用HTTPS服务")
	flag.BoolVar(&httpsEnabled, "enable-https", false, "启用HTTPS服务")
	flag.IntVar(&httpPort, "port", 8080, "指定HTTP服务器端口")
	flag.IntVar(&httpPort, "p", 8080, "指定HTTP服务器端口")
	flag.IntVar(&httpsPort, "https-port", 8443, "指定HTTPS服务器端口")
	flag.StringVar(&certDir, "cert-dir", "./cert", "SSL证书目录")
	flag.BoolVar(&socks5Enabled, "socks5", false, "启用SOCKS5代理服务")
	flag.BoolVar(&socks5Enabled, "enable-socks5", false, "启用SOCKS5代理服务")
	flag.IntVar(&socks5Port, "socks5-port", 1080, "指定SOCKS5代理服务器端口")
	flag.BoolVar(&proxyEnabled, "proxy", false, "启用HTTP代理服务")
	flag.BoolVar(&proxyEnabled, "enable-proxy", false, "启用HTTP代理服务")
	flag.IntVar(&proxyPort, "proxy-port", 10808, "指定HTTP代理服务器端口")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助信息
	if showHelp {
		showHelpInfo()
		return
	}

	// 创建web目录（如果不存在）
	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		err := os.Mkdir(webDir, 0755)
		if err != nil {
			log.Fatalf("无法创建web目录: %v", err)
		}
	}

	// 检查并创建默认页面
	createDefaultPageIfNeeded(webDir, uploadEnabled)

	// 处理静态文件（HTML, JS等）
	//fileServer := http.FileServer(http.Dir(webDir))
	//http.Handle("/", fileServer)
	// 创建一个文件服务器，用于处理静态文件
	// 注意：http.Dir 会自动处理路径安全，阻止目录遍历
	fileServer := http.FileServer(http.Dir(webDir))

	// 使用 http.HandleFunc 注册一个通用的处理函数
	// 这个处理函数会优先处理 "/" 路径，然后将其他路径转发给文件服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果请求的是根路径 "/"
		if r.URL.Path == "/" {
			serveIndexOrAbout(w, r, webDir)
			return // 处理完毕，不再继续
		}

		// 如果不是根路径，则使用文件服务器处理请求
		// 这样做是为了让像 /style.css, /images/logo.png 这样的请求能够被正常服务
		fileServer.ServeHTTP(w, r)
	})

	// 添加状态API端点
	http.HandleFunc("/api/status", statusHandler)

	// 添加文件浏览API端点
	http.HandleFunc("/api/files", filesListHandler)

	// 根据参数决定是否启用文件浏览
	if filesEnabled {
		http.HandleFunc("/files", filesBrowseHandler)
		fmt.Println("✅ 文件浏览功能已启用")
	} else {
		http.HandleFunc("/files", filesDisabledHandler)
		fmt.Println("🔒 文件浏览功能已禁用 (使用 -files 参数启用)")
	}

	// 根据参数决定是否启用文件上传
	if uploadEnabled {
		http.HandleFunc("/upload", uploadHandler)
		fmt.Println("✅ 文件上传功能已启用")
	} else {
		http.HandleFunc("/upload", uploadDisabledHandler)
		fmt.Println("🔒 文件上传功能已禁用 (使用 -upload 参数启用)")
	}

	// 根据参数决定是否启用WebDAV服务
	if webdavEnabled {
		setupWebDAVHandler()
		if webdavReadonly {
			fmt.Printf("✅ WebDAV服务已启用 (只读模式) - 目录: %s\n", webdavDir)
		} else {
			fmt.Printf("✅ WebDAV服务已启用 (读写模式) - 目录: %s\n", webdavDir)
		}
	} else {
		http.HandleFunc("/webdav", webdavDisabledHandler)
		fmt.Println("🔒 WebDAV服务已禁用 (使用 -webdav 参数启用)")
	}

	// 根据参数决定是否启用SOCKS5代理服务
	if socks5Enabled {
		if err := startSOCKS5Server(socks5Port); err != nil {
			log.Fatalf("启动SOCKS5代理服务失败: %v", err)
		}
		fmt.Printf("✅ SOCKS5代理服务已启用 - 端口: %d\n", socks5Port)
	} else {
		http.HandleFunc("/socks5", socks5DisabledHandler)
		fmt.Println("🔒 SOCKS5代理服务已禁用 (使用 -socks5 参数启用)")
	}

	// 根据参数决定是否启用HTTP代理服务
	if proxyEnabled {
		if err := startHTTPProxyServer(proxyPort); err != nil {
			log.Fatalf("启动HTTP代理服务失败: %v", err)
		}
		fmt.Printf("✅ HTTP代理服务已启用 - 端口: %d\n", proxyPort)
	} else {
		http.HandleFunc("/proxy", httpProxyDisabledHandler)
		fmt.Println("🔒 HTTP代理服务已禁用 (使用 -proxy 参数启用)")
	}

	// 启动服务器
	startServers()
}
