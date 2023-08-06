CREATE DATABASE testing;
USE testing;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'TODO',
	`createdAt` timestamp NOT NULL COMMENT 'TODO',
	`email` varchar(100) NOT NULL COMMENT 'TODO',
	PRIMARY KEY (`id`)
) COMMENT='TODO';
