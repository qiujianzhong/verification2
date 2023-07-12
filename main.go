package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	io "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db = &sql.DB{}
var conf DBConf
var conf_ []string

var file_locker sync.Mutex //config file locker

type DBConf struct {
	Name       string   `json:"name"`
	Ip         string   `json:"ip"`
	Port       string   `json:"port"`
	DriverName string   `json:"driver_name"`
	Conn_dev1  []string `json:"conn_dev1"`
	Conn_dev3  []string `json:"conn_dev3"`
	Conn_dev5  []string `json:"conn_dev5"`
}

type Users struct {
	Info     string
	Uid      string
	Mail     string
	Phone    string
	Nickname string
	Tier     string
	Verified string
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
	http.HandleFunc("/query", query)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/verified", verified)
	http.HandleFunc("/tier", tier)
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/useradd", useradd)
	http.HandleFunc("/price", price)

	conf, ok := LoadConfig("./server.json")
	if !ok {
		fmt.Println("load config failed")
		return
	}
	fmt.Println(conf)
	for _, sqlConn := range conf.Conn_dev3 {
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

	err := http.ListenAndServe(conf.Ip+":"+conf.Port, nil) //监听端口
	if err != nil {
		log.Fatal("listenandserver", err)
	}
}

//主页
func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, _ := template.ParseFiles("static/index.html") //读取首页的html
	m := make(map[string]interface{})
	m["Name"] = conf.Name
	log.Println(t.Execute(w, m))
}

//判断查询类型
func types(query1 string) string {
	types := ""
	if strings.Contains(query1, "@") {
		types = "email"
	} else {
		if _, err := strconv.Atoi(query1); err == nil {
			if len(query1) == 11 {
				types = "phone"
			} else {
				types = "uid"
			}
		} else {
			types = "nickname"
		}
	}
	return types

}

func GetRandomString(l int) string {
	str1 := ""
	for i := 0; i < l; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(9)
		str1 = str1 + strconv.Itoa(r)
	}
	// fmt.Println(str1)
	return string(str1)
}

//useradd
func useradd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var user1 string
	var phone1 string
	var mail1 string

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if len(r.Form["mail"]) != 0 {
		mail1 = r.Form["mail"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入mail'}")
	}

	if len(r.Form["phone"][0]) > 11 {
		fmt.Fprintf(w, "{'Info':'请输入mail'}")
	} else if len(r.Form["phone"][0]) == 11 {
		phone1 = r.Form["phone"][0]
	} else if len(r.Form["phone"][0]) == 10 {
		phone1 = "1" + (r.Form["phone"][0])
	} else if len(r.Form["phone"][0]) == 0 {
		phone1 = "13" + GetRandomString(9-len(r.Form["phone"][0]))
	} else {
		phone1 = (r.Form["phone"][0]) + GetRandomString(11-len(r.Form["phone"][0]))
	}

	user1 = r.Form["user"][0] + phone1
	mail1 = phone1 + r.Form["mail"][0]

	if dev == "dev1" {
		conf_ = conf.Conn_dev1
	}
	if dev == "dev3" {
		conf_ = conf.Conn_dev3
	}
	if dev == "dev5" {
		conf_ = conf.Conn_dev5
	}

	// insert_cmd := `INSERT INTO gateio.user (nickname,email,password,fundpass,login2,verified,deposref,timest,invested,flag,phone,receiveadminemails,firstname,lastname,loggedin,loggedin_timest,tier_timest,tier,regtime,photoid,iprestrict,language,ip,ip_reg,admin,ref_uid,email_verified,country,is_sub,sub_status,main_uid,sub_remark,point_type,is_broker,ref_ratio,ref_ratio_str,sub_type,sub_website_id,agency_type,type) VALUES ("` + user1 + `","` + mail1 + `","2608a4e7704861106de689e850dab18e","873017f8c6513904a4f8ee5bf07b803b",1,1,"81054449","2025-07-25 17:10:27",5,0,"` + phone1 + `","No","测试工具","",1,"2022-09-27 17:57:40","2023-11-15 21:15:53",4,0,"362322199510203424",1,"cn","112.65.61.127","127.0.0.1",0,0,1,"83",0,0,0,"",0,0,0,"3.1-0.75",0,0,0,0)`
	insert_cmd := `INSERT INTO gateio.user (nickname,email,password,fundpass,login2,verified,deposref,timest,invested,flag,phone,receiveadminemails,firstname,lastname,loggedin,loggedin_timest,tier_timest,tier,regtime,photoid,iprestrict,language,ip,ip_reg,admin,ref_uid,email_verified,country,is_sub,sub_status,main_uid,sub_remark,point_type,is_broker,ref_ratio,ref_ratio_str,sub_type,sub_website_id,agency_type,type) 
	VALUES (?,?,"2608a4e7704861106de689e850dab18e","873017f8c6513904a4f8ee5bf07b803b",1,1,"81054449","2025-07-25 17:10:27",5,0, ? ,"No","测试工具","",1,"2022-09-27 17:57:40","2023-11-15 21:15:53",4,0,"362322199510203424",1,"cn","112.65.61.127","127.0.0.1",0,0,1,"83",0,0,0,"",0,0,0,"3.1-0.75",0,0,0,0)`

	for _, sqlConn := range conf_ {
		db, _ := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		_, err := db.Exec(insert_cmd, user1, mail1, phone1)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{'Info':'创建会员失败'}")
		} else {
			log.Println("创建会员成功：用户名:" + user1 + " 手机:" + phone1 + " 邮箱:" + mail1)
			fmt.Fprintf(w, "创建会员成功："+user1+" 手机:"+phone1+" 邮箱:"+mail1+" 登录密码: 1234qwer 支付密码: 123123")
		}

	}
}

