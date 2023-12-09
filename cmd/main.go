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

type UserInfoUpdate struct {
	Avatar        string `json:"avatar"`
	BackgroundUrl string `json:"backgroundUrl"`
	Intro         string `json:"intro"`
	Name          string `json:"name"`
	UserId        int    `json:"userId"`
}

type PostItem struct {
	Id           int    `json:"postId"`
	Content      string `json:"content"`
	UserId       int    `json:"userId"`
	UserName     string `json:"userName"`
	UserAvatar   string `json:"userAvatarPath"`
	CreateTime   string `json:"createTime"`
	LikeCount    int    `json:"likedCount"`
	ReadCount    int    `json:"readCount"`
	CommentCount int    `json:"commentCount"`
	CategoryId   int    `json:"CategoryId"`
	CategoryName string `json:"categoryName"`
}

type PostList struct {
	Post []PostItem
}

type ReplyItem struct {
	Id             int    `json:"id"`
	CommentId      int    `json:"commentId"`
	CreateTime     string `json:"createTime"`
	Content        string `json:"content"`
	LikeCount      int    `json:"likeCount"`
	FromUserId     int    `json:"fromUserId"`
	FromUserName   string `json:"fromUserName"`
	FromUserAvatar string `json:"fromUserAvatarPath"`
	ToUserId       int    `json:"toUserId"`
	ToUserName     string `json:"toUserName"`
	ToUserAvatar   string `json:"toUserAvatarPath"`
}

type ReplyList struct {
	Reply []ReplyItem
}

type CommentItem struct {
	Id         int    `json:"id"`
	PostId     int    `json:"postId"`
	UserId     int    `json:"userId"`
	UserAvatar string `json:"userAvatarPath"`
	UserName   string `json:"userName"`
	Comment    string `json:"comment"`
	CreateTime string `json:"createTime"`
	LikeCount  int    `json:"likeCount"`
	ReplyCount int    `json:"replyCount"`
	Reply      string `json:"reply"`
}

type CommentList struct {
	Comment []CommentItem
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

func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is update")
	r.ParseForm()
	id := r.Form["id"][0]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res UserInfoUpdate
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	sql := "update user_info set name = ?,avatar = ?,intro = ? where id = ?"
	_, err = db.Exec(sql, res.Name, res.Avatar, res.Intro, id)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
	}
	json.NewEncoder(w).Encode(data)
}

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is get")
	r.ParseForm()
	limit := 10
	sql := "SELECT * FROM post_view ORDER BY create_time LIMIT ?"
	rows, err := db.Query(sql, limit)
	if err != nil {
		log.Fatal(err)
	}
	var res []PostItem
	for rows.Next() {
		var p PostItem
		if err := rows.Scan(&p.Id, &p.Content, &p.UserId, &p.UserName, &p.UserAvatar, &p.CreateTime, &p.LikeCount, &p.ReadCount, &p.CommentCount, &p.CategoryId, &p.CategoryName); err != nil {
			log.Fatal(err)
		}
		res = append(res, p)
	}
	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"list":   res,
			"pageNo": 1,
			"total":  1,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func postGetOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("this is postGetOne")
	r.ParseForm()
	id := r.Form["id"][0]
	fmt.Println(id)
	// read_count ++
	sql := "update post_list set read_count = read_count + 1 where id = ?"
	_, err := db.Exec(sql, id)
	if err != nil {
		log.Fatal(err)
	}

	// get post
	sql = "select * from post_view where id = ?"

	var res PostItem
	err = db.QueryRow(sql, id).Scan(&res.Id, &res.Content, &res.UserId, &res.UserName, &res.UserAvatar, &res.CreateTime, &res.LikeCount, &res.ReadCount, &res.CommentCount, &res.CategoryId, &res.CategoryName)
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

func getComment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postId := r.Form["postId"][0]

	// get post
	sql := "select * from comment_view where post_id = ?"
	rows, err := db.Query(sql, postId)
	if err != nil {
		log.Fatal(err)
	}
	var res []CommentItem
	for rows.Next() {
		// get comment
		var p CommentItem
		if err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.UserAvatar, &p.UserName, &p.Comment, &p.CreateTime, &p.LikeCount, &p.ReplyCount); err != nil {
			log.Fatal(err)
		}

		// get comment's reply
		sql := "select * from reply_view where comment_id = ?"
		rows2, err := db.Query(sql, p.Id)
		if err != nil {
			log.Fatal(err)
		}
		var reply ReplyList
		for rows2.Next() {
			var r ReplyItem
			if err := rows2.Scan(&r.Id, &r.CommentId, &r.CreateTime, &r.Content, &r.LikeCount, &r.FromUserId, &r.FromUserName, &r.FromUserAvatar, &r.ToUserId, &r.ToUserName, &r.ToUserAvatar); err != nil {
				log.Fatal(err)
			}
			reply.Reply = append(reply.Reply, r)
		}

		// encode json to string
		encodedReply, err := json.Marshal(reply.Reply)
		if err != nil {
			log.Fatal(err)
		}
		p.Reply = string(encodedReply)

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"code": 0,
		"data": res,
	}

	json.NewEncoder(w).Encode(data)
}
func postLike(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postId := r.Form["postId"][0]
	// get post
	sql := "update post_list set like_count = like_count + 1 where id = ?"
	_, err := db.Exec(sql, postId)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}

func commentLike(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	commentId := r.Form["id"][0]
	level2CommentId := r.Form["level2CommentId"][0]
	level := r.Form["level"][0]
	if level == "1" {
		sql := "update comment_list set like_count = like_count + 1 where id = ?"
		_, err := db.Exec(sql, commentId)
		if err != nil {
			log.Fatal(err)
		}
	}
	if level == "2" {
		sql := "update reply_list set like_count = like_count + 1 where id = ?"
		_, err := db.Exec(sql, level2CommentId)
		if err != nil {
			log.Fatal(err)
		}
	}
	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/sk/users/login", login)        // 登录
	http.HandleFunc("/sk/users/insert", insert)      // 新用户信息插入
	http.HandleFunc("/sk/users/getOne", getOne)      // 个人信息页
	http.HandleFunc("/sk/users/update", update)      // 老用户信息修改
	http.HandleFunc("/sk/post/get", get)             // 帖子列表获取
	http.HandleFunc("/sk/post/getOne", postGetOne)   //获取单个帖子
	http.HandleFunc("/sk/comment/get", getComment)   //获取评论列表
	http.HandleFunc("/sk/post/like", postLike)       //帖子点赞
	http.HandleFunc("/sk/comment/like", commentLike) //评论点赞

	err := http.ListenAndServe(":13001", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
