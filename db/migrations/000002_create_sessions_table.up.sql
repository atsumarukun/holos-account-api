CREATE TABLE IF NOT EXISTS `sessions` (
  `account_id` CHAR(36) NOT NULL COMMENT "アカウントID",
  `token` CHAR(32) NOT NULL COMMENT "トークン",
  `expires_at` DATETIME (6) NOT NULL COMMENT "有効期限",
  PRIMARY KEY (`account_id`),
  CONSTRAINT `fk_sessions_account_id` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  UNIQUE `uq_sessions_token` (`token`)
);
