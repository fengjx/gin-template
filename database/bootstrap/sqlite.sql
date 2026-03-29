CREATE TABLE IF NOT EXISTS sys_schema_info (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  schema_version INTEGER NOT NULL,
  initialized_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sys_users (
  uid INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL,
  status TEXT NOT NULL,
  display_name TEXT NOT NULL,
  email_verified INTEGER NOT NULL DEFAULT 0,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sys_refresh_tokens (
  id TEXT PRIMARY KEY,
  uid INTEGER NOT NULL,
  token_hash TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  revoked_at DATETIME,
  user_agent TEXT,
  client_ip TEXT,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_options (
  id TEXT PRIMARY KEY,
  option_key TEXT NOT NULL UNIQUE,
  option_value TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  is_public INTEGER NOT NULL DEFAULT 0,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sys_oauth_bindings (
  id TEXT PRIMARY KEY,
  uid INTEGER NOT NULL,
  provider TEXT NOT NULL,
  provider_user_id TEXT NOT NULL,
  provider_username TEXT NOT NULL DEFAULT '',
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(provider, provider_user_id),
  FOREIGN KEY(uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_email_verifications (
  id TEXT PRIMARY KEY,
  uid INTEGER NOT NULL,
  email TEXT NOT NULL,
  token_hash TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  used_at DATETIME,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_password_resets (
  id TEXT PRIMARY KEY,
  uid INTEGER NOT NULL,
  email TEXT NOT NULL,
  token_hash TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  used_at DATETIME,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(uid) REFERENCES sys_users(uid)
);

CREATE TABLE IF NOT EXISTS sys_files (
  id TEXT PRIMARY KEY,
  uid INTEGER NOT NULL,
  original_name TEXT NOT NULL,
  storage_name TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size INTEGER NOT NULL,
  path TEXT NOT NULL,
  utime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(uid) REFERENCES sys_users(uid)
);

CREATE TRIGGER IF NOT EXISTS trg_sys_schema_info_utime
AFTER UPDATE ON sys_schema_info
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_schema_info
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_users_utime
AFTER UPDATE ON sys_users
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_users
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE uid = OLD.uid;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_refresh_tokens_utime
AFTER UPDATE ON sys_refresh_tokens
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_refresh_tokens
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_options_utime
AFTER UPDATE ON sys_options
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_options
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_oauth_bindings_utime
AFTER UPDATE ON sys_oauth_bindings
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_oauth_bindings
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_email_verifications_utime
AFTER UPDATE ON sys_email_verifications
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_email_verifications
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_password_resets_utime
AFTER UPDATE ON sys_password_resets
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_password_resets
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_sys_files_utime
AFTER UPDATE ON sys_files
FOR EACH ROW
WHEN NEW.utime = OLD.utime
BEGIN
  UPDATE sys_files
  SET utime = STRFTIME('%Y-%m-%d %H:%M:%f', 'now')
  WHERE id = OLD.id;
END;

INSERT INTO sys_schema_info (schema_version, initialized_at, utime, ctime)
SELECT 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM sys_schema_info);

INSERT INTO sys_options (id, option_key, option_value, description, is_public, utime, ctime)
VALUES
  ('notice', 'notice', '欢迎使用 gin-template', '系统公告', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('about', 'about', 'Gin + React 同构脚手架', '关于信息', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('pprof_url', 'pprof_url', '/debug/pprof/', 'pprof 监控地址', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(option_key) DO NOTHING;
