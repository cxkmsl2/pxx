-- PXX 商品库存服务数据库
CREATE DATABASE IF NOT EXISTS pxx_item DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE pxx_item;

-- 商品表
CREATE TABLE IF NOT EXISTS products (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    seller_id      BIGINT UNSIGNED NOT NULL,
    title          VARCHAR(128) NOT NULL,
    `desc`         TEXT,
    category       VARCHAR(32) NOT NULL DEFAULT '',
    tag            VARCHAR(32) NOT NULL DEFAULT '',
    price          BIGINT NOT NULL DEFAULT 0,
    original_price BIGINT NOT NULL DEFAULT 0,
    images         TEXT,
    `condition`    TINYINT NOT NULL DEFAULT 1,
    campus         VARCHAR(64) NOT NULL DEFAULT '',
    dormitory      VARCHAR(64) NOT NULL DEFAULT '',
    stock          INT NOT NULL DEFAULT 1 COMMENT '库存数量',
    view_count     INT NOT NULL DEFAULT 0,
    like_count     INT NOT NULL DEFAULT 0,
    status         TINYINT NOT NULL DEFAULT 1 COMMENT '1上架 2锁定 0下架',
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_seller (seller_id),
    INDEX idx_category (category),
    INDEX idx_campus (campus)
) ENGINE=InnoDB;

-- 租赁物品表
CREATE TABLE IF NOT EXISTS rental_items (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    owner_id     BIGINT UNSIGNED NOT NULL,
    title        VARCHAR(128) NOT NULL,
    `desc`       TEXT,
    cover_image  VARCHAR(512) NOT NULL DEFAULT '',
    images       TEXT,
    daily_price  BIGINT NOT NULL DEFAULT 0,
    weekly_price BIGINT NOT NULL DEFAULT 0,
    deposit      BIGINT NOT NULL DEFAULT 0,
    campus       VARCHAR(64) NOT NULL DEFAULT '',
    status       TINYINT NOT NULL DEFAULT 1,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_id)
) ENGINE=InnoDB;

-- 以物换物表
CREATE TABLE IF NOT EXISTS barter_items (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    title      VARCHAR(128) NOT NULL,
    `desc`     TEXT,
    images     TEXT,
    want_item  VARCHAR(256) NOT NULL DEFAULT '',
    status     TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id)
) ENGINE=InnoDB;

-- 订阅账号表
CREATE TABLE IF NOT EXISTS subscriptions (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    owner_id       BIGINT UNSIGNED NOT NULL,
    brand          VARCHAR(64) NOT NULL DEFAULT '',
    title          VARCHAR(128) NOT NULL,
    `desc`         TEXT,
    account_cipher TEXT COMMENT 'AES加密的账号',
    pass_cipher    TEXT COMMENT 'AES加密的密码',
    price_per_day  BIGINT NOT NULL DEFAULT 0,
    price_per_week BIGINT NOT NULL DEFAULT 0,
    status         TINYINT NOT NULL DEFAULT 1 COMMENT '1空闲 2使用中 3过期',
    rented_by      BIGINT UNSIGNED NULL,
    rent_start     DATETIME NULL,
    rent_end       DATETIME NULL,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner (owner_id),
    INDEX idx_brand (brand),
    INDEX idx_status (status)
) ENGINE=InnoDB;

-- 拼团表
CREATE TABLE IF NOT EXISTS group_buys (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title          VARCHAR(128) NOT NULL,
    `desc`         TEXT,
    cover_image    VARCHAR(512) NOT NULL DEFAULT '',
    price_per_unit BIGINT NOT NULL DEFAULT 0,
    min_people     INT NOT NULL DEFAULT 2,
    current_people INT NOT NULL DEFAULT 0,
    total_stock    INT NOT NULL DEFAULT 0,
    sold_count     INT NOT NULL DEFAULT 0,
    start_at       DATETIME NOT NULL,
    end_at         DATETIME NOT NULL,
    status         TINYINT NOT NULL DEFAULT 1,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- TCC 控制表
CREATE TABLE IF NOT EXISTS tcc_logs (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tx_id       VARCHAR(64) NOT NULL,
    branch_id   BIGINT UNSIGNED NOT NULL,
    action_type VARCHAR(16) NOT NULL COMMENT 'LOCK_STOCK/DEDUCT_STOCK/UNLOCK_STOCK',
    status      VARCHAR(16) NOT NULL COMMENT 'TRYING/LOCKED/CONFIRMED/CANCELLED/FAILED',
    product_id  BIGINT UNSIGNED NOT NULL DEFAULT 0,
    payload     JSON,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_tx_branch (tx_id, branch_id)
) ENGINE=InnoDB;
