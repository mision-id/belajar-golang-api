package main

import (
	"encoding/json"
	"fmt"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DbConn string `mapstructure:"DB_CONNECTION"`
	APIKey string `mapstructure:"API_KEY`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DbConn: viper.GetString("DB_CONNECTION"),
		APIKey: viper.GetString("API_KEY"),
	}

	//Setup Database
	db, err := database.InitDB(config.DbConn)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server Running on ", addr)

	apiKeyMidlleware := middlewares.APIKey(config.APIKey)

	//Injection Endpoint /api/v1/products
	productRepo := repositories.NewProductRepository(db)
	ProductService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(ProductService)

	//Injection Endpoint /api/v1/category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	//Injection Endpoint /api/v1/products
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	//Injection Endpoint /api/v1/report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Endpoint route /api/v1/products
	http.HandleFunc("/api/v1/products", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/v1/products/", middlewares.CORS(middlewares.Logger(apiKeyMidlleware(productHandler.HandlerProductsByID))))

	// Endpoint route /api/v1/categories
	http.HandleFunc("/api/v1/categories", middlewares.CORS(middlewares.Logger(categoryHandler.HandleCategories)))
	http.HandleFunc("/api/v1/categories/", middlewares.CORS(middlewares.Logger(apiKeyMidlleware(categoryHandler.HandleCategoriesByID))))

	// Endpoint route /api/v1/checkout
	http.HandleFunc("/api/v1/checkout", middlewares.CORS(middlewares.Logger(apiKeyMidlleware(transactionHandler.Checkout))))

	// Endpoint route /api/v1/report query param
	http.HandleFunc("/api/v1/report", middlewares.CORS(middlewares.Logger(reportHandler.ReportByRange)))
	// Endpoint route /api/v1/report/today
	http.HandleFunc("/api/v1/report/", middlewares.CORS(middlewares.Logger(apiKeyMidlleware(reportHandler.DailyReport))))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		//set jadi konsensus JSON
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{
			"status":  "Ok",
			"message": "API Running",
		})
		if err != nil {
			http.Error(w, "Failed to  encode health response", http.StatusInternalServerError)
		}
	})

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Server Failed to Start")
	}

}
