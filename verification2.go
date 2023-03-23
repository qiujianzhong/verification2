package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db = &sql.DB{}

//可以改为不同环境的链接 可以相同
var db1 = "root:123456@tcp(192.168.1.200:3306)/micro_credit?timeout=5s&readTimeout=2s"
var db2 = "root:123456@tcp(192.168.1.200:3306)/micro_credit?timeout=5s&readTimeout=2s"
var db3 = "root:123456@tcp(192.168.1.200:3306)/micro_credit?timeout=5s&readTimeout=2s"
var db4 = "root:123456@tcp(192.168.1.200:3306)/micro_credit?timeout=5s&readTimeout=2s"
var db5 = "root:123456@tcp(192.168.1.200:3306)/micro_credit?timeout=5s&readTimeout=2s"

func main() {
	//路由
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/code", code)
	http.HandleFunc("/sms", sms)
	http.HandleFunc("/getmobile", getmobile)
	http.HandleFunc("/setmobile", setmobile)

	err := http.ListenAndServe(":8889", nil) //监听端口
	if err != nil {
		log.Fatal("listenandserver", err)
	}
}

//主页
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("index.html") //读取首页的html
	log.Println(t.Execute(w, nil))
}

//读取验证码
func code(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var mobile string
	var code string
	var code1 string
	var code2 string
	var code3 string
	var code4 string
	var code5 string
	str_time := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time1 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time2 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time3 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time4 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time5 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")

	if len(r.Form["mobile"]) != 0 {
		mobile = r.Form["mobile"][0]
	}

	db, _ = sql.Open("mysql", db1)
	defer db.Close()
	db.QueryRow("SELECT  `code`,createTime FROM `t_wd_sms_log` WHERE `code` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?   ORDER BY id DESC ", mobile).Scan(&code1, &str_time1)

	db, _ = sql.Open("mysql", db2)
	defer db.Close()
	db.QueryRow("SELECT  `code`,createTime FROM `t_wd_sms_log` WHERE `code` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&code2, &str_time2)

	db, _ = sql.Open("mysql", db3)
	defer db.Close()
	db.QueryRow("SELECT  `code`,createTime FROM `t_wd_sms_log` WHERE `code` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&code3, &str_time3)

	db, _ = sql.Open("mysql", db4)
	defer db.Close()
	db.QueryRow("SELECT  `code`,createTime FROM `t_wd_sms_log` WHERE `code` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&code4, &str_time4)

	db, _ = sql.Open("mysql", db5)
	defer db.Close()
	db.QueryRow("SELECT  `code`,createTime FROM `t_wd_sms_log` WHERE `code` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&code5, &str_time5)

	//fmt.Println(str_time1,str_time2,str_time3)
	//fmt.Println(code1,code2,code3)

	if str_time1 > str_time2 {
		code = code1
		str_time = str_time1
	} else {
		code = code2
		str_time = str_time2
	}
	if str_time < str_time3 {
		code = code3
	}
	if str_time < str_time4 {
		code = code4
		str_time = str_time4
	}
	if str_time < str_time5 {
		code = code5
		//str_time = str_time5
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
	var sms1 string
	var sms2 string
	var sms3 string
	var sms4 string
	var sms5 string

	str_time := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time1 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time2 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time3 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time4 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")
	str_time5 := time.Unix(1389058332, 0).Format("2006-01-02 15:04:05")

	if len(r.Form["mobile"]) != 0 {
		mobile = r.Form["mobile"][0]
	}

	db, _ = sql.Open("mysql", db1)
	defer db.Close()
	db.QueryRow("SELECT  `smsContent`,createTime FROM `t_wd_sms_log` WHERE `smsContent` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?   ORDER BY id DESC ", mobile).Scan(&sms1, &str_time1)

	db, _ = sql.Open("mysql", db2)
	defer db.Close()
	db.QueryRow("SELECT  `smsContent`,createTime FROM `t_wd_sms_log` WHERE `smsContent` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&sms2, &str_time2)

	db, _ = sql.Open("mysql", db3)
	defer db.Close()
	db.QueryRow("SELECT  `smsContent`,createTime FROM `t_wd_sms_log` WHERE `smsContent` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&sms3, &str_time3)

	db, _ = sql.Open("mysql", db4)
	defer db.Close()
	db.QueryRow("SELECT  `smsContent`,createTime FROM `t_wd_sms_log` WHERE `smsContent` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&sms4, &str_time4)

	db, _ = sql.Open("mysql", db5)
	defer db.Close()
	db.QueryRow("SELECT  `smsContent`,createTime FROM `t_wd_sms_log` WHERE `smsContent` IS NOT NULL AND TIMEDIFF(NOW(), createTime) < '00:10:00.00.000000' AND mobile= ?  ORDER BY id DESC ", mobile).Scan(&sms5, &str_time5)

	//fmt.Println(str_time1,str_time2,str_time3)
	//fmt.Println(sms1,sms2,sms3)

	if str_time1 > str_time2 {
		sms = sms1
		str_time = str_time1
	} else {
		sms = sms2
		str_time = str_time2
	}
	if str_time < str_time3 {
		sms = sms3
	}
	if str_time < str_time4 {
		sms = sms4
		str_time = str_time4
	}
	if str_time < str_time5 {
		sms = sms5
		//str_time = str_time5
	}

	if sms != "" {
		log.Println("mobile:" + mobile + ",sms:" + sms)
		fmt.Fprintf(w, sms)
	} else {
		fmt.Fprintf(w, "亲，没有查到短信♥")
	}

}

//读取订单对应的登录手机号
func getmobile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var orderid string
	var mobile string

	if len(r.Form["orderid"]) != 0 {
		orderid = r.Form["orderid"][0]
	} else {
		fmt.Fprintf(w, "请输入订单号♥")
	}

	db, _ = sql.Open("mysql", db1)
	defer db.Close()
	db.QueryRow("SELECT mobile FROM `credit_finger_bank`.`fg_credit_user` WHERE user_id in (SELECT channel_user_id FROM `credit_finger_bank`.`fg_market_order` WHERE `order_id` = ? ) ", orderid).Scan(&mobile)

	if mobile != "" {
		//log.Println("mobile:"mobile+",sms:"+sms)
		fmt.Fprintf(w, mobile)
	} else {
		fmt.Fprintf(w, "未查到订单♥")
	}

}

//修改订单对应的登录手机号
func setmobile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var orderid string
	var mobile string

	if len(r.Form["orderid"]) != 0 {
		orderid = r.Form["orderid"][0]
	} else {
		fmt.Fprintf(w, "请输入订单号♥")
	}
	if len(r.Form["mobile"]) != 0 {
		mobile = r.Form["mobile"][0]
	} else {
		fmt.Fprintf(w, "请输入手机号♥")
	}

	db, _ = sql.Open("mysql", db1)
	defer db.Close()

	_, err := db.Exec("UPDATE `credit_finger_bank`.`fg_credit_user` SET mobile = ? WHERE user_id in (SELECT channel_user_id FROM `credit_finger_bank`.`fg_market_order` WHERE `order_id` = ? )", mobile, orderid)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "更新失败")
	} else {
		log.Println("更新订单:" + orderid + ",mobile:" + mobile)
		fmt.Fprintf(w, "更新成功")
	}

}
