package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Data struct {
	Images   []Images `json:"images"`
	Tooltips Tooltips `json:"tooltips"`
}
type Images struct {
	Urlbase       string   `json:"urlbase"`
	Enddate       string   `json:"enddate"`
	Startdate     string   `json:"startdate"`
	Quiz          string   `json:"quiz"`
	Top           int      `json:"top"`
	Fullstartdate string   `json:"fullstartdate"`
	Copyrightlink string   `json:"copyrightlink"`
	URL           string   `json:"url"`
	Hs            []string `json:"hs"`
	Title         string   `json:"title"`
	Drk           int      `json:"drk"`
	Wp            bool     `json:"wp"`
	Bot           int      `json:"bot"`
	Copyright     string   `json:"copyright"`
	Hsh           string   `json:"hsh"`
}
type Tooltips struct {
	Previous string `json:"previous"`
	Loading  string `json:"loading"`
	Next     string `json:"next"`
	Walle    string `json:"walle"`
	Walls    string `json:"walls"`
}

//CLI帮助
func inCliValue() string {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p",
			Value: "8084",
			Usage: "server port",
		},
	}
	p := ""
	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			fmt.Println("您想做什么？")
		}
		if c.String("p") != "" {
			p = c.String("p")
		} else {
			p = "8084"
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if p == "" {
		os.Exit(0)
	}
	return p
}

func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func GetBingImgData(url string) (Data, error) {
	response, err := http.Get(url)
	if err != nil {
		return Data{}, err
	} else if response.StatusCode != http.StatusOK {
		return Data{}, errors.New(fmt.Sprint(http.StatusServiceUnavailable))
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	previous := gjson.Get(string(body), "tooltips.previous").Str
	loading := gjson.Get(string(body), "tooltips.loading").Str
	next := gjson.Get(string(body), "tooltips.next").Str
	walle := gjson.Get(string(body), "tooltips.walle").Str
	walls := gjson.Get(string(body), "tooltips.walls").Str

	result := gjson.Get(string(body), "images")
	var d Data
	result.ForEach(
		func(key, value gjson.Result) bool {
			d.Images = append(d.Images, Images{
				Urlbase:       gjson.Get(value.String(), "urlbase").Str,
				Enddate:       gjson.Get(value.String(), "enddate").Str,
				Startdate:     gjson.Get(value.String(), "startdate").Str,
				Quiz:          gjson.Get(value.String(), "quiz").Str,
				Top:           int(gjson.Get(value.String(), "top").Int()),
				Fullstartdate: gjson.Get(value.String(), "fullstartdate").Str,
				Copyrightlink: gjson.Get(value.String(), "copyrightlink").Str,
				URL:           gjson.Get(value.String(), "url").Str,
				Hs:            []string{gjson.Get(value.String(), "hs").Str},
				Title:         gjson.Get(value.String(), "title").Str,
				Drk:           int(gjson.Get(value.String(), "drk").Int()),
				Wp:            gjson.Get(value.String(), "wp").Bool(),
				Bot:           int(gjson.Get(value.String(), "bot").Int()),
				Copyright:     gjson.Get(value.String(), "copyright").Str,
				Hsh:           gjson.Get(value.String(), "Hsh").Str,
			})
			return true
			// 继续迭代
		})

	data := Data{
		Images: d.Images,
		Tooltips: Tooltips{
			Previous: previous,
			Loading:  loading,
			Next:     next,
			Walle:    walle,
			Walls:    walls,
		},
	}
	return data, nil
}

func APIRoot(c *gin.Context) {
	resultOne, err := GetBingImgData("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=8")
	if err != nil {
		log.Println(err)
		c.Status(http.StatusServiceUnavailable)
	}
	resultTwo, err := GetBingImgData("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=7&n=8")
	if err != nil {
		log.Println(err)
		c.Status(http.StatusServiceUnavailable)
	}
	for i := 1; i <= len(resultTwo.Images)-1; i++ {
		resultOne.Images = append(resultOne.Images, resultTwo.Images[i])
	}
	for i := 0; i <= len(resultOne.Images)-1; i++ {
		resultOne.Images[i].Urlbase = "//cn.bing.com" + resultOne.Images[i].Urlbase
		resultOne.Images[i].Quiz = "//cn.bing.com" + resultOne.Images[i].Quiz
		resultOne.Images[i].URL = "//cn.bing.com" + resultOne.Images[i].URL
	}
	c.JSON(http.StatusOK, resultOne)
}

func main() {
	port := inCliValue()
	fmt.Println("Running at port " + port + " ...")
	fmt.Println("Other options input -h")

	lmt := tollbooth.NewLimiter(30, nil)
	// Set a custom message.
	lmt.SetMessage("请求过于频繁，请稍后重试")

	gin.SetMode(gin.ReleaseMode) //启动生产环境
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	// 记录日志到文件。
	logFile, err := os.OpenFile("gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Println("日志写入错误", err)
	}
	//及时关闭file句柄
    defer logFile.Close()
	gin.DefaultWriter = io.MultiWriter(logFile)
	router := gin.Default()
	router.Use(Cors())
	router.Use(LimitHandler(lmt))

	// API分组(RESTFULL)以及版本控制
	API := router.Group("/")
	{
		API.GET("/", APIRoot)
		API.POST("/", APIRoot)
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		// 服务连接
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("listen:", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
