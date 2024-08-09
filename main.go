package main

import (
	"github.com/altRush/go-gin-pokemon-showroom-services/models"
	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/store/all", getAllStoredPokemons)
	router.GET("/store/:pokemonStoreId", getPokemonByStoreIdFromStore)

	router.POST("/store", addPokemonToStore)

	router.Run("localhost:8080")
}

func addPokemonToStore(c *gin.Context) {

	add_pokemon_to_store, err := models.AddPokemonToStore(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, add_pokemon_to_store)
}

func getPokemonByStoreIdFromStore(c *gin.Context) {
	pokemon, err := models.GetPokemonByStoreIdFromStore(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, pokemon)

}

func getAllStoredPokemons(c *gin.Context) {
	allStoredPokemon, err := models.GetAllStoredPokemons(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, allStoredPokemon)
}
