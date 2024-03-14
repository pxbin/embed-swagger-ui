package openapiv3

type options struct {
	basePath    string `json:"basePath"`       // Base URL to docs.
	Title       string `json:"title"`          // Title of index file.
	SwaggerJSON string `json:"swaggerJsonUrl"` // URL to openapi.json/swagger.json document specification.

	// InternalBasePath is used to override BasePath if external
	// url differs from internal one.
	InternalBasePath string `json:"-"`

	ShowTopBar         bool              `json:"showTopBar"`         // Show navigation top bar, hidden by default.
	HideCurl           bool              `json:"hideCurl"`           // Hide curl code snippet.
	JSONEditor         bool              `json:"jsonEditor"`         // Enable visual json editor support (experimental, can fail with complex schemas).
	PreAuthorizeAPIKey map[string]string `json:"preAuthorizeApiKey"` // Map of security name to key value.

	// SettingsUI contains keys and plain javascript values of SwaggerUIBundle configuration.
	// Overrides default values.
	// See https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/ for available options.
	SettingsUI map[string]string `json:"-"`

	LocalOpenAPIFile string `json:"-"` // Local openapi file path.
}

type HandlerOption func(opt *options)

// // WithBasePath sets base URL to docs.
// func WithBasePath(basePath string) HandlerOption {
// 	return func(opt *options) {
// 		opt.basePath = basePath
// 	}
// }

// WithTitle sets title of index file.
func WithTitle(title string) HandlerOption {
	return func(opt *options) {
		opt.Title = title
	}
}

// WithSwaggerJSON sets URL to openapi.json/swagger.json document specification.
func WithSwaggerJSON(url string) HandlerOption {
	return func(opt *options) {
		opt.SwaggerJSON = url
	}
}

// WithInternalBasePath sets internal base URL to docs.
func WithInternalBasePath(basePath string) HandlerOption {
	return func(opt *options) {
		opt.InternalBasePath = basePath
	}
}

// WithShowTopBar sets whether to show navigation top bar.
func WithShowTopBar(show bool) HandlerOption {
	return func(opt *options) {
		opt.ShowTopBar = show
	}
}

// WithHideCurl sets whether to hide curl code snippet.
func WithHideCurl(hide bool) HandlerOption {
	return func(opt *options) {
		opt.HideCurl = hide
	}
}

// WithJSONEditor sets whether to enable visual json editor support.
func WithJSONEditor(enable bool) HandlerOption {
	return func(opt *options) {
		opt.JSONEditor = enable
	}
}

// WithPreAuthorizeAPIKey sets map of security name to key value.
func WithPreAuthorizeAPIKey(securityName, key string) HandlerOption {
	return func(opt *options) {
		if opt.PreAuthorizeAPIKey == nil {
			opt.PreAuthorizeAPIKey = make(map[string]string)
		}
		opt.PreAuthorizeAPIKey[securityName] = key
	}
}

// WithSettingsUI sets keys and plain javascript values of SwaggerUIBundle configuration.
// Overrides default values.
// See https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/ for available options.
func WithSettingsUI(settings map[string]string) HandlerOption {
	return func(opt *options) {
		opt.SettingsUI = settings
	}
}

func WithLocalFile(filePath string) HandlerOption {
	return func(opt *options) {
		opt.LocalOpenAPIFile = filePath
	}
}
