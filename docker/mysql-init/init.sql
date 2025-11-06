-- 为了可重复执行，先按依赖关系安全删除
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `user_permanent_files`;
DROP TABLE IF EXISTS `user_files`;
DROP TABLE IF EXISTS `user_action_logs`;
DROP TABLE IF EXISTS `membership_orders`;
DROP TABLE IF EXISTS `membership_info`;
DROP TABLE IF EXISTS `chat_history`;
DROP TABLE IF EXISTS `users`;
SET FOREIGN_KEY_CHECKS = 1;

-- 1) 用户表（被引用端）
CREATE TABLE `users` (
  `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `full_name` VARCHAR(100) DEFAULT NULL,
  `email` VARCHAR(100) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `security_question1` VARCHAR(255) NOT NULL,
  `security_answer1_hash` VARCHAR(255) NOT NULL,
  `security_question2` VARCHAR(255) NOT NULL,
  `security_answer2_hash` VARCHAR(255) NOT NULL,
  `reset_token` VARCHAR(255) DEFAULT NULL,
  `reset_token_expires_at` TIMESTAMP NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `phone_number` (`phone_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2) 会话历史（外键列统一为 INT UNSIGNED）
CREATE TABLE `chat_history` (
  `chat_id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `record_id` VARCHAR(64) NOT NULL,
  `content_url` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`chat_id`),
  UNIQUE KEY `record_id` (`record_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `chat_history_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 3) 会员信息（已是 INT UNSIGNED，保持不变并补充 ON UPDATE CASCADE 以一致化）
CREATE TABLE `membership_info` (
  `membership_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `start_date` DATE NOT NULL,
  `expire_date` DATE NOT NULL,
  `status` ENUM('active','expired') DEFAULT 'active',
  `last_updated` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`membership_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `membership_info_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 4) 会员订单（外键列改为 INT UNSIGNED）
CREATE TABLE `membership_orders` (
  `order_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `purchase_date` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `duration_months` INT NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL,
  `payment_method` ENUM('alipay','wechat','card','other') DEFAULT 'other',
  PRIMARY KEY (`order_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `membership_orders_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 5) 用户操作日志（若需要和 users 关联，类型也统一；并加外键）
CREATE TABLE `user_action_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `action_type` VARCHAR(64) NOT NULL,
  `resource_type` VARCHAR(64) NOT NULL,
  `resource_id` VARCHAR(128) NOT NULL,
  `request_id` VARCHAR(128) DEFAULT NULL,
  `ip_addr` VARCHAR(64) DEFAULT NULL,
  `user_agent` VARCHAR(512) DEFAULT NULL,
  `extra_json` JSON DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`,`created_at`),
  KEY `idx_action_type` (`action_type`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  CONSTRAINT `user_action_logs_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 6) 临时用户文件（外键列改为 INT UNSIGNED）
CREATE TABLE `user_files` (
  `file_id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `filename` VARCHAR(255) NOT NULL,
  `filepath` VARCHAR(255) NOT NULL,
  `size` BIGINT NOT NULL,
  `uploaded_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`file_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `user_files_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 7) 永久用户文件（原来没有外键且 user_id 为 BIGINT；统一为 INT UNSIGNED 并加外键）
CREATE TABLE `user_permanent_files` (
  `file_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `filename` VARCHAR(255) NOT NULL,
  `filepath` TEXT NOT NULL,
  `size` BIGINT NOT NULL,
  `uploaded_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`file_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_time` (`uploaded_at`),
  CONSTRAINT `user_permanent_files_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `password_reset_tokens` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `token` VARCHAR(64) NOT NULL,
  `expire_at` DATETIME NOT NULL,
  `used` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_token` (`token`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `password_reset_tokens_ibfk_1`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
