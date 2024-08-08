package types

import "github.com/lib/pq"

type PokemonTypes struct {
	Type_name string `json:"type_name"`
	Url       string `json:"url"`
}

type Pokemon_profile_with_types struct {
	Name             string         `json:"name"`
	Url              string         `json:"url"`
	Sprite           string         `json:"sprite"`
	Types            []PokemonTypes `json:"types"`
	Pokemon_store_id int            `json:"pokemon_store_id"`
	Trainer_id       string         `json:"trainer_id"`
}

type Pokemon_profile_db struct {
	Name             string         `db:"name"`
	Url              string         `db:"url"`
	Sprite           string         `db:"sprite"`
	Types            pq.StringArray `db:"types"`
	Pokemon_store_id int            `db:"pokemon_store_id"`
	Trainer_id       *string        `db:"trainer_id"`
}
