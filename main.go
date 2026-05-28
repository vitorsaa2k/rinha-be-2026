package main

import (
	"gin-test/internal/routes"
	"gin-test/internal/store"

	"github.com/gofiber/fiber/v3"
)

/* func addAlbum(c fiber.Ctx) {
	var newAlbum album
	c.BindJSON(&newAlbum)
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusOK, newAlbum)
} */

func main() {
	app := fiber.New()
	/* first := utils.GetLimit(transaction.Amount, store.Normalizer.Max_amount)
	second := utils.GetLimit(transaction.Hour, store.Normalizer.Max_hour)
	third := utils.GetLimit(transaction.Customer_avg_amount, store.Normalizer.Max_avg)
	vector := []float64{first, second, third}
	result, err := utils.SearchInVector(vector, 3, store.Dataset)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if result.IsPossibleFraud {
		fmt.Println("This is possibly a fraud")
	} else {
		fmt.Println("This is possibly not a fraud")
	}
	fmt.Println(result)
	fmt.Println(transaction)
	fmt.Println(first, second, third) */
	routes.RegisterRoutes(app)
	app.Listen(":8090", fiber.ListenConfig{
		DisableStartupMessage: true,
	})

}

func init() {
	store.LoadReferencesGziped(100000)
	store.LoadNormalizer()
}
