package openapi

import (
	"errors"

	projectassets "gin-template"
	"gin-template/pkg/errs"
	yaml "go.yaml.in/yaml/v3"
)

func ValidateEmbeddedSpec() error {
	data, err := projectassets.ReadOpenAPI()
	if err != nil {
		return errs.Wrap(err, "读取内嵌 OpenAPI 规范失败")
	}
	return ValidateYAML(data)
}

func ValidateYAML(data []byte) error {
	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return errs.Wrap(err, "解析 OpenAPI YAML 失败")
	}
	version, _ := doc["openapi"].(string)
	if version == "" {
		return errs.WithStack(errors.New("openapi yaml 缺少 openapi 版本"))
	}
	info, _ := doc["info"].(map[string]any)
	if info == nil {
		return errs.WithStack(errors.New("openapi yaml 缺少 info 节点"))
	}
	title, _ := info["title"].(string)
	if title == "" {
		return errs.WithStack(errors.New("openapi yaml 缺少 info.title"))
	}
	return nil
}
