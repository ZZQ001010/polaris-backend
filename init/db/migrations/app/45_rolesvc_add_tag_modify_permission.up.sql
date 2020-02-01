INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 160, 0, 39, 'Permission.Pro.Tag.Modify', '编辑标签', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);

Update ppm_rol_permission_operation set `lang_code` = 'Permission.Pro.Tag.Remove' where id = 151;
