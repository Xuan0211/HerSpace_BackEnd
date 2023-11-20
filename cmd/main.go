package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

/*
| user_info | CREATE TABLE `user_info` (

	`id` int(11) NOT NULL AUTO_INCREMENT,
	`name` char(255) DEFAULT NULL,
	`avatar` char(255) DEFAULT NULL,
	`wechat_code` char(255) DEFAULT NULL,
	`created_time` datetime DEFAULT NULL,
	`followed_count` int(11) DEFAULT NULL,
	`following_count` int(11) DEFAULT NULL,
	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 |
*/
type UserInfo struct {
	Id             int
	Name           sql.NullString `json:"name"`
	Avatar         sql.NullString
	WechatCode     string
	CreatedTime    sql.NullString
	FollowedCount  sql.NullInt32
	FollowingCount sql.NullInt32
	Intro          sql.NullString `json:"intro"`
}

/*
avatar:string
backgroundUrl:string
email:string
intro:string
level:int
name:string
password:string
phone:string
userId:string
*/
type UserInfoAdd struct {
	Avatar        string `json:"avatar"`
	BackgroundUrl string `json:"backgroundUrl"`
	Email         string `json:"email"`
	Intro         string `json:"intro"`
	Level         int    `json:"level"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	UserId        int    `json:"userId"`
}

func init() {
	fmt.Println("init")
	var err error
	db, err = sql.Open("mysql", "HerSpace:88516098@/HerSpace")
	if err != nil {
		panic(err.Error())
	}
}
func sql_user_info_query_by_wechatcode(WechatCode string) (ret UserInfo, ret_err error) {
	sql := "select * from user_info where wechat_code = '" + WechatCode + "'"
	var res UserInfo
	err := db.QueryRow(sql).Scan(&res.Id, &res.Name, &res.Avatar, &res.WechatCode, &res.CreatedTime, &res.FollowedCount, &res.FollowingCount, &res.Intro)
	if err != nil {
		return res, err
	} else {
		return res, nil
	}
}
func sql_user_info_insert_wechatcode(WechatCode string) (ret_err error) {
	sql := "insert into user_info (wechat_code) values (?)"
	_, err := db.Exec(sql, WechatCode)
	return err
}

func sql_login(WechatCode string) (ret UserInfo, ret_err error) {
	res, err := sql_user_info_query_by_wechatcode(WechatCode)
	if err != nil {
		sql_user_info_insert_wechatcode(WechatCode)
		res, err = sql_user_info_query_by_wechatcode(WechatCode)
		return res, err
	} else {
		return res, nil
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//fmt.Println(r.Form["code"])
	res, err := sql_login(r.Form["code"][0])
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
		"data": res,
	}
	json.NewEncoder(w).Encode(data)
}

func insert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is insert")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res UserInfoAdd
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	sql := "update user_info set name = ?,avatar = ?,intro = ? where id = ?"
	_, err = db.Exec(sql, res.Name, res.Avatar, res.Intro, res.UserId)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
	}
	json.NewEncoder(w).Encode(data)
}

func getOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is getOne")
	r.ParseForm()
	id := r.Form["id"][0]

	sql := "select * from user_info where id = ?"
	var res UserInfo
	err := db.QueryRow(sql, id).Scan(&res.Id, &res.Name, &res.Avatar, &res.WechatCode, &res.CreatedTime, &res.FollowedCount, &res.FollowingCount, &res.Intro)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
		"data": res,
	}
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/sk/users/login", login)   // 设置访问的路由
	http.HandleFunc("/sk/users/insert", insert) // 设置访问的路由
	http.HandleFunc("/sk/users/getOne", getOne) // 设置访问的路由
	err := http.ListenAndServe(":13001", nil)   // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
