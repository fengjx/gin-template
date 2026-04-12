package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	appHTTP "gin-template/internal/app/http"
	sysfileStore "gin-template/internal/store/sysfile"
	sysuserStore "gin-template/internal/store/sysuser"
	"gin-template/t/testkit"

	"go.yaml.in/yaml/v3"
	"gorm.io/gorm"
)

type operationCase struct {
	name      string
	operation string
	run       func(t *testing.T, h *testkit.Harness)
}

type extraCase struct {
	name string
	run  func(t *testing.T, h *testkit.Harness)
}

type suite struct {
	name       string
	operations []operationCase
	extras     []extraCase
}

func TestOpenAPIIntegration(t *testing.T) {
	suites := []suite{
		{
			name:       "system",
			operations: systemOperationCases(),
		},
		{
			name:       "auth",
			operations: authOperationCases(),
			extras:     authExtraCases(),
		},
		{
			name:       "users",
			operations: userOperationCases(),
			extras:     userExtraCases(),
		},
		{
			name:       "files",
			operations: fileOperationCases(),
			extras:     fileExtraCases(),
		},
		{
			name:       "options",
			operations: optionOperationCases(),
			extras:     optionExtraCases(),
		},
	}

	assertOperationCoverage(t, suites)

	for _, suite := range suites {
		suite := suite
		t.Run(suite.name, func(t *testing.T) {
			for _, tc := range suite.operations {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					h := testkit.NewHarness(t)
					tc.run(t, h)
				})
			}

			for _, tc := range suite.extras {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					h := testkit.NewHarness(t)
					tc.run(t, h)
				})
			}
		})
	}
}

