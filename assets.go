package projectassets

import (
	"embed"
	"io/fs"
)

//go:embed api/openapi/openapi.yaml
var OpenAPIFS embed.FS

//go:embed database/bootstrap/*.sql
var BootstrapFS embed.FS

//go:embed admin/dist/*
var AdminDistFS embed.FS

func ReadOpenAPI() ([]byte, error) {
	return OpenAPIFS.ReadFile("api/openapi/openapi.yaml")
}

func ReadBootstrapSQL(driver string) ([]byte, error) {
	return BootstrapFS.ReadFile("database/bootstrap/" + driver + ".sql")
}

func AdminAssets() (fs.FS, error) {
	return fs.Sub(AdminDistFS, "admin/dist")
}
