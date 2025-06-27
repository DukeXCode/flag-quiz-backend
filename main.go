package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "/app/data/database.db")
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

		var countries []map[string]any
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
			SelectedCountry int   `json:"selectedCountry" binding:"required"`
			CorrectCountry  int   `json:"correctCountry" binding:"required"`
			IsCorrect       *bool `json:"isCorrect" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := db.Exec("INSERT INTO answer (selected_country, correct_country, is_correct) VALUES (?, ?, ?)", json.SelectedCountry, json.CorrectCountry, json.IsCorrect)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, _ := res.LastInsertId()
		c.JSON(http.StatusCreated, gin.H{"id": id, "selectedCountry": json.SelectedCountry, "correctCountry": json.CorrectCountry, "isCorrect": json.IsCorrect})
	})

	r.GET("/answers/wrong/countries/:id", func(c *gin.Context) {
		id := c.Param("id")
		rows, err := db.Query("SELECT selected_country FROM answer WHERE correct_country = ? AND is_correct = FALSE", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		country_ids := []int{}
		for rows.Next() {
			var answer_id int
			if err := rows.Scan(&answer_id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			country_ids = append(country_ids, answer_id)
		}

		randomCountries := selectRandomCountries(country_ids)
		c.JSON(http.StatusOK, randomCountries)
	})
}

func selectRandomCountries(countries []int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedItems := []int{}
	usedIds := make(map[int]bool)
	var numItemsToSelect int

	numItemsToSelect = min(len(countries), 3)

	for len(selectedItems) < numItemsToSelect {
		randomIndex := r.Intn(len(countries))
		selectedId := countries[randomIndex]

		if _, ok := usedIds[selectedId]; !ok {
			selectedItems = append(selectedItems, selectedId)
			usedIds[selectedId] = true
		}
	}

	return selectedItems
}
