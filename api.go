package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	DB_USER := os.Getenv("CRONDBUSER")
	DB_PASS := os.Getenv("CRONDBPASS")
	DB_HOST := os.Getenv("DBHOST")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/cronager?parseTime=true"

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
	type Cronjob struct {
		Color string `json:"color"`
	}

	router := gin.Default()
	// Add API handlers here

	// GET a cronjob
	router.GET("/color", func(c *gin.Context) {
		var (
			cronjob Cronjob
			result  gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, name, crondef, command, description, active, logtime from jobs where id = ?;", id)
		err = row.Scan(&cronjob.Id, &cronjob.Name, &cronjob.Cron_def, &cronjob.Command, &cronjob.Description, &cronjob.Active, &cronjob.Logtime)
		if err != nil {
			// if no results, send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": cronjob,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	router.PUT("/job", func(c *gin.Context) {
		var cronjob Cronjob
		c.BindJSON(&cronjob)

		stmt, err := db.Prepare("update jobs set name = ?, crondef = ?, command = ?, description = ?, active = ? where id = ?;")

		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(cronjob.Name, cronjob.Cron_def, cronjob.Command, cronjob.Description, cronjob.Active, cronjob.Id)
		if err != nil {
			fmt.Print(err.Error())
		}

		defer stmt.Close()
		c.JSON(http.StatusOK, cronjob)
	})

	router.OPTIONS("/job", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.OPTIONS("/jobs", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.Use(cors.Default())

	router.Run(":3000")
}
