package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.db"
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group("/rest/v1")
	routingV1(v1, db)
	v2 := r.Group("/rest/v2")
	routingV2(v2, db)
	r.Run(":8080")
}

func routingV1(r *gin.RouterGroup, db *sql.DB) {
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

func routingV2(r *gin.RouterGroup, db *sql.DB) {
	r.GET("/questions", func(c *gin.Context) {
		countStr := c.DefaultQuery("count", "10")
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 {
			count = 10
		}

		rows, err := db.Query("SELECT id, name, iso2 FROM country")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		countryMap := make(map[int]gin.H)
		var countries []gin.H
		for rows.Next() {
			var id int
			var name string
			var iso2 string
			if err := rows.Scan(&id, &name, &iso2); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c := gin.H{"id": id, "name": name, "iso2": iso2}
			countryMap[id] = c
			countries = append(countries, c)
		}

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		usedIds := make(map[int]bool)
		questions := []gin.H{}
		numQuestions := count
		if numQuestions > len(countries) {
			numQuestions = len(countries)
		}

		for len(questions) < numQuestions {
			ci := rng.Intn(len(countries))
			correctCountry := countries[ci]
			correctId := correctCountry["id"].(int)

			if usedIds[correctId] {
				continue
			}
			usedIds[correctId] = true

			histRows, err := db.Query("SELECT selected_country FROM answer WHERE correct_country = ? AND is_correct = FALSE", correctId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			histWrong := []int{}
			for histRows.Next() {
				var sc int
				if err := histRows.Scan(&sc); err != nil {
					histRows.Close()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				histWrong = append(histWrong, sc)
			}
			histRows.Close()

			if len(histWrong) > 1 {
				rng.Shuffle(len(histWrong), func(i, j int) {
					histWrong[i], histWrong[j] = histWrong[j], histWrong[i]
				})
			}

			options := []gin.H{correctCountry}
			histUsed := make(map[int]bool)
			histCount := 0
			for _, wid := range histWrong {
				if histCount >= 2 {
					break
				}
				if wid != correctId && !histUsed[wid] {
					options = append(options, countryMap[wid])
					histUsed[wid] = true
					histCount++
				}
			}

			otherIds := []int{}
			for _, country := range countries {
				id := country["id"].(int)
				if id != correctId && !histUsed[id] {
					otherIds = append(otherIds, id)
				}
			}
			rng.Shuffle(len(otherIds), func(i, j int) {
				otherIds[i], otherIds[j] = otherIds[j], otherIds[i]
			})

			for len(options) < 4 && len(otherIds) > 0 {
				nextId := otherIds[0]
				otherIds = otherIds[1:]
				options = append(options, countryMap[nextId])
			}

			rng.Shuffle(len(options), func(i, j int) {
				options[i], options[j] = options[j], options[i]
			})

			questions = append(questions, gin.H{
				"correctCountry": correctCountry,
				"options":        options,
			})
		}

		c.JSON(http.StatusOK, questions)
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
