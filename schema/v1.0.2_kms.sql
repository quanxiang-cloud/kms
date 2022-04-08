DROP TABLE IF EXISTS customer_auth_method;
CREATE TABLE `customer_auth_method` (
    `id`            VARCHAR(64)                 COMMENT 'unique id',
    `owner`         VARCHAR(64)                 COMMENT 'owner id',
    `owner_name`    VARCHAR(64)                 COMMENT 'owner name',
    `name`          VARCHAR(64)     NOT NULL    COMMENT 'name of method',
    `title`         VARCHAR(64)     NOT NULL    COMMENT 'title of method',
    `description`   VARCHAR(255),
    `image`         VARCHAR(255)    NOT NULL    COMMENT 'path of docker image',
    `create_at`     BIGINT(20)                  COMMENT 'create time',
    `update_at`     BIGINT(20)                  COMMENT 'update time',
    `delete_at`     BIGINT(20)                  COMMENT 'delete time',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `customer_auth_path`;
CREATE TABLE `customer_auth_path` (
    `id`            VARCHAR(64)                 COMMENT 'unique id',
    `method_id`     VARCHAR(64)     NOT NULL    COMMENT 'method_id',
    `namespace`     VARCHAR(64)     NOT NULL    COMMENT 'namespace path',
    `create_at`     BIGINT(20)                  COMMENT 'create time',
    `update_at`     BIGINT(20)                  COMMENT 'update time',
    `delete_at`     BIGINT(20)                  COMMENT 'delete time',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;