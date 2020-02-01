ALTER TABLE `ppm_org_user_config`
ADD COLUMN `default_project_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '' AFTER `daily_project_report_message_status`;
