package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	io "io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db = &sql.DB{}
var conf DBConf

var file_locker sync.Mutex //config file locker

type DBConf struct {
	Name       string   `json:"name"`
	Port       string   `json:"port"`
	DriverName string   `json:"driver_name"`
	SQLCode    string   `json:"sql_code"`
	SQSms      string   `json:"sql_sms"`
	SQLConn    []string `json:"sql_conn"`
}

func LoadConfig(filename string) (DBConf, bool) {
	file_locker.Lock()
	data, err := io.ReadFile(filename) //read config file
	file_locker.Unlock()
	if err != nil {
		fmt.Println("read json file error")
		return conf, false
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		fmt.Println("unmarshal json file error")
		return conf, false
	}
	return conf, true
}

func main() {
	//路由
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/code", code)
	http.HandleFunc("/sms", sms)

	conf, ok := LoadConfig("./server.json")
	if !ok {
		fmt.Println("load config failed")
		return
	}
	fmt.Println(conf)
	for _, sqlConn := range conf.SQLConn {
		db, err := sql.Open(conf.DriverName, sqlConn)
		if err != nil {
			panic("cb sql connect error:" + err.Error())
		}
		err = db.Ping()
		if err != nil {
			panic("cb sql can not ping:" + err.Error())
		}
	}
	fmt.Println("ALL DB CONF ARE CORRECT")

	err := http.ListenAndServe(":"+conf.Port, nil) //监听端口
	if err != nil {
		log.Fatal("listenandserver", err)
	}
}

//主页
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("index.html") //读取首页的html
	m := make(map[string]interface{})
	m["Name"] = conf.Name
	log.Println(t.Execute(w, m))
}

//读取验证码
func code(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var mobile string
	var code string
	var code_temp string
	var str_time = time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	var str_time_temp = time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")

	if len(r.Form["mobile"]) != 0 {
		mobile = r.Form["mobile"][0]
	} else {
		fmt.Fprintf(w, "亲，没有查到验证码♥")
	}

	for _, sqlConn := range conf.SQLConn {
		db, err := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		db.QueryRow(conf.SQLCode, mobile).Scan(&code_temp, &str_time_temp)
		if err != nil {
			fmt.Println("can not open db => ", sqlConn)
			fmt.Println("please check out your server.json")
		}
		//fmt.Println(" " + str_time_temp + " " + code_temp)
		if str_time_temp > str_time {
			code = code_temp
			str_time = str_time_temp
		}

	}
	if code != "" {
		log.Println("mobile:" + mobile + ",code:" + code)
		fmt.Fprintf(w, code)
	} else {
		fmt.Fprintf(w, "亲，没有查到验证码♥")
	}

}

//读取短信
func sms(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var mobile string
	var sms string
	var sms_temp string
	var str_time = time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	var str_time_temp = time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")

	if len(r.Form["mobile"]) != 0 {
		mobile = r.Form["mobile"][0]
	} else {
		fmt.Fprintf(w, "亲，没有查到短信♥")
	}

	for _, sqlConn := range conf.SQLConn {
		db, err := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		db.QueryRow(conf.SQSms, mobile).Scan(&sms_temp, &str_time_temp)
		if err != nil {
			fmt.Println("can not open db => ", sqlConn)
			fmt.Println("please check out your server.json")
		}
		//fmt.Println(" " + str_time_temp + " " + code_temp)
		if str_time_temp > str_time {
			sms = sms_temp
			str_time = str_time_temp
		}

	}
	if sms != "" {
		log.Println("mobile:" + mobile + ",sms:" + sms)
		fmt.Fprintf(w, sms)
	} else {
		fmt.Fprintf(w, "亲，没有查到短信♥")
	}

}
