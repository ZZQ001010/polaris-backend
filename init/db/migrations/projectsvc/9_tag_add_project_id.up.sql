ALTER TABLE `ppm_pri_tag`
ADD COLUMN `project_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '项目id' AFTER `org_id`,
ADD INDEX `index_ppm_pri_tag_project_id`(`project_id`) USING BTREE;

ALTER TABLE `ppm_pri_tag` ADD COLUMN `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER `create_time`;

ALTER TABLE `ppm_pri_tag` ADD COLUMN `updator` bigint(20) NOT NULL DEFAULT '0' AFTER `create_time`;

