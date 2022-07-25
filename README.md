# 一个简单的图书管理

## 	1、简介

### 		此项目是基于Gin的一个简单的图书管理，前后端不分离，功能简单，可以克隆过去体验，缺点是代码写的很乱，并且冗		余，重复的没有优化，当作练手的小项目了

## 	2、技术选型

### 		前端：LayUi框架套壳

### 		后端：Golang、Gin

### 		数据库：mongodb

## 	3、启动方式

### 		首先，需要golang的环境，需要mongodb数据库支持

### 		其次，需要更改建立数据库、集合，有mongodb之后可以在navicat或者DataGrip里执行以下sql语句

```sql
use Book;
db.createCollection("books")
db.books.insertOne({
    name:"时间简史",
    index:"0A340",
    classify:"自然类基本学科",
    type:"图书",
    num:12
})
db.createCollection("lends")
db.lends.insertOne({
    bookname: "时间简史",
    index: "0A340",
    name:"张三",
    studentid:"2452652",
    academy:"航空学院",
    startDate:"2022-07-01",
    returnDate:"2022-07-08",
    status:"学生",
    seashell:"无"
})
db.createCollection("admins")
db.admins.insertOne({
    username:"101",
    password:"123456",
    nickname:"图书管理人员",
    sex:"男",
    age:"22",
})
db.createCollection("times")
db.times.insertOne({
    username:"101",
    time:"2022-07-25 00:20:14"
})
```

### 		最后进入main.go所在的文件夹下执行以下命令

​		

```go
x:\xxx\xxx>go build main.go

x:\xxx\xxx>main.exe
```

### 最后打开浏览器访问http://localhost:8082/，没有设置注册功能，（账号：101，密码：123456），可以自己在Book数据库的admin集合中添加，输入后登录即可
