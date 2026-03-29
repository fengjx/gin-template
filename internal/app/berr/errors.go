package berr

import "net/http"

const (
	StatusOK = 0
)

const (
	// 系统级状态码：100000-199999。
	StatusSystemBase          = 100000
	StatusBadRequest          = 100001
	StatusUnauthorized        = 100002
	StatusForbidden           = 100003
	StatusNotFound            = 100004
	StatusConflict            = 100005
	StatusTooManyRequests     = 100006
	StatusInternalServerError = 100007

	StatusInvalidRequest           = 100101
	StatusInvalidToken             = 100102
	StatusRequireLogin             = 100103
	StatusRequireAdmin             = 100104
	StatusRequireRoot              = 100105
	StatusOpenAPISpecUnavailable   = 100201
	StatusFrontendAssetUnavailable = 100202
)

var (
	// 系统级错误定义。
	ErrInternalServerError      = New(http.StatusInternalServerError, StatusInternalServerError, "服务内部错误")
	ErrResourceNotFound         = New(http.StatusNotFound, StatusNotFound, "资源不存在")
	ErrTooManyRequests          = New(http.StatusTooManyRequests, StatusTooManyRequests, "请求过于频繁")
	ErrInvalidRequest           = New(http.StatusBadRequest, StatusInvalidRequest, "请求参数无效")
	ErrInvalidToken             = New(http.StatusUnauthorized, StatusInvalidToken, "令牌无效")
	ErrRequireLogin             = New(http.StatusUnauthorized, StatusRequireLogin, "请先登录")
	ErrRequireAdmin             = New(http.StatusForbidden, StatusRequireAdmin, "需要管理员权限")
	ErrRequireRoot              = New(http.StatusForbidden, StatusRequireRoot, "需要 root 权限")
	ErrOpenAPISpecUnavailable   = New(http.StatusInternalServerError, StatusOpenAPISpecUnavailable, "OpenAPI 规范不可用")
	ErrFrontendAssetUnavailable = New(http.StatusInternalServerError, StatusFrontendAssetUnavailable, "前端资源不可用")
)



const (
	// Auth 业务状态码：200000-209999。
	StatusAuthBase                           = 200000
	StatusTurnstileVerifyFailed              = 200001
	StatusAuthFieldsRequired                 = 200002
	StatusInvalidPasswordResetToken          = 200003
	StatusInvalidEmailVerificationToken      = 200004
	StatusInvalidCredentials                 = 200005
	StatusUserDisabled                       = 200006
	StatusCreatePasswordResetTokenFailed     = 200007
	StatusPasswordUpdateFailed               = 200008
	StatusCreateEmailVerificationTokenFailed = 200009
	StatusUpdateEmailVerificationFailed      = 200010
	StatusGitHubOAuthDisabled                = 200011
	StatusWeChatOAuthDisabled                = 200012
)

var (
	// Auth 业务错误定义。
	ErrTurnstileVerifyFailed              = New(http.StatusBadRequest, StatusTurnstileVerifyFailed, "人机校验失败")
	ErrAuthFieldsRequired                 = New(http.StatusBadRequest, StatusAuthFieldsRequired, "用户名、邮箱、密码不能为空")
	ErrInvalidPasswordResetToken          = New(http.StatusBadRequest, StatusInvalidPasswordResetToken, "重置令牌无效")
	ErrInvalidEmailVerificationToken      = New(http.StatusBadRequest, StatusInvalidEmailVerificationToken, "验证令牌无效")
	ErrInvalidCredentials                 = New(http.StatusUnauthorized, StatusInvalidCredentials, "账号或密码错误")
	ErrUserDisabled                       = New(http.StatusForbidden, StatusUserDisabled, "用户已禁用")
	ErrGitHubOAuthDisabled                = New(http.StatusBadRequest, StatusGitHubOAuthDisabled, "GitHub OAuth 未启用")
	ErrWeChatOAuthDisabled                = New(http.StatusBadRequest, StatusWeChatOAuthDisabled, "微信 OAuth 未启用")
	ErrCreatePasswordResetTokenFailed     = New(http.StatusInternalServerError, StatusCreatePasswordResetTokenFailed, "创建重置令牌失败")
	ErrPasswordUpdateFailed               = New(http.StatusInternalServerError, StatusPasswordUpdateFailed, "密码更新失败")
	ErrCreateEmailVerificationTokenFailed = New(http.StatusInternalServerError, StatusCreateEmailVerificationTokenFailed, "创建验证令牌失败")
	ErrUpdateEmailVerificationFailed      = New(http.StatusInternalServerError, StatusUpdateEmailVerificationFailed, "更新邮箱状态失败")
)

const (
	// User 业务状态码：210000-219999。
	StatusUserBase                  = 210000
	StatusInvalidUserID             = 210001
	StatusInvalidRole               = 210002
	StatusInvalidUserStatus         = 210003
	StatusEmailRequired             = 210004
	StatusCannotDeleteCurrentUser   = 210005
	StatusUserNotFound              = 210006
	StatusUsernameExists            = 210007
	StatusEmailExists               = 210008
	StatusOnlyRootCanCreateRootUser = 210009
	StatusOnlyRootCanModifyRootUser = 210010
	StatusOnlyRootCanGrantRootRole  = 210011
	StatusOnlyRootCanDeleteRootUser = 210012
	StatusUpdateProfileFailed       = 210013
	StatusQueryUsersFailed          = 210014
	StatusCheckUsernameFailed       = 210015
	StatusCheckEmailFailed          = 210016
	StatusPasswordProcessFailed     = 210017
	StatusCreateUserFailed          = 210018
	StatusUpdateUserFailed          = 210019
	StatusDeleteUserFailed          = 210020
)

