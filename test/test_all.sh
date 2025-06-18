#!/bin/bash

echo "========================================"
echo "SWeb Complete Test Suite"
echo "SWeb 完整测试套件"
echo "========================================"
echo

echo "Please ensure the server is started with all features enabled:"
echo "请确保服务器已启动并启用了所有功能："
echo "./sweb -upload -files -webdav -https -socks5 -proxy"
echo

echo "Starting comprehensive tests..."
echo "开始综合测试..."
echo

show_menu() {
    echo "========================================"
    echo "Test Menu / 测试菜单"
    echo "========================================"
    echo "1. Run all tests / 运行所有测试"
    echo "2. Basic functionality test / 基础功能测试"
    echo "3. SOCKS5 proxy test / SOCKS5代理测试"
    echo "4. HTTP/HTTPS proxy test / HTTP/HTTPS代理测试"
    echo "5. Debug connection test / 调试连接测试"
    echo "6. Quick SOCKS5 test / SOCKS5快速测试"
    echo "0. Exit / 退出"
    echo "========================================"
}

basic_test() {
    echo
    echo "========================================"
    echo "1. Basic Functionality Test"
    echo "1. 基础功能测试"
    echo "========================================"
    echo
    echo "Testing server status and basic features..."
    echo "测试服务器状态和基础功能..."
    go run test_all.go
    echo
    read -p "Press Enter to continue / 按回车键继续..."
}

socks5_test() {
    echo
    echo "========================================"
    echo "2. SOCKS5 Proxy Test"
    echo "2. SOCKS5代理测试"
    echo "========================================"
    echo
    echo "Compiling SOCKS5 test..."
    echo "编译SOCKS5测试..."
    go build -o test_socks5_temp test_socks5.go
    if [ $? -ne 0 ]; then
        echo "✗ SOCKS5 test compilation failed"
        echo "✗ SOCKS5测试编译失败"
        read -p "Press Enter to continue / 按回车键继续..."
        return
    fi

    echo "Running SOCKS5 proxy test..."
    echo "运行SOCKS5代理测试..."
    ./test_socks5_temp
    rm -f test_socks5_temp
    echo
    read -p "Press Enter to continue / 按回车键继续..."
}

proxy_test() {
    echo
    echo "========================================"
    echo "3. HTTP/HTTPS Proxy Test"
    echo "3. HTTP/HTTPS代理测试"
    echo "========================================"
    echo
    echo "Compiling HTTP proxy test..."
    echo "编译HTTP代理测试..."
    go build -o test_proxy_temp test_proxy.go
    if [ $? -ne 0 ]; then
        echo "✗ HTTP proxy test compilation failed"
        echo "✗ HTTP代理测试编译失败"
        read -p "Press Enter to continue / 按回车键继续..."
        return
    fi

    echo "Running HTTP/HTTPS proxy test..."
    echo "运行HTTP/HTTPS代理测试..."
    ./test_proxy_temp
    rm -f test_proxy_temp
    echo
    read -p "Press Enter to continue / 按回车键继续..."
}

debug_test() {
    echo
    echo "========================================"
    echo "4. Debug Connection Test"
    echo "4. 调试连接测试"
    echo "========================================"
    echo
    echo "Compiling debug test..."
    echo "编译调试测试..."
    go build -o debug_socks5_temp debug_socks5.go
    if [ $? -ne 0 ]; then
        echo "✗ Debug test compilation failed"
        echo "✗ 调试测试编译失败"
        read -p "Press Enter to continue / 按回车键继续..."
        return
    fi

    echo "Running debug connection test..."
    echo "运行调试连接测试..."
    ./debug_socks5_temp
    rm -f debug_socks5_temp
    echo
    read -p "Press Enter to continue / 按回车键继续..."
}

quick_test() {
    echo
    echo "========================================"
    echo "5. Quick SOCKS5 Test"
    echo "5. SOCKS5快速测试"
    echo "========================================"
    echo
    echo "Compiling main program..."
    echo "编译主程序..."
    cd ..
    go build -o sweb_temp main.go upload.go files.go webdav.go server.go utils.go socks5.go proxy.go
    if [ $? -ne 0 ]; then
        echo "✗ Main program compilation failed"
        echo "✗ 主程序编译失败"
        cd test
        read -p "Press Enter to continue / 按回车键继续..."
        return
    fi
    cd test

    echo "✓ Compilation successful"
    echo "✓ 编译成功"
    echo
    echo "Please run in another terminal: ../sweb_temp -socks5"
    echo "请在另一个终端运行: ../sweb_temp -socks5"
    echo "Then press Enter to continue testing..."
    echo "然后按回车键继续测试..."
    read

    echo
    echo "Testing SOCKS5 connection..."
    echo "测试SOCKS5连接..."
    go run debug_socks5.go

    rm -f ../sweb_temp
    echo
    read -p "Press Enter to continue / 按回车键继续..."
}

all_tests() {
    echo "========================================"
    echo "Running All Tests / 运行所有测试"
    echo "========================================"
    basic_test
    socks5_test
    proxy_test
    debug_test
}

# Main menu loop
while true; do
    show_menu
    read -p "Please select an option / 请选择选项 (0-6): " choice
    
    case $choice in
        0)
            break
            ;;
        1)
            all_tests
            ;;
        2)
            basic_test
            ;;
        3)
            socks5_test
            ;;
        4)
            proxy_test
            ;;
        5)
            debug_test
            ;;
        6)
            quick_test
            ;;
        *)
            echo "Invalid option / 无效选项"
            ;;
    esac
done

echo
echo "========================================"
echo "Test Suite Completed!"
echo "测试套件完成！"
echo "========================================"
echo
echo "Summary / 总结:"
echo "- Server status: http://localhost:8080/api/status"
echo "- 服务器状态: http://localhost:8080/api/status"
echo "- SOCKS5 proxy: localhost:1080"
echo "- SOCKS5代理: localhost:1080"
echo "- HTTP proxy: localhost:10808"
echo "- HTTP代理: localhost:10808"
echo
echo "Browser configuration tips:"
echo "浏览器配置提示："
echo "1. For SOCKS5: Set SOCKS5 proxy to 127.0.0.1:1080"
echo "1. SOCKS5: 设置SOCKS5代理为 127.0.0.1:1080"
echo "2. For HTTP/HTTPS: Set HTTP proxy to 127.0.0.1:10808"
echo "2. HTTP/HTTPS: 设置HTTP代理为 127.0.0.1:10808"
echo
echo "Thank you for using SWeb!"
echo "感谢使用 SWeb！"
echo
