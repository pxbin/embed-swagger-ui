package main

import (
	"log"
	"net/http"

	"github.com/pxbin/embed-swagger-ui/openapiv3"
)

func main() {
	router := openapiv3.NewHandler(
		openapiv3.WithTitle("Petstore"),
		openapiv3.WithLocalFile("./openapi.yaml"),
		// openapiv3.WithSwaggerJSON("https://petstore3.swagger.io/api/v3/openapi.json"),
	)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
	}

	log.Fatal(srv.ListenAndServe())
}
