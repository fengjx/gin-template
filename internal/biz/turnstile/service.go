package turnstile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gin-template/internal/app/config"
	"gin-template/pkg/errs"
)

func Verify(ctx context.Context, token, remoteIP string) error {
	cfg := config.Get()
	if !cfg.Turnstile.Enabled {
		return nil
	}
	if strings.TrimSpace(token) == "" {
		return errs.WithStack(errors.New("turnstile token 不能为空"))
	}
	form := url.Values{}
	form.Set("secret", cfg.Turnstile.SecretKey)
	form.Set("response", token)
	form.Set("remoteip", remoteIP)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://challenges.cloudflare.com/turnstile/v0/siteverify", strings.NewReader(form.Encode()))
	if err != nil {
		return errs.Wrap(err, "创建 turnstile 校验请求失败")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errs.Wrap(err, "发送 turnstile 校验请求失败")
	}
	defer resp.Body.Close()

	var body struct {
		Success bool     `json:"success"`
		Errors  []string `json:"error-codes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return errs.Wrap(err, "解析 turnstile 校验响应失败")
	}
	if !body.Success {
		return errs.WithStack(fmt.Errorf("turnstile 校验失败: %v", body.Errors))
	}
	return nil
}
