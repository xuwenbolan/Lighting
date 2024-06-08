package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var sourcePath = "./TestRes"

func main() {
	// 创建一个日志文件
	file, err := os.OpenFile("Lighting.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	// 设置日志的输出目标为文件
	log.SetOutput(file)

	// 如果需要同时将日志输出到控制台，可以使用多重写入器
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	// 检查命令行参数是否正确
	if len(os.Args) != 2 {
		log.Println("Usage: Lighting.exe <file or directory path>")
		return
	}

	// 获取文件或文件夹路径
	sourcePath = os.Args[1]

	if sourcePath == "" {
		return
	}

	log.Println("Selected file or directory:", sourcePath)

	http.HandleFunc("/lighting", downloadHandler)

	log.Println("Server started on :80")
	localIP, err := getLocalIP()
	if err != nil {
		log.Println("Failed to get local IP:", err)
	} else {
		log.Println("http://" + localIP + "/lighting")
	}
	log.Fatal(http.ListenAndServe(":80", nil))

}
