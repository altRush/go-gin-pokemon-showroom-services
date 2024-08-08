package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Types struct {
	Type_name string `json:"type_name"`
	Url       string `json:"url"`
}

type Pokemon_profile_with_types struct {
	Name             string  `json:"name"`
	Url              string  `json:"url"`
	Sprite           string  `json:"sprite"`
	Types            []Types `json:"types"`
	Pokemon_store_id int     `json:"pokemon_store_id"`
	Trainer_id       string  `json:"trainer_id"`
}

type Pokemon_profile_db struct {
	Name             string         `db:"name"`
	Url              string         `db:"url"`
	Sprite           string         `db:"sprite"`
	Types            pq.StringArray `db:"types"`
	Pokemon_store_id int            `db:"pokemon_store_id"`
	Trainer_id       *string        `db:"trainer_id"`
}

func getAllStoredPokemons(c *gin.Context) {
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

	var pokemons []Pokemon_profile_with_types

	for rows.Next() {
		var pkms Pokemon_profile_db
		err := rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.Pokemon_store_id, &pkms.Trainer_id)
		if err != nil {
			log.Fatalln(err)
		}
		var pkmnsWithTypes Pokemon_profile_with_types

		typesJson, err := json.Marshal(pkms.Types)

		if err != nil {
			log.Fatalln(err)
		}

		foo := string(typesJson)
		bar := strings.Trim(foo, "[]")
		baz := strings.Replace(bar, "\"", "'", 4)

		unnestSql := fmt.Sprintf("select t.* from unnest(array[%s]) type_name_s left join types t on t.type_name = type_name_s", baz)

		typesRows, err := db.Query(unnestSql)

		if err != nil {
			panic(err)
		}

		var types []Types

		for typesRows.Next() {
			var t Types
			err := typesRows.Scan(&t.Type_name, &t.Url)
			if err != nil {
				log.Fatalln(err)
			}

			types = append(types, t)
		}

		var trainer_id string
		if pkms.Trainer_id != nil {
			trainer_id = *pkms.Trainer_id
		}

		pkmnsWithTypes = Pokemon_profile_with_types{
			Name:             pkms.Name,
			Url:              pkms.Url,
			Sprite:           pkms.Sprite,
			Pokemon_store_id: pkms.Pokemon_store_id,
			Trainer_id:       trainer_id,
			Types:            types,
		}
		pokemons = append(pokemons, pkmnsWithTypes)

		typesRows.Close()
	}

	rows.Close()

	db.Close()

	c.IndentedJSON(http.StatusOK, pokemons)
}

func main() {
	router := gin.Default()
	router.GET("/store/all", getAllStoredPokemons)

	router.Run("localhost:8080")
}