var (
	// User 业务错误定义。
	ErrInvalidUserID             = New(http.StatusBadRequest, StatusInvalidUserID, "用户标识无效")
	ErrInvalidRole               = New(http.StatusBadRequest, StatusInvalidRole, "角色无效")
	ErrInvalidUserStatus         = New(http.StatusBadRequest, StatusInvalidUserStatus, "状态无效")
	ErrEmailRequired             = New(http.StatusBadRequest, StatusEmailRequired, "邮箱不能为空")
	ErrCannotDeleteCurrentUser   = New(http.StatusBadRequest, StatusCannotDeleteCurrentUser, "不能删除当前登录用户")
	ErrUserNotFound              = New(http.StatusNotFound, StatusUserNotFound, "用户不存在")
	ErrUsernameExists            = New(http.StatusConflict, StatusUsernameExists, "用户名已存在")
	ErrEmailExists               = New(http.StatusConflict, StatusEmailExists, "邮箱已存在")
	ErrOnlyRootCanCreateRootUser = New(http.StatusForbidden, StatusOnlyRootCanCreateRootUser, "仅 root 可创建 root 用户")
	ErrOnlyRootCanModifyRootUser = New(http.StatusForbidden, StatusOnlyRootCanModifyRootUser, "仅 root 可修改 root 用户")
	ErrOnlyRootCanGrantRootRole  = New(http.StatusForbidden, StatusOnlyRootCanGrantRootRole, "仅 root 可授予 root 权限")
	ErrOnlyRootCanDeleteRootUser = New(http.StatusForbidden, StatusOnlyRootCanDeleteRootUser, "仅 root 可删除 root 用户")
	ErrUpdateProfileFailed       = New(http.StatusInternalServerError, StatusUpdateProfileFailed, "更新个人信息失败")
	ErrQueryUsersFailed          = New(http.StatusInternalServerError, StatusQueryUsersFailed, "查询用户失败")
	ErrCheckUsernameFailed       = New(http.StatusInternalServerError, StatusCheckUsernameFailed, "检查用户名失败")
	ErrCheckEmailFailed          = New(http.StatusInternalServerError, StatusCheckEmailFailed, "检查邮箱失败")
	ErrPasswordProcessFailed     = New(http.StatusInternalServerError, StatusPasswordProcessFailed, "密码处理失败")
	ErrCreateUserFailed          = New(http.StatusInternalServerError, StatusCreateUserFailed, "创建用户失败")
	ErrUpdateUserFailed          = New(http.StatusInternalServerError, StatusUpdateUserFailed, "更新用户失败")
	ErrDeleteUserFailed          = New(http.StatusInternalServerError, StatusDeleteUserFailed, "删除用户失败")
)

const (
	// File 业务状态码：220000-229999。
	StatusFileBase               = 220000
	StatusMissingUploadFile      = 220001
	StatusFileNotFound           = 220002
	StatusQueryFilesFailed       = 220003
	StatusOpenUploadFileFailed   = 220004
	StatusCreateUploadDirFailed  = 220005
	StatusCreateTargetFileFailed = 220006
	StatusWriteTargetFileFailed  = 220007
	StatusSaveFileRecordFailed   = 220008
	StatusUploadFileFailed       = 220009
	StatusDeleteFileFailed       = 220010
)


var (
	// File 业务错误定义。
	ErrMissingUploadFile      = New(http.StatusBadRequest, StatusMissingUploadFile, "缺少上传文件")
	ErrFileNotFound           = New(http.StatusNotFound, StatusFileNotFound, "文件不存在")
	ErrQueryFilesFailed       = New(http.StatusInternalServerError, StatusQueryFilesFailed, "查询文件失败")
	ErrOpenUploadFileFailed   = New(http.StatusInternalServerError, StatusOpenUploadFileFailed, "打开上传文件失败")
	ErrCreateUploadDirFailed  = New(http.StatusInternalServerError, StatusCreateUploadDirFailed, "创建上传目录失败")
	ErrCreateTargetFileFailed = New(http.StatusInternalServerError, StatusCreateTargetFileFailed, "创建目标文件失败")
	ErrWriteTargetFileFailed  = New(http.StatusInternalServerError, StatusWriteTargetFileFailed, "写入文件失败")
	ErrSaveFileRecordFailed   = New(http.StatusInternalServerError, StatusSaveFileRecordFailed, "保存文件记录失败")
	ErrUploadFileFailed       = New(http.StatusInternalServerError, StatusUploadFileFailed, "上传文件失败")
	ErrDeleteFileFailed       = New(http.StatusInternalServerError, StatusDeleteFileFailed, "删除文件失败")
)


const (
	// Option 业务状态码：230000-239999。
	StatusOptionBase          = 230000
	StatusAboutNotFound       = 230001
	StatusNoticeNotFound      = 230002
	StatusOptionNotFound      = 230003
	StatusGetOptionsFailed    = 230004
	StatusUpdateOptionFailed  = 230005
	StatusPprofURLUnavailable = 230006
)


var (
	// Option 业务错误定义。
	ErrAboutNotFound       = New(http.StatusNotFound, StatusAboutNotFound, "about 不存在")
	ErrNoticeNotFound      = New(http.StatusNotFound, StatusNoticeNotFound, "notice 不存在")
	ErrOptionNotFound      = New(http.StatusNotFound, StatusOptionNotFound, "配置项不存在")
	ErrGetOptionsFailed    = New(http.StatusInternalServerError, StatusGetOptionsFailed, "获取配置失败")
	ErrUpdateOptionFailed  = New(http.StatusInternalServerError, StatusUpdateOptionFailed, "更新配置失败")
	ErrPprofURLUnavailable = New(http.StatusNotFound, StatusPprofURLUnavailable, "Pprof 地址不可用")
)
