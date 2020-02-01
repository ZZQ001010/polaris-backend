ALTER TABLE `ppm_org_user`
MODIFY COLUMN `source_platform` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `password_salt`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `source_platform`;

ALTER TABLE `ppm_org_department`
ADD COLUMN `source_platform` varchar(32) NOT NULL DEFAULT '' AFTER `is_hide`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `is_hide`;

ALTER TABLE `ppm_org_department_out_info`
ADD COLUMN `source_platform` varchar(32) NOT NULL DEFAULT '' AFTER `department_id`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `department_id`;

ALTER TABLE `ppm_org_organization`
ADD COLUMN `source_platform` varchar(32) NOT NULL DEFAULT '' AFTER `scale`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `scale`;

ALTER TABLE `ppm_org_organization_out_info`
ADD COLUMN `source_platform` varchar(32) NOT NULL DEFAULT '' AFTER `out_org_id`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `out_org_id`;

ALTER TABLE `ppm_org_user_out_info`
ADD COLUMN `source_platform` varchar(32) NOT NULL DEFAULT '' AFTER `user_id`,
MODIFY COLUMN `source_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `user_id`;

