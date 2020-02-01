ALTER TABLE `ppm_org_user_config`
    ADD COLUMN `daily_project_report_message_status` tinyint(4) NOT NULL DEFAULT 2 AFTER `relation_message_status`;