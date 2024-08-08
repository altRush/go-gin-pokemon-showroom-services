package main

import (
	store_pokemons "go-gin-pokemon-showroom-services/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/store/all", store_pokemons.GetAllStoredPokemons)

	router.Run("localhost:8080")
}
