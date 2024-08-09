package store_pokemons

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/altRush/go-gin-pokemon-showroom-services/types"
	"github.com/altRush/go-gin-pokemon-showroom-services/utils"

	"github.com/gin-gonic/gin"
)

func AddPokemonToStore(c *gin.Context) {
	var pokemonProfile types.Pokemon_profile
	c.BindJSON(&pokemonProfile)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	typesString := utils.ConvertDbArrayToUnnestArrayString(pokemonProfile.Types)

	_, dbErr := db.Exec("INSERT INTO public.stored_pokemons (name, url, sprite, types, trainer_id) VALUES ($1, $2, $3, ARRAY[$4], $5 ::uuid)", pokemonProfile.Name, pokemonProfile.Url, pokemonProfile.Sprite, typesString, pokemonProfile.Trainer_id)

	if dbErr != nil {
		panic(dbErr)
	}

	add_pokemon_to_store := types.Add_pokemon_to_store{Result: "Added to store"}
	c.JSON(http.StatusCreated, add_pokemon_to_store)
}

func GetPokemonByStoreIdFromStore(c *gin.Context) {
	pokemonStoreId := c.Param("pokemonStoreId")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	var pokemon types.Pokemon_profile_from_db

	row := db.QueryRow("SELECT * FROM stored_pokemons WHERE pokemon_store_id = $1", pokemonStoreId)

	scanErr := row.Scan(&pokemon.Name, &pokemon.Url, &pokemon.Sprite, &pokemon.Types, &pokemon.Pokemon_store_id, &pokemon.Trainer_id)

	if scanErr != nil {
		panic(scanErr)
	}

	var pkmnsWithTypes types.Pokemon_profile_from_db_with_types
	typesString := utils.ConvertDbArrayToUnnestArrayString(pokemon.Types)

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
	if pokemon.Trainer_id != nil {
		trainer_id = *pokemon.Trainer_id
	}

	pkmnsWithTypes = types.Pokemon_profile_from_db_with_types{
		Name:             pokemon.Name,
		Url:              pokemon.Url,
		Sprite:           pokemon.Sprite,
		Pokemon_store_id: pokemon.Pokemon_store_id,
		Trainer_id:       trainer_id,
		Types:            pokemonTypes,
	}

	typesRows.Close()

	c.JSON(http.StatusOK, pkmnsWithTypes)

}

func GetAllStoredPokemons(c *gin.Context) {
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

	var pokemons []types.Pokemon_profile_from_db_with_types
	for rows.Next() {
		var pkms types.Pokemon_profile_from_db
		err := rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.Pokemon_store_id, &pkms.Trainer_id)
		if err != nil {
			log.Fatalln(err)
		}
		var pkmnsWithTypes types.Pokemon_profile_from_db_with_types
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

		pkmnsWithTypes = types.Pokemon_profile_from_db_with_types{
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
