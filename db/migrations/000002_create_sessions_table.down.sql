ALTER TABLE `sessions`
DROP FOREIGN KEY `fk_sessions_account_id`;

ALTER TABLE `sessions`
DROP INDEX `uq_sessions_token`;

DROP TABLE IF EXISTS `sessions`;
