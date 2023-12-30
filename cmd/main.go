package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
	CircleId     int    `json:"CategoryId"`
	CircleName   string `json:"categoryName"`
	IsLike       int    `json:"isLiked"`
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
	IsLike         int    `json:"isLiked"`
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
	IsLike     int    `json:"isLiked"`
}

type CommentList struct {
	Comment []CommentItem
}

type PostAdd struct {
	Content     string   `json:"content"`
	UserId      int      `json:"userId"`
	CircleId    int      `json:"categoryId"`
	IsPublished string   `json:"isPublished"`
	Urls        []string `json:"urls"`
}

type CategoryItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CircleItem struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Intro   string `json:"intro"`
	Avatar  string `json:"avatar"`
	Type    int    `json:"type"`
	IsAudit bool   `json:"isAudit"`
	Top     int    `json:"top"`
}

type HotCircleItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Intro    string `json:"intro"`
	Avatar   string `json:"avatar"`
	Type     int    `json:"type"`
	IsAudit  int    `json:"isAudit"`
	Top      int    `json:"top"`
	IsFollow int    `json:"isFollow"`
}

func DBinit() {
	var err error
	db, err = sql.Open("mysql", "HerSpace:88516098@/HerSpace")
	if err != nil {
		panic(err.Error())
	}
	// 读取queries.sql文件
	sqlFile, err := ioutil.ReadFile("../sql/init.sql")
	if err != nil {
		fmt.Println("无法读取SQL文件:", err)
		return
	}

	// 将SQL语句拆分为多个语句
	initSQL := strings.Split(string(sqlFile), ";")

	// 执行每个SQL语句
	for _, sql := range initSQL {
		if sql == "" {
			break
		}
		_, err := db.Exec(sql)
		if err != nil {
			fmt.Println("无法执行SQL语句:", err)
			return
		}
	}
	fmt.Println("DB init success!")
}

func init() {
	fmt.Println("init")
	DBinit()
}

