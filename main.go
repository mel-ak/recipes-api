package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mel-ak/recipes-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// --- General API Information Annotations (from your input) ---
// @title Recipes API
// @description This is a sample recipes API. You can find out more about the API at https://github.com/mel-ak/recipes-api
// @version 1.0.0
// @host localhost:8080
// @BasePath /
// @contact.name Melak Sisay
// @contact.url https://me-port-sigma.vercel.app/
// @contact.email melakesisay@gmail.com
// @accept json
// @produce json
// @schemes http

// --- Data Model Annotation ---

// Recipe represents a food recipe.
// swagger:model
type Recipe struct {
	// the id for the recipe
	// read only: true
	// example: 60e9d1f37e19a6b83f3e69f8
	ID string `json:"id"`
	// The name of the recipe
	// required: true
	// min length: 3
	// example: Spaghetti Carbonara
	Nmae string `json:"name"`
	// A list of tags for categorization (e.g., Italian, Pasta)
	// example: ["Italian", "Dinner"]
	Tags []string `json:"tags"`
	// A list of ingredients needed for the recipe
	// required: true
	Ingredients []string `json:"ingredients"`
	// Step-by-step instructions
	// required: true
	Instructions []string `json:"instructions`
	// The date the recipe was published
	PublishedAt time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

// NewRecipeHandler handles the request to create a new recipe.
// @Summary Create a new recipe
// @Description Creates a new recipe and stores it in the database.
// @Tags recipes
// @ID newRecipe
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe object to be created"
// @Success 201 {object} Recipe "Successfully created the recipe"
// @Failure 400 {object} string "Invalid input or validation failed"
// @Router /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler handles the request to list all recipes.
// @Summary List all recipes
// @Description Gets a comprehensive list of all recipes currently available.
// @Tags recipes
// @ID listRecipes // Added ID tag based on user input
// @Accept json
// @Produce json
// @Success 200 {array} Recipe "Successfully retrieved list of recipes"
// @Failure 500 {object} string "Internal Server Error"
// @Router /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler handles the request to update an existing recipe.
// @Summary Update an existing recipe
// @Description Updates the details of a recipe specified by its ID.
// @Tags recipes
// @ID updateRecipe
// @Accept json
// @Produce json
// @Param id path string true "ID of the recipe"
// @Param recipe body Recipe true "Recipe object that needs to be updated"
// @Success 200 {object} Recipe "Successful operation, returns the updated recipe"
// @Failure 400 {object} string "Invalid input or validation failed"
// @Failure 404 {object} string "Invalid recipe ID / Recipe not found"
// @Router /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipe.ID = id

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)

}

// DeleteRecipeHandler handles the request to delete a recipe by ID.
// @Summary Delete a recipe
// @Description Deletes a recipe based on the provided ID.
// @Tags recipes
// @ID deleteRecipe
// @Produce json
// @Param id path string true "ID of the recipe to be deleted"
// @Success 204 "No content, deletion successful"
// @Failure 404 {object} string "Recipe ID not found"
// @Router /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})

}

// SearchRecipesHandler handles the request to search recipes by tag.
// @Summary Search recipes
// @Description Search for recipes based on tags or other criteria.
// @Tags recipes
// @ID searchRecipes
// @Accept json
// @Produce json
// @Param tag query string false "Tag to filter recipes by"
// @Success 200 {array} Recipe "List of matching recipes"
// @Router /recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false

		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
