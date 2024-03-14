package main

import (
	"log"
	"net/http"

	"github.com/pxbin/embed-swagger-ui/openapiv3"
)

func main() {

	http.Handle("/docs/", openapiv3.NewHandler(
		openapiv3.WithBasePath("/docs/"),
		openapiv3.WithTitle("Petstore"),
		openapiv3.WithSwaggerJSON("https://petstore3.swagger.io/api/v3/openapi.json"),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
