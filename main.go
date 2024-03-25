package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openfoodfacts/openfoodfacts-go"
)

type ProductInfo struct {
	ProductName     string        `json:"product_name"`
	NutrientLevels  NutrientLevel `json:"nutrient_levels"`
	ProductQuantity string        `json:"product_quantity"`
	Ingredients     string        `json:"ingredients"`
	NutriScore      string        `json:"nutri_score"`
}
type NutrientLevel struct {
	Sugars       string `json:"sugars"`
	SaturatedFat string `json:"saturated-fat"`
	Fat          string `json:"fat"`
	Salt         string `json:"salt"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	params := r.URL.Query()
	paramValue := params.Get("param")
	api := openfoodfacts.NewClient("world", "", "")

	// Handle the received parameter
	product, err := api.Product(paramValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productInfo := ProductInfo{
		ProductName:     product.ProductName,
		NutrientLevels:  NutrientLevel(product.NutrientLevels),
		ProductQuantity: product.Quantity,
		Ingredients:     product.IngredientsText,
		NutriScore:      product.NutritionGrades,
	}

	// Convert struct to JSON
	productJSON, err := json.Marshal(productInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(productJSON)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
