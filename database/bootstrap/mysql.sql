CREATE TABLE IF NOT EXISTS sys_schema_info (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  schema_version BIGINT NOT NULL,
  initialized_at DATETIME NOT NULL DEFAULT NOW(),
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间'
);

CREATE TABLE IF NOT EXISTS sys_users (
  uid BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(128) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  display_name VARCHAR(128) NOT NULL,
  email_verified TINYINT(1) NOT NULL DEFAULT 0,
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间'
);

CREATE TABLE IF NOT EXISTS sys_refresh_tokens (
  id VARCHAR(36) PRIMARY KEY,
  uid BIGINT NOT NULL,
  token_hash VARCHAR(255) NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked_at DATETIME NULL,
  user_agent VARCHAR(255) NULL,
  client_ip VARCHAR(64) NULL,
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  CONSTRAINT fk_sys_refresh_tokens_user FOREIGN KEY (uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_options (
  id VARCHAR(64) PRIMARY KEY,
  option_key VARCHAR(128) NOT NULL UNIQUE,
  option_value TEXT NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  is_public TINYINT(1) NOT NULL DEFAULT 0,
  type VARCHAR(32) NOT NULL DEFAULT 'string',
  status VARCHAR(32) NOT NULL DEFAULT 'online',
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间'
);

CREATE TABLE IF NOT EXISTS sys_oauth_bindings (
  id VARCHAR(36) PRIMARY KEY,
  uid BIGINT NOT NULL,
  provider VARCHAR(32) NOT NULL,
  provider_user_id VARCHAR(128) NOT NULL,
  provider_username VARCHAR(128) NOT NULL DEFAULT '',
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  UNIQUE KEY uniq_sys_oauth_bindings_provider (provider, provider_user_id),
  CONSTRAINT fk_sys_oauth_bindings_user FOREIGN KEY (uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_email_verifications (
  id VARCHAR(36) PRIMARY KEY,
  uid BIGINT NOT NULL,
  email VARCHAR(128) NOT NULL,
  token_hash VARCHAR(255) NOT NULL,
  expires_at DATETIME NOT NULL,
  used_at DATETIME NULL,
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  CONSTRAINT fk_sys_email_verifications_user FOREIGN KEY (uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_password_resets (
  id VARCHAR(36) PRIMARY KEY,
  uid BIGINT NOT NULL,
  email VARCHAR(128) NOT NULL,
  token_hash VARCHAR(255) NOT NULL,
  expires_at DATETIME NOT NULL,
  used_at DATETIME NULL,
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  CONSTRAINT fk_sys_password_resets_user FOREIGN KEY (uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_files (
  id VARCHAR(36) PRIMARY KEY,
  uid BIGINT NOT NULL,
  original_name VARCHAR(255) NOT NULL,
  storage_name VARCHAR(255) NOT NULL,
  content_type VARCHAR(128) NOT NULL,
  size BIGINT NOT NULL,
  path VARCHAR(1024) NOT NULL,
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  CONSTRAINT fk_sys_files_user FOREIGN KEY (uid) REFERENCES sys_users(uid)
);

INSERT INTO sys_schema_info (schema_version, initialized_at, utime, ctime)
SELECT 3, NOW(), NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_schema_info);

INSERT INTO sys_options (id, option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'notice', 'notice', '欢迎使用 gin-template', '系统公告', 1, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'notice');

INSERT INTO sys_options (id, option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'about', 'about', 'Gin + React 同构脚手架', '关于信息', 1, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'about');

INSERT INTO sys_options (id, option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'pprof_url', 'pprof_url', '/debug/pprof/', 'pprof 监控地址', 0, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'pprof_url');
