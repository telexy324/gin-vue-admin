CREATE TABLE `application_servers`
(
    `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`      datetime(0) NULL DEFAULT NULL,
    `updated_at`      datetime(0) NULL DEFAULT NULL,
    `deleted_at`      datetime(0) NULL DEFAULT NULL,
    `hostname`        varchar(100) NOT NULL DEFAULT '' COMMENT '机器名称',
    `architecture`    tinyint(2) NOT NULL DEFAULT '0' COMMENT '架构 1 x86,2 arm',
    `manage_ip`       varchar(15)  NOT NULL DEFAULT '' COMMENT '管理ip',
    `os`              tinyint(2) NOT NULL DEFAULT '0' COMMENT '系统 1 redhat,2 suse,3 centos,4 kylin',
    `os_version`      varchar(50)  NOT NULL DEFAULT '' COMMENT '系统版本',
    `system_id`       bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '所属系统id',
    `application_ids` text                  DEFAULT '' COMMENT '部署应用ids',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_server_relations`
(
    `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`      datetime(0) NULL DEFAULT NULL,
    `updated_at`      datetime(0) NULL DEFAULT NULL,
    `deleted_at`      datetime(0) NULL DEFAULT NULL,
    `start_server_id` bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '源节点id',
    `end_server_id`   bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '目的节点id',
    `end_server_url`  varchar(255) NOT NULL DEFAULT '' COMMENT '目的url',
    `relation`        varchar(100) NOT NULL DEFAULT '' COMMENT '调用关系',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_systems`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT '系统名称',
    `position`   tinyint(2) NOT NULL DEFAULT '0' COMMENT '节点是否系统外 0 系统外, 1 系统内',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `admins`
(
    `id`            bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`    datetime(0) NULL DEFAULT NULL,
    `updated_at`    datetime(0) NULL DEFAULT NULL,
    `deleted_at`    datetime(0) NULL DEFAULT NULL,
    `name`          varchar(100) NOT NULL DEFAULT '' COMMENT '姓名',
    `mobile`        varchar(30)  NOT NULL DEFAULT '' COMMENT '电话',
    `email`         varchar(50)  NOT NULL DEFAULT '' COMMENT '邮箱',
    `department_id` bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '部门id',
    PRIMARY KEY (`id`),
    UNIQUE (`name`),
    UNIQUE (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `departments`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
    `parent_id`  bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '上级部门id',
    `leader_id`  bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '领导id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `applications`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `type`       tinyint(2) NOT NULL DEFAULT '0' COMMENT '应用类型 0 未定义 1 数据库 2 缓存 3 应用 4 存储 5 负载均衡 6 备份 7 反向代理',
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT '应用名称',
    `version`    varchar(100) NOT NULL DEFAULT '' COMMENT '版本',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_system_admins`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `system_id`  bigint UNSIGNED NOT NULL COMMENT '系统id',
    `admin_id`   bigint UNSIGNED NOT NULL COMMENT '管理员id',
    `is_prime`   tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否主管 0 否 1 是',
    PRIMARY KEY (`id`),
    UNIQUE (`system_id`, `admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;