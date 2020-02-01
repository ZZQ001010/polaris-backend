ALTER TABLE `ppm_pri_tag`
    ADD COLUMN `bg_style` varchar(8) NOT NULL DEFAULT '' AFTER `name_pinyin`,
    ADD COLUMN `font_style` varchar(8) NOT NULL DEFAULT '' AFTER `bg_style`;