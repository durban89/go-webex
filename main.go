package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, nil))
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

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
