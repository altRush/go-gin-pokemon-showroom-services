package main

import (
	"database/sql"
	"net/http"

	"github.com/altRush/go-gin-pokemon-showroom-services/models"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/store/all", getAllStoredPokemons)
	router.GET("/store/:pokemonStoreId", getPokemonByStoreIdFromStore)

	router.POST("/store", addPokemonToStore)

	router.Run("localhost:8080")
}

func addPokemonToStore(c *gin.Context) {
	result, err := models.AddPokemonToStore(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func getPokemonByStoreIdFromStore(c *gin.Context) {
	pokemon, err := models.GetPokemonByStoreIdFromStore(c)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "pokemon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

func getAllStoredPokemons(c *gin.Context) {
	allStoredPokemon, err := models.GetAllStoredPokemons(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, allStoredPokemon)
}
