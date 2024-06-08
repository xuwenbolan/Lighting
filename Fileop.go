package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// 压缩目录到文件或http.ResponseWriter
func zipDir(source string, target interface{}) error {
	var zipWriter *zip.Writer
	var zipFile *os.File
	var err error

	switch t := target.(type) {
	case string:
		zipFile, err = os.Create(t)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter = zip.NewWriter(zipFile)
	case http.ResponseWriter:
		zipWriter = zip.NewWriter(t)
	default:
		return fmt.Errorf("unsupported target type")
	}
	defer zipWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// 创建zip文件内的路径
		zipPath := filepath.ToSlash(relPath)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

// 压缩单个文件到文件或http.ResponseWriter
func zipFile(source string, target interface{}) error {
	var zipWriter *zip.Writer
	var zipFile *os.File
	var err error

	switch t := target.(type) {
	case string:
		zipFile, err = os.Create(t)
		if err != nil {
			return err
		}
		defer zipFile.Close()
		zipWriter = zip.NewWriter(zipFile)
	case http.ResponseWriter:
		zipWriter = zip.NewWriter(t)
	default:
		return fmt.Errorf("unsupported target type")
	}
	defer zipWriter.Close()

	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(source))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// 入口函数，根据SourcePath是文件还是目录，选择相应的压缩方法
func compress(source string, target interface{}) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return zipDir(source, target)
	} else {
		return zipFile(source, target)
	}
}
