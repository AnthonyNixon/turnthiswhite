package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func getCurrentNumberFromDB() (v int) {
	DB_USER := os.Getenv("TTSDBUSER")
	DB_PASS := os.Getenv("TTSDBPASS")
	DB_HOST := os.Getenv("DBHOST")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/turnthiswhite"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err.Error())
	}
	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	var currentValue = 0

	row := db.QueryRow("select value from variables where name = 'currentnum';")
	err = row.Scan(&currentValue)
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Printf("Starting at: %d\n", currentValue)
	return currentValue
}

var currentNumber = 0
var rates map[string]int
var RATE_LIMIT = 10

func startSync() {
	DB_USER := os.Getenv("TTSDBUSER")
	DB_PASS := os.Getenv("TTSDBPASS")
	DB_HOST := os.Getenv("DBHOST")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/turnthiswhite"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	for range time.Tick(time.Minute) {
		fmt.Println("Syncing DB...")
		stmt, err := db.Prepare("update `variables` SET value = ? WHERE name = 'currentnum';")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(currentNumber)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func rateClear() {
	for range time.Tick(time.Second) {
		fmt.Println("Clearing Rates...")
		rates = make(map[string]int)
	}
}

func main() {
	os.Setenv("APIPID", strconv.Itoa(os.Getpid()))
	fmt.Println(os.Getpid())

	file, err := os.Create("api.pid")
	if err != nil {
		log.Fatal("Cannot create pid file", err)
	}

	fmt.Fprintf(file, strconv.Itoa(os.Getpid()))

	file.Close()

	currentNumber = getCurrentNumberFromDB()
	go startSync()
	go rateClear()

	type Turnthiswhite struct {
		Number       int    `json:"number"`
		Color        string `json:"color"`
		InverseColor string `json:"inverseColor"`
	}

	router := gin.Default()
	// Add API handlers here

	// GET a cronjob
	router.GET("/color", func(c *gin.Context) {
		var turnthiswhite Turnthiswhite

		turnthiswhite.Number = currentNumber
		turnthiswhite.Color = strings.Replace(fmt.Sprintf("#%6x", currentNumber), " ", "0", -1)
		turnthiswhite.InverseColor = strings.Replace(fmt.Sprintf("#%6x", (16777215-currentNumber)), " ", "0", -1)

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, turnthiswhite)
		fmt.Println(rates)
	})

	router.PUT("/lighter", func(c *gin.Context) {
		var turnthiswhite Turnthiswhite
		var remoteAddr = c.Request.RemoteAddr

		rates[remoteAddr] += 1
		if rates[remoteAddr] <= RATE_LIMIT {
			currentNumber++
			if currentNumber > 16777215 {
				currentNumber = 16777215
			}
			turnthiswhite.Number = currentNumber
			turnthiswhite.Color = strings.Replace(fmt.Sprintf("#%6x", currentNumber), " ", "0", -1)
			turnthiswhite.InverseColor = strings.Replace(fmt.Sprintf("#%6x", (16777215-currentNumber)), " ", "0", -1)

			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusOK, turnthiswhite)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusTooManyRequests, turnthiswhite)
		}

	})

	router.PUT("/darker", func(c *gin.Context) {
		var turnthiswhite Turnthiswhite
		var remoteAddr = c.Request.RemoteAddr

		rates[remoteAddr] += 1
		if rates[remoteAddr] <= RATE_LIMIT {
			currentNumber--
			if currentNumber < 0 {
				currentNumber = 0
			}
			turnthiswhite.Number = currentNumber
			turnthiswhite.Color = strings.Replace(fmt.Sprintf("#%6x", currentNumber), " ", "0", -1)
			turnthiswhite.InverseColor = strings.Replace(fmt.Sprintf("#%6x", (16777215-currentNumber)), " ", "0", -1)

			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusOK, turnthiswhite)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusTooManyRequests, turnthiswhite)
		}
	})

	router.OPTIONS("/color", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,PUT")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.OPTIONS("/lighter", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,PUT")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.OPTIONS("/darker", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,PUT")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.Use(cors.Default())

	router.Run(":3000")
}
