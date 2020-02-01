CREATE TABLE if not exists  `ppm_pri_tag` (
   `id` bigint(20) NOT NULL ,
   `org_id` bigint(20) NOT NULL DEFAULT '0' ,
   `name` varchar(128) COLLATE utf8mb4_bin NOT NULL DEFAULT '' ,
   `name_pinyin` varchar(128) COLLATE utf8mb4_bin NOT NULL DEFAULT '' ,
   `creator` bigint(20) NOT NULL DEFAULT '0' ,
   `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ,
   `version` int(11) NOT NULL DEFAULT '1' ,
   `is_delete` tinyint(4) NOT NULL DEFAULT '2' ,
   PRIMARY KEY (`id`),
   KEY `index_ppm_pri_tag_name` (`name`),
   KEY `index_ppm_pri_tag_org_id` (`org_id`),
   KEY `index_ppm_pri_tag_name_pinyin` (`name_pinyin`),
   KEY `index_ppm_pri_tag_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE if not exists  `ppm_pri_issue_tag` (
 `id` bigint(20) NOT NULL ,
 `org_id` bigint(20) NOT NULL DEFAULT '0' ,
 `project_id` bigint(20) NOT NULL DEFAULT '0' ,
 `issue_id` bigint(20) NOT NULL DEFAULT '0' ,
 `tag_id` bigint(20) NOT NULL DEFAULT '0' ,
 `tag_name` varchar(128) COLLATE utf8mb4_bin NOT NULL DEFAULT '' ,
 `creator` bigint(20) NOT NULL DEFAULT '0' ,
 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ,
 `updator` bigint(20) NOT NULL DEFAULT '0' ,
 `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
 `version` int(11) NOT NULL DEFAULT '1' ,
 `is_delete` tinyint(4) NOT NULL DEFAULT '2' ,
 PRIMARY KEY (`id`),
 KEY `index_ppm_pri_issue_tag_project_id` (`project_id`),
 KEY `index_ppm_pri_issue_tag_org_id` (`org_id`),
 KEY `index_ppm_pri_issue_tag_issue_id` (`issue_id`),
 KEY `index_ppm_pri_issue_tag_tag_name` (`tag_name`),
 KEY `index_ppm_pri_issue_tag_tag_id` (`tag_id`),
 KEY `index_ppm_pri_issue_relation_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


ALTER TABLE `ppm_pri_issue_detail` DROP COLUMN `tags`;

ALTER TABLE `ppm_pri_issue` ADD COLUMN `sort` int(11) NOT NULL DEFAULT 0 AFTER `end_time`;

ALTER TABLE `ppm_pri_issue`
    DROP INDEX `index_ppm_pri_issue_project_id_project_objct_type_id`,
    ADD INDEX `index_ppm_pri_issue_project_id_project_objct_type_id_sort`(`project_id`, `project_object_type_id`, `sort`) USING BTREE;