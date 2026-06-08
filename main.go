package main

import (
	"fmt"
	"gin-test/internal/handlers"
	"gin-test/internal/ivf"
	"gin-test/internal/store"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("Building IVF index...")
	ivf.Clusters = ivf.BuildIVFIndexStreamed(store.References, 4096, 1)
	fmt.Println("Initializing Server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	if err := fasthttp.ListenAndServe(":"+port, handlers.FastHTTPHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

}

func init() {
	fmt.Println("Loading references")
	store.LoadReferencesGziped(10000)
	store.LoadNormalizer()
}
