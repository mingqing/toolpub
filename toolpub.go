package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
)

type Login struct {
	User   string `form:"user" json:"user" binding:"required"`
	Passwd string `form:"password" json:"password" binding:"required"`
}

func main() {
	//r := gin.Default()
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Static("/tmp", "/tmp")

	//r.LoadHTMLGlob("templates/*")
	tmpl := template.New("tmpl").Delims("{{%", "%}}")
	tmpl.ParseFiles(
		"templates/base-header.html",
		"templates/base-nav.html",
		"templates/base-footer.html")
	_, err := tmpl.ParseGlob("templates/*/*")
	if err != nil {
		fmt.Println("new template err:", err)
	}
	r.SetHTMLTemplate(tmpl)

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user/index.html", gin.H{
			"title": "Main Title",
		})
	})

	r.GET("/help/ie", BrowserUpgrade)
	r.GET("/insert", InsertDatas)

	v1 := r.Group("/v1")
	{
		v1.GET("/test1", V1test1)
		v1.GET("/test2", V1test2)
	}

	r.Run()
}

func InsertDatas(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func V1test1(c *gin.Context) {
	influx, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.5.17.224:8086",
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		fmt.Println("err:", err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "testdb",
		Precision: "s",
	})
	if err != nil {
		fmt.Println("err:", err)
	}

	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10,
		"system": 53,
		"user":   46,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		fmt.Println("err:", err)
	}

	bp.AddPoint(pt)

	influx.Write(bp)

	writePoints(influx)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
func V1test2(c *gin.Context) {
	influx, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.5.17.224:8086",
		Username: "admin",
		Password: "admin",
	})

	if err != nil {
		fmt.Println("err:", err)
	}

	writePoints(influx)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
func BrowserUpgrade(c *gin.Context) {
	c.HTML(http.StatusOK, "help/ie-upgrade.html", gin.H{})
}

func writePoints(clnt client.Client) {
	sampleSize := 10
	rand.Seed(42)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "systemstats",
		Precision: "us",
	})

	for i := 0; i < sampleSize; i++ {
		regions := []string{"us-west1", "us-west2", "us-west3", "us-east1"}
		tags := map[string]string{
			"cpu":    "cpu-total",
			"host":   fmt.Sprintf("host%d", rand.Intn(1000)),
			"region": regions[rand.Intn(len(regions))],
		}

		idle := rand.Float64() * 100
		fields := map[string]interface{}{
			"idle": idle,
			"busy": 100 - idle,
		}

		pt, err := client.NewPoint(
			"cpu_usage",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			fmt.Println("err:", err)
		}

		bp.AddPoint(pt)
	}

	err := clnt.Write(bp)
	if err != nil {
		fmt.Println(err)
	}
}