func sql_user_info_query_by_wechatcode(WechatCode string) (ret UserInfo, ret_err error) {
	sql := "select * from user_info where wechat_code = '" + WechatCode + "'"
	var res UserInfo
	err := db.QueryRow(sql).Scan(&res.Id, &res.Name, &res.Avatar, &res.Intro, &res.WechatCode, &res.CreatedTime, &res.FollowedCount, &res.FollowingCount)
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
	err := db.QueryRow(sql, id).Scan(&res.Id, &res.Name, &res.Avatar, &res.Intro, &res.WechatCode, &res.CreatedTime, &res.FollowedCount, &res.FollowingCount)
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

func postGet(w http.ResponseWriter, r *http.Request) {
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
		err := rows.Scan(&p.Id, &p.Content, &p.UserId, &p.UserName, &p.UserAvatar, &p.CreateTime, &p.LikeCount, &p.ReadCount, &p.CommentCount, &p.CircleId, &p.CircleName)
		if err != nil {
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
	err = db.QueryRow(sql, id).Scan(&res.Id, &res.Content, &res.UserId, &res.UserName, &res.UserAvatar, &res.CreateTime, &res.LikeCount, &res.ReadCount, &res.CommentCount, &res.CircleId, &res.CircleName)
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

func postGetCate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sql := "select * from post_view where circle_id = ? ORDER BY create_time"
	rows, err := db.Query(sql, r.Form["id"][0])
	if err != nil {
		log.Fatal(err)
	}
	var res []PostItem
	for rows.Next() {
		var p PostItem
		err := rows.Scan(&p.Id, &p.Content, &p.UserId, &p.UserName, &p.UserAvatar, &p.CreateTime, &p.LikeCount, &p.ReadCount, &p.CommentCount, &p.CircleId, &p.CircleName)
		if err != nil {
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

func commentGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postId := r.Form["postId"][0]
	userId := r.Form["userId"][0]

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

			// get isLike
			sql = "SELECT EXISTS (SELECT * FROM reply_like_list WHERE from_user_id = ? AND reply_id = ?) AS result;"
			err = db.QueryRow(sql, userId, r.Id).Scan(&r.IsLike)
			if err != nil {
				log.Fatal(err)
			}

			reply.Reply = append(reply.Reply, r)
		}
		if reply.Reply != nil {
			// encode json to string
			encodedReply, err := json.Marshal(reply.Reply)
			if err != nil {
				log.Fatal(err)
			}
			p.Reply = string(encodedReply)
		}

		// get isLike
		sql = "SELECT EXISTS (SELECT * FROM comment_like_list WHERE from_user_id = ? AND comment_id = ?) AS result;"
		err = db.QueryRow(sql, userId, p.Id).Scan(&p.IsLike)
		if err != nil {
			log.Fatal(err)
		}

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
	userId := r.Form["userId"][0]
	status := r.Form["status"][0]
	postId := r.Form["postId"][0]
	if status == "1" {
		sql := "insert into post_like_list(from_user_id, post_id) values(?,?)"
		_, err := db.Exec(sql, userId, postId)
		if err != nil {
			log.Fatal(err)
		}
		sql = "update post_list set like_count = like_count + 1 where id = ?"
		_, err = db.Exec(sql, postId)
		if err != nil {
			log.Fatal(err)
		}
	}
	if status == "2" {
		sql := "delete from post_like_list where from_user_id = ? and post_id = ?"
		_, err := db.Exec(sql, userId, postId)
		if err != nil {
			log.Fatal(err)
		}
		sql = "update post_list set like_count = like_count - 1 where id = ?"
		_, err = db.Exec(sql, postId)
		if err != nil {
			log.Fatal(err)
		}
	}

	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}

func commentLike(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	isLike := r.Form["isLike"][0]
	commentId := r.Form["id"][0]
	level2CommentId := r.Form["level2CommentId"][0]
	level := r.Form["level"][0]
	userId := r.Form["userId"][0]
	if isLike == "1" {
		if level == "1" {
			sql := "insert into comment_like_list(comment_id, from_user_id) values(?, ?)"
			_, err := db.Exec(sql, commentId, userId)
			if err != nil {
				log.Fatal(err)
			}
			sql = "update comment_list set like_count = like_count + 1 where id = ?"
			_, err = db.Exec(sql, commentId)
			if err != nil {
				log.Fatal(err)
			}
		}
		if level == "2" {
			sql := "insert into reply_like_list(reply_id, from_user_id) values(?, ?)"
			_, err := db.Exec(sql, level2CommentId, userId)
			if err != nil {
				log.Fatal(err)
			}
			sql = "update reply_list set like_count = like_count + 1 where id = ?"
			_, err = db.Exec(sql, level2CommentId)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if isLike == "0" {
		if level == "1" {
			sql := "delete from comment_like_list where comment_id = ? and from_user_id = ?"
			_, err := db.Exec(sql, commentId, userId)
			if err != nil {
				log.Fatal(err)
			}
			sql = "update comment_list set like_count = like_count - 1 where id = ?"
			_, err = db.Exec(sql, commentId)
			if err != nil {
				log.Fatal(err)
			}
		}
		if level == "2" {
			sql := "delete from reply_like_list where reply_id = ? and from_user_id = ?"
			_, err := db.Exec(sql, level2CommentId, userId)
			if err != nil {
				log.Fatal(err)
			}
			sql = "update reply_list set like_count = like_count - 1 where id = ?"
			_, err = db.Exec(sql, level2CommentId)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}
func commentCo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postId := r.Form["postId"][0]
	content := r.Form["content"][0]
	fromUserId := r.Form["fromUserId"][0]
	toUserId := r.Form["toUserId"][0]
	//isAT := r.Form["isAT"][0]
	//atUserId := r.Form["atUserId"][0]
	createTime := time.Now()
	level := r.Form["level"][0]
	if level == "1" {
		sql := "insert into comment_list(post_id, user_id, comment, create_time, like_count) values(? ,? ,? ,?, 0)"
		_, err := db.Exec(sql, postId, fromUserId, content, createTime.Format("200601021504"))
		if err != nil {
			log.Fatal(err)
		}
	}
	if level == "2" {
		commentId := r.Form["commentId"][0]
		sql := "insert into reply_list(comment_id, from_user_id, content, create_time, to_user_id, like_count) values(?, ?, ?, ?, ?, 0)"
		_, err := db.Exec(sql, commentId, fromUserId, content, createTime.Format("200601021504"), toUserId)
		if err != nil {
			log.Fatal(err)
		}
	}
	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}

func postAdd(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res PostAdd
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	if res.CircleId == 0 {
		res.CircleId = 1
	}

	createTime := time.Now()

	sql := "insert into post_list(user_id, content, circle_id, create_time, like_count, read_count) values(?, ?, ?, ?, 0, 0)"
	_, err = db.Exec(sql, res.UserId, res.Content, res.CircleId, createTime.Format("200601021504"))
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	data := map[string]interface{}{
		"code": 0,
	}

	json.NewEncoder(w).Encode(data)
}
func getAllCategoryType() (ret []CategoryItem) {
	sql := "SELECT * FROM category_info"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	var res []CategoryItem
	for rows.Next() {
		var p CategoryItem
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, p)
	}
	return res
}
func circleGetCate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sql := "SELECT * FROM circle_list"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	var res []CircleItem
	for rows.Next() {
		var p CircleItem
		err := rows.Scan(&p.Id, &p.Name, &p.Intro, &p.Avatar, &p.Type)
		if err != nil {
			log.Fatal(err)
		}

		// get isAudit
		sql = "SELECT EXISTS (SELECT * FROM circle_audit_list WHERE user_id = ? AND circle_id = ?) AS result;"
		err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&p.IsAudit)
		if err != nil {
			log.Fatal(err)
		}

		// get isTop
		sql = "SELECT EXISTS (SELECT * FROM circle_top_list WHERE user_id = ? AND circle_id = ?) AS result;"
		var isTop bool
		err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isTop)
		if err != nil {
			log.Fatal(err)
		}

		if isTop {
			p.Top = 1
		} else {
			p.Top = 0
		}

		res = append(res, p)
	}
	fmt.Println(res)

	typeList := getAllCategoryType()
	fmt.Println(typeList)

	data := map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"list":     res,
			"typeList": typeList,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func categoryGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := getAllCategoryType()
	fmt.Println(res)

	data := map[string]interface{}{
		"code": 0,
		"data": res,
	}
	json.NewEncoder(w).Encode(data)
}

func circleGetHot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sql := "SELECT * FROM circle_list"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	var res []HotCircleItem
	for rows.Next() {
		var p HotCircleItem
		err := rows.Scan(&p.Id, &p.Name, &p.Intro, &p.Avatar, &p.Type)
		if err != nil {
			log.Fatal(err)
		}

		p.IsAudit = 1
		// get isFollow
		var isFollow bool
		sql = "SELECT EXISTS (SELECT * FROM circle_audit_list WHERE user_id = ? AND circle_id = ?) AS result;"
		err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isFollow)
		if err != nil {
			log.Fatal(err)
		}
		if isFollow {
			p.IsFollow = 1
		} else {
			p.IsFollow = 0
		}

		// get isTop
		sql = "SELECT EXISTS (SELECT * FROM circle_top_list WHERE user_id = ? AND circle_id = ?) AS result;"
		var isTop bool
		err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isTop)
		if err != nil {
			log.Fatal(err)
		}
		if isTop {
			p.Top = 1
		} else {
			p.Top = 0
		}

		res = append(res, p)
	}
	fmt.Println(res)

	typeList := getAllCategoryType()
	fmt.Println(typeList)

	data := map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"list":     res,
			"typeList": typeList,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func circleGetType(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	type_id := r.Form["type"][0]
	var sql string
	var res []HotCircleItem
	if type_id == "0" {
		sql = "SELECT * FROM circle_list"
		rows, err := db.Query(sql)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var p HotCircleItem
			err := rows.Scan(&p.Id, &p.Name, &p.Intro, &p.Avatar, &p.Type)
			if err != nil {
				log.Fatal(err)
			}

			p.IsAudit = 1
			// get isFollow
			var isFollow bool
			sql = "SELECT EXISTS (SELECT * FROM circle_audit_list WHERE user_id = ? AND circle_id = ?) AS result;"
			err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isFollow)
			if err != nil {
				log.Fatal(err)
			}
			if isFollow {
				p.IsFollow = 1
			} else {
				p.IsFollow = 0
			}

			// get isTop
			sql = "SELECT EXISTS (SELECT * FROM circle_top_list WHERE user_id = ? AND circle_id = ?) AS result;"
			var isTop bool
			err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isTop)
			if err != nil {
				log.Fatal(err)
			}

			if isTop {
				p.Top = 1
			} else {
				p.Top = 0
			}
			res = append(res, p)
		}
	} else {
		sql = "SELECT * FROM circle_list WHERE type_id = ?"
		rows, err := db.Query(sql, type_id)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			var p HotCircleItem
			err := rows.Scan(&p.Id, &p.Name, &p.Intro, &p.Avatar, &p.Type)
			if err != nil {
				log.Fatal(err)
			}

			p.IsAudit = 1
			// get isFollow
			var isFollow bool
			sql = "SELECT EXISTS (SELECT * FROM circle_audit_list WHERE user_id = ? AND circle_id = ?) AS result;"
			err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isFollow)
			if err != nil {
				log.Fatal(err)
			}
			if isFollow {
				p.IsFollow = 1
			} else {
				p.IsFollow = 0
			}

			// get isTop
			sql = "SELECT EXISTS (SELECT * FROM circle_top_list WHERE user_id = ? AND circle_id = ?) AS result;"
			var isTop bool
			err = db.QueryRow(sql, r.Form["userId"][0], p.Id).Scan(&isTop)
			if err != nil {
				log.Fatal(err)
			}

			if isTop {
				p.Top = 1
			} else {
				p.Top = 0
			}
			res = append(res, p)
		}
	}

	fmt.Println(res)
	typeList := getAllCategoryType()
	fmt.Println(typeList)
	data := map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"list":     res,
			"typeList": typeList,
		},
	}
	json.NewEncoder(w).Encode(data)
}

func circleAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sql := "insert into circle_list(name, intro, avatar, type_id) values(?, ?, ?, ?)"
	_, err := db.Exec(sql, r.Form["name"][0], r.Form["intro"][0], "undefined", r.Form["typeId"][0])
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{
		"code": 0,
	}
	json.NewEncoder(w).Encode(data)
}

func circleFollow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := r.Form["userId"][0]
	cateId := r.Form["caId"][0]
	status := r.Form["status"][0]
	if status == "1" {
		sql := "insert into circle_audit_list(user_id, circle_id) values(?, ?)"
		_, err := db.Exec(sql, userId, cateId)
		if err != nil {
			log.Fatal(err)
		}
	}
	if status == "2" {
		sql := "delete from circle_audit_list where user_id = ? and circle_id = ?"
		_, err := db.Exec(sql, userId, cateId)
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
	http.HandleFunc("/sk/users/login", login)           // 登录
	http.HandleFunc("/sk/users/insert", insert)         // 新用户信息插入
	http.HandleFunc("/sk/users/getOne", getOne)         // 个人信息页
	http.HandleFunc("/sk/users/update", update)         // 老用户信息修改
	http.HandleFunc("/sk/post/get", postGet)            // 帖子列表获取
	http.HandleFunc("/sk/post/getOne", postGetOne)      // 获取单个帖子
	http.HandleFunc("/sk/post/insert", postAdd)         // 添加帖子
	http.HandleFunc("/sk/post/cate", postGetCate)       // 获取圈子中的帖子
	http.HandleFunc("/sk/post/like", postLike)          // 帖子点赞
	http.HandleFunc("/sk/comment/get", commentGet)      // 获取评论列表
	http.HandleFunc("/sk/comment/like", commentLike)    // 评论点赞
	http.HandleFunc("/sk/comment/co", commentCo)        // 评论帖子
	http.HandleFunc("/sk/category/cate", circleGetCate) // 获取关注的圈子
	http.HandleFunc("/sk/category/hot", circleGetHot)   // 获取热门圈子
	http.HandleFunc("/sk/category/get", circleGetType)  // 获取指定类型圈子
	http.HandleFunc("/sk/category/insert", circleAdd)   // 添加圈子
	http.HandleFunc("/sk/category/type", categoryGet)   // 获取圈子类型
	http.HandleFunc("/sk/category/addCa", circleFollow) //关注/取消关注圈子

	err := http.ListenAndServe(":13001", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
