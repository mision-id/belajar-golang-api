package main

import (
	"encoding/json"
	"fmt"

	"kasir-api/database"
	"kasir-api/handlers"
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
}

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

	//Endpoint route /api/v1/products
	productRepo := repositories.NewProductRepository(db)
	ProductService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(ProductService)

	//Endpoint route /api/v1/category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Endpoint route /api/v1/products
	http.HandleFunc("/api/v1/products", productHandler.HandleProducts)
	http.HandleFunc("/api/v1/products/", productHandler.HandlerProductsByID)

	// Endpoint route /api/v1/categories
	http.HandleFunc("/api/v1/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/v1/categories/", categoryHandler.HandleCategoriesByID)

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
