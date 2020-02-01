ALTER TABLE `ppm_org_organization`
ADD COLUMN `code` varchar(64) NOT NULL DEFAULT '' COMMENT '编号' AFTER `name`,
ADD INDEX `index_ppm_org_organization_code`(`code`) USING BTREE;
