ALTER TABLE `ppm_pro_project_detail`
    ADD COLUMN `is_sync_out_calendar` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_enable_work_hours`;