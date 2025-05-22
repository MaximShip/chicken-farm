package controller

import (
	"net/http"
	"strconv"

	"chicken-farm/internal/model"
	"chicken-farm/internal/service"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	employeeService *service.EmployeeService
}

func NewEmployeeController(employeeService *service.EmployeeService) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
	}
}

func (c *EmployeeController) RegisterRoutes(router *gin.Engine) {
	employees := router.Group("/api/employees")
	{
		employees.GET("", c.GetAllEmployees)
		employees.GET("/:id", c.GetEmployeeByID)
		employees.POST("", c.CreateEmployee)
		employees.PUT("/:id", c.UpdateEmployee)
		employees.DELETE("/:id", c.DeleteEmployee)
		employees.GET("/:id/chicken-count", c.GetEmployeeChickenCount)
		employees.GET("/chicken-counts", c.GetAllEmployeeChickenCounts)
		employees.GET("/:id/egg-count", c.GetEmployeeEggCount)
		employees.GET("/egg-counts", c.GetAllEmployeeEggCounts)
	}
}

func (c *EmployeeController) GetAllEmployees(ctx *gin.Context) {
	employees, err := c.employeeService.GetAllEmployees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employees)
}

func (c *EmployeeController) GetEmployeeByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	employee, err := c.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var employee model.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.employeeService.CreateEmployee(&employee); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, employee)
}

func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var employee model.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.ID = uint(id)
	if err := c.employeeService.UpdateEmployee(&employee); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.employeeService.DeleteEmployee(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "employee deleted successfully"})
}

func (c *EmployeeController) GetEmployeeChickenCount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	count, err := c.employeeService.GetEmployeeChickenCount(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"employee_id": id, "chicken_count": count})
}

func (c *EmployeeController) GetAllEmployeeChickenCounts(ctx *gin.Context) {
	counts, err := c.employeeService.GetAllEmployeeChickenCounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, counts)
}

func (c *EmployeeController) GetEmployeeEggCount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	count, err := c.employeeService.GetEmployeeEggCount(uint(id), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"employee_id": id,
		"start_date":  startDate,
		"end_date":    endDate,
		"egg_count":   count,
	})
}

func (c *EmployeeController) GetAllEmployeeEggCounts(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	counts, err := c.employeeService.GetAllEmployeeEggCounts(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"start_date": startDate,
		"end_date":   endDate,
		"counts":     counts,
	})
}
