package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chicken-farm/internal/controller"
	"chicken-farm/internal/model"
	"chicken-farm/internal/repository"
	"chicken-farm/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	db                 *gorm.DB
	router             *gin.Engine
	chickenController  *controller.ChickenController
	employeeController *controller.EmployeeController
	reportController   *controller.ReportController
}

func (suite *TestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(
		&model.Chicken{},
		&model.Employee{},
		&model.EmployeeCage{},
		&model.Farm{},
		&model.Cage{},
		&model.ConfigParam{},
	)
	suite.Require().NoError(err)

	suite.db = db

	chickenRepo := repository.NewChickenRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	farmRepo := repository.NewFarmRepository(db)

	chickenService := service.NewChickenService(chickenRepo, farmRepo)
	employeeService := service.NewEmployeeService(employeeRepo, farmRepo)
	reportService := service.NewReportService(chickenRepo, employeeRepo, farmRepo)

	suite.chickenController = controller.NewChickenController(chickenService)
	suite.employeeController = controller.NewEmployeeController(employeeService)
	suite.reportController = controller.NewReportController(reportService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	suite.chickenController.RegisterRoutes(router)
	suite.employeeController.RegisterRoutes(router)
	suite.reportController.RegisterRoutes(router)
	suite.router = router

	suite.seedTestData()
}

func (suite *TestSuite) seedTestData() {
	cages := []model.Cage{
		{Number: 1},
		{Number: 2},
		{Number: 3},
	}
	for _, cage := range cages {
		suite.db.Create(&cage)
	}

	chickens := []model.Chicken{
		{CageID: 1, Weight: 2.5, Age: 12, EggPerMonth: 25, Breed: "Леггорн"},
		{CageID: 2, Weight: 3.0, Age: 18, EggPerMonth: 22, Breed: "Род-Айленд"},
	}
	for _, chicken := range chickens {
		suite.db.Create(&chicken)
	}

	employees := []model.Employee{
		{FullName: "Иванов Иван Иванович", PassportData: "1234 567890", Salary: 50000},
		{FullName: "Петров Петр Петрович", PassportData: "2345 678901", Salary: 45000},
	}
	for _, employee := range employees {
		suite.db.Create(&employee)
	}

	employeeCages := []model.EmployeeCage{
		{EmployeeID: 1, CageID: 1},
		{EmployeeID: 2, CageID: 2},
	}
	for _, ec := range employeeCages {
		suite.db.Create(&ec)
	}
}

func (suite *TestSuite) TestGetAllChickens() {
	req, _ := http.NewRequest("GET", "/api/chickens", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var chickens []model.Chicken
	err := json.Unmarshal(w.Body.Bytes(), &chickens)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), chickens, 2)
}

func (suite *TestSuite) TestGetChickenByID() {
	req, _ := http.NewRequest("GET", "/api/chickens/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var chicken model.Chicken
	err := json.Unmarshal(w.Body.Bytes(), &chicken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), chicken.ID)
	assert.Equal(suite.T(), "Леггорн", chicken.Breed)
}

