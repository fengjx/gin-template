ALTER TABLE sys_options
  ADD COLUMN type VARCHAR(32) NOT NULL DEFAULT 'string' AFTER is_public,
  ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'online' AFTER type;

UPDATE sys_options
SET type = 'string'
WHERE type = '';

UPDATE sys_options
SET status = 'online'
WHERE status = '';

UPDATE sys_schema_info
SET schema_version = 3,
    utime = NOW()
WHERE schema_version < 3;
