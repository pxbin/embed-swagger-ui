package openapiv3

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"sort"
	"strings"
)

// handler handles swagger UI request.
type handler struct {
	*options

	ConfigJSON template.JS

	tpl          *template.Template
	staticServer http.Handler
}

// NewHandler creates HTTP handler for Swagger UI.
func NewHandler(handlerOpts ...HandlerOption) http.Handler {
	opts := &options{}

	for _, o := range handlerOpts {
		o(opts)
	}

	if opts.BasePath != "" {
		opts.BasePath = strings.TrimSuffix(opts.BasePath, "/") + "/"
	}

	h := &handler{
		options: opts,
	}

	js, err := json.Marshal(h.options)
	if err != nil {
		panic(err)
	}

	h.ConfigJSON = template.JS(js) //nolint:gosec // Data is well-formed.

	err = h.LoadIndexTpl()
	if err != nil {
		panic(err)
	}

	// // vsfgen assets.go:
	// h.staticServer = http.StripPrefix(h.BasePath, http.FileServer(swagger_ui.Assets))

	// Note:
	// The patterns are interpreted relative to the package directory containing
	// the source file when use go:embed directive.
	stripped, err := fs.Sub(uiAssets, "swagger_ui/dist")
	if err != nil {
		panic(err)
	}

	h.staticServer = http.StripPrefix(h.BasePath, http.FileServer(http.FS(stripped)))

	return h
}

// IndexTpl creates page template.
//
//nolint:funlen // The template is long.
func (h *handler) LoadIndexTpl() error {
	settings := map[string]string{
		"url":         "url",
		"dom_id":      "'#swagger-ui'",
		"deepLinking": "true",
		"presets": `[
				SwaggerUIBundle.presets.apis,
				SwaggerUIStandalonePreset
			]`,
		"plugins": `[
				SwaggerUIBundle.plugins.DownloadUrl
			]`,
		"layout":                   `"StandaloneLayout"`,
		"showExtensions":           "true",
		"showCommonExtensions":     "true",
		"validatorUrl":             `"https://validator.swagger.io/validator"`,
		"defaultModelsExpandDepth": "1", // Hides schemas, override with value "1" in Config.SettingsUI to show schemas.
		`onComplete`: `function() {
                if (cfg.preAuthorizeApiKey) {
                    for (var name in cfg.preAuthorizeApiKey) {
                        ui.preauthorizeApiKey(name, cfg.preAuthorizeApiKey[name]);
                    }
                }

                var dom = document.querySelector('.scheme-container select');
                for (var key in dom) {
                    if (key.startsWith("__reactInternalInstance$")) {
                        var compInternals = dom[key]._currentElement;
                        var compWrapper = compInternals._owner;
                        compWrapper._instance.setScheme(window.location.protocol.slice(0,-1));
                    }
                }
            }`,
	}

	for k, v := range h.options.SettingsUI {
		settings[k] = v
	}

	settingsStr := make([]string, 0, len(settings))
	for k, v := range settings {
		settingsStr = append(settingsStr, "\t\t\t"+k+": "+v)
	}

	sort.Strings(settingsStr)

	indexHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }} - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <link rel="stylesheet" type="text/css" href="index.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }

        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }

        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>

<body>
<div id="swagger-ui"></div>

<script src="./swagger-ui-bundle.js" charset="UTF-8"> </script>
<script src="./swagger-ui-standalone-preset.js" charset="UTF-8"> </script>

<script>
    window.onload = function () {
        var conf = {{ .ConfigJSON }};
        var url = conf.swaggerJsonUrl;
        if (!url.startsWith("https://") && !url.startsWith("http://")) {
           url = window.location.protocol + "//" + window.location.host + url;
        }

        // Build a system
        var settings = {
` + strings.Join(settingsStr, ",\n") + `
        };

        if (!conf.showTopBar) {
            settings.plugins.push(() => {return {components: {Topbar: () => () => null}}});
        }

        if (conf.hideCurl) {
            settings.plugins.push(() => {return {wrapComponents: {curl: () => () => null}}});
        }

        window.ui = SwaggerUIBundle(settings);
    }
</script>
</body>
</html>
`

	tpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		return err
	}

	h.tpl = tpl
	return nil
}

// ServeHTTP implements http.Handler interface to handle swagger UI request.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSuffix(r.URL.Path, "/") != strings.TrimSuffix(h.BasePath, "/") && h.staticServer != nil {
		h.staticServer.ServeHTTP(w, r)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := h.tpl.Execute(w, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
