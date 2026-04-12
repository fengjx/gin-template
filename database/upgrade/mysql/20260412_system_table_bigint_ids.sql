DROP TABLE IF EXISTS sys_refresh_tokens;
DROP TABLE IF EXISTS sys_options;
DROP TABLE IF EXISTS sys_oauth_bindings;
DROP TABLE IF EXISTS sys_email_verifications;
DROP TABLE IF EXISTS sys_password_resets;
DROP TABLE IF EXISTS sys_files;

CREATE TABLE IF NOT EXISTS sys_refresh_tokens (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
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
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  option_key VARCHAR(128) NOT NULL,
  option_value TEXT NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  is_public TINYINT(1) NOT NULL DEFAULT 0,
  type VARCHAR(32) NOT NULL DEFAULT 'string',
  status VARCHAR(32) NOT NULL DEFAULT 'online',
  utime DATETIME NOT NULL DEFAULT NOW() ON UPDATE NOW() COMMENT '更新时间',
  ctime DATETIME NOT NULL DEFAULT NOW() COMMENT '创建时间',
  CONSTRAINT uk_o UNIQUE (option_key)
);

CREATE TABLE IF NOT EXISTS sys_oauth_bindings (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
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
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
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
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
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
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
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

INSERT INTO sys_options (option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'notice', '欢迎使用 gin-template', '系统公告', 1, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'notice');

INSERT INTO sys_options (option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'about', 'Gin + React 同构脚手架', '关于信息', 1, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'about');

INSERT INTO sys_options (option_key, option_value, description, is_public, type, status, utime, ctime)
SELECT 'pprof_url', '/debug/pprof/', 'pprof 监控地址', 0, 'string', 'online', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM sys_options WHERE option_key = 'pprof_url');

UPDATE sys_schema_info
SET schema_version = 4,
    utime = NOW()
WHERE schema_version < 4;
