# 定义项目名称，通常是你的模块名或你想生成的可执行文件的名称
PROJECT_NAME := sweb

# 源文件列表，包含所有需要编译的Go文件
SOURCE_FILES := ./main.go ./upload.go ./files.go ./webdav.go ./server.go ./utils.go ./socks5.go ./proxy.go

# 定义输出目录
BUILD_DIR := ./bin

.PHONY: all clean windows linux

all: windows linux

windows:
	@echo "Building for Windows..."
	mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(PROJECT_NAME).exe $(SOURCE_FILES)
	@echo "Windows build complete: $(BUILD_DIR)/windows/$(PROJECT_NAME).exe"

linux:
	@echo "Building for Linux..."
	mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(PROJECT_NAME) $(SOURCE_FILES)
	@echo "Linux build complete: $(BUILD_DIR)/linux/$(PROJECT_NAME)"

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete."