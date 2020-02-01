
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 'web', 'web主站', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 'ding', '钉钉', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO ppm_bas_source_channel ( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 3, 'wechat', '微信', '', 1, 1, now( ), 1, now( ), 1, 2 );
INSERT INTO ppm_bas_source_channel( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 4, 'demo', '示例', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_source_channel( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 5, 'invite', '用户邀请', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_bas_source_channel( `id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 6, 'lark-xyjh2019', 'Lark巡洋计划2019', '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (1, 'PayLevel.Trial', '试用级别', 1024000, 100, 0, 0, 1296000, 2, 0, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (2, 'PayLevel.Free', '免费版', 1024000, 10, 0, 0, 0, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (3, 'PayLevel.Professional', '专业版', 1024000000, 100, 25600, 25600, 31536000, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_bas_pay_level`( `id`, `lang_code`, `name`, `storage`, `member_count`, `price`, `member_price`, `duration`, `is_show`, `sort`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES (4, 'PayLevel.Vip', 'VIP版', 1024000000, 1000, 102400, 25600, 31536000, 1, 1, 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_org_organization( `id`, `name`, `web_site`, `industry_id`, `scale`, `source_channel`, `country_id`, `province_id`, `city_id`, `address`, `logo_url`, `resource_id`, `owner`, `is_authenticated`, `remark`, `init_status`, `status`, `is_show`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, '北极星组织', '', 0, '', 'web', 0, 0, 0, '', '', 0, 1, 1, '', 3, 1, 0, 1, now(), 1, now(), 1, 2);


INSERT INTO ppm_org_user( `id`, `org_id`, `name`, `login_name`, `login_name_edit_count`, `email`, `mobile`, `birthday`, `sex`, `password`, `password_salt`, `source_channel`, `language`, `motto`, `last_login_ip`, `last_login_time`, `login_fail_count`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 1, '管理员', 'admin', 99, 'admin@admin.com', '12345678901', '1970-01-01 00:00:00', 99, '2eab55be9dfa1ba036d9024baf2c60c8', '8a66572aac6911e99080784f439212b0', 'web', 'zh-CN', '', '', '1970-01-01 00:00:00', 0, 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_org_user_organization( `id`, `org_id`, `user_id`, `check_status`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_tem_team( `id`, `org_id`, `name`, `nick_name`, `owner`, `department_id`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, '北极星组织团队', '北极星组织团队', 1, 0, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_tem_user_team( `id`, `org_id`, `team_id`, `user_id`, `relation_type`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
