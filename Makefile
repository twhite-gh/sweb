# 定义项目名称，通常是你的模块名或你想生成的可执行文件的名称
PROJECT_NAME := sweb

# 源文件列表，包含所有需要编译的Go文件
SOURCE_FILES := ./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go

# 定义输出目录
BUILD_DIR := ./bin

# 编译标志
LDFLAGS := -ldflags "-s -w"

.PHONY: all clean windows windows-32 windows-64 linux linux-32 linux-64 macos macos-intel macos-arm64

all: windows linux macos

# Windows builds
windows: windows-64 windows-32

windows-64:
	@echo "Building for Windows (64-bit)..."
	mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(PROJECT_NAME)-windows-amd64.exe $(SOURCE_FILES)
	@echo "Windows 64-bit build complete: $(BUILD_DIR)/windows/$(PROJECT_NAME)-windows-amd64.exe"

windows-32:
	@echo "Building for Windows (32-bit)..."
	mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(PROJECT_NAME)-windows-386.exe $(SOURCE_FILES)
	@echo "Windows 32-bit build complete: $(BUILD_DIR)/windows/$(PROJECT_NAME)-windows-386.exe"

# Linux builds
linux: linux-64 linux-32

linux-64:
	@echo "Building for Linux (64-bit)..."
	mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(PROJECT_NAME)-linux-amd64 $(SOURCE_FILES)
	@echo "Linux 64-bit build complete: $(BUILD_DIR)/linux/$(PROJECT_NAME)-linux-amd64"

linux-32:
	@echo "Building for Linux (32-bit)..."
	mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(PROJECT_NAME)-linux-386 $(SOURCE_FILES)
	@echo "Linux 32-bit build complete: $(BUILD_DIR)/linux/$(PROJECT_NAME)-linux-386"

# macOS builds
macos: macos-intel macos-arm64

macos-intel:
	@echo "Building for macOS (Intel)..."
	mkdir -p $(BUILD_DIR)/macos
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos/$(PROJECT_NAME)-darwin-amd64 $(SOURCE_FILES)
	@echo "macOS Intel build complete: $(BUILD_DIR)/macos/$(PROJECT_NAME)-darwin-amd64"

macos-arm64:
	@echo "Building for macOS (Apple Silicon)..."
	mkdir -p $(BUILD_DIR)/macos
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos/$(PROJECT_NAME)-darwin-arm64 $(SOURCE_FILES)
	@echo "macOS Apple Silicon build complete: $(BUILD_DIR)/macos/$(PROJECT_NAME)-darwin-arm64"

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete."