SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- PXX 二手交易平台 数据库初始化
CREATE DATABASE IF NOT EXISTS pxx DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE pxx;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    open_id     VARCHAR(64) NOT NULL DEFAULT '',
    nickname    VARCHAR(64) NOT NULL DEFAULT '',
    avatar      VARCHAR(512) NOT NULL DEFAULT '',
    phone       VARCHAR(20) NOT NULL DEFAULT '',
    campus      VARCHAR(64) NOT NULL DEFAULT '',
    department  VARCHAR(64) NOT NULL DEFAULT '',
    dormitory   VARCHAR(64) NOT NULL DEFAULT '',
    is_verified TINYINT(1) NOT NULL DEFAULT 0,
    balance     BIGINT NOT NULL DEFAULT 0,
    role        VARCHAR(20) NOT NULL DEFAULT 'user',
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_open_id (open_id)
) ENGINE=InnoDB;

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
    view_count     INT NOT NULL DEFAULT 0,
    like_count     INT NOT NULL DEFAULT 0,
    status         TINYINT NOT NULL DEFAULT 1,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_seller (seller_id),
    INDEX idx_category (category),
    INDEX idx_campus (campus)
) ENGINE=InnoDB;

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no      VARCHAR(32) NOT NULL DEFAULT '',
    buyer_id      BIGINT UNSIGNED NOT NULL,
    seller_id     BIGINT UNSIGNED NOT NULL,
    product_id    BIGINT UNSIGNED NOT NULL,
    amount        BIGINT NOT NULL DEFAULT 0,
    status        TINYINT NOT NULL DEFAULT 1,
    pay_method    VARCHAR(16) NOT NULL DEFAULT '',
    paid_at       DATETIME NULL,
    delivered_at  DATETIME NULL,
    completed_at  DATETIME NULL,
    cancelled_at  DATETIME NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_order_no (order_no),
    INDEX idx_buyer (buyer_id),
    INDEX idx_seller (seller_id)
) ENGINE=InnoDB;

-- 帖子表
CREATE TABLE IF NOT EXISTS posts (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    type       TINYINT NOT NULL DEFAULT 1,
    title      VARCHAR(128) NOT NULL,
    content    TEXT,
    images     TEXT,
    tags       VARCHAR(256) NOT NULL DEFAULT '',
    product_id BIGINT UNSIGNED NULL,
    view_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    status     TINYINT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_type (type)
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

-- 拼团订单表
CREATE TABLE IF NOT EXISTS group_buy_orders (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    group_buy_id BIGINT UNSIGNED NOT NULL,
    user_id      BIGINT UNSIGNED NOT NULL,
    is_leader    TINYINT(1) NOT NULL DEFAULT 0,
    quantity     INT NOT NULL DEFAULT 1,
    amount       BIGINT NOT NULL DEFAULT 0,
    status       TINYINT NOT NULL DEFAULT 1,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_group (group_buy_id),
    INDEX idx_user (user_id)
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

-- 租赁订单表
CREATE TABLE IF NOT EXISTS rental_orders (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    rental_id    BIGINT UNSIGNED NOT NULL,
    renter_id    BIGINT UNSIGNED NOT NULL,
    start_date   DATETIME NOT NULL,
    end_date     DATETIME NOT NULL,
    total_amount BIGINT NOT NULL DEFAULT 0,
    deposit      BIGINT NOT NULL DEFAULT 0,
    status       TINYINT NOT NULL DEFAULT 1,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rental (rental_id),
    INDEX idx_renter (renter_id)
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

-- 置换订单表
CREATE TABLE IF NOT EXISTS barter_orders (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no   VARCHAR(32) NOT NULL DEFAULT '',
    user_a_id  BIGINT UNSIGNED NOT NULL,
    user_b_id  BIGINT UNSIGNED NOT NULL,
    item_a_id  BIGINT UNSIGNED NOT NULL,
    item_b_id  BIGINT UNSIGNED NOT NULL,
    status     TINYINT NOT NULL DEFAULT 1,
    meet_place VARCHAR(128) NOT NULL DEFAULT '',
    meet_time  DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_order_no (order_no)
) ENGINE=InnoDB;

-- 用户行为埋点表
CREATE TABLE IF NOT EXISTS user_behaviors (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    target_id   BIGINT UNSIGNED NOT NULL,
    target_type VARCHAR(32) NOT NULL DEFAULT '',
    action      VARCHAR(16) NOT NULL DEFAULT '',
    duration    INT NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_action (user_id, action),
    INDEX idx_target (target_type, target_id)
) ENGINE=InnoDB;
