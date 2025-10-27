-- 数据库迁移脚本：从AUTO_INCREMENT切换到雪花算法ID
-- 警告：此脚本会删除现有数据，请在执行前备份数据库！

-- ========== 1. 删除外键约束 ==========
ALTER TABLE `comments` DROP FOREIGN KEY `fk_comments_topic_id`;
ALTER TABLE `comments` DROP FOREIGN KEY `fk_comments_user_id`;
ALTER TABLE `topics` DROP FOREIGN KEY `fk_topics_user_id`;
ALTER TABLE `votes` DROP FOREIGN KEY `fk_votes_user_id`;
ALTER TABLE `votes` DROP FOREIGN KEY `fk_votes_topic_id`;

-- ========== 2. 清空所有表（可选，保留数据请跳过此步） ==========
-- 如果你想保留现有数据，请注释掉下面的TRUNCATE语句
-- 注意：保留数据需要手动迁移ID
TRUNCATE TABLE `comments`;
TRUNCATE TABLE `votes`;
TRUNCATE TABLE `topics`;
TRUNCATE TABLE `users`;

-- ========== 3. 修改ID列类型并移除AUTO_INCREMENT ==========

-- 修改 users 表
ALTER TABLE `users` 
    MODIFY COLUMN `id` BIGINT NOT NULL COMMENT '用户ID (使用雪花算法生成)';

-- 修改 topics 表
ALTER TABLE `topics` 
    MODIFY COLUMN `id` BIGINT NOT NULL COMMENT '话题ID (使用雪花算法生成)',
    MODIFY COLUMN `user_id` BIGINT NOT NULL COMMENT '发布者ID';

-- 修改 votes 表
ALTER TABLE `votes` 
    MODIFY COLUMN `id` BIGINT NOT NULL COMMENT '投票ID (使用雪花算法生成)',
    MODIFY COLUMN `user_id` BIGINT NOT NULL COMMENT '用户ID',
    MODIFY COLUMN `topic_id` BIGINT NOT NULL COMMENT '话题ID';

-- 修改 comments 表
ALTER TABLE `comments` 
    MODIFY COLUMN `id` BIGINT NOT NULL COMMENT '评论ID (使用雪花算法生成)',
    MODIFY COLUMN `topic_id` BIGINT NOT NULL COMMENT '话题ID',
    MODIFY COLUMN `user_id` BIGINT NOT NULL COMMENT '用户ID',
    MODIFY COLUMN `parent_id` BIGINT DEFAULT NULL COMMENT '父评论ID（用于回复）';

-- ========== 4. 重新添加外键约束 ==========
ALTER TABLE `topics` 
    ADD CONSTRAINT `fk_topics_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

ALTER TABLE `votes` 
    ADD CONSTRAINT `fk_votes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    ADD CONSTRAINT `fk_votes_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE;

ALTER TABLE `comments` 
    ADD CONSTRAINT `fk_comments_topic_id` FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE,
    ADD CONSTRAINT `fk_comments_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

-- ========== 5. 验证修改 ==========
-- 查看表结构确认修改成功
SHOW CREATE TABLE `users`;
SHOW CREATE TABLE `topics`;
SHOW CREATE TABLE `votes`;
SHOW CREATE TABLE `comments`;

-- 完成！现在可以启动应用程序，使用雪花算法生成ID