func systemOperationCases() []operationCase {
	return []operationCase{
		{
			name:      "GET /system/status returns ok",
			operation: "GET /system/status",
			run: func(t *testing.T, h *testkit.Harness) {
				var resp appHTTP.SystemStatusEnvelope
				h.NewSession(t).JSON(t, http.MethodGet, "/system/status", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Status != "ok" {
					t.Fatalf("expected system status ok, got %q", resp.Data.Status)
				}
			},
		},
		{
			name:      "GET /system/about returns default option",
			operation: "GET /system/about",
			run: func(t *testing.T, h *testkit.Harness) {
				var resp appHTTP.OptionValueEnvelope
				h.NewSession(t).JSON(t, http.MethodGet, "/system/about", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Value != "Gin + React 同构脚手架" {
					t.Fatalf("unexpected about value: %q", resp.Data.Value)
				}
			},
		},
		{
			name:      "GET /system/notice returns default option",
			operation: "GET /system/notice",
			run: func(t *testing.T, h *testkit.Harness) {
				var resp appHTTP.OptionValueEnvelope
				h.NewSession(t).JSON(t, http.MethodGet, "/system/notice", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Value != "欢迎使用 gin-template" {
					t.Fatalf("unexpected notice value: %q", resp.Data.Value)
				}
			},
		},
	}
}

func authOperationCases() []operationCase {
	return []operationCase{
		{
			name:      "POST /auth/register creates user and session",
			operation: "POST /auth/register",
			run: func(t *testing.T, h *testkit.Harness) {
				var resp appHTTP.AuthEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/register", map[string]any{
					"username": "register_user",
					"email":    "register_user@example.com",
					"password": "secret123",
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.AccessToken == "" {
					t.Fatal("expected access token in register response")
				}
				if resp.Data.User.Username != "register_user" {
					t.Fatalf("unexpected username: %q", resp.Data.User.Username)
				}
				if resp.Data.User.Role != sysuserStore.RoleUser {
					t.Fatalf("unexpected role: %q", resp.Data.User.Role)
				}
			},
		},
		{
			name:      "POST /auth/login returns admin token",
			operation: "POST /auth/login",
			run: func(t *testing.T, h *testkit.Harness) {
				var resp appHTTP.AuthEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/login", map[string]any{
					"identifier": testkit.DefaultAdminUsername,
					"password":   testkit.DefaultAdminPassword,
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.AccessToken == "" {
					t.Fatal("expected access token in login response")
				}
				if resp.Data.User.Role != sysuserStore.RoleAdmin {
					t.Fatalf("expected admin role, got %q", resp.Data.User.Role)
				}
			},
		},
		{
			name:      "POST /auth/logout revokes current session",
			operation: "POST /auth/logout",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				var resp appHTTP.MessageEnvelope
				admin.JSON(t, http.MethodPost, "/auth/logout", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Message == "" {
					t.Fatal("expected logout message")
				}
			},
		},
		{
			name:      "POST /auth/refresh issues a new access token",
			operation: "POST /auth/refresh",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var resp appHTTP.AuthEnvelope
				admin.JSON(t, http.MethodPost, "/auth/refresh", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.AccessToken == "" {
					t.Fatal("expected access token in refresh response")
				}
				if resp.Data.User.Username != testkit.DefaultAdminUsername {
					t.Fatalf("unexpected refreshed user: %q", resp.Data.User.Username)
				}
			},
		},
		{
			name:      "POST /auth/password/request-reset returns debug token in dev",
			operation: "POST /auth/password/request-reset",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "reset_user", "reset_user@example.com", "secret123")

				var resp appHTTP.MessageEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/password/request-reset", map[string]any{
					"token": "unused-for-request-reset",
					"email": user.User.Email,
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.DebugToken == nil || *resp.Data.DebugToken == "" {
					t.Fatal("expected debug token for password reset in dev mode")
				}
			},
		},
		{
			name:      "POST /auth/password/reset updates password with debug token",
			operation: "POST /auth/password/reset",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "password_reset_user", "password_reset_user@example.com", "secret123")

				var requestReset appHTTP.MessageEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/password/request-reset", map[string]any{
					"token": "unused-for-request-reset",
					"email": user.User.Email,
				}, http.StatusOK, &requestReset)
				if requestReset.Data.DebugToken == nil || *requestReset.Data.DebugToken == "" {
					t.Fatal("expected debug token for password reset")
				}

				var resetResp appHTTP.MessageEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/password/reset", map[string]any{
					"token":        *requestReset.Data.DebugToken,
					"new_password": "new-secret123",
				}, http.StatusOK, &resetResp)

				assertOKEnvelope(t, resetResp.Status, resetResp.Msg)
				if resetResp.Data.Message == "" {
					t.Fatal("expected password reset message")
				}

				relogged := h.Login(t, user.User.Email, "new-secret123")
				if relogged.User.Email != user.User.Email {
					t.Fatalf("expected to login with new password, got %q", relogged.User.Email)
				}
			},
		},
		{
			name:      "POST /auth/email/send-verification returns debug token in dev",
			operation: "POST /auth/email/send-verification",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "verify_sender", "verify_sender@example.com", "secret123")

				var resp appHTTP.MessageEnvelope
				user.JSON(t, http.MethodPost, "/auth/email/send-verification", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.DebugToken == nil || *resp.Data.DebugToken == "" {
					t.Fatal("expected debug token for email verification in dev mode")
				}
			},
		},
		{
			name:      "POST /auth/email/verify marks email as verified",
			operation: "POST /auth/email/verify",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "verify_user", "verify_user@example.com", "secret123")

				var sendResp appHTTP.MessageEnvelope
				user.JSON(t, http.MethodPost, "/auth/email/send-verification", nil, http.StatusOK, &sendResp)
				if sendResp.Data.DebugToken == nil || *sendResp.Data.DebugToken == "" {
					t.Fatal("expected debug token for email verification")
				}

				var verifyResp appHTTP.MessageEnvelope
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/email/verify", map[string]any{
					"token": *sendResp.Data.DebugToken,
				}, http.StatusOK, &verifyResp)

				assertOKEnvelope(t, verifyResp.Status, verifyResp.Msg)

				var me appHTTP.UserEnvelope
				user.JSON(t, http.MethodGet, "/users/me", nil, http.StatusOK, &me)
				assertOKEnvelope(t, me.Status, me.Msg)
				if !me.Data.EmailVerified {
					t.Fatal("expected verified email after /auth/email/verify")
				}
			},
		},
	}
}