func (suite *TestSuite) TestCreateChicken() {
	newChicken := model.Chicken{
		CageID:      3,
		Weight:      2.8,
		Age:         15,
		EggPerMonth: 28,
		Breed:       "Нью-Гемпшир",
	}

	jsonData, _ := json.Marshal(newChicken)
	req, _ := http.NewRequest("POST", "/api/chickens", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdChicken model.Chicken
	err := json.Unmarshal(w.Body.Bytes(), &createdChicken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Нью-Гемпшир", createdChicken.Breed)
	assert.Equal(suite.T(), uint(3), createdChicken.CageID)
}

func (suite *TestSuite) TestCreateChickenInvalidCage() {
	newChicken := model.Chicken{
		CageID:      999,
		Weight:      2.8,
		Age:         15,
		EggPerMonth: 28,
		Breed:       "Нью-Гемпшир",
	}

	jsonData, _ := json.Marshal(newChicken)
	req, _ := http.NewRequest("POST", "/api/chickens", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *TestSuite) TestCreateChickenOccupiedCage() {
	newChicken := model.Chicken{
		CageID:      1,
		Weight:      2.8,
		Age:         15,
		EggPerMonth: 28,
		Breed:       "Нью-Гемпшир",
	}

	jsonData, _ := json.Marshal(newChicken)
	req, _ := http.NewRequest("POST", "/api/chickens", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *TestSuite) TestUpdateChicken() {
	updatedChicken := model.Chicken{
		ID:          1,
		CageID:      1,
		Weight:      3.0,
		Age:         13,
		EggPerMonth: 30,
		Breed:       "Леггорн Улучшенный",
	}

	jsonData, _ := json.Marshal(updatedChicken)
	req, _ := http.NewRequest("PUT", "/api/chickens/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var chicken model.Chicken
	err := json.Unmarshal(w.Body.Bytes(), &chicken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Леггорн Улучшенный", chicken.Breed)
	assert.Equal(suite.T(), 30, chicken.EggPerMonth)
}

func (suite *TestSuite) TestDeleteChicken() {
	req, _ := http.NewRequest("DELETE", "/api/chickens/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/api/chickens/1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestGetAllEmployees() {
	req, _ := http.NewRequest("GET", "/api/employees", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var employees []model.Employee
	err := json.Unmarshal(w.Body.Bytes(), &employees)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), employees, 2)
}

func (suite *TestSuite) TestGetEmployeeByID() {
	req, _ := http.NewRequest("GET", "/api/employees/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var employee model.Employee
	err := json.Unmarshal(w.Body.Bytes(), &employee)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint(1), employee.ID)
	assert.Equal(suite.T(), "Иванов Иван Иванович", employee.FullName)
}

func (suite *TestSuite) TestCreateEmployee() {
	newEmployee := model.Employee{
		FullName:     "Сидоров Сидор Сидорович",
		PassportData: "3456 789012",
		Salary:       55000,
		Cages:        []uint{3},
	}

	jsonData, _ := json.Marshal(newEmployee)
	req, _ := http.NewRequest("POST", "/api/employees", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdEmployee model.Employee
	err := json.Unmarshal(w.Body.Bytes(), &createdEmployee)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Сидоров Сидор Сидорович", createdEmployee.FullName)
}

func (suite *TestSuite) TestCreateEmployeeInvalidCage() {
	newEmployee := model.Employee{
		FullName:     "Сидоров Сидор Сидорович",
		PassportData: "3456 789012",
		Salary:       55000,
		Cages:        []uint{999},
	}

	jsonData, _ := json.Marshal(newEmployee)
	req, _ := http.NewRequest("POST", "/api/employees", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *TestSuite) TestUpdateEmployee() {
	updatedEmployee := model.Employee{
		ID:           1,
		FullName:     "Иванов Иван Петрович",
		PassportData: "1234 567890",
		Salary:       60000,
		Cages:        []uint{1, 3},
	}

	jsonData, _ := json.Marshal(updatedEmployee)
	req, _ := http.NewRequest("PUT", "/api/employees/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var employee model.Employee
	err := json.Unmarshal(w.Body.Bytes(), &employee)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Иванов Иван Петрович", employee.FullName)
	assert.Equal(suite.T(), 60000.0, employee.Salary)
}

func (suite *TestSuite) TestDeleteEmployee() {
	req, _ := http.NewRequest("DELETE", "/api/employees/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/api/employees/1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestGetLowProductivityChickens() {
	req, _ := http.NewRequest("GET", "/api/reports/low-productivity-chickens", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var chickens []model.Chicken
	err := json.Unmarshal(w.Body.Bytes(), &chickens)
	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), []model.Chicken{}, chickens)
}

func (suite *TestSuite) TestGetMostProductiveChicken() {
	req, _ := http.NewRequest("GET", "/api/reports/most-productive-chicken", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var stats map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), stats, "chicken_id")
	assert.Contains(suite.T(), stats, "cage_id")
	assert.Contains(suite.T(), stats, "egg_per_month")
}

func (suite *TestSuite) TestGetEmployeeChickenCounts() {
	req, _ := http.NewRequest("GET", "/api/reports/employee-chicken-counts", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var stats []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), stats, 2)
}

func (suite *TestSuite) TestGetNonExistentChicken() {
	req, _ := http.NewRequest("GET", "/api/chickens/999", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestGetNonExistentEmployee() {
	req, _ := http.NewRequest("GET", "/api/employees/999", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestInvalidChickenData() {
	invalidData := `{"cage_id": "invalid", "weight": -1}`
	req, _ := http.NewRequest("POST", "/api/chickens", bytes.NewBufferString(invalidData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func BenchmarkGetAllChickens(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Chicken{}, &model.Cage{})

	cage := model.Cage{Number: 1}
	db.Create(&cage)

	for i := 0; i < 100; i++ {
		chicken := model.Chicken{
			CageID:      1,
			Weight:      2.5,
			Age:         12,
			EggPerMonth: 25,
			Breed:       "TestBreed",
		}
		db.Create(&chicken)
	}

	chickenRepo := repository.NewChickenRepository(db)
	farmRepo := repository.NewFarmRepository(db)
	chickenService := service.NewChickenService(chickenRepo, farmRepo)
	chickenController := controller.NewChickenController(chickenService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	chickenController.RegisterRoutes(router)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/chickens", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkCreateChicken(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Chicken{}, &model.Cage{})

	cage := model.Cage{Number: 1}
	db.Create(&cage)

	chickenRepo := repository.NewChickenRepository(db)
	farmRepo := repository.NewFarmRepository(db)
	chickenService := service.NewChickenService(chickenRepo, farmRepo)
	chickenController := controller.NewChickenController(chickenService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	chickenController.RegisterRoutes(router)

	chicken := model.Chicken{
		CageID:      1,
		Weight:      2.5,
		Age:         12,
		EggPerMonth: 25,
		Breed:       "BenchmarkBreed",
	}
	jsonData, _ := json.Marshal(chicken)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/chickens", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TearDownTest() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *TestSuite) TestChickenBusinessLogic() {
	chickenRepo := repository.NewChickenRepository(suite.db)
	farmRepo := repository.NewFarmRepository(suite.db)
	chickenService := service.NewChickenService(chickenRepo, farmRepo)

	invalidChicken := &model.Chicken{
		CageID:      999,
		Weight:      2.5,
		Age:         12,
		EggPerMonth: 25,
		Breed:       "Test",
	}

	err := chickenService.CreateChicken(invalidChicken)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cage not found")
}

func (suite *TestSuite) TestEmployeeBusinessLogic() {
	employeeRepo := repository.NewEmployeeRepository(suite.db)
	farmRepo := repository.NewFarmRepository(suite.db)
	employeeService := service.NewEmployeeService(employeeRepo, farmRepo)

	invalidEmployee := &model.Employee{
		FullName:     "Test Employee",
		PassportData: "9999 999999",
		Salary:       50000,
		Cages:        []uint{999},
	}

	err := employeeService.CreateEmployee(invalidEmployee)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cage not found")
}
