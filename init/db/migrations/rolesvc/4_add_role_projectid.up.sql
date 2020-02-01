ALTER TABLE `ppm_rol_role`
    ADD COLUMN `project_id` bigint(0) NOT NULL DEFAULT 0 AFTER `lang_code`,
ADD INDEX `index_ppm_rol_role_project_id`(`project_id`) USING BTREE;