//修改市场价格
func price(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var market string
	var price string

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "dev错误")
	}

	if len(r.Form["market"]) != 0 {
		market = r.Form["market"][0]
	} else {
		fmt.Fprintf(w, "market错误")
	}

	if len(r.Form["price"]) != 0 {
		price = r.Form["price"][0]
	} else {
		fmt.Fprintf(w, "price错误")
	}

	rpc_url := "http://172.16.226.191:42345/usdt/contracts/" + market + "/index_price"

	if dev == "dev1" {
		rpc_url = "http://172.16.226.191:42345/usdt/contracts/" + market + "/index_price"
	}
	if dev == "dev3" {
		rpc_url = "http://172.16.226.191:42345/usdt/contracts/" + market + "/index_price"
	}
	if dev == "dev5" {
		rpc_url = "http://172.16.226.191:42345/usdt/contracts/" + market + "/index_price"
	}

	//json序列化
	post := "price=" + price

	fmt.Println(rpc_url, "post", post)

	var jsonStr = []byte(post)
	//fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	if strings.Contains(string(body), "ok") {
		log.Println("执行成功:market:" + market + ",price:" + price)
		fmt.Fprintf(w, "执行成功")
	} else {
		log.Println(err)
		fmt.Fprintf(w, "执行失败"+string(body))

	}

}

