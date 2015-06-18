package api

// ConfigSet POST /v1/apps/<app id>/config/
type ConfigSet struct {
	Values map[string]string `json:"values"`
}
