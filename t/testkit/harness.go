package testkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gin-template/internal/app/bootstrap"
	"gin-template/internal/app/config"
	"gin-template/internal/app/db"
	appEnv "gin-template/internal/app/env"
	appHTTP "gin-template/internal/app/http"
	appLog "gin-template/internal/app/log"
	"gin-template/internal/app/security"
	appService "gin-template/internal/service"
	sysuserStore "gin-template/internal/store/sysuser"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin"

	DefaultRootUsername = "root"
	DefaultRootPassword = "root"
	DefaultRootEmail    = "root@example.com"
)

// UploadFile 描述 multipart 上传时的文件内容。
type UploadFile struct {
	FieldName   string
	Filename    string
	ContentType string
	Data        []byte
}

// Harness 封装 OpenAPI 集成测试的运行环境。
type Harness struct {
	handler   http.Handler
	baseURL   *url.URL
	repoRoot  string
	tempRoot  string
	dbPath    string
	uploadDir string
}

// Session 表示一个带独立 cookie jar 的测试客户端。
type Session struct {
	Harness     *Harness
	Client      *http.Client
	AccessToken string
	User        appHTTP.User
}

// NewHarness 启动一套隔离的 sqlite + httptest 服务，供单个测试用例使用。
func NewHarness(t *testing.T) *Harness {
	t.Helper()

	repoRoot := findRepoRoot(t)
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	h := &Harness{
		repoRoot: repoRoot,
		tempRoot: t.TempDir(),
	}
	h.dbPath = filepath.Join(h.tempRoot, "integration.db")
	h.uploadDir = filepath.Join(h.tempRoot, "uploads")

	t.Cleanup(func() {
		appLog.Sync()
		appLog.ResetForTest()
		db.ResetForTest()
		config.ResetForTest()
		appEnv.ResetForTest()
		appService.ResetOptionServiceForTest()
		_ = os.Chdir(originalWD)
	})

	appLog.ResetForTest()
	db.ResetForTest()
	config.ResetForTest()
	appEnv.ResetForTest()
	appService.ResetOptionServiceForTest()

	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("change directory to repo root: %v", err)
	}

	t.Setenv("APP_ENV", "dev")
	t.Setenv("APP_DATABASE_SQLITE_PATH", h.dbPath)
	t.Setenv("APP_STORAGE_LOCAL_DIR", h.uploadDir)
	t.Setenv("APP_RATE_LIMIT_ENABLED", "false")
	t.Setenv("APP_TURNSTILE_ENABLED", "false")
	t.Setenv("APP_MAIL_ENABLED", "false")
	t.Setenv("APP_AUTH_JWT_SECRET", "integration-secret")

	if err := config.Load(); err != nil {
		t.Fatalf("load config: %v", err)
	}

	if err := os.Chdir(h.tempRoot); err != nil {
		t.Fatalf("change directory to temp root: %v", err)
	}

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		t.Fatalf("create upload dir: %v", err)
	}

	if err := bootstrap.EnsureSystemOptions(context.Background()); err != nil {
		t.Fatalf("ensure system options: %v", err)
	}
	if err := bootstrap.EnsureDefaultAdmin(context.Background()); err != nil {
		t.Fatalf("ensure default admin: %v", err)
	}
	if err := seedRootUser(context.Background()); err != nil {
		t.Fatalf("seed root user: %v", err)
	}

	engine, err := appHTTP.NewEngine()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	h.handler = engine
	h.baseURL, err = url.Parse("http://api.integration.test")
	if err != nil {
		t.Fatalf("parse base url: %v", err)
	}
	return h
}

// APIURL 返回带 `/api/v1` 前缀的完整地址。
func (h *Harness) APIURL(path string) string {
	if strings.HasPrefix(path, "/") {
		return h.baseURL.String() + "/api/v1" + path
	}
	return h.baseURL.String() + "/api/v1/" + path
}

// TempRoot 返回当前测试用例的临时根目录。
func (h *Harness) TempRoot() string {
	return h.tempRoot
}

// DBPath 返回隔离 sqlite 文件路径。
func (h *Harness) DBPath() string {
	return h.dbPath
}

// UploadDir 返回隔离上传目录。
func (h *Harness) UploadDir() string {
	return h.uploadDir
}

// NewSession 创建一个匿名测试客户端。
func (h *Harness) NewSession(t *testing.T) *Session {
	t.Helper()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("new cookie jar: %v", err)
	}

	client := &http.Client{
		Jar:       jar,
		Timeout:   10 * time.Second,
		Transport: &handlerTransport{handler: h.handler},
	}
	return &Session{
		Harness: h,
		Client:  client,
	}
}

// RegisterUser 通过公开注册接口创建普通用户，并保留其 cookie / token。
func (h *Harness) RegisterUser(t *testing.T, username, email, password string) *Session {
	t.Helper()

	session := h.NewSession(t)
	var resp appHTTP.AuthEnvelope
	session.JSON(t, http.MethodPost, "/auth/register", map[string]any{
		"username": username,
		"email":    email,
		"password": password,
	}, http.StatusOK, &resp)
	session.AccessToken = resp.Data.AccessToken
	session.User = resp.Data.User
	return session
}

