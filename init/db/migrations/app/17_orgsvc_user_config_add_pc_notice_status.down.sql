ALTER TABLE `ppm_org_user_config`
DROP COLUMN `default_project_object_type_id`,
DROP COLUMN `pc_notice_open_status`,
DROP COLUMN `pc_issue_remind_message_status`,
DROP COLUMN `pc_org_message_status`,
DROP COLUMN `pc_project_message_status`,
DROP COLUMN `pc_comment_at_message_status`,
MODIFY COLUMN `daily_report_message_status` tinyint(4) NOT NULL DEFAULT '2',
MODIFY COLUMN `remind_message_status` tinyint(4) NOT NULL DEFAULT '2';