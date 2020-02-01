INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (18, 'Upload', '上传', '上传文件（附件）', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (19, 'Download', '下载', '下载文件（附件）', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (20, 'Remove', '移除', '移除', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (21, 'ModifyPermission', '编辑角色权限', '编辑角色权限', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (22, 'CreateFolder', '创建文件夹', '创建文件夹', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (23, 'ModifyFolder', '编辑文件夹', '编辑文件夹', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO `ppm_rol_operation` (`id`, `code`, `name`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete`)
    VALUES (24, 'DeleteFolder', '删除文件夹', '删除文件夹', 1, 1, now(), 1, now(), 1, 2);



INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 38, 0, 'Permission.Pro.File', 'File', '文件管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/File', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 39, 0, 'Permission.Pro.Tag', 'Tag', '标签管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Tag', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 40, 0, 'Permission.Pro.Attachment', 'Attachment', '附件管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Attachment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 41, 0, 'Permission.Pro.Member', 'Member', '项目成员管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Member', '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 142, 0, 38, 'Permission.Pro.File.Upload', '上传文件', 'Upload', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 143, 0, 38, 'Permission.Pro.File.Download', '下载文件', 'Download', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 144, 0, 38, 'Permission.Pro.File.Modify', '编辑文件', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 145, 0, 38, 'Permission.Pro.File.Delete', '删除文件', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 146, 0, 38, 'Permission.Pro.File.CreateFolder', '创建文件夹', 'CreateFolder', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 147, 0, 38, 'Permission.Pro.File.ModifyFolder', '修改文件夹', 'ModifyFolder', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 148, 0, 38, 'Permission.Pro.File.DeleteFolder', '删除文件夹', 'DeleteFolder', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 149, 0, 39, 'Permission.Pro.Tag.Create', '创建标签', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 150, 0, 39, 'Permission.Pro.Tag.Delete', '删除标签', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 151, 0, 39, 'Permission.Pro.Tag.Remove.', '移除标签', 'Remove', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 152, 0, 40, 'Permission.Pro.Attachment.Upload', '上传附件', 'Upload', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 153, 0, 40, 'Permission.Pro.Attachment.Download', '下载附件', 'Download', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 154, 0, 40, 'Permission.Pro.Attachment.Delete', '删除文件', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 155, 0, 41, 'Permission.Pro.Member.Bind', '添加成员', 'Bind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 156, 0, 41, 'Permission.Pro.Member.Unbind', '移除成员', 'Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 157, 0, 33, 'Permission.Pro.Member.Create', '创建项目角色', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 158, 0, 33, 'Permission.Pro.Member.Delete', '删除项目角色', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 159, 0, 33, 'Permission.Pro.Member.ModifyPermission', '编辑角色权限', 'ModifyPermission', '', 1, 1, now(), 1, now(), 1, 2);

Update ppm_rol_permission_operation set is_delete = 1 where id = 102;
Update ppm_rol_permission_operation set permission_id = 8 where id = 35;
Update ppm_rol_permission_operation set permission_id = 41, `remark` = `name`, `name` = '修改成员角色' where id = 129;
Update ppm_rol_permission set `name` = '任务栏管理' where id = 15;
Update ppm_rol_permission_operation set `remark` = `name`, `name` = '查看任务栏' where id = 50;
Update ppm_rol_permission_operation set `remark` = `name`, `name` = '编辑任务栏' where id = 51;
Update ppm_rol_permission_operation set `remark` = `name`, `name` = '创建任务栏' where id = 52;
Update ppm_rol_permission_operation set `remark` = `name`, `name` = '删除任务栏' where id = 53;
