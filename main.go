package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Type struct {
	Type_name string `json:"type_name"`
	Url       string `json:"url"`
}

type Pokemon_profile_with_types struct {
	Name             string `json:"name"`
	Url              string `json:"url"`
	Sprite           string `json:"sprite"`
	Type             Type   `json:"type"`
	Pokemon_store_id int    `json:"pokemon_store_id"`
	Trainer_id       string `json:"trainer_id"`
}

type Pokemon_profile_db struct {
	Name             string         `db:"name"`
	Url              string         `db:"url"`
	Sprite           string         `db:"sprite"`
	Types            pq.StringArray `db:"types"`
	Pokemon_store_id string         `db:"pokemon_store_id"`
	Trainer_id       *string        `db:"trainer_id"`
}

func getStoredPokemonById(c *gin.Context) {
	godotenv.Load()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM stored_pokemons")
	if err != nil {
		panic(err)
	}

	var pokemons []Pokemon_profile_db

	for rows.Next() {
		var pkms Pokemon_profile_db
		err := rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.Pokemon_store_id, &pkms.Trainer_id)
		if err != nil {
			log.Fatalln(err)
		}
		pokemons = append(pokemons, pkms)
	}

	rows.Close()

	db.Close()

	c.IndentedJSON(http.StatusOK, pokemons)
}

func main() {
	router := gin.Default()
	router.GET("/store", getStoredPokemonById)

	router.Run("localhost:8080")
}
