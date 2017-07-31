package main

import (
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"strings"
	//"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//DB_USER := os.Getenv("CRONDBUSER")
	//DB_PASS := os.Getenv("CRONDBPASS")
	//DB_HOST := os.Getenv("DBHOST")
	//
	//dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/cronager?parseTime=true"
	//
	//db, err := sql.Open("mysql", dsn)
	//if err != nil {
	//	fmt.Print(err.Error())
	//}
	//defer db.Close()
	//// make sure our connection is available
	//err = db.Ping()
	//if err != nil {
	//	fmt.Print(err.Error())
	//}
	type Turnthiswhite struct {
		Number int    `json:"number"`
		Color  string `json:"color"`
	}

	currentNumber := 0

	router := gin.Default()
	// Add API handlers here

	// GET a cronjob
	router.GET("/color", func(c *gin.Context) {
		var turnthiswhite Turnthiswhite

		currentNumber++
		turnthiswhite.Number = currentNumber
		turnthiswhite.Color = strings.Replace(fmt.Sprintf("#%-6x", currentNumber), " ", "0", -1)

		c.JSON(http.StatusOK, turnthiswhite)
	})

	router.PUT("/color", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{})
	})

	router.OPTIONS("/color", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.Use(cors.Default())

	router.Run(":3000")
}
