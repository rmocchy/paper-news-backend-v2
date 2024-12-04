
-- +migrate Up
ALTER TABLE `papers`
ADD CONSTRAINT `uq_article_id` UNIQUE (`article_id`);

-- +migrate Down
ALTER TABLE `papers`
DROP INDEX `uq_article_id`;
