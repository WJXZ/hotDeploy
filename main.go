package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	upload_path string = "./upload/"
)

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func getAllFile(pathname string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	var s []string
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			s = append(s, fi.Name())
		}
	}
	return s, nil
}

//上传
func uploadHandle(w http.ResponseWriter, r *http.Request) {
	//从请求当中判断方法
	if r.Method == "GET" {
		files, _ := getAllFile(upload_path)
		files_str := strings.Join(files, ",")

		io.WriteString(w, "<html><head><title>上传</title></head>"+
			"<body><form action='#' method=\"post\" enctype=\"multipart/form-data\">"+
			"<label>文件列表</label>: "+files_str+"<br/><br/>    "+
			"<label>上传脚本</label>: "+
			"<input type=\"file\" name='file'  /><br/><br/>    "+
			"<label><input type=\"submit\" value=\"上传脚本文件\"/></label></form></body></html>")
	} else {
		//获取文件内容 要这样获取
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//创建文件
		fW, err := os.Create(upload_path + head.Filename)
		if err != nil {
			fmt.Println("文件创建失败")
			return
		}
		defer fW.Close()
		_, err = io.Copy(fW, file)
		if err != nil {
			fmt.Println("文件保存失败")
			return
		}
		io.WriteString(w, "上传成功!")
	}
}

func reLaunch(s string) (info string) {
	info = ""
	if IsExist(upload_path + s) {
		cmd := exec.Command("sh", upload_path+s)
		err := cmd.Start()
		if err != nil {
			info = "脚本错误"
		}
		cmd.Wait()
	} else {
		info = "文件不存在"
	}
	return
}

func deployHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	script := r.Form.Get("id")
	err := reLaunch(script)
	res := "deploy success"
	if err != "" {
		res = err
	}
	io.WriteString(w, "<h1>"+res+"</h1>")
}

func main() {
	http.HandleFunc("/upload", uploadHandle)
	http.HandleFunc("/deploy", deployHandle)
	err := http.ListenAndServe(":88", nil)
	if err != nil {
		fmt.Println("服务器启动失败")
		return
	}
	fmt.Println("服务器启动成功")
}
