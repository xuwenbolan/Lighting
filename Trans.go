package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

// 是否压缩
var isCompress = false

// GetLocalIP 返回本机的一个非环回 IPv4 地址。
func getLocalIP() (string, error) {
	// 创建一个UDP连接，以此获取本地IP地址
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr)
	info, err := os.Stat(sourcePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	if info.IsDir() || isCompress {
		// 如果是文件夹，或者是文件且需要压缩
		w.Header().Set("Content-Disposition", "attachment; filename=\""+info.Name()+".zip\"")
		w.Header().Set("Content-Type", "application/zip")
		err = compress(sourcePath, w)
		if err != nil {
			http.Error(w, "Could not compress file", http.StatusInternalServerError)
			return
		}
	} else {
		// 如果是文件且不需要压缩，直接发送文件
		file, err := os.Open(sourcePath)
		if err != nil {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", "attachment; filename=\""+info.Name()+"\"")
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Error while sending file.", http.StatusInternalServerError)
			return
		}
	}
}
