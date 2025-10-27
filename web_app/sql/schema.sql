-- Blossom 论坛数据库表结构

-- ========== 用户表 ==========
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT NOT NULL COMMENT '用户ID (使用雪花算法生成)',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ========== 话题表 ==========
CREATE TABLE IF NOT EXISTS `topics` (
    `id` BIGINT NOT NULL COMMENT '话题ID (使用雪花算法生成)',
    `user_id` BIGINT NOT NULL COMMENT '发布者ID',
    `title` VARCHAR(200) NOT NULL COMMENT '话题标题',
    `content` TEXT NOT NULL COMMENT '话题内容',
    `category` VARCHAR(20) NOT NULL COMMENT '分类：tech/design/discuss/share/product',
    `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    `dislike_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点踩数',
    `comment_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数',
    `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_category` (`category`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_like_count` (`like_count`),
    CONSTRAINT `fk_topics_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='话题表';

-- ========== 投票表 ==========
CREATE TABLE IF NOT EXISTS `votes` (
    `id` BIGINT NOT NULL COMMENT '投票ID (使用雪花算法生成)',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `topic_id` BIGINT NOT NULL COMMENT '话题ID',
    `vote_type` TINYINT NOT NULL COMMENT '投票类型：1=点赞，-1=点踩',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_topic` (`user_id`, `topic_id`),
    KEY `idx_topic_id` (`topic_id`),
    CONSTRAINT `fk_votes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_votes_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='投票表';

-- ========== 评论表 ==========
CREATE TABLE IF NOT EXISTS `comments` (
    `id` BIGINT NOT NULL COMMENT '评论ID (使用雪花算法生成)',
    `topic_id` BIGINT NOT NULL COMMENT '话题ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `parent_id` BIGINT DEFAULT NULL COMMENT '父评论ID（用于回复）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_topic_id` (`topic_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    CONSTRAINT `fk_comments_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comments_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ========== 插入测试数据 ==========
-- 注意：由于使用雪花算法生成ID，测试数据需要通过应用程序API插入
-- 或手动指定有效的雪花算法ID

-- 示例：手动插入测试用户（需要使用有效的雪花算法ID）
-- INSERT INTO `users` (`id`, `username`, `email`, `password`) VALUES
-- (1000000000000001, '技术极客', 'tech@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbxqYYGabdz0XN3o6'),
-- (1000000000000002, '设计师小王', 'design@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbxqYYGabdz0XN3o6'),
-- (1000000000000003, '前端开发者', 'frontend@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lbxqYYGabdz0XN3o6');

