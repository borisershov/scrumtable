CREATE TABLE `settings` (
  `tlgrm_chat_id` bigint NOT NULL,
  `current_date` varchar(15) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `issues` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tlgrm_chat_id` bigint NOT NULL,
  `created_at` varchar(15) NOT NULL,
  `date` varchar(15) NOT NULL,
  `done` tinyint(1) NOT NULL DEFAULT '0',
  `text` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sprint_issues` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tlgrm_chat_id` bigint NOT NULL,
  `date` varchar(15) NOT NULL,
  `goal` tinyint(1) NOT NULL DEFAULT '0',
  `done` tinyint(1) NOT NULL DEFAULT '0',
  `text` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
