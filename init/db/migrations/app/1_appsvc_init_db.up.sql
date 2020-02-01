/*==============================================================*/
/* Table: ppm_bas_app_info                                      */
/*==============================================================*/
CREATE TABLE `ppm_bas_app_info` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `name` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `code` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `secret1` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `secret2` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `owner` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `check_status` tinyint(4) NOT NULL DEFAULT '1',
  `remark` varchar(512) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `creator` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint(20) NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int(11) NOT NULL DEFAULT '1',
  `is_delete` tinyint(4) NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_bas_app_info_create_time` (`create_time`),
  KEY `index_ppm_bas_app_info_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;