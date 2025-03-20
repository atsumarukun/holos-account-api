CREATE TABLE IF NOT EXISTS `sessions` (
  `id` CHAR(36) NOT NULL COMMENT "ID",
  `account_id` CHAR(36) NOT NULL COMMENT "アカウントID",
  `token` CHAR(32) NOT NULL COMMENT "トークン",
  `expires_at` DATETIME (6) NOT NULL COMMENT "有効期限",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  `deleted_at` DATETIME (6) COMMENT "削除日時",
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_sessions_account_id` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  UNIQUE `uq_sessions_account_id` (`account_id`),
  UNIQUE `uq_sessions_token` (`token`)
);
