package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
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
