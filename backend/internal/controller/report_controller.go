package controller

import (
	"net/http"

	"chicken-farm/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

func (c *ReportController) RegisterRoutes(router *gin.Engine) {
	reports := router.Group("/api/reports")
	{
		reports.GET("/egg-stats", c.GetTotalEggStats)
		reports.GET("/employee-egg-stats", c.GetEmployeeEggStats)
		reports.GET("/low-productivity-chickens", c.GetLowProductivityChickens)
		reports.GET("/most-productive-chicken", c.GetMostProductiveChickenStats)
		reports.GET("/employee-chicken-counts", c.GetEmployeeChickenCountStats)
	}
}

func (c *ReportController) GetTotalEggStats(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	stats, err := c.reportService.GetTotalEggStats(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"start_date": startDate,
		"end_date":   endDate,
		"total_eggs": stats.TotalEggs,
		"total_cost": stats.TotalCost,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ReportController) GetEmployeeEggStats(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	stats, err := c.reportService.GetEmployeeEggStats(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"start_date": startDate,
		"end_date":   endDate,
		"stats":      stats,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ReportController) GetLowProductivityChickens(ctx *gin.Context) {
	chickens, err := c.reportService.GetLowProductivityChickens()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chickens)
}

func (c *ReportController) GetMostProductiveChickenStats(ctx *gin.Context) {
	stats, err := c.reportService.GetMostProductiveChickenStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

func (c *ReportController) GetEmployeeChickenCountStats(ctx *gin.Context) {
	stats, err := c.reportService.GetEmployeeChickenCountStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
