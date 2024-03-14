# embed-swagger-ui

## How To Use

install:
```shell
go get -u github.com/pxbin/embed-swagger-ui
```

use:
```go
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

```

## How To Generate Embedded Assets
use `vfsgen` to generate embedded assets:
```shell
cd openapiv3/swagger_ui/vfsgen && go run generate.go 
```

- update openapiv3/handler.go: L50-L51
```go 
// vsfgen assets.go:
h.staticServer = http.StripPrefix(h.BasePath, http.FileServer(swagger_ui.Assets))
```
- update openapiv3/handler.go: L50-L61 打开注释或使用以下代码来替换
```go 
// vsfgen assets.go:
h.staticServer = http.StripPrefix(h.BasePath, http.FileServer(swagger_ui.Assets))
```
- update openapiv3/handler.go: L53-L61 移除或注释掉
```go 
	// // Note:
	// // The patterns are interpreted relative to the package directory containing
	// // the source file when use go:embed directive.
	// stripped, err := fs.Sub(uiAssets, "swagger_ui/dist")
	// if err != nil {
	// 	panic(err)
	// }

	// h.staticServer = http.StripPrefix(h.BasePath, http.FileServer(http.FS(stripped)))
```

or

use `go:embed` to generate embedded assets:
```go
import "embed"

//go:embed swagger_ui/dist/*
var uiAssets embed.FS
```
Note: 
- `go:embed` is only available in go1.16 or later.
- "The patterns are interpreted relative to the package directory containing the source file"

因为上面的语句得到的 uiAssets 中每一个文件都是以 swagger_ui/dist/ 为根目录，它们的路径都是以 swagger_ui/dist/ 开头，比如 index.html，实际是 swagger_ui/dist/index.html。
使用 fs.Sub() 将 swagger_ui/dist 从 uiAssets 根目录路径中 strip 掉。
```go 
func Assets() http.FileSystem {
	stripped, err := fs.Sub(uiAssets, "swagger_ui/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(stripped)
}
```

## References

- [swagger-ui - github](https://github.com/swagger-api/swagger-ui)
- [Swagger Open API Specification 2.0 and 3.0 in Go](https://kecci.medium.com/swagger-open-api-specification-2-0-and-3-0-in-go-c1f05b51a595)

- [Serve SwaggerUI within your Golang application](https://ribice.medium.com/serve-swaggerui-within-your-golang-application-5486748a5ed4)
- [go-kratos swagger-api](https://github.com/go-kratos/swagger-api)
- [tx7do kratos-swagger-api](https://github.com/tx7do/kratos-swagger-ui)
- [vfsgen](https://github.com/shurcooL/vfsgen)
- [vfsgen-sample](https://github.com/qtopie/vfsgen-sample)

- [Embed design](https://github.com/golang/proposal/blob/master/design/draft-embed.md)
- [Embed examples](https://pkg.go.dev/embed)
- [Embedded Swagger UI for Go](https://github.com/swaggest/swgui)
- [Embedding Vue.js Apps in Go](https://hackandsla.sh/posts/2021-06-18-embed-vuejs-in-go/)

- [Tutorial: Developing a RESTful API with Go, JSON Schema validation and OpenAPI docs](https://dev.to/vearutop/tutorial-developing-a-restful-api-with-go-json-schema-validation-and-openapi-docs-2490)