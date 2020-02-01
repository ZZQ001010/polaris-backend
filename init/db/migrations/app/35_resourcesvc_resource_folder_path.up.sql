ALTER TABLE `ppm_res_folder`
    ADD COLUMN `path` varchar(512) NOT NULL DEFAULT '' AFTER `file_type`,
    ADD INDEX `index_ppm_res_folder_path`(`path`) USING BTREE;