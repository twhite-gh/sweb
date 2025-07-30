package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 获取当前可执行文件的目录
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execDir := filepath.Dir(execPath)
	
	// 构建src目录中的实际程序
	srcDir := filepath.Join(execDir, "src")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = srcDir
	cmd.Args = append(cmd.Args, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	err = cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
