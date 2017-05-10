package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // 解析参数 默认是不会解析的
	fmt.Println(r.Form) //服务器端输出信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("schema", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("value:", strings.Join(v, ""))
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)

		// log.Println(t.Execute(w, nil))
	} else {
		// 表单处理
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		// 表单验证 - 必填字端
		if len(r.Form["username"][0]) == 0 {
			fmt.Println("Username can not empty.")
		}

		// form validate - number
		getint, err := strconv.Atoi(r.Form["username"][0])
		if err != nil {
			fmt.Println("not number")
		}

		if getint > 1000 {
			fmt.Println("number is so larger.")
		}

		// 防止多次递交表单 - 示例
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}

		// 预防跨站脚本 - 示例
		fmt.Println("username:", template.HTMLEscapeString(r.Form["username"][0]))
		fmt.Println("password:", template.HTMLEscapeString(r.Form["password"][0]))

		// template.HTMLEscape(w, []byte(r.Form["username"][0]))

		// t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		// err = t.ExecuteTemplate(w, "T", r.Form["username"][0])

		t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		err = t.ExecuteTemplate(w, "T", template.HTML(r.Form["username"][0]))
	}

}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method:", r.Method)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		fmt.Fprint(w, "%v", handler.Header)

		f, err := os.OpenFile("./static/tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在./static/tmp/目录

		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

	} else if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("views/upload.gtpl")
		t.Execute(w, token)
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