//deposit
func deposit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var query1 string
	var coin string
	var amount string

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if len(r.Form["query1"]) != 0 {
		query1 = r.Form["query1"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入'}")
	}

	if len(r.Form["coin"]) != 0 {
		coin = r.Form["coin"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入coin'}")
	}
	if len(r.Form["amount"]) != 0 {
		amount = r.Form["amount"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入amount'}")
	}

	// query_type := types(query1)

	rpc_url := "http://118.178.105.202:8401/spotme"

	if dev == "dev1" {
		rpc_url = "http://118.178.105.202:8401/spotme"
	}
	if dev == "dev3" {
		rpc_url = "http://118.178.105.202:8401/spotme"
	}
	if dev == "dev5" {
		rpc_url = "http://118.178.105.202:8401/spotme"
	}

	//json序列化
	post := `{"id": 1515752473250, "method": "balance.update","params": [` + query1 + `, "` + coin + `", "deposit", ` + strconv.FormatInt(time.Now().Unix(), 10) + `, "` + amount + `" , {}, "force"]}`

	fmt.Println(rpc_url, "post", post)

	var jsonStr = []byte(post)
	//fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", rpc_url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	if strings.Contains(string(body), "success") {
		log.Println("充值成功:" + query1 + ",coin:" + coin + ",amount:" + amount)
		fmt.Fprintf(w, "充值成功")
	} else {
		log.Println(err)
		fmt.Fprintf(w, "充值失败"+string(body))

	}

}

//tier
func tier(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var query1 string
	var tier string

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if len(r.Form["query1"]) != 0 {
		query1 = r.Form["query1"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入'}")
	}

	if len(r.Form["tier"]) != 0 {
		tier = r.Form["tier"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入tier'}")
	}

	query_type := types(query1)

	if dev == "dev1" {
		conf_ = conf.Conn_dev1
	}
	if dev == "dev3" {
		conf_ = conf.Conn_dev3
	}
	if dev == "dev5" {
		conf_ = conf.Conn_dev5
	}

	for _, sqlConn := range conf_ {
		db, _ := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		_, err := db.Exec("UPDATE `gateio`.`user` SET `tier` = ? ,`tier_timest` = '2030-10-10 00:00:00'  WHERE  "+query_type+" = ? ", tier, query1)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{'Info':'调整会员等级失败'}")
		} else {
			log.Println("调整会员等级成功:" + query_type + ":" + query1)
			fmt.Fprintf(w, "{'Info':'调整会员等级成功'}")
		}

	}
}

//verified
func verified(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var query1 string

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if len(r.Form["query1"]) != 0 {
		query1 = r.Form["query1"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入'}")
	}

	query_type := types(query1)

	if dev == "dev1" {
		conf_ = conf.Conn_dev1
	}
	if dev == "dev3" {
		conf_ = conf.Conn_dev3
	}
	if dev == "dev5" {
		conf_ = conf.Conn_dev5
	}

	for _, sqlConn := range conf_ {
		db, _ := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		_, err := db.Exec("UPDATE `gateio`.`user` SET  photoid = '362322199510203424' ,`verified` = '1' WHERE  "+query_type+" = ? ", query1)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{'Info':'d认证失败'}")
		} else {
			log.Println("认证成功:" + query_type + ":" + query1)
			fmt.Fprintf(w, "{'Info':'认证成功'}")
		}

	}
}

//重置密码
func reset(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var query1 string

	if len(r.Form["query1"]) != 0 {
		query1 = r.Form["query1"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入'}")
	}

	query_type := types(query1)

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if dev == "dev1" {
		conf_ = conf.Conn_dev1
	}
	if dev == "dev3" {
		conf_ = conf.Conn_dev3
	}
	if dev == "dev5" {
		conf_ = conf.Conn_dev5
	}

	for _, sqlConn := range conf_ {
		db, _ := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		_, err := db.Exec("UPDATE `gateio`.`user` SET  fundpass = '873017f8c6513904a4f8ee5bf07b803b' ,`password` = '2608a4e7704861106de689e850dab18e' WHERE  "+query_type+" = ? ", query1)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{'Info':'d更新失败'}")
		} else {
			log.Println("更新成功:" + query_type + ":" + query1)
			fmt.Fprintf(w, "{'Info':'更新成功，登录密码 1234qwer 支付密码 123123'}")
		}

	}
}

//查询会员
func query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var dev string
	var query1 string

	var uid string
	var mail string
	var phone string
	var nickname string
	var tier string
	var verified string

	var SQL = ""

	if len(r.Form["query1"]) != 0 {
		query1 = r.Form["query1"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'请输入'}")
	}
	SQL = types(query1)

	SQL = "select uid,email,phone,nickname,tier,verified from `gateio`.`user` where " + SQL + " = ? limit 1"

	if len(r.Form["dev"]) != 0 {
		dev = r.Form["dev"][0]
	} else {
		fmt.Fprintf(w, "{'Info':'dev选择错误'}")
	}

	if dev == "dev1" {
		conf_ = conf.Conn_dev1
	}
	if dev == "dev3" {
		conf_ = conf.Conn_dev3
	}
	if dev == "dev5" {
		conf_ = conf.Conn_dev5
	}

	for _, sqlConn := range conf_ {
		// fmt.Println("%s %s", SQL, query1)
		db, err := sql.Open(conf.DriverName, sqlConn)
		defer db.Close()
		db.QueryRow(SQL, query1).Scan(&uid, &mail, &phone, &nickname, &tier, &verified)
		if err != nil {
			fmt.Println("can not open db => ", sqlConn)
			fmt.Println("please check out your server.json")
		}

	}
	if uid != "" {
		user := Users{"ok", uid, mail, phone, nickname, tier, verified}
		fmt.Println("%v: \n", user) //查询的会员
		// JSON format:
		js, _ := json.Marshal(user)

		fmt.Fprintf(w, string(js))
	} else {
		fmt.Fprintf(w, "{'Info':'未查询到会员'}")
	}

}
