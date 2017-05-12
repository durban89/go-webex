package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mysqlEx(w http.ResponseWriter, r *http.Request) {
	// CREATE TABLE `userinfo` (
	// 	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	// 	`username` VARCHAR(64) NULL DEFAULT NULL,
	// 	`departname` VARCHAR(64) NULL DEFAULT NULL,
	// 	`created` DATE NULL DEFAULT NULL,
	// 	PRIMARY KEY (`uid`)
	// );

	// CREATE TABLE `userdetail` (
	// 	`uid` INT(10) NOT NULL DEFAULT '0',
	// 	`intro` TEXT NULL,
	// 	`profile` TEXT NULL,
	// 	PRIMARY KEY (`uid`)
	// )
	db, err := sql.Open("mysql", "root:123456@/test?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	//update data
	stmt, err = db.Prepare("update userinfo set username=? where uid = ?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//query data
	rows, err := db.Query("select * from userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)

	}

	// delete data
	stmt, err = db.Prepare("delete from userinfo where uid = ?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	db.Close()
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/mysqlEx", mysqlEx)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
