CREATE TABLE `access_keys`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `name`       varchar(255) NOT NULL DEFAULT '',
    `type`       varchar(255) NOT NULL DEFAULT '',
    `project_id` bigint UNSIGNED NOT NULL DEFAULT '0',
    `key`        text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `environments`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `name`       varchar(255) NOT NULL DEFAULT '',
    `project_id` bigint UNSIGNED NOT NULL DEFAULT '0',
    `password`   varchar(255) NOT NULL DEFAULT '',
    `json`       text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `inventorys`
(
    `id`            bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`    datetime(0) NULL DEFAULT NULL,
    `updated_at`    datetime(0) NULL DEFAULT NULL,
    `deleted_at`    datetime(0) NULL DEFAULT NULL,
    `project_id`    bigint UNSIGNED NOT NULL DEFAULT '0',
    `type`          varchar(255) NOT NULL DEFAULT '',
    `ssh_key_id`    bigint UNSIGNED NOT NULL DEFAULT '0',
    `become_key_id` bigint UNSIGNED NOT NULL DEFAULT '0',
    `inventory`     text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `projects`
(
    `id`                 bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`         datetime(0) NULL DEFAULT NULL,
    `updated_at`         datetime(0) NULL DEFAULT NULL,
    `deleted_at`         datetime(0) NULL DEFAULT NULL,
    `created`            datetime     not null,
    `name`               varchar(255) not null,
    `alert`              BOOLEAN      NOT NULL DEFAULT FALSE,
    `alert_chat`         varchar(30)  NOT NULL DEFAULT '',
    `max_parallel_tasks` int(11) not null default 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `project_users`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `project_id` bigint UNSIGNED NOT NULL DEFAULT '0',
    `user_id`    bigint UNSIGNED NOT NULL DEFAULT '0',
    `admin`      boolean not null default false,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `environments`
(
    `id`               bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`       datetime(0) NULL DEFAULT NULL,
    `updated_at`       datetime(0) NULL DEFAULT NULL,
    `deleted_at`       datetime(0) NULL DEFAULT NULL,
    `project_id`       bigint UNSIGNED NOT NULL DEFAULT '0',
    `template_id`      int primary key,
    `cron_format`      varchar(255) not null,
    `last_commit_hash` varchar(40)  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `tasks`
(
    `id`             bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`     datetime(0) NULL DEFAULT NULL,
    `updated_at`     datetime(0) NULL DEFAULT NULL,
    `deleted_at`     datetime(0) NULL DEFAULT NULL,
    `template_id`    int          NOT NULL DEFAULT '0',
    `status`         varchar(255) NOT NULL DEFAULT '',
    `playbook`       varchar(255) NOT NULL DEFAULT '',
    `environment`    text,
    `created`        datetime null,
    `start`          datetime null,
    `end`            datetime null,
    `arguments`      text,
    `project_id`     bigint UNSIGNED NOT NULL DEFAULT '0',
    `message`        varchar(250) NOT NULL DEFAULT '',
    `version`        varchar(20)  NOT NULL DEFAULT '',
    `commit_hash`    varchar(40)  NOT NULL DEFAULT '',
    `commit_message` varchar(100) NOT NULL DEFAULT '',
    `build_task_id`  bigint UNSIGNED NOT NULL DEFAULT '0',
    `debug`          BOOLEAN      NOT NULL DEFAULT FALSE,
    `dry_run`        BOOLEAN      NOT NULL DEFAULT FALSE,
    `user_id`        bigint UNSIGNED NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `task_outputs`
(
    `id`         bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(0) NULL DEFAULT NULL,
    `updated_at` datetime(0) NULL DEFAULT NULL,
    `deleted_at` datetime(0) NULL DEFAULT NULL,
    `task_id`    int          NOT NULL DEFAULT '0',
    `task`       varchar(255) NOT NULL DEFAULT '',
    `time`       datetime     not null,
    `output`     text,

    unique (`task_id`, `time`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create table `templates`
(
    `id`                          bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`                  datetime(0) NULL DEFAULT NULL,
    `updated_at`                  datetime(0) NULL DEFAULT NULL,
    `deleted_at`                  datetime(0) NULL DEFAULT NULL,
    `project_id`                  bigint UNSIGNED NOT NULL DEFAULT '0',
    `inventory_id`                bigint UNSIGNED NOT NULL DEFAULT '0',
    `environment_id`              bigint UNSIGNED NOT NULL DEFAULT '0',
    `playbook`                    varchar(255) NOT NULL DEFAULT '',
    `arguments`                   text,
    `description`                 text,
    `become_key_id`               bigint UNSIGNED NOT NULL DEFAULT '0',
    `vault_key_id`                bigint UNSIGNED NOT NULL DEFAULT '0',
    `view_id`                     int          references `project__view` (id) on delete set null,
    `survey_vars`                 text,
    `autorun`                     boolean               default false,
    `allow_override_args_in_task` bool         not null default false,
    `name`                        varchar(100) not null,
    `suppress_success_alerts`     bool         not null default false,
    `build_template_id`           int          not null,
    `start_version`               varchar(100) not null,
    `type`                        varchar(40)  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;