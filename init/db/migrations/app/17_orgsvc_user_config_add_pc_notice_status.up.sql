ALTER TABLE `ppm_org_user_config`
ADD COLUMN `default_project_object_type_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '默认工作栏id' AFTER `default_project_id`,
ADD COLUMN `pc_notice_open_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'pc桌面通知开关状态, 2否, 1是' AFTER `default_project_object_type_id`,
ADD COLUMN `pc_issue_remind_message_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'pc任务提醒状态, 2否, 1是' AFTER `pc_notice_open_status`,
ADD COLUMN `pc_org_message_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'pc组织相关推送状态, 2否, 1是' AFTER `pc_issue_remind_message_status`,
ADD COLUMN `pc_project_message_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'pc项目相关推送状态, 2否, 1是' AFTER `pc_org_message_status`,
ADD COLUMN `pc_comment_at_message_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'pc评论相关推送状态, 2否, 1是' AFTER `pc_project_message_status`,
MODIFY COLUMN `daily_report_message_status` tinyint(4) NOT NULL DEFAULT '1',
MODIFY COLUMN `remind_message_status` tinyint(4) NOT NULL DEFAULT '1';