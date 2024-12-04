
-- +migrate Up
CREATE TABLE IF NOT EXISTS `papers` (
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `article_id`  VARCHAR(255) NOT NULL,
    `title`       VARCHAR(255) NOT NULL,
    `abstract`    TEXT NOT NULL,
    `abstract_jp` TEXT NOT NULL,
    `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS `papers`;
