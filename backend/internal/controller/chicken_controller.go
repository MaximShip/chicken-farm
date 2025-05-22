package controller

import (
	"net/http"
	"strconv"

	"chicken-farm/internal/model"
	"chicken-farm/internal/service"

	"github.com/gin-gonic/gin"
)

type ChickenController struct {
	chickenService *service.ChickenService
}

func NewChickenController(chickenService *service.ChickenService) *ChickenController {
	return &ChickenController{
		chickenService: chickenService,
	}
}

func (c *ChickenController) RegisterRoutes(router *gin.Engine) {
	chickens := router.Group("/api/chickens")
	{
		chickens.GET("", c.GetAllChickens)
		chickens.GET("/:id", c.GetChickenByID)
		chickens.POST("", c.CreateChicken)
		chickens.PUT("/:id", c.UpdateChicken)
		chickens.DELETE("/:id", c.DeleteChicken)
		chickens.GET("/low-productivity", c.GetChickensWithLowProductivity)
		chickens.GET("/most-productive", c.GetMostProductiveChicken)
		chickens.GET("/avg-eggs", c.GetAvgEggsByWeightAndAge)
	}
}

func (c *ChickenController) GetAllChickens(ctx *gin.Context) {
	chickens, err := c.chickenService.GetAllChickens()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chickens)
}

func (c *ChickenController) GetChickenByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	chicken, err := c.chickenService.GetChickenByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "chicken not found"})
		return
	}

	ctx.JSON(http.StatusOK, chicken)
}

func (c *ChickenController) CreateChicken(ctx *gin.Context) {
	var chicken model.Chicken
	if err := ctx.ShouldBindJSON(&chicken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.chickenService.CreateChicken(&chicken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, chicken)
}

func (c *ChickenController) UpdateChicken(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var chicken model.Chicken
	if err := ctx.ShouldBindJSON(&chicken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chicken.ID = uint(id)
	if err := c.chickenService.UpdateChicken(&chicken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chicken)
}

func (c *ChickenController) DeleteChicken(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.chickenService.DeleteChicken(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "chicken deleted successfully"})
}

func (c *ChickenController) GetChickensWithLowProductivity(ctx *gin.Context) {
	chickens, err := c.chickenService.GetChickensWithLowProductivity()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chickens)
}

func (c *ChickenController) GetMostProductiveChicken(ctx *gin.Context) {
	chicken, err := c.chickenService.GetMostProductiveChicken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chicken)
}

func (c *ChickenController) GetAvgEggsByWeightAndAge(ctx *gin.Context) {
	weightStr := ctx.Query("weight")
	ageStr := ctx.Query("age")

	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid weight"})
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid age"})
		return
	}

	avgEggs, err := c.chickenService.GetAvgEggsByWeightAndAge(weight, age)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"weight": weight, "age": age, "avg_eggs": avgEggs})
}
