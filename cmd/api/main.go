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
	indexPath := os.Getenv("INDEX_PATH")
	if indexPath == "" {

		indexPath = "public/out_quantized.bin"
	}
	ivf.LoadPartitions()
	//store.LoadReferencesGzipedPath("./public/references.json.gz", 600000)
	/* if _, err := os.Stat(indexPath); err == nil {
		if _, err := ivf.LoadQuantizedIndexFromFile(indexPath); err != nil {
			log.Fatalf("Failed to load quantized index: %v", err)
		}
	} else {
		fmt.Println("WARN: quantized index file not found, running without index")
	} */
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
	fmt.Println("Loading normalizer and MCC risk")
	store.LoadNormalizer()
	store.LoadMccRisk()
}
