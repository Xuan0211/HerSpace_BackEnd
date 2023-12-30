
CREATE TABLE IF NOT EXISTS category_info (
  id int PRIMARY KEY AUTO_INCREMENT, 
  name char(255) NOT NULL
);

INSERT INTO category_info (name)
SELECT '生活娱乐'
WHERE NOT EXISTS (
    SELECT 1 FROM category_info
);

Create table IF NOT EXISTS circle_list(
  id int AUTO_INCREMENT PRIMARY KEY,
  name char(255) NOT NULL,
  intro char(255) NOT NULL,
  avatar char(255) DEFAULT '' NOT NULL,
  type_id int NOT NULL,
  foreign key (type_id) references category_info(id)
);

INSERT INTO circle_list(name, intro, avatar, type_id)
SELECT '美好的一天', '爸爸妈妈去上班，我去幼儿园', 'undefined', 1
WHERE NOT EXISTS (
    SELECT 1 FROM circle_list
);

CREATE TABLE IF NOT EXISTS circle_audit_list(
  id int AUTO_INCREMENT PRIMARY KEY,
  circle_id int NOT NULL,
  user_id int NOT NULL,
  FOREIGN KEY (`circle_id`) REFERENCES `circle_list`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`)
);

CREATE TABLE IF NOT EXISTS circle_top_list(
  id int AUTO_INCREMENT PRIMARY KEY,
  circle_id int NOT NULL,
  user_id int NOT NULL,
  FOREIGN KEY (`circle_id`) REFERENCES `circle_list`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`)
);

CREATE TABLE IF NOT EXISTS `user_info` (
  `id` int PRIMARY KEY AUTO_INCREMENT, 
  `name` char(255) DEFAULT '' NOT NULL,
  `avatar` char(255) DEFAULT '' NOT NULL,
  `intro` char(255)  DEFAULT '' NOT NULL,
  `wechat_code` char(255) DEFAULT '' NOT NULL,
  `created_time` char(255) DEFAULT '' NOT NULL,
  `followed_count` int(11) DEFAULT 0 NOT NULL,
  `following_count` int(11) DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS `post_list` (
  `id` int AUTO_INCREMENT PRIMARY KEY, 
  `content` char(255) DEFAULT '' NOT NULL,
  `user_id` int NOT NULL,
  `create_time` char(255) DEFAULT '' NOT NULL,
  `like_count` int NOT NULL,
  `read_count` int NOT NULL,
  `circle_id` int NOT NULL,
   FOREIGN KEY (`user_id`) REFERENCES `user_info` (`id`),
   FOREIGN KEY (`circle_id`) REFERENCES `circle_list` (`id`)
);

create table IF NOT EXISTS comment_list(
	id int AUTO_INCREMENT PRIMARY KEY, 
	comment char(255) NOT NULL,
	post_id int NOT NULL,
	foreign key (post_id) references post_list(id),
	user_id int NOT NULL, 
	foreign key (user_id) references user_info(id),
	create_time char(255) NOT NULL,
	like_count int NOT NULL
);

create table IF NOT EXISTS reply_list(
	id int AUTO_INCREMENT PRIMARY KEY, 
	like_count int NOT NULL,
	create_time char(255) NOT NULL,
	comment_id int NOT NULL,
	content char(255) NOT NULL,
	from_user_id int NOT NULL, 
	to_user_id int NOT NULL,
	foreign key (from_user_id) references user_info(id),
	foreign key (to_user_id) references user_info(id),
	foreign key (comment_id) references comment_list(id)
);

CREATE OR REPLACE
VIEW comment_view
(id,post_id,user_id,user_avatar,user_name,comment,create_time,like_count,reply_count)
as
SELECT comment.id, comment.post_id, user.id, user.avatar, user.name,comment.comment,comment.create_time, comment.like_count, COALESCE(COUNT(reply.comment_id), 0)
FROM comment_list comment
JOIN user_info user ON comment.user_id = user.id
LEFT JOIN reply_list reply ON comment.id = reply.comment_id
GROUP BY comment.id;

CREATE OR REPLACE
VIEW reply_view
(id,comment_id,create_time,content,like_count,from_user_id,from_user_name,from_user_avatar,to_user_id,to_user_name,to_user_avatar)
as 
SELECT reply.id,reply.comment_id,reply.create_time, reply.content, reply.like_count, reply.from_user_id, from_user.name, from_user.avatar, reply.to_user_id, to_user.name, to_user.avatar
FROM user_info from_user, user_info to_user, reply_list reply
WHERE reply.to_user_id = to_user.id AND reply.from_user_id = from_user.id;

CREATE OR REPLACE 
VIEW post_view
(id,content,user_id,user_name,user_avatar,create_time,like_count,read_count,comment_count,circle_id,circle_name)
AS
SELECT post.id, post.content, post.user_id, user.name, user.avatar, post.create_time, post.like_count, post.read_count, COUNT(comment.id), circ.id, circ.name
FROM post_list post
JOIN user_info user ON post.user_id = user.id
JOIN circle_list circ ON circ.id = post.circle_id
LEFT JOIN comment_list comment ON comment.post_id = post.id
GROUP BY post.id;

create table IF NOT EXISTS post_like_list(
	id int AUTO_INCREMENT PRIMARY KEY, 
	post_id int NOT NULL,
	from_user_id int NOT NULL, 
	foreign key (post_id) references post_list(id),
	foreign key (from_user_id) references user_info(id)
);

create table IF NOT EXISTS comment_like_list(
	id int AUTO_INCREMENT PRIMARY KEY, 
	comment_id int NOT NULL,
	from_user_id int NOT NULL, 
	foreign key (comment_id) references comment_list(id),
	foreign key (from_user_id) references user_info(id)
);

create table IF NOT EXISTS reply_like_list(
	id int AUTO_INCREMENT PRIMARY KEY, 
	reply_id int NOT NULL,
	from_user_id int NOT NULL, 
	foreign key (reply_id) references reply_list(id),
	foreign key (from_user_id) references user_info(id)
);