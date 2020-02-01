/*==============================================================*/
/* Table: ppm_bas_object_id                                     */
/*==============================================================*/
create table if not exists ppm_tak_message
(
  `id` bigint(20) NOT NULL,
  `org_id` bigint(20) NOT NULL DEFAULT '0',
  `topic` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `project_id` bigint(20) NOT NULL DEFAULT '0',
  `issue_id` bigint(20) NOT NULL DEFAULT '0',
  `trends_id` bigint(20) NOT NULL DEFAULT '0',
  `info` varchar(4096) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `content` text COLLATE utf8mb4_bin,
  `fail_count` int(11) NOT NULL DEFAULT '0',
  `fail_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `fail_msg` varchar(1024) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `finish_status` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `finish_msg` varchar(1024) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `start_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `creator` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint(20) NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `version` int(11) NOT NULL DEFAULT '1',
  `is_delete` tinyint(4) NOT NULL DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `index_ppm_tak_message_issue_id` (`issue_id`),
  KEY `index_ppm_tak_message_project_id` (`project_id`),
  KEY `index_ppm_tak_message_org_id` (`org_id`),
  KEY `index_ppm_tak_message_create_time` (`create_time`),
  KEY `index_ppm_tak_message_trends_id` (`trends_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

