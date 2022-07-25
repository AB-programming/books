package main

import (
	"Book/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Middle(confiem bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if confiem {
			cookie, err := c.Cookie("confirm")
			if err != nil {
				fmt.Println(err)
			}
			if cookie == "yes" {
				c.Next()
			} else {
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("html/login.tmpl", "html/main.tmpl", "html/lend.tmpl", "html/lendlist.tmpl",
		"html/search.tmpl", "html/add.tmpl", "html/error.tmpl", "html/admin.tmpl")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		admins := sql.FindAdmin()
		for _, value := range admins {
			if value.Username == username && value.Password == password {
				now := time.Now()
				time := now.Format("2006-01-02 15:04:05")
				c.SetCookie("confirm", "yes", 0, "/", "localhost", false, true)
				c.SetCookie("username", username, 0, "/", "localhost", false, true)
				sql.Addtime(username, time)
				c.Redirect(http.StatusMovedPermanently, "http://localhost:8082/main")
				return
			}
		}
		c.HTML(http.StatusOK, "error.tmpl", gin.H{})
	})

	r.GET("/main", Middle(true), func(c *gin.Context) {
		username, err := c.Cookie("username")
		if err != nil {
			fmt.Println(err)
		}
		admin := sql.FindOneAdmin(username)
		books := sql.Find()
		times := sql.Findtime(username)
		c.HTML(http.StatusOK, "main.tmpl", gin.H{
			"lists": books,
			"admin": admin,
			"times": times,
		})
	})
	r.GET("/lend", Middle(true), func(c *gin.Context) {
		c.HTML(http.StatusOK, "lend.tmpl", gin.H{})
	})

	r.GET("/lendlist", Middle(true), func(c *gin.Context) {
		results := sql.Findlend()
		c.HTML(http.StatusOK, "lendlist.tmpl", gin.H{
			"results": results,
		})
	})

	r.GET("/search", Middle(true), func(c *gin.Context) {
		results := sql.Find()
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"tables": results,
		})
	})

	r.POST("/result", Middle(true), func(c *gin.Context) {
		classify := c.PostForm("classify")
		mytype := c.PostForm("type")
		name := c.PostForm("name")
		results := sql.Findbook(classify, mytype, name)
		c.HTML(http.StatusOK, "search.tmpl", gin.H{
			"tables": results,
		})
	})

	r.GET("/add", Middle(true), func(c *gin.Context) {
		c.HTML(http.StatusOK, "add.tmpl", gin.H{})
	})

	r.POST("/addresult", Middle(true), func(c *gin.Context) {
		bookname := c.PostForm("bookname")
		index := c.PostForm("index")
		mytype := c.PostForm("type")
		classify := c.PostForm("classify")
		num, err := strconv.ParseInt(c.PostForm("num"), 10, 32)
		if err != nil {
			panic(err)
		}
		sql.Addbook(bookname, index, classify, mytype, int32(num))
	})

	r.POST("/insert", Middle(true), func(c *gin.Context) {
		lend := make(map[string]string)
		lend["bookname"] = c.PostForm("BookName")
		lend["index"] = c.PostForm("Index")
		lend["name"] = c.PostForm("Name")
		lend["studentid"] = c.PostForm("StudentId")
		lend["academy"] = c.PostForm("Academy")
		lend["startdate"] = c.PostForm("StartDate")
		lend["returndate"] = c.PostForm("ReturnDate")
		lend["status"] = c.PostForm("Status")
		lend["seashell"] = c.PostForm("Seashell")

		//调用sql包下的Addlend方法增加一条借阅记录
		sql.Addlend(lend)
		//调用sql包下的Updatebook的方法同步更新books集合中书的所剩数量
		sql.Updatebook(lend["index"])
	})

	r.GET("/delete", Middle(true), func(c *gin.Context) {
		id := c.Query("id")
		fmt.Println("id是" + id)
		sql.Deletebook(id)
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8082/main")
	})

	r.GET("/return", func(c *gin.Context) {
		studentid := c.Query("studentid")
		index := c.Query("index")
		sql.ReturnId(studentid, index)
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8082/lendlist/")
	})

	r.GET("/exit", Middle(true), func(c *gin.Context) {
		c.SetCookie("confirm", "yes", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8082/")
	})

	r.GET("/admin", Middle(true), func(c *gin.Context) {
		username, err := c.Cookie("username")
		if err != nil {
			fmt.Println(err)
		}
		admin := sql.FindOneAdmin(username)
		c.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"admin": admin,
		})
	})

	r.Run(":8082") // 监听并在 0.0.0.0:8081 上启动服务
}