// Login 使用用户名或邮箱登录，并保留其 cookie / token。
func (h *Harness) Login(t *testing.T, identifier, password string) *Session {
	t.Helper()

	session := h.NewSession(t)
	var resp appHTTP.AuthEnvelope
	session.JSON(t, http.MethodPost, "/auth/login", map[string]any{
		"identifier": identifier,
		"password":   password,
	}, http.StatusOK, &resp)
	session.AccessToken = resp.Data.AccessToken
	session.User = resp.Data.User
	return session
}

// LoginAdmin 登录默认管理员。
func (h *Harness) LoginAdmin(t *testing.T) *Session {
	t.Helper()
	return h.Login(t, DefaultAdminUsername, DefaultAdminPassword)
}

// LoginRoot 登录测试夹具中的 root 账号。
func (h *Harness) LoginRoot(t *testing.T) *Session {
	t.Helper()
	return h.Login(t, DefaultRootUsername, DefaultRootPassword)
}

// JSON 发送 JSON 请求，并在成功后解码响应。
func (s *Session) JSON(t *testing.T, method, path string, body any, wantStatus int, target any) []byte {
	t.Helper()

	var reader io.Reader
	contentType := ""
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal json body: %v", err)
		}
		reader = bytes.NewReader(payload)
		contentType = "application/json"
	}
	return s.Harness.doRequest(t, s.Client, s.AccessToken, method, path, reader, contentType, wantStatus, target)
}

// Raw 发送原始请求体，适合构造非法 JSON 等场景。
func (s *Session) Raw(t *testing.T, method, path string, body io.Reader, contentType string, wantStatus int, target any) []byte {
	t.Helper()
	return s.Harness.doRequest(t, s.Client, s.AccessToken, method, path, body, contentType, wantStatus, target)
}

// Multipart 发送 multipart/form-data 请求，可选携带文件。
func (s *Session) Multipart(t *testing.T, method, path string, fields map[string]string, file *UploadFile, wantStatus int, target any) []byte {
	t.Helper()

	var payload bytes.Buffer
	writer := multipart.NewWriter(&payload)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write multipart field %q: %v", key, err)
		}
	}

	if file != nil {
		part, err := writer.CreateFormFile(file.FieldName, file.Filename)
		if err != nil {
			t.Fatalf("create multipart file %q: %v", file.Filename, err)
		}
		if _, err := part.Write(file.Data); err != nil {
			t.Fatalf("write multipart file %q: %v", file.Filename, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	return s.Harness.doRequest(t, s.Client, s.AccessToken, method, path, bytes.NewReader(payload.Bytes()), writer.FormDataContentType(), wantStatus, target)
}

func (h *Harness) doRequest(t *testing.T, client *http.Client, accessToken, method, path string, body io.Reader, contentType string, wantStatus int, target any) []byte {
	t.Helper()

	req, err := http.NewRequest(method, h.APIURL(path), body)
	if err != nil {
		t.Fatalf("new request %s %s: %v", method, path, err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body %s %s: %v", method, path, err)
	}

	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s: expected status %d, got %d, body=%s", method, path, wantStatus, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if target != nil && len(bytes.TrimSpace(respBody)) > 0 {
		if err := json.Unmarshal(respBody, target); err != nil {
			t.Fatalf("unmarshal response %s %s: %v, body=%s", method, path, err, strings.TrimSpace(string(respBody)))
		}
	}

	return respBody
}

type handlerTransport struct {
	handler http.Handler
}

func (t *handlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		_ = req.Body.Close()
	}

	clone := httptest.NewRequest(req.Method, req.URL.String(), bytes.NewReader(bodyBytes))
	clone.Header = req.Header.Clone()
	clone.Host = req.Host
	clone.RemoteAddr = "127.0.0.1:12345"

	recorder := httptest.NewRecorder()
	t.handler.ServeHTTP(recorder, clone)

	resp := recorder.Result()
	resp.Request = req
	return resp, nil
}

func findRepoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	current := wd
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			t.Fatalf("cannot find repo root from %s", wd)
		}
		current = parent
	}
}

func seedRootUser(ctx context.Context) error {
	if _, err := sysuserStore.ByUsername(ctx, DefaultRootUsername); err == nil {
		return nil
	} else if !sysuserStore.IsNotFound(err) {
		return fmt.Errorf("check root user existence: %w", err)
	}

	hash, err := security.HashPassword(DefaultRootPassword)
	if err != nil {
		return fmt.Errorf("hash root password: %w", err)
	}

	rootUser := sysuserStore.New(DefaultRootUsername, DefaultRootEmail, hash)
	rootUser.Role = sysuserStore.RoleRoot
	rootUser.EmailVerified = true
	if err := sysuserStore.Create(ctx, rootUser); err != nil {
		return fmt.Errorf("create root user: %w", err)
	}
	return nil
}
