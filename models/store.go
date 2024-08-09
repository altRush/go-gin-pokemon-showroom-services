package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type PokemonTypes struct {
	Type_name string `json:"type_name"`
	Url       string `json:"url"`
}

type Pokemon_profile_from_db_with_types struct {
	Name             string         `json:"name"`
	Url              string         `json:"url"`
	Sprite           string         `json:"sprite"`
	Types            []PokemonTypes `json:"types"`
	Pokemon_store_id int            `json:"pokemon_store_id"`
	Trainer_id       string         `json:"trainer_id"`
}

type Pokemon_profile_from_db struct {
	Name             string         `db:"name"`
	Url              string         `db:"url"`
	Sprite           string         `db:"sprite"`
	Types            pq.StringArray `db:"types"`
	Pokemon_store_id int            `db:"pokemon_store_id"`
	Trainer_id       *string        `db:"trainer_id"`
}

type Pokemon_profile struct {
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Sprite     string   `json:"sprite"`
	Types      []string `json:"types"`
	Trainer_id string   `json:"trainer_id"`
}

type Add_pokemon_to_store struct {
	Result string `json:"result"`
}

func AddPokemonToStore(c *gin.Context) (Add_pokemon_to_store, error) {
	var pokemonProfile Pokemon_profile
	c.BindJSON(&pokemonProfile)

	typesString := convertDbArrayToUnnestArrayString(pokemonProfile.Types)

	_, dbErr := db.Exec("INSERT INTO public.stored_pokemons (name, url, sprite, types, trainer_id) VALUES ($1, $2, $3, ARRAY["+typesString+"], $4 ::uuid)", pokemonProfile.Name, pokemonProfile.Url, pokemonProfile.Sprite, pokemonProfile.Trainer_id)

	if dbErr != nil {
		return Add_pokemon_to_store{}, dbErr
	}

	add_pokemon_to_store := Add_pokemon_to_store{Result: "Added to store"}

	return add_pokemon_to_store, nil
}

func GetPokemonByStoreIdFromStore(c *gin.Context) (Pokemon_profile_from_db_with_types, error) {
	pokemonStoreId := c.Param("pokemonStoreId")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	var pokemon Pokemon_profile_from_db

	row := db.QueryRow("SELECT * FROM stored_pokemons WHERE pokemon_store_id = $1", pokemonStoreId)

	scanErr := row.Scan(&pokemon.Name, &pokemon.Url, &pokemon.Sprite, &pokemon.Types, &pokemon.Pokemon_store_id, &pokemon.Trainer_id)

	if scanErr != nil {
		panic(scanErr)
	}

	var pkmnsWithTypes Pokemon_profile_from_db_with_types
	typesString := convertDbArrayToUnnestArrayString(pokemon.Types)

	unnestSql := fmt.Sprintf("select t.* from unnest(array[%s]) type_name_s left join types t on t.type_name = type_name_s", typesString)

	typesRows, err := db.Query(unnestSql)

	if err != nil {
		return Pokemon_profile_from_db_with_types{}, err
	}

	var pokemonTypes []PokemonTypes

	for typesRows.Next() {
		var t PokemonTypes
		err := typesRows.Scan(&t.Type_name, &t.Url)
		if err != nil {
			return Pokemon_profile_from_db_with_types{}, err
		}

		pokemonTypes = append(pokemonTypes, t)
	}

	var trainer_id string
	if pokemon.Trainer_id != nil {
		trainer_id = *pokemon.Trainer_id
	}

	pkmnsWithTypes = Pokemon_profile_from_db_with_types{
		Name:             pokemon.Name,
		Url:              pokemon.Url,
		Sprite:           pokemon.Sprite,
		Pokemon_store_id: pokemon.Pokemon_store_id,
		Trainer_id:       trainer_id,
		Types:            pokemonTypes,
	}

	typesRows.Close()

	return pkmnsWithTypes, nil

}

func GetAllStoredPokemons(c *gin.Context) ([]Pokemon_profile_from_db_with_types, error) {

	rows, err := db.Query("SELECT * FROM stored_pokemons")
	if err != nil {
		return []Pokemon_profile_from_db_with_types{}, err
	}

	var pokemons []Pokemon_profile_from_db_with_types
	var typesRows *sql.Rows

	for rows.Next() {
		var pkms Pokemon_profile_from_db
		err = rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.Pokemon_store_id, &pkms.Trainer_id)
		if err != nil {
			return []Pokemon_profile_from_db_with_types{}, err
		}
		var pkmnsWithTypes Pokemon_profile_from_db_with_types
		typesString := convertDbArrayToUnnestArrayString(pkms.Types)

		fmt.Println(typesString)

		unnestSql := fmt.Sprintf("select t.* from unnest(array[%s]) type_name_s left join types t on t.type_name = type_name_s", typesString)

		typesRows, err = db.Query(unnestSql)

		if err != nil {
			return []Pokemon_profile_from_db_with_types{}, err
		}

		var pokemonTypes []PokemonTypes

		for typesRows.Next() {
			var t PokemonTypes
			err = typesRows.Scan(&t.Type_name, &t.Url)
			if err != nil {
				return []Pokemon_profile_from_db_with_types{}, err
			}

			pokemonTypes = append(pokemonTypes, t)
		}

		var trainer_id string
		if pkms.Trainer_id != nil {
			trainer_id = *pkms.Trainer_id
		}

		pkmnsWithTypes = Pokemon_profile_from_db_with_types{
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

	return pokemons, nil
}
