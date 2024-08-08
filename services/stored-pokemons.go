package store_pokemons

import (
	"database/sql"
	"fmt"
	"go-gin-pokemon-showroom-services/types"
	"go-gin-pokemon-showroom-services/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GetAllStoredPokemons(c *gin.Context) {
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

	var pokemons []types.Pokemon_profile_with_types

	for rows.Next() {
		var pkms types.Pokemon_profile_db
		err := rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.Pokemon_store_id, &pkms.Trainer_id)
		if err != nil {
			log.Fatalln(err)
		}
		var pkmnsWithTypes types.Pokemon_profile_with_types

		typesString := utils.ConvertDbArrayToUnnestArrayString(pkms.Types)

		unnestSql := fmt.Sprintf("select t.* from unnest(array[%s]) type_name_s left join types t on t.type_name = type_name_s", typesString)

		typesRows, err := db.Query(unnestSql)

		if err != nil {
			panic(err)
		}

		var pokemonTypes []types.PokemonTypes

		for typesRows.Next() {
			var t types.PokemonTypes
			err := typesRows.Scan(&t.Type_name, &t.Url)
			if err != nil {
				log.Fatalln(err)
			}

			pokemonTypes = append(pokemonTypes, t)
		}

		var trainer_id string
		if pkms.Trainer_id != nil {
			trainer_id = *pkms.Trainer_id
		}

		pkmnsWithTypes = types.Pokemon_profile_with_types{
			Name:             pkms.Name,
			Url:              pkms.Url,
			Sprite:           pkms.Sprite,
			Pokemon_store_id: pkms.Pokemon_store_id,
			Trainer_id:       trainer_id,
			Types:            pokemonTypes,
		}
		pokemons = append(pokemons, pkmnsWithTypes)

		typesRows.Close()
	}

	rows.Close()

	db.Close()

	c.JSON(http.StatusOK, pokemons)
}
