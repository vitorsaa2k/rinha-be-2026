package main

import (
	"gin-test/internal/routes"
	"gin-test/models"

	"github.com/gofiber/fiber/v3"
)

var transaction = models.TransactionSctruct{Amount: 12500, Hour: 22, Customer_avg_amount: 4800}

var albums = []models.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func home(c fiber.Ctx) error {

	return c.JSON(albums)
}

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

	app.Get("/", home)
	routes.RegisterRoutes(app)
	app.Listen(":8090")

}
