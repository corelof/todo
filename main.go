package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fileList := []string{}
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && !isBinaryFile(path) {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, path := range fileList {
		f, _ := os.Open(path)
		scanner := bufio.NewScanner(f)
		ln := 1
		for scanner.Scan() {
			line := scanner.Text()
			var ctd = false
			ctd = contains(line)
			if ctd {
				stat, err := f.Stat()
				if err != nil {
					log.Println(err)
				}
				for len(line) > 0 && line[0] == '	' {
					line = line[1:]
				}
				fmt.Printf("%s:%d - %s\n", stat.Name(), ln, line)
			}
			ln++
		}
		f.Close()
	}
}

func isBinaryFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return true
	}
	defer file.Close()
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		log.Println(err)
		return true
	}
	contentType := http.DetectContentType(buffer[:n])
	return (contentType == "application/octet-stream")
}

func contains(s string) bool {
	str := "TODO"
	n := len(s)
	for i := 0; i < n-3; i++ {
		sub := strings.ToUpper(s[i : i+4])
		if sub == str {
			return true
		}
	}
	return false
}
