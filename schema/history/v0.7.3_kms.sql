ALTER TABLE customer_secret_key DROP INDEX idx_global_name;
ALTER TABLE customer_secret_key ADD UNIQUE idx_global_key (service,key_id) USING BTREE;