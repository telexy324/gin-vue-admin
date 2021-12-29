CREATE TABLE `application_server`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `hostname`     varchar(100) NOT NULL DEFAULT '' COMMENT '机器名称',
    `architecture` int(2) NOT NULL DEFAULT '0' COMMENT '架构 1 x86,2 arm',
    `manage_ip`    varchar(15)  NOT NULL DEFAULT '' COMMENT '管理ip',
    `os`           int(2) NOT NULL DEFAULT '0' COMMENT '系统 1 redhat,2 suse,3 centos,4 kylin',
    `os_version`   varchar(50)  NOT NULL DEFAULT '' COMMENT '系统版本',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `sys_authority_menus`
(
    `application_server_id`           bigint UNSIGNED NOT NULL,
    `sys_authority_authority_id` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    PRIMARY KEY (`sys_base_menu_id`, `sys_authority_authority_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;