
-- +migrate Up
CREATE INDEX `idx_created_at` ON `papers` (`created_at`);

-- +migrate Down
DROP INDEX `idx_created_at` ON `papers`;