func authExtraCases() []extraCase {
	return []extraCase{
		{
			name: "POST /auth/register returns 409 on duplicate username or email",
			run: func(t *testing.T, h *testkit.Harness) {
				h.RegisterUser(t, "duplicate_user", "duplicate_user@example.com", "secret123")

				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/register", map[string]any{
					"username": "duplicate_user",
					"email":    "another@example.com",
					"password": "secret123",
				}, http.StatusConflict, &problem)
				assertProblemEnvelope(t, problem)

				h.NewSession(t).JSON(t, http.MethodPost, "/auth/register", map[string]any{
					"username": "another_user",
					"email":    "duplicate_user@example.com",
					"password": "secret123",
				}, http.StatusConflict, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /auth/login returns 400 on malformed json",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).Raw(t, http.MethodPost, "/auth/login", bytes.NewBufferString("{"), "application/json", http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /auth/login returns 401 on invalid credentials",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/login", map[string]any{
					"identifier": testkit.DefaultAdminUsername,
					"password":   "wrong-password",
				}, http.StatusUnauthorized, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /auth/email/send-verification returns 401 for anonymous user",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/email/send-verification", nil, http.StatusUnauthorized, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /auth/password/reset returns 400 on invalid token",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodPost, "/auth/password/reset", map[string]any{
					"token":        "not-a-real-token",
					"new_password": "new-secret123",
				}, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
	}
}

