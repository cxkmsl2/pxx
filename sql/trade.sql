-- PXX 交易服务数据库
CREATE DATABASE IF NOT EXISTS pxx_trade DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE pxx_trade;

-- 商品订单表
CREATE TABLE IF NOT EXISTS orders (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no      VARCHAR(32) NOT NULL DEFAULT '',
    buyer_id      BIGINT UNSIGNED NOT NULL,
    seller_id     BIGINT UNSIGNED NOT NULL,
    product_id    BIGINT UNSIGNED NOT NULL,
    amount        BIGINT NOT NULL DEFAULT 0,
    status        TINYINT NOT NULL DEFAULT 1 COMMENT '1待支付 2已支付 3已发货 4已完成 5已取消',
    tx_id         VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'TCC事务ID',
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

-- 拼团订单表
CREATE TABLE IF NOT EXISTS group_buy_orders (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    group_buy_id BIGINT UNSIGNED NOT NULL,
    user_id      BIGINT UNSIGNED NOT NULL,
    is_leader    TINYINT(1) NOT NULL DEFAULT 0,
    quantity     INT NOT NULL DEFAULT 1,
    amount       BIGINT NOT NULL DEFAULT 0,
    status       TINYINT NOT NULL DEFAULT 1,
    tx_id        VARCHAR(64) NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_group (group_buy_id),
    INDEX idx_user (user_id)
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
    tx_id        VARCHAR(64) NOT NULL DEFAULT '',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rental (rental_id),
    INDEX idx_renter (renter_id)
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
    tx_id      VARCHAR(64) NOT NULL DEFAULT '',
    meet_place VARCHAR(128) NOT NULL DEFAULT '',
    meet_time  DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_order_no (order_no)
) ENGINE=InnoDB;

-- 订阅订单表
CREATE TABLE IF NOT EXISTS sub_orders (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no       VARCHAR(32) NOT NULL DEFAULT '',
    subscription_id BIGINT UNSIGNED NOT NULL,
    buyer_id       BIGINT UNSIGNED NOT NULL,
    days           INT NOT NULL DEFAULT 0,
    amount         BIGINT NOT NULL DEFAULT 0,
    status         TINYINT NOT NULL DEFAULT 1 COMMENT '1租用中 2已完成 3已投诉',
    tx_id          VARCHAR(64) NOT NULL DEFAULT '',
    dispute_reason TEXT,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_order_no (order_no),
    INDEX idx_subscription (subscription_id)
) ENGINE=InnoDB;

-- TCC 全局事务表
CREATE TABLE IF NOT EXISTS tcc_transactions (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tx_id      VARCHAR(64) NOT NULL UNIQUE,
    status     VARCHAR(16) NOT NULL DEFAULT 'TRYING' COMMENT 'TRYING/CONFIRMED/CANCELLED',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;
