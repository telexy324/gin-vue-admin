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
    `system_id`    bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '所属系统id',
    `app_ids`      text COMMENT '部署应用ids',
    `ssh_port`     int(5) NOT NULL DEFAULT '0' COMMENT 'ssh端口',
    `display_name` varchar(100) NOT NULL DEFAULT '' COMMENT '展示名称',
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
    `network`    text COMMENT 'cidr networks',
    `position`   tinyint(2) NOT NULL DEFAULT '0' COMMENT '系统位置 0 未知, 1 月坛, 2 昌平, 3 丰台, 4 珠海, 5 西安',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_admins`
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

CREATE TABLE `application_departments`
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

CREATE TABLE `application_apps`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `type`       tinyint(2) NOT NULL DEFAULT '0' COMMENT '应用类型 0 未定义 1 数据库 2 缓存 3 web中间件 4 存储 5 负载均衡 6 备份 7 反向代理 8 队列 9 搜索引擎',
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
    PRIMARY KEY (`id`),
    UNIQUE (`system_id`, `admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_system_relations`
(
    `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`      datetime(0) NULL DEFAULT NULL,
    `updated_at`      datetime(0) NULL DEFAULT NULL,
    `deleted_at`      datetime(0) NULL DEFAULT NULL,
    `start_system_id` bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '源系统id',
    `end_system_id`   bigint UNSIGNED NOT NULL DEFAULT '0' COMMENT '目的系统id',
    `end_system_url`  varchar(255) NOT NULL DEFAULT '' COMMENT '目的url',
    `relation`        varchar(100) NOT NULL DEFAULT '' COMMENT '调用关系',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `application_task_schedules`
(
    `id`               bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`       datetime(0) NULL DEFAULT NULL,
    `updated_at`       datetime(0) NULL DEFAULT NULL,
    `deleted_at`       datetime(0) NULL DEFAULT NULL,
    `template_id`      bigint UNSIGNED NOT NULL DEFAULT '0',
    `cron_format`      varchar(255) not null,
    `last_commit_hash` varchar(40)  NOT NULL DEFAULT '',
    `valid`            tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 禁用 1 启用',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `application_tasks`
(
    `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`     datetime(0) NULL DEFAULT NULL,
    `updated_at`     datetime(0) NULL DEFAULT NULL,
    `deleted_at`     datetime(0) NULL DEFAULT NULL,
    `template_id`    int          NOT NULL DEFAULT '0',
    `status`         varchar(255) NOT NULL DEFAULT '',
    `system_user_id` bigint UNSIGNED NOT NULL DEFAULT '0',
    `begin_time`     datetime(0) NULL DEFAULT NULL,
    `end_time`       datetime(0) NULL DEFAULT NULL,
    `set_task_id`    bigint UNSIGNED NOT NULL DEFAULT '0',
    `file_download`  varchar(255) NOT NULL DEFAULT '',
    `system_id`      bigint UNSIGNED NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `application_task_outputs`
(
    `id`          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(0) NULL DEFAULT NULL,
    `updated_at`  datetime(0) NULL DEFAULT NULL,
    `deleted_at`  datetime(0) NULL DEFAULT NULL,
    `task_id`     int         NOT NULL DEFAULT '0',
    `record_time` datetime    not null,
    `output`      text,
    `manage_ip`   varchar(15) NOT NULL DEFAULT '' COMMENT '管理ip',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `application_task_templates`
(
    `id`              bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`      datetime(0) NULL DEFAULT NULL,
    `updated_at`      datetime(0) NULL DEFAULT NULL,
    `deleted_at`      datetime(0) NULL DEFAULT NULL,
    `name`            varchar(100) not null,
    `description`     text,
    target_server_ids text,
    `mode`            tinyint(2) NOT NULL DEFAULT '0' COMMENT '运行方式 0 未定义 1 命令 2 脚本',
    `command`         text,
    `script_path`     varchar(255) NOT NULL DEFAULT '',
    `last_task_id`    int          NOT NULL DEFAULT '0',
    `sys_user`        varchar(30)  NOT NULL DEFAULT 'root',
    `system_id`       bigint UNSIGNED NOT NULL DEFAULT '0',
    `execute_type`    tinyint(2) NOT NULL DEFAULT '0' COMMENT '模板类型 1 普通 2 日志提取',
    `log_path`        varchar(255) NOT NULL DEFAULT '',
    `script_hash`     char(32)     NOT NULL DEFAULT '',
    `log_output`      tinyint(2) NOT NULL DEFAULT '0' COMMENT '1 直接下载 2 上传服务器',
    `log_dst`         varchar(255) NOT NULL DEFAULT '',
    `dst_server_id`   int          NOT NULL DEFAULT '0',
    `secret_id`       int          NOT NULL DEFAULT '0',
    `shell_type`      int          NOT NULL DEFAULT '0' COMMENT '1 sh 2 bash',
    `shell_vars`      varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_system_sys_admins`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `system_id`  bigint UNSIGNED NOT NULL COMMENT '系统id',
    `admin_id`   bigint UNSIGNED NOT NULL COMMENT '管理员id',
    PRIMARY KEY (`id`),
    UNIQUE (`system_id`, `admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_system_edit_relations`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `system_id`  bigint UNSIGNED NOT NULL COMMENT '系统id',
    `relation`   text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_task_template_sets`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `system_id`  bigint UNSIGNED NOT NULL COMMENT '系统id',
    `name`       varchar(100) not null,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_task_template_set_templates`
(
    `id`          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(0) NULL DEFAULT NULL,
    `updated_at`  datetime(0) NULL DEFAULT NULL,
    `deleted_at`  datetime(0) NULL DEFAULT NULL,
    `set_id`      bigint UNSIGNED NOT NULL COMMENT '模板集id',
    `template_id` bigint UNSIGNED NOT NULL COMMENT '模板id',
    `seq`         int(4) NOT NULL COMMENT '执行顺序',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `application_set_tasks`
(
    `id`               bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`       datetime(0) NULL DEFAULT NULL,
    `updated_at`       datetime(0) NULL DEFAULT NULL,
    `deleted_at`       datetime(0) NULL DEFAULT NULL,
    `set_id`           bigint UNSIGNED NOT NULL DEFAULT '0',
    `system_user_id`   bigint UNSIGNED NOT NULL DEFAULT '0',
    `current_task_id`  bigint UNSIGNED NOT NULL DEFAULT '0',
    `total_steps`      int(4) NOT NULL DEFAULT '0',
    `current_step`     int(4) NOT NULL DEFAULT '0',
    `templates_string` text,
    `tasks_string`     text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_upload_servers`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `hostname`   varchar(100) NOT NULL DEFAULT '' COMMENT '机器名称',
    `manage_ip`  varchar(15)  NOT NULL DEFAULT '' COMMENT '管理ip',
    `mode`       tinyint(2) NOT NULL DEFAULT '0' COMMENT '运行方式 0 未定义 1 ssh 2 ftp',
    `port`       int(5) NOT NULL DEFAULT '0' COMMENT '端口',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `application_upload_secrets`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `server_id`  bigint UNSIGNED NOT NULL COMMENT '上传服务器id',
    `name`       varchar(100) NOT NULL DEFAULT '' COMMENT '用户名',
    `password`   varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;