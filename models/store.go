package models

import (
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type PokemonType struct {
	TypeName string `json:"type_name"`
	Url      string `json:"url"`
}

type PokemonProfileFromDBWithTypes struct {
	Name           string        `json:"name"`
	Url            string        `json:"url"`
	Sprite         string        `json:"sprite"`
	Types          []PokemonType `json:"types"`
	PokemonStoreID int           `json:"pokemon_store_id"`
	TrainerID      string        `json:"trainer_id"`
}

type PokemonProfileFromDB struct {
	Name           string         `db:"name"`
	Url            string         `db:"url"`
	Sprite         string         `db:"sprite"`
	Types          pq.StringArray `db:"types"`
	PokemonStoreID int            `db:"pokemon_store_id"`
	TrainerID      *string        `db:"trainer_id"`
}

type PokemonProfile struct {
	Name      string   `json:"name"`
	Url       string   `json:"url"`
	Sprite    string   `json:"sprite"`
	Types     []string `json:"types"`
	TrainerID string   `json:"trainer_id"`
}

type AddResult struct {
	Result string `json:"result"`
}

func AddPokemonToStore(c *gin.Context) (AddResult, error) {
	var pokemonProfile PokemonProfile
	if err := c.BindJSON(&pokemonProfile); err != nil {
		return AddResult{}, err
	}

	_, dbErr := db.Exec(
		"INSERT INTO public.stored_pokemons (name, url, sprite, types, trainer_id) VALUES ($1, $2, $3, $4, $5 ::uuid)",
		pokemonProfile.Name, pokemonProfile.Url, pokemonProfile.Sprite, pq.Array(pokemonProfile.Types), pokemonProfile.TrainerID,
	)
	if dbErr != nil {
		return AddResult{}, dbErr
	}

	return AddResult{Result: "Added to store"}, nil
}

func GetPokemonByStoreIdFromStore(c *gin.Context) (PokemonProfileFromDBWithTypes, error) {
	pokemonStoreID := c.Param("pokemonStoreId")

	var pokemon PokemonProfileFromDB

	row := db.QueryRow(
		"SELECT name, url, sprite, types, pokemon_store_id, trainer_id FROM stored_pokemons WHERE pokemon_store_id = $1",
		pokemonStoreID,
	)
	scanErr := row.Scan(&pokemon.Name, &pokemon.Url, &pokemon.Sprite, &pokemon.Types, &pokemon.PokemonStoreID, &pokemon.TrainerID)
	if scanErr != nil {
		return PokemonProfileFromDBWithTypes{}, scanErr
	}

	unnestSql := "select t.* from unnest($1::text[]) type_name_s left join types t on t.type_name = type_name_s"
	typesRows, err := db.Query(unnestSql, pq.Array(pokemon.Types))
	if err != nil {
		return PokemonProfileFromDBWithTypes{}, err
	}
	defer typesRows.Close()

	var pokemonTypes []PokemonType
	for typesRows.Next() {
		var t PokemonType
		err := typesRows.Scan(&t.TypeName, &t.Url)
		if err != nil {
			return PokemonProfileFromDBWithTypes{}, err
		}
		pokemonTypes = append(pokemonTypes, t)
	}

	var trainerID string
	if pokemon.TrainerID != nil {
		trainerID = *pokemon.TrainerID
	}

	pkmnsWithTypes := PokemonProfileFromDBWithTypes{
		Name:           pokemon.Name,
		Url:            pokemon.Url,
		Sprite:         pokemon.Sprite,
		PokemonStoreID: pokemon.PokemonStoreID,
		TrainerID:      trainerID,
		Types:          pokemonTypes,
	}

	return pkmnsWithTypes, nil
}

func GetAllStoredPokemons(c *gin.Context) ([]PokemonProfileFromDBWithTypes, error) {
	rows, err := db.Query("SELECT name, url, sprite, types, pokemon_store_id, trainer_id FROM stored_pokemons")
	if err != nil {
		return []PokemonProfileFromDBWithTypes{}, err
	}
	defer rows.Close()

	var pokemons []PokemonProfileFromDBWithTypes

	for rows.Next() {
		var pkms PokemonProfileFromDB
		err = rows.Scan(&pkms.Name, &pkms.Url, &pkms.Sprite, &pkms.Types, &pkms.PokemonStoreID, &pkms.TrainerID)
		if err != nil {
			return []PokemonProfileFromDBWithTypes{}, err
		}

		unnestSql := "select t.* from unnest($1::text[]) type_name_s left join types t on t.type_name = type_name_s"
		typesRows, err := db.Query(unnestSql, pq.Array(pkms.Types))
		if err != nil {
			return []PokemonProfileFromDBWithTypes{}, err
		}

		var pokemonTypes []PokemonType
		for typesRows.Next() {
			var t PokemonType
			err = typesRows.Scan(&t.TypeName, &t.Url)
			if err != nil {
				typesRows.Close()
				return []PokemonProfileFromDBWithTypes{}, err
			}
			pokemonTypes = append(pokemonTypes, t)
		}
		typesRows.Close()

		var trainerID string
		if pkms.TrainerID != nil {
			trainerID = *pkms.TrainerID
		}

		pkmnsWithTypes := PokemonProfileFromDBWithTypes{
			Name:           pkms.Name,
			Url:            pkms.Url,
			Sprite:         pkms.Sprite,
			PokemonStoreID: pkms.PokemonStoreID,
			TrainerID:      trainerID,
			Types:          pokemonTypes,
		}
		pokemons = append(pokemons, pkmnsWithTypes)
	}

	return pokemons, nil
}
