CREATE TABLE IF NOT EXISTS `accounts` (
  `id` CHAR(36) NOT NULL COMMENT "ID",
  `name` VARCHAR(24) NOT NULL COMMENT "アカウント名",
  `password` VARCHAR(60) NOT NULL COMMENT "パスワード",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  `deleted_at` DATETIME (6) COMMENT "削除日時",
  PRIMARY KEY (`id`),
  UNIQUE `uq_accounts_name` (`name`)
);
