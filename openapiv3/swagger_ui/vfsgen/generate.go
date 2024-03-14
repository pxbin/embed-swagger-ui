//go:build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("../dist")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "swagger_ui",
		BuildTags:    "!swagger_ui",
		Filename:     "../assets.go",
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
