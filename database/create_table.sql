CREATE TABLE `application_servers`
(
    `id`           bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`   datetime(0) NULL DEFAULT NULL,
    `updated_at`   datetime(0) NULL DEFAULT NULL,
    `deleted_at`   datetime(0) NULL DEFAULT NULL,
    `hostname`     varchar(100) NOT NULL DEFAULT '' COMMENT '机器名称',
    `architecture` tinyint(2) NOT NULL DEFAULT '0' COMMENT '架构 1 x86,2 arm',
    `manage_ip`    varchar(15)  NOT NULL DEFAULT '' COMMENT '管理ip',
    `os`           tinyint(2) NOT NULL DEFAULT '0' COMMENT '系统 1 redhat,2 suse,3 centos,4 kylin',
    `os_version`   varchar(50)  NOT NULL DEFAULT '' COMMENT '系统版本',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
