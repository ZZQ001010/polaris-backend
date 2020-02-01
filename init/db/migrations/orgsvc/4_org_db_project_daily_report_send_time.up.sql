ALTER TABLE `ppm_orc_config`
    ADD COLUMN `project_daily_report_send_time` varchar(8) NOT NULL DEFAULT '18:00' COMMENT '项目日报发送时间' AFTER `remind_send_time`;