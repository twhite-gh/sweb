# 定义项目名称，通常是你的模块名或你想生成的可执行文件的名称
PROJECT_NAME := sweb

# 源文件入口，通常是包含 main 函数的文件
MAIN_FILE := ./main.go

# 定义输出目录
BUILD_DIR := ./bin

.PHONY: all clean windows linux

all: windows linux

windows:
	@echo "Building for Windows..."
	mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(PROJECT_NAME).exe $(MAIN_FILE)
	@echo "Windows build complete: $(BUILD_DIR)/windows/$(PROJECT_NAME).exe"

linux:
	@echo "Building for Linux..."
	mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "Linux build complete: $(BUILD_DIR)/linux/$(PROJECT_NAME)"

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete."