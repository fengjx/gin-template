ALTER TABLE sys_options ADD COLUMN type TEXT NOT NULL DEFAULT 'string';
ALTER TABLE sys_options ADD COLUMN status TEXT NOT NULL DEFAULT 'online';

UPDATE sys_options
SET type = 'string'
WHERE TRIM(type) = '';

UPDATE sys_options
SET status = 'online'
WHERE TRIM(status) = '';

UPDATE sys_schema_info
SET schema_version = 3,
    utime = CURRENT_TIMESTAMP
WHERE schema_version < 3;
