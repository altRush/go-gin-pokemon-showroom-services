package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	Name   string `db:"name"`
	Url    string `db:"url"`
	Sprite string `db:"sprite"`
}

var stored_pokemon = Pokemon_profile_with_types{
	Name:   "test name",
	Url:    "test url",
	Sprite: "test sprite",
	Type: Type{
		Type_name: "test type name",
		Url:       "test url",
	},
	Pokemon_store_id: 0,
	Trainer_id:       "test trainer id",
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

	rows, err := db.Query("SELECT name, url, sprite FROM stored_pokemons")
	if err != nil {
		panic(err)
	}

	var (
		name   string
		url    string
		sprite string
	)

	for rows.Next() {
		err := rows.Scan(&name, &url, &sprite)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%#v\n", name, "%#v\n", url, "%#v\n", sprite)
	}

	rows.Close()

	db.Close()

	c.IndentedJSON(http.StatusOK, stored_pokemon)
}

func main() {
	router := gin.Default()
	router.GET("/store", getStoredPokemonById)

	router.Run("localhost:8080")
}
