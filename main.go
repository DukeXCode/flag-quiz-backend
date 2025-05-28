package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/rest/v1")
	routing(v1, db)
	r.Run(":8080")
}

func routing(r *gin.RouterGroup, db *sql.DB) {
	r.GET("/countries", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, iso2 FROM country")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var countries []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			var iso2 string
			if err := rows.Scan(&id, &name, &iso2); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			countries = append(countries, gin.H{"id": id, "name": name, "iso2": iso2})
		}
		c.JSON(http.StatusOK, countries)
	})

	r.POST("/answers", func(c *gin.Context) {
		var json struct {
			SelectedCountry int  `json:"selectedCountry" binding:"required"`
			CorrectCountry  int  `json:"correctCountry" binding:"required"`
			IsCorrect       bool `json:"isCorrect" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("%#v", json)
		res, err := db.Exec("INSERT INTO answer (selected_country, correct_country, is_correct) VALUES (?, ?, ?)", json.SelectedCountry, json.CorrectCountry, json.IsCorrect)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, _ := res.LastInsertId()
		c.JSON(http.StatusCreated, gin.H{"id": id, "selectedCountry": json.SelectedCountry, "correctCountry": json.CorrectCountry, "isCorrect": json.IsCorrect})
	})
}
