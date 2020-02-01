ALTER TABLE `ppm_cmm_area` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_cities` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_continents` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_countries` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_industry` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_mobile_prefix` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_regions` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;
ALTER TABLE `ppm_cmm_states` 
ADD COLUMN `is_delete` tinyint(4) NOT NULL DEFAULT 2 AFTER `is_default`;