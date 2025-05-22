package main

import (
	"log"
	"net/http"

	"chicken-farm/internal/controller"
	"chicken-farm/internal/model"
	"chicken-farm/internal/repository"
	"chicken-farm/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := migrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := seedData(db); err != nil {
		log.Fatal("Failed to seed data:", err)
	}

	chickenRepo := repository.NewChickenRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	farmRepo := repository.NewFarmRepository(db)

	chickenService := service.NewChickenService(chickenRepo, farmRepo)
	employeeService := service.NewEmployeeService(employeeRepo, farmRepo)
	reportService := service.NewReportService(chickenRepo, employeeRepo, farmRepo)

	chickenController := controller.NewChickenController(chickenService)
	employeeController := controller.NewEmployeeController(employeeService)
	reportController := controller.NewReportController(reportService)

	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 часов
	}
	router.Use(cors.New(config))

	router.Use(func(c *gin.Context) {
		log.Printf("Method: %s, Path: %s, Origin: %s", c.Request.Method, c.Request.URL.Path, c.Request.Header.Get("Origin"))
		c.Next()
	})

	router.OPTIONS("/api/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusOK)
	})

	chickenController.RegisterRoutes(router)
	employeeController.RegisterRoutes(router)
	reportController.RegisterRoutes(router)

	router.Static("/static", "./web/build/static")
	router.StaticFile("/", "./web/build/index.html")

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("chicken_farm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Chicken{},
		&model.Employee{},
		&model.EmployeeCage{},
		&model.Farm{},
		&model.Cage{},
		&model.ConfigParam{},
	)
}

// начальные данные
func seedData(db *gorm.DB) error {
	var cageCount int64
	if err := db.Model(&model.Cage{}).Count(&cageCount).Error; err != nil {
		return err
	}

	if cageCount > 0 {
		return nil
	}

	cages := []model.Cage{
		{Number: 1},
		{Number: 2},
		{Number: 3},
		{Number: 4},
		{Number: 5},
	}

	for _, cage := range cages {
		if err := db.Create(&cage).Error; err != nil {
			return err
		}
	}

	chickens := []model.Chicken{
		{CageID: 1, Weight: 2.5, Age: 12, EggPerMonth: 25, Breed: "Леггорн"},
		{CageID: 2, Weight: 3.0, Age: 18, EggPerMonth: 22, Breed: "Род-Айленд"},
		{CageID: 3, Weight: 2.8, Age: 15, EggPerMonth: 28, Breed: "Нью-Гемпшир"},
	}

	for _, chicken := range chickens {
		if err := db.Create(&chicken).Error; err != nil {
			return err
		}
	}

	employees := []model.Employee{
		{FullName: "Иванов Иван Иванович", PassportData: "1234 567890", Salary: 50000},
		{FullName: "Петров Петр Петрович", PassportData: "2345 678901", Salary: 45000},
	}

	for _, employee := range employees {
		if err := db.Create(&employee).Error; err != nil {
			return err
		}
	}

	employeeCages := []model.EmployeeCage{
		{EmployeeID: 1, CageID: 1},
		{EmployeeID: 1, CageID: 2},
		{EmployeeID: 2, CageID: 3},
		{EmployeeID: 2, CageID: 4},
		{EmployeeID: 2, CageID: 5},
	}

	for _, ec := range employeeCages {
		if err := db.Create(&ec).Error; err != nil {
			return err
		}
	}

	eggPrice := model.ConfigParam{
		Key:   "egg_price",
		Value: "10.0",
	}

	if err := db.Create(&eggPrice).Error; err != nil {
		return err
	}

	return nil
}
