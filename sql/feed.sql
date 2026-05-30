-- PXX 社区信息服务数据库
CREATE DATABASE IF NOT EXISTS pxx_feed DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE pxx_feed;

-- 帖子表
CREATE TABLE IF NOT EXISTS posts (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    type       TINYINT NOT NULL DEFAULT 1 COMMENT '1求购悬赏 2好物种草',
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

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    from_user_id BIGINT UNSIGNED NOT NULL,
    to_user_id   BIGINT UNSIGNED NOT NULL,
    content     TEXT,
    msg_type    VARCHAR(16) NOT NULL DEFAULT 'text',
    read_at     DATETIME NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_from (from_user_id),
    INDEX idx_to (to_user_id)
) ENGINE=InnoDB;

-- 跑腿任务表
CREATE TABLE IF NOT EXISTS tasks (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    publisher_id BIGINT UNSIGNED NOT NULL,
    taker_id     BIGINT UNSIGNED NULL,
    type         VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'food/express/print/book',
    title        VARCHAR(128) NOT NULL,
    `desc`       TEXT,
    from_place   VARCHAR(128) NOT NULL DEFAULT '',
    to_place     VARCHAR(128) NOT NULL DEFAULT '',
    reward       BIGINT NOT NULL DEFAULT 0,
    status       TINYINT NOT NULL DEFAULT 1 COMMENT '1待接单 2已接单 3已完成 4已取消',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_publisher (publisher_id),
    INDEX idx_taker (taker_id),
    INDEX idx_type (type)
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

-- 仲裁表
CREATE TABLE IF NOT EXISTS dispute_cases (
    id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id         BIGINT UNSIGNED NOT NULL,
    plaintiff_id     BIGINT UNSIGNED NOT NULL,
    defendant_id     BIGINT UNSIGNED NOT NULL,
    title            VARCHAR(256) NOT NULL DEFAULT '',
    plaintiff_desc   TEXT,
    defendant_desc   TEXT,
    plaintiff_votes  INT NOT NULL DEFAULT 0,
    defendant_votes  INT NOT NULL DEFAULT 0,
    status           TINYINT NOT NULL DEFAULT 1 COMMENT '0待审核 1投票中 2已结案',
    winner           VARCHAR(16) NOT NULL DEFAULT '',
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order (order_id),
    INDEX idx_plaintiff (plaintiff_id),
    INDEX idx_defendant (defendant_id)
) ENGINE=InnoDB;

-- 仲裁投票表
CREATE TABLE IF NOT EXISTS case_votes (
    id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    case_id   BIGINT UNSIGNED NOT NULL,
    voter_id  BIGINT UNSIGNED NOT NULL,
    vote_for  VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'plaintiff/defendant',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_case (case_id),
    INDEX idx_voter (voter_id)
) ENGINE=InnoDB;
