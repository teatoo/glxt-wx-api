package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/weigj/go-odbc/driver"
)

var port = ":8086" //端口
var queryHandler *sql.DB
var softvareVersion = "V2.01"

func main() {
	//读取本地局域网IP
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		temp, _ := inter.Addrs()
		for _, addr := range temp {
			if addr.String() != "0.0.0.0" {
				fmt.Println("Server local Ip address:" + addr.String() + port)
			}
		}
	}
	//建立数据库连接
	db, err := sql.Open("odbc", "DSN=glxt;uid=jsuser;pwd=jamesy;")
	if err != nil {
		fmt.Println("Connecting Error", err)
		return
	}
	//defer db.Close()

	queryHandler = db
	//defer stmt.Close()
	//切换运行模式到发布状态
	gin.SetMode(gin.ReleaseMode)
	//路由设置
	r := gin.Default()

	//	r.GET("/ping", func(c *gin.Context) {
	//		c.String(http.StatusOK, "pong")
	//	})
	r.Use(static.Serve("/", static.LocalFile("", false)))
	r.Use(static.Serve("/static", static.LocalFile("/static", false)))

	//	r.Static("/static", "./static")
	//	r.Static("/", "./static")
	//	r.Static("/files", "./files")
	//	r.GET("/favicon.ico", func(c *gin.Context) {
	//		c.File("favicon.ico")
	//	})
	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
		//		c.File("index.html?timestamp=" + time.Now().String())
	})
	r.GET("/periodInfo", QueryPeriod)
	//	r.POST("/file", Upload)
	//	r.DELETE("/file", Delete)
	//显示版本信息
	fmt.Println("继教学时查询微信服务程序", softvareVersion, "版本")
	//启动服务
	s := &http.Server{
		Addr:           port,
		Handler:        r,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

type ListPeriod struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Name    string `json:"name"`
	Credita int    `json:"credita"`
	Creditb int    `json:"creditb"`
}

func QueryPeriod(c *gin.Context) {
	var idCard string
	var inputTeacherName string
	var queryStartTime string
	var queryEndTime string

	idCard = strings.TrimSpace(c.Query("idCard"))
	inputTeacherName = strings.TrimSpace(c.Query("queryName"))

	queryStartTime = c.Query("startTime") + "-09-01"
	queryEndTime = c.Query("endtTime") + "-08-31"

	stmt, err := queryHandler.Prepare("SELECT c.starttime as starttime,c.endtime as endtime,c.ClassRoom as ClassRoom,p.credita as credita,p.creditb as creditb,p.TeacherName as TeacherName FROM info_period p INNER JOIN info_class c ON  p.ClassId = c.ClassId where  p.idcard=?  and  c.starttime >=? and starttime <? order by starttime")
	if err != nil {
		fmt.Println("Prepare Error", err)
		return
	}

	rows, err := stmt.Query(idCard, queryStartTime, queryEndTime)
	if err != nil {
		fmt.Println("Query Error", err)
		return
	}
	defer rows.Close()

	var starttime string
	var endtime string
	var ClassRoom string
	var credita int
	var creditb int
	var TeacherName string

	lp := make([]ListPeriod, 0)

	//	totalCredita = 0
	//	totalCreditb = 0
	//	totalTrainTimes = 0
	for rows.Next() {
		_ = rows.Scan(&starttime, &endtime, &ClassRoom, &credita, &creditb, &TeacherName)
		if strings.TrimSpace(TeacherName) == strings.TrimSpace(inputTeacherName) {
			fmt.Println(starttime, "|", endtime, "|", ClassRoom, "|", credita, "|", creditb, "|", TeacherName)
			//			totalCredita = totalCredita + credita
			//			totalCreditb = totalCreditb + creditb
			//			totalTrainTimes = totalTrainTimes + 1
			var p ListPeriod
			p.Start = starttime
			p.End = endtime
			p.Name = ClassRoom
			p.Credita = credita
			p.Creditb = creditb
			lp = append(lp, p)
		}
	}
	//	lp := make([]ListPeriod, 0)
	//	var p ListPeriod
	//	p.Start = "2012-10-10"
	//	p.End = "2012-10-11"
	//	p.Name = "测试班级"
	//	p.Credita = 10
	//	p.Creditb = 15
	//	lp = append(lp, p)
	//	var p2 ListPeriod
	//	p2.Start = "2013-10-10"
	//	p2.End = "2013-10-11"
	//	p2.Name = "测试班级2"
	//	p2.Credita = 110
	//	p2.Creditb = 125
	//	lp = append(lp, p2)
	c.JSON(http.StatusOK, lp)
	//	fmt.Println("共计培训",totalTrainTimes,"次,灵活性学时:",totalCreditb,",规范性学时:",totalCredita)
}
