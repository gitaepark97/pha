CREATE TABLE `user` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `phone_number` char(11) UNIQUE NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `session` (
  `id` varchar(36) PRIMARY KEY,
  `user_id` bigint NOT NULL,
  `refresh_token` varchar(300) NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(45) NOT NULL,
  `is_blocked` tinyint(1) NOT NULL DEFAULT 0,
  `expired_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `session` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE;