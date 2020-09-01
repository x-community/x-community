CREATE TABLE `user` (
	`id` INT UNSIGNED AUTO_INCREMENT,
	`username` VARCHAR (32),
	`email` VARCHAR (256),
	`salt` VARCHAR (10),
	`password` VARCHAR (128),
	`actived` boolean,
	`active_code` VARCHAR (32),
	`last_login_ip` VARCHAR (64),
	`last_login_time` datetime,
	`created_at` DATETIME NULL,
	`actived_at` DATETIME NULL,
	`updated_at` DATETIME NULL,
	PRIMARY KEY (`id`)
);

CREATE UNIQUE INDEX `uix_user_username` ON `user` (`username`);

CREATE UNIQUE INDEX `uix_user_email` ON `user` (`email`);

CREATE TABLE `user_relation` (
	`id` BIGINT AUTO_INCREMENT,
	`user_id` INT UNSIGNED,
	`fellow_user_id` INT UNSIGNED,
	`created_at` DATETIME NULL,
	PRIMARY KEY (`id`)
);

CREATE INDEX idx_fellow_uid ON `user_relation`(fellow_user_id);

CREATE UNIQUE INDEX `uix_user_relation` ON `user_relation` (`user_id`, `fellow_user_id`);