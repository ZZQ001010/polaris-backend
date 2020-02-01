ALTER TABLE `ppm_tre_notice`
    ADD COLUMN `relation_type` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `status`,
    ADD COLUMN `ext` varchar(4096) COLLATE utf8mb4_bin NOT NULL DEFAULT '' AFTER `relation_type`;