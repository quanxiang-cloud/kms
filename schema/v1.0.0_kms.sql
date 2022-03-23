ALTER TABLE secret_key MODIFY COLUMN create_at BIGINT(20) COMMENT 'create time';
ALTER TABLE secret_key MODIFY COLUMN update_at BIGINT(20) COMMENT 'update time';
ALTER TABLE secret_key MODIFY COLUMN delete_at BIGINT(20) COMMENT 'delete time';

ALTER TABLE customer_secret_key MODIFY COLUMN create_at BIGINT(20) COMMENT 'create time';
ALTER TABLE customer_secret_key MODIFY COLUMN update_at BIGINT(20) COMMENT 'update time';
ALTER TABLE customer_secret_key MODIFY COLUMN delete_at BIGINT(20) COMMENT 'delete time';
ALTER TABLE customer_secret_key ADD COLUMN parsed INT(11) NOT NULL COMMENT '1 parsed 0 not parse' AFTER active;
-- ALTER TABLE customer_secret_key DROP INDEX idx_global_name;
-- ALTER TABLE customer_secret_key ADD UNIQUE idx_global_key (`service`,`key_id`);

ALTER TABLE secret_key_config DROP COLUMN key_num;
ALTER TABLE secret_key_config DROP COLUMN expiry;
ALTER TABLE secret_key_config DROP COLUMN expire_at;
ALTER TABLE secret_key_config ADD COLUMN config_content VARCHAR(64) COMMENT 'config content' AFTER owner_name;
ALTER TABLE secret_key_config MODIFY COLUMN create_at BIGINT(20) COMMENT 'create time' AFTER config_content;
ALTER TABLE secret_key_config ADD COLUMN update_at BIGINT(20) COMMENT 'update time' AFTER create_at;
ALTER TABLE secret_key_config ADD COLUMN delete_at BIGINT(20) COMMENT 'delete time' AFTER update_at;
REPLACE INTO `secret_key_config` VALUES ('1', 'system', '系统', '{"keyNum": "5"}', unix_timestamp(NOW())*1000, unix_timestamp(NOW())*1000, NULL)

