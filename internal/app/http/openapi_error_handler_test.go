package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"gin-template/internal/middleware"
)

type noopServer struct{}

func (noopServer) PostAuthEmailSendVerification(c *gin.Context) { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthEmailVerify(c *gin.Context)           { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthLogin(c *gin.Context)                 { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthLogout(c *gin.Context)                { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthPasswordRequestReset(c *gin.Context)  { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthPasswordReset(c *gin.Context)         { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthRefresh(c *gin.Context)               { c.Status(http.StatusNoContent) }
func (noopServer) PostAuthRegister(c *gin.Context)              { c.Status(http.StatusNoContent) }
func (noopServer) GetFiles(c *gin.Context)                      { c.Status(http.StatusNoContent) }
func (noopServer) GetFilesSearch(c *gin.Context, _ GetFilesSearchParams) {
	c.Status(http.StatusNoContent)
}
func (noopServer) PostFilesUpload(c *gin.Context) { c.Status(http.StatusNoContent) }

//nolint:revive // 生成的 openapi 接口要求保持 DeleteFilesId 这种命名。
func (noopServer) DeleteFilesId(c *gin.Context, _ string) { c.Status(http.StatusNoContent) }

//nolint:revive // 生成的 openapi 接口要求保持 GetFilesId 这种命名。
func (noopServer) GetFilesId(c *gin.Context, _ string)       { c.Status(http.StatusNoContent) }
func (noopServer) GetOptions(c *gin.Context)                 { c.Status(http.StatusNoContent) }
func (noopServer) PostOptions(c *gin.Context)                { c.Status(http.StatusNoContent) }
func (noopServer) PutOptions(c *gin.Context)                 { c.Status(http.StatusNoContent) }
func (noopServer) GetSystemAbout(c *gin.Context)             { c.Status(http.StatusNoContent) }
func (noopServer) GetSystemNotice(c *gin.Context)            { c.Status(http.StatusNoContent) }
func (noopServer) GetSystemStatus(c *gin.Context)            { c.Status(http.StatusNoContent) }
func (noopServer) GetUsers(c *gin.Context, _ GetUsersParams) { c.Status(http.StatusNoContent) }
func (noopServer) PostUsers(c *gin.Context)                  { c.Status(http.StatusNoContent) }
func (noopServer) GetUsersMe(c *gin.Context)                 { c.Status(http.StatusNoContent) }
func (noopServer) PutUsersMe(c *gin.Context)                 { c.Status(http.StatusNoContent) }
func (noopServer) GetUsersSearch(c *gin.Context, _ GetUsersSearchParams) {
	c.Status(http.StatusNoContent)
}

//nolint:revive // 生成的 openapi 接口要求保持 DeleteUsersId 这种命名。
func (noopServer) DeleteUsersId(c *gin.Context, _ string) { c.Status(http.StatusNoContent) }

//nolint:revive // 生成的 openapi 接口要求保持 GetUsersId 这种命名。
func (noopServer) GetUsersId(c *gin.Context, _ string) { c.Status(http.StatusNoContent) }

//nolint:revive // 生成的 openapi 接口要求保持 PutUsersId 这种命名。
func (noopServer) PutUsersId(c *gin.Context, _ string) { c.Status(http.StatusNoContent) }

func TestOpenAPIErrorHandlerReturnsProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.Trace())
	RegisterHandlersWithOptions(engine, noopServer{}, GinServerOptions{
		ErrorHandler: OpenAPIErrorHandler,
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/users?page=abc", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}

	var problem struct {
		Status  int             `json:"status"`
		Msg     string          `json:"msg"`
		Details string          `json:"details"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatalf("unmarshal problem: %v", err)
	}
	if problem.Msg != "请求参数无效" {
		t.Fatalf("expected msg 请求参数无效, got %s", problem.Msg)
	}
	if !strings.Contains(problem.Details, "Invalid format for parameter page") {
		t.Fatalf("expected details contains parse error, got %s", problem.Details)
	}
	if strings.TrimSpace(string(problem.Data)) != "null" {
		t.Fatalf("expected null data, got %s", string(problem.Data))
	}
	if strings.TrimSpace(recorder.Header().Get("X-Trace-Id")) == "" {
		t.Fatal("expected trace id in openapi error response header")
	}
}
