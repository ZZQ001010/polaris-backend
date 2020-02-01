
ALTER TABLE `ppm_res_resource`
    ADD COLUMN `file_type` tinyint(4) NOT NULL DEFAULT 0 AFTER `md5`;