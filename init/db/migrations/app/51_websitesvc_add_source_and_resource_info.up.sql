ALTER TABLE `ppm_wst_contact`
ADD COLUMN `source` varchar(8) NOT NULL DEFAULT '' COMMENT '' AFTER `intention`,
ADD COLUMN `resource_info` varchar(2048) NOT NULL DEFAULT '' COMMENT '' AFTER `source`;