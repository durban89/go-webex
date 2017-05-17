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

	"github.com/durban.zhang/webex/helpers/session"
	_ "github.com/durban.zhang/webex/helpers/session/providers/memory"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
	sess := globalSessions.SessionStart(w, r)
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
		fmt.Println(err)

		// session
		sess.Set("username", r.Form["username"][0])
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
	stmt, err := db.Prepare("INSERT INTO userinfo(username,departname, created) values (?,?,?)")
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

func sqliteEx(w http.ResponseWriter, r *http.Request) {
	// CREATE TABLE `userinfo` (
	// 	`uid` INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	`username` VARCHAR(64) NULL,
	// 	`departname` VARCHAR(64) NULL,
	// 	`created` DATE NULL
	// );

	// CREATE TABLE `userdeatail` (
	// 	`uid` INT(10) NULL,
	// 	`intro` TEXT NULL,
	// 	`profile` TEXT NULL,
	// 	PRIMARY KEY (`uid`)
	// );

	db, err := sql.Open("sqlite3", "./data/test.db")
	checkErr(err)

	// //插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username,departname, created) values (?,?,?)")
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

func postgresqlEx(w http.ResponseWriter, r *http.Request) {
	// CREATE TABLE userinfo
	// (
	// 	uid serial NOT NULL,
	// 	username character varying(100) NOT NULL,
	// 	departname character varying(500) NOT NULL,
	// 	Created date,
	// 	CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
	// )
	// WITH (OIDS=FALSE);

	// CREATE TABLE userdeatail
	// (
	// 	uid integer,
	// 	intro character varying(100),
	// 	profile character varying(100)
	// )
	// WITH(OIDS=FALSE);

	db, err := sql.Open("postgres", "user=root password=123456 dbname=test sslmode=disable")
	checkErr(err)

	// insert data
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) VALUES ($1,$2,$3) RETURNING uid")
	checkErr(err)

	res, err := stmt.Exec("durban", "研发部门", "2016-05-12")
	checkErr(err)
	fmt.Println(res)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username, departname, created) VALUES ($1,$2,$3) RETURNING uid;", "durban", "研发部门", "2016-05-12").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("最后插入的ID", lastInsertId)

	// update data
	stmt, err = db.Prepare("update userinfo set username = $1 where uid = $2")
	checkErr(err)

	res, err = stmt.Exec("durban2", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)

	// query data
	rows, err := db.Query("SELECT * FROM userinfo")
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
	stmt, err = db.Prepare("delete from userinfo where uid = $1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}

func cookieEx(w http.ResponseWriter, r *http.Request) {
	// set cookie
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "username", Value: "durban", Expires: expiration}
	http.SetCookie(w, &cookie)
	// get cookie
	cookie1, _ := r.Cookie("username")
	fmt.Fprintln(w, cookie1)

	fmt.Fprintln(w, "All Cookie ========== ")
	// get all cookie
	for _, cookie2 := range r.Cookies() {
		fmt.Fprintln(w, cookie2.Name)
	}
}

func sessionEx(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	fmt.Println("SessionStart==")
	fmt.Println(sess)
	fmt.Fprintln(w, sess.Get("username"))
	// t.Execute(w, sess.Get("countnum"))
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/mysqlEx", mysqlEx)
	http.HandleFunc("/sqliteEx", sqliteEx)
	http.HandleFunc("/postgresqlEx", postgresqlEx)
	http.HandleFunc("/cookieEx", cookieEx)
	http.HandleFunc("/sessionEx", sessionEx)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}
