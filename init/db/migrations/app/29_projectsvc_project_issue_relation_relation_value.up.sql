ALTER TABLE `ppm_pro_project_relation`
    ADD COLUMN `relation_code` varchar(128) NOT NULL DEFAULT '' AFTER `relation_id`;

ALTER TABLE `ppm_pri_issue_relation`
    ADD COLUMN `relation_code` varchar(128) NOT NULL DEFAULT '' AFTER `relation_id`;