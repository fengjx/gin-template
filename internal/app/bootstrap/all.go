package bootstrap

import (
	// 注册认证相关路由与初始化逻辑。
	_ "gin-template/internal/biz/auth"
	// 注册邮件相关依赖。
	_ "gin-template/internal/biz/email"
	// 注册文件管理模块。
	_ "gin-template/internal/biz/file"
	// 注册 GitHub OAuth 模块。
	_ "gin-template/internal/biz/github"
	// 注册系统配置与状态模块。
	_ "gin-template/internal/biz/option"
	// 注册 Turnstile 能力。
	_ "gin-template/internal/biz/turnstile"
	// 注册用户管理模块。
	_ "gin-template/internal/biz/user"
	// 注册微信 OAuth 模块。
	_ "gin-template/internal/biz/wechat"
)
