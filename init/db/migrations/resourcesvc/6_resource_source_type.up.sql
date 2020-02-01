ALTER TABLE `ppm_res_resource`
ADD COLUMN `source_type` tinyint(4)
NOT NULL DEFAULT '0' AFTER `file_type`;