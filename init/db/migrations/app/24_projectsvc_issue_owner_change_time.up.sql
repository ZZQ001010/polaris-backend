ALTER TABLE `ppm_pri_issue` 
ADD COLUMN `owner_change_time` datetime(0) NOT NULL DEFAULT '1970-01-01 00:00:00' AFTER `owner`,
ADD INDEX `index_ppm_pri_issue_owner_change_time`(`owner_change_time`) USING BTREE;