ALTER TABLE `ppm_org_user_organization`
ADD COLUMN `audit_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '审核时间' AFTER `status`;

ALTER TABLE `ppm_org_user_organization`
ADD COLUMN `auditor_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '审核人id' AFTER `status`;

ALTER TABLE `ppm_org_user_organization`
ADD COLUMN `status_change_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '状态变更时间' AFTER `status`;

ALTER TABLE `ppm_org_user_organization`
ADD COLUMN `status_changer_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '状态变更人id' AFTER `status`;