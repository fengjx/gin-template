package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"gin-template/pkg/errs"
)

type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type WeChatUser struct {
	OpenID string `json:"openid"`
	Name   string `json:"nickname"`
}

func GitHubAuthURL(clientID, redirectURL, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURL)
	q.Set("scope", "read:user user:email")
	q.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + q.Encode()
}

func ExchangeGitHubToken(ctx context.Context, clientID, clientSecret, code string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", errs.Wrap(err, "创建 GitHub access token 请求失败")
	}
	q := req.URL.Query()
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", code)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errs.Wrap(err, "请求 GitHub access token 失败")
	}
	defer resp.Body.Close()

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", errs.Wrap(err, "解析 GitHub access token 响应失败")
	}
	if body.AccessToken == "" {
		return "", errs.WithStack(errors.New("github oauth failed"))
	}
	return body.AccessToken, nil
}

func FetchGitHubUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, errs.Wrap(err, "创建 GitHub 用户信息请求失败")
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(err, "请求 GitHub 用户信息失败")
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errs.Wrap(err, "解析 GitHub 用户信息失败")
	}
	return &user, nil
}

func WeChatAuthURL(appID, redirectURL, state string) string {
	q := url.Values{}
	q.Set("appid", appID)
	q.Set("redirect_uri", redirectURL)
	q.Set("response_type", "code")
	q.Set("scope", "snsapi_login")
	q.Set("state", state)
	return "https://open.weixin.qq.com/connect/qrconnect?" + q.Encode() + "#wechat_redirect"
}

func ExchangeWeChatToken(ctx context.Context, appID, secret, code string) (string, string, error) {
	q := url.Values{}
	q.Set("appid", appID)
	q.Set("secret", secret)
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weixin.qq.com/sns/oauth2/access_token?"+q.Encode(), nil)
	if err != nil {
		return "", "", errs.Wrap(err, "创建微信 access token 请求失败")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", errs.Wrap(err, "请求微信 access token 失败")
	}
	defer resp.Body.Close()

	var body struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", "", errs.Wrap(err, "解析微信 access token 响应失败")
	}
	if body.AccessToken == "" || body.OpenID == "" {
		return "", "", errs.WithStack(fmt.Errorf("wechat oauth failed: %s", body.ErrMsg))
	}
	return body.AccessToken, body.OpenID, nil
}

func FetchWeChatUser(ctx context.Context, accessToken, openID string) (*WeChatUser, error) {
	q := url.Values{}
	q.Set("access_token", accessToken)
	q.Set("openid", openID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weixin.qq.com/sns/userinfo?"+q.Encode(), nil)
	if err != nil {
		return nil, errs.Wrap(err, "创建微信用户信息请求失败")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(err, "请求微信用户信息失败")
	}
	defer resp.Body.Close()

	var user WeChatUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errs.Wrap(err, "解析微信用户信息失败")
	}
	return &user, nil
}