func userOperationCases() []operationCase {
	return []operationCase{
		{
			name:      "GET /users/me returns current user",
			operation: "GET /users/me",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "me_user", "me_user@example.com", "secret123")

				var resp appHTTP.UserEnvelope
				user.JSON(t, http.MethodGet, "/users/me", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Email != "me_user@example.com" {
					t.Fatalf("unexpected me email: %q", resp.Data.Email)
				}
			},
		},
		{
			name:      "PUT /users/me updates display name",
			operation: "PUT /users/me",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "update_me_user", "update_me_user@example.com", "secret123")

				var resp appHTTP.UserEnvelope
				user.JSON(t, http.MethodPut, "/users/me", map[string]any{
					"display_name": "新的昵称",
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.DisplayName != "新的昵称" {
					t.Fatalf("expected display name to be updated, got %q", resp.Data.DisplayName)
				}
			},
		},
		{
			name:      "GET /users lists users for admin",
			operation: "GET /users",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				var resp appHTTP.UserListEnvelope
				admin.JSON(t, http.MethodGet, "/users?page=1&page_size=20", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if len(resp.Data.Items) < 2 {
					t.Fatalf("expected at least admin and root users, got %d", len(resp.Data.Items))
				}
			},
		},
		{
			name:      "POST /users creates managed user for admin",
			operation: "POST /users",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				var resp appHTTP.UserEnvelope
				admin.JSON(t, http.MethodPost, "/users", map[string]any{
					"username":       "managed_user",
					"email":          "managed_user@example.com",
					"password":       "secret123",
					"display_name":   "Managed User",
					"role":           sysuserStore.RoleUser,
					"status":         sysuserStore.StatusActive,
					"email_verified": true,
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Username != "managed_user" {
					t.Fatalf("unexpected username: %q", resp.Data.Username)
				}
				if !resp.Data.EmailVerified {
					t.Fatal("expected managed user to be marked as verified")
				}
			},
		},
		{
			name:      "GET /users/search filters users for admin",
			operation: "GET /users/search",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				createManagedUser(t, admin, "search_target", "search_target@example.com")

				var resp appHTTP.UserListEnvelope
				admin.JSON(t, http.MethodGet, "/users/search?q=search_target", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Total < 1 {
					t.Fatalf("expected search result, got total=%d", resp.Data.Total)
				}
				if !containsUser(resp.Data.Items, "search_target") {
					t.Fatal("expected search result to include search_target")
				}
			},
		},
		{
			name:      "GET /users/{id} returns target user",
			operation: "GET /users/{id}",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				created := createManagedUser(t, admin, "lookup_target", "lookup_target@example.com")

				var resp appHTTP.UserEnvelope
				admin.JSON(t, http.MethodGet, fmt.Sprintf("/users/%d", created.Uid), nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Username != "lookup_target" {
					t.Fatalf("unexpected user lookup result: %q", resp.Data.Username)
				}
			},
		},
		{
			name:      "PUT /users/{id} updates target user",
			operation: "PUT /users/{id}",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				created := createManagedUser(t, admin, "update_target", "update_target@example.com")

				var resp appHTTP.UserEnvelope
				admin.JSON(t, http.MethodPut, fmt.Sprintf("/users/%d", created.Uid), map[string]any{
					"display_name":   "Updated Target",
					"email":          "updated_target@example.com",
					"status":         sysuserStore.StatusLocked,
					"email_verified": true,
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.DisplayName != "Updated Target" {
					t.Fatalf("expected updated display name, got %q", resp.Data.DisplayName)
				}
				if resp.Data.Status != sysuserStore.StatusLocked {
					t.Fatalf("expected locked status, got %q", resp.Data.Status)
				}
			},
		},
		{
			name:      "DELETE /users/{id} removes target user",
			operation: "DELETE /users/{id}",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				created := createManagedUser(t, admin, "delete_target", "delete_target@example.com")

				var resp appHTTP.MessageEnvelope
				admin.JSON(t, http.MethodDelete, fmt.Sprintf("/users/%d", created.Uid), nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Message == "" {
					t.Fatal("expected delete user message")
				}

				_, err := sysuserStore.ByUID(context.Background(), created.Uid)
				if !sysuserStore.IsNotFound(err) {
					t.Fatalf("expected deleted user to be absent, got err=%v", err)
				}
			},
		},
	}
}

func userExtraCases() []extraCase {
	return []extraCase{
		{
			name: "GET /users returns 401 for anonymous user",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodGet, "/users", nil, http.StatusUnauthorized, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /users returns 403 for regular user",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "forbidden_user", "forbidden_user@example.com", "secret123")

				var problem appHTTP.Problem
				user.JSON(t, http.MethodGet, "/users", nil, http.StatusForbidden, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /users returns 400 on malformed json",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.Raw(t, http.MethodPost, "/users", bytes.NewBufferString("{"), "application/json", http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /users/{id} returns 400 on invalid id",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodGet, "/users/not-a-number", nil, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /users/{id} returns 404 when user is missing",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodGet, "/users/999999", nil, http.StatusNotFound, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "DELETE /users/{id} returns 400 when deleting current user",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodDelete, fmt.Sprintf("/users/%d", admin.User.Uid), nil, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
	}
}

func fileOperationCases() []operationCase {
	return []operationCase{
		{
			name:      "GET /files lists uploaded files for admin",
			operation: "GET /files",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				uploadFixtureFile(t, admin, "list-files.txt", []byte("list files body"))

				var resp appHTTP.FileListEnvelope
				admin.JSON(t, http.MethodGet, "/files", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Total < 1 {
					t.Fatalf("expected at least one uploaded file, got %d", resp.Data.Total)
				}
			},
		},
		{
			name:      "GET /files/search filters uploaded files",
			operation: "GET /files/search",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				uploadFixtureFile(t, admin, "search-report.txt", []byte("search files body"))

				var resp appHTTP.FileListEnvelope
				admin.JSON(t, http.MethodGet, "/files/search?q=search-report", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Total < 1 {
					t.Fatalf("expected file search result, got total=%d", resp.Data.Total)
				}
				if !containsFile(resp.Data.Items, "search-report.txt") {
					t.Fatal("expected search result to include search-report.txt")
				}
			},
		},
		{
			name:      "GET /files/{id} returns uploaded file",
			operation: "GET /files/{id}",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				uploaded := uploadFixtureFile(t, admin, "lookup-file.txt", []byte("lookup file body"))

				var resp appHTTP.FileEnvelope
				admin.JSON(t, http.MethodGet, "/files/"+strconv.FormatInt(uploaded.Id, 10), nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Id != uploaded.Id {
					t.Fatalf("expected file id %d, got %d", uploaded.Id, resp.Data.Id)
				}
			},
		},
		{
			name:      "POST /files/upload stores file metadata and content",
			operation: "POST /files/upload",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var resp appHTTP.FileEnvelope
				admin.Multipart(t, http.MethodPost, "/files/upload", nil, &testkit.UploadFile{
					FieldName: "file",
					Filename:  "upload-proof.txt",
					Data:      []byte("upload proof"),
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Path == "" {
					t.Fatal("expected uploaded file path")
				}
				if !strings.HasPrefix(resp.Data.Path, h.UploadDir()) {
					t.Fatalf("expected uploaded path under %s, got %s", h.UploadDir(), resp.Data.Path)
				}
				if _, err := os.Stat(resp.Data.Path); err != nil {
					t.Fatalf("expected uploaded file to exist on disk: %v", err)
				}
			},
		},
		{
			name:      "DELETE /files/{id} removes metadata and file content",
			operation: "DELETE /files/{id}",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)
				uploaded := uploadFixtureFile(t, admin, "delete-file.txt", []byte("delete me"))

				var resp appHTTP.MessageEnvelope
				admin.JSON(t, http.MethodDelete, "/files/"+strconv.FormatInt(uploaded.Id, 10), nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Message == "" {
					t.Fatal("expected delete file message")
				}
				if _, err := os.Stat(uploaded.Path); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected file to be deleted from disk, got err=%v", err)
				}
				if _, err := sysfileStore.ByID(context.Background(), uploaded.Id); !errors.Is(err, gorm.ErrRecordNotFound) {
					t.Fatalf("expected file metadata to be deleted, got err=%v", err)
				}
			},
		},
	}
}

func fileExtraCases() []extraCase {
	return []extraCase{
		{
			name: "GET /files returns 401 for anonymous user",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodGet, "/files", nil, http.StatusUnauthorized, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /files returns 403 for regular user",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "file_user", "file_user@example.com", "secret123")

				var problem appHTTP.Problem
				user.JSON(t, http.MethodGet, "/files", nil, http.StatusForbidden, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /files/upload returns 400 when file is missing",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.Multipart(t, http.MethodPost, "/files/upload", nil, nil, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /files/{id} returns 404 when file is missing",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodGet, "/files/999999", nil, http.StatusNotFound, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
	}
}

func optionOperationCases() []operationCase {
	return []operationCase{
		{
			name:      "GET /options returns all system options for admin",
			operation: "GET /options",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var resp appHTTP.OptionListEnvelope
				admin.JSON(t, http.MethodGet, "/options", nil, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if len(resp.Data) < 3 {
					t.Fatalf("expected bootstrap options, got %d", len(resp.Data))
				}
				if !containsOption(resp.Data, "about") || !containsOption(resp.Data, "notice") {
					t.Fatal("expected about and notice options to exist")
				}
				about := findOption(resp.Data, "about")
				if about == nil {
					t.Fatal("expected to find about option")
				}
				if about.Type != "string" || about.Status != "online" {
					t.Fatalf("expected about option to expose type/status, got %#v", about)
				}
			},
		},
		{
			name:      "POST /options creates option for admin",
			operation: "POST /options",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var resp appHTTP.OptionEnvelope
				admin.JSON(t, http.MethodPost, "/options", map[string]any{
					"key":         "site_profile",
					"value":       "{\"name\":\"gin-template\"}",
					"description": "站点配置",
					"is_public":   true,
					"type":        "json",
					"status":      "online",
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.Type != "json" || resp.Data.Status != "online" {
					t.Fatalf("unexpected created option metadata: %#v", resp.Data)
				}
			},
		},
		{
			name:      "PUT /options updates option for admin",
			operation: "PUT /options",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var resp appHTTP.OptionEnvelope
				admin.JSON(t, http.MethodPut, "/options", map[string]any{
					"key":         "about",
					"value":       "新的关于信息",
					"description": "关于信息",
					"is_public":   true,
					"type":        "string",
					"status":      "online",
				}, http.StatusOK, &resp)

				assertOKEnvelope(t, resp.Status, resp.Msg)
				if resp.Data.OptionValue != "新的关于信息" {
					t.Fatalf("expected updated option value, got %q", resp.Data.OptionValue)
				}

				var about appHTTP.OptionValueEnvelope
				h.NewSession(t).JSON(t, http.MethodGet, "/system/about", nil, http.StatusOK, &about)
				assertOKEnvelope(t, about.Status, about.Msg)
				if about.Data.Value != "新的关于信息" {
					t.Fatalf("expected public about endpoint to reflect updated value, got %q", about.Data.Value)
				}
			},
		},
	}
}

func optionExtraCases() []extraCase {
	return []extraCase{
		{
			name: "GET /options returns 401 for anonymous user",
			run: func(t *testing.T, h *testkit.Harness) {
				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodGet, "/options", nil, http.StatusUnauthorized, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "GET /options returns 403 for regular user",
			run: func(t *testing.T, h *testkit.Harness) {
				user := h.RegisterUser(t, "option_user", "option_user@example.com", "secret123")

				var problem appHTTP.Problem
				user.JSON(t, http.MethodGet, "/options", nil, http.StatusForbidden, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /options returns 409 when option already exists",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodPost, "/options", map[string]any{
					"key":         "about",
					"value":       "duplicate",
					"description": "duplicate",
					"is_public":   true,
					"type":        "string",
					"status":      "online",
				}, http.StatusConflict, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "POST /options returns 400 when json is invalid",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodPost, "/options", map[string]any{
					"key":         "broken_json",
					"value":       "{invalid",
					"description": "broken",
					"is_public":   false,
					"type":        "json",
					"status":      "online",
				}, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "PUT /options returns 404 when option is missing",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodPut, "/options", map[string]any{
					"key":         "missing_option",
					"value":       "missing",
					"description": "missing",
					"is_public":   false,
					"type":        "string",
					"status":      "online",
				}, http.StatusNotFound, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "PUT /options returns 400 when type is invalid",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var problem appHTTP.Problem
				admin.JSON(t, http.MethodPut, "/options", map[string]any{
					"key":         "about",
					"value":       "bad type",
					"description": "bad type",
					"is_public":   true,
					"type":        "yaml",
					"status":      "online",
				}, http.StatusBadRequest, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
		{
			name: "PUT /options offline value is hidden from public reader",
			run: func(t *testing.T, h *testkit.Harness) {
				admin := h.LoginAdmin(t)

				var updateResp appHTTP.OptionEnvelope
				admin.JSON(t, http.MethodPut, "/options", map[string]any{
					"key":         "about",
					"value":       "offline about",
					"description": "offline about",
					"is_public":   true,
					"type":        "string",
					"status":      "offline",
				}, http.StatusOK, &updateResp)
				assertOKEnvelope(t, updateResp.Status, updateResp.Msg)

				var problem appHTTP.Problem
				h.NewSession(t).JSON(t, http.MethodGet, "/system/about", nil, http.StatusNotFound, &problem)
				assertProblemEnvelope(t, problem)
			},
		},
	}
}

func assertOperationCoverage(t *testing.T, suites []suite) {
	t.Helper()

	specOperations := loadOpenAPIOperations(t)
	registered := make(map[string]string, len(specOperations))

	for _, suite := range suites {
		for _, tc := range suite.operations {
			if tc.operation == "" {
				t.Fatalf("suite %q case %q has empty operation", suite.name, tc.name)
			}
			if previous, exists := registered[tc.operation]; exists {
				t.Fatalf("duplicate operation registration for %q in %q and %q", tc.operation, previous, tc.name)
			}
			registered[tc.operation] = tc.name
		}
	}

	for operation := range specOperations {
		if _, ok := registered[operation]; !ok {
			t.Fatalf("missing integration case for OpenAPI operation %q", operation)
		}
	}

	for operation := range registered {
		if _, ok := specOperations[operation]; !ok {
			t.Fatalf("registered operation %q is not present in openapi.yaml", operation)
		}
	}
}

func loadOpenAPIOperations(t *testing.T) map[string]struct{} {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(findRepoRoot(t), "api", "openapi", "openapi.yaml"))
	if err != nil {
		t.Fatalf("read openapi.yaml: %v", err)
	}

	var spec struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		t.Fatalf("unmarshal openapi.yaml: %v", err)
	}

	operations := make(map[string]struct{})
	for path, methods := range spec.Paths {
		for method := range methods {
			normalizedMethod := strings.ToUpper(method)
			if !isHTTPMethod(normalizedMethod) {
				continue
			}
			operations[normalizedMethod+" "+path] = struct{}{}
		}
	}
	return operations
}

func isHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func assertOKEnvelope(t *testing.T, status int, msg string) {
	t.Helper()

	if status != 0 {
		t.Fatalf("expected success status 0, got %d", status)
	}
	if msg != "ok" {
		t.Fatalf("expected success msg \"ok\", got %q", msg)
	}
}

func assertProblemEnvelope(t *testing.T, problem appHTTP.Problem) {
	t.Helper()

	if problem.Status == 0 {
		t.Fatal("expected non-zero problem status")
	}
	if strings.TrimSpace(problem.Msg) == "" {
		t.Fatal("expected problem msg")
	}
	if strings.TrimSpace(problem.Details) == "" {
		t.Fatal("expected problem details")
	}
}

func createManagedUser(t *testing.T, admin *testkit.Session, username, email string) appHTTP.User {
	t.Helper()

	var resp appHTTP.UserEnvelope
	admin.JSON(t, http.MethodPost, "/users", map[string]any{
		"username":       username,
		"email":          email,
		"password":       "secret123",
		"display_name":   strings.ToUpper(username),
		"role":           sysuserStore.RoleUser,
		"status":         sysuserStore.StatusActive,
		"email_verified": false,
	}, http.StatusOK, &resp)
	assertOKEnvelope(t, resp.Status, resp.Msg)
	return resp.Data
}

func uploadFixtureFile(t *testing.T, admin *testkit.Session, filename string, content []byte) appHTTP.File {
	t.Helper()

	var resp appHTTP.FileEnvelope
	admin.Multipart(t, http.MethodPost, "/files/upload", nil, &testkit.UploadFile{
		FieldName: "file",
		Filename:  filename,
		Data:      content,
	}, http.StatusOK, &resp)
	assertOKEnvelope(t, resp.Status, resp.Msg)
	return resp.Data
}

func containsUser(items []appHTTP.User, username string) bool {
	for _, item := range items {
		if item.Username == username {
			return true
		}
	}
	return false
}

func containsFile(items []appHTTP.File, originalName string) bool {
	for _, item := range items {
		if item.OriginalName == originalName {
			return true
		}
	}
	return false
}

func containsOption(items []appHTTP.Option, key string) bool {
	return findOption(items, key) != nil
}

func findOption(items []appHTTP.Option, key string) *appHTTP.Option {
	for i := range items {
		if items[i].OptionKey == key {
			return &items[i]
		}
	}
	return nil
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
