INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 34, 0, 'Permission.Pro.Test', 'Test', '项目测试管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Test', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 35, 0, 'Permission.Pro.Test.TestApp', 'TestApp', '测试应用管理', 34, 3, '/Org/{org_id}/Pro/{pro_id}/Test/TestApp', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 36, 0, 'Permission.Pro.Test.TestDevice', 'TestDevice', '测试设备管理', 34, 3, '/Org/{org_id}/Pro/{pro_id}/Test/TestDevice', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 37, 0, 'Permission.Pro.Test.TestReport', 'TestReport', '测试报告管理', 34, 3, '/Org/{org_id}/Pro/{pro_id}/Test/TestReport', '', 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 130, 0, 35, 'Permission.Pro.Test.TestApp.View', '查看测试应用', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 131, 0, 35, 'Permission.Pro.Test.TestApp.Create', '创建测试应用', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 132, 0, 35, 'Permission.Pro.Test.TestApp.Modify', '编辑测试应用', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 133, 0, 35, 'Permission.Pro.Test.TestApp.Delete', '删除测试应用', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 134, 0, 36, 'Permission.Pro.Test.TestDevice.View', '查看测试设备', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 135, 0, 36, 'Permission.Pro.Test.TestDevice.Create', '创建测试设备', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 136, 0, 36, 'Permission.Pro.Test.TestDevice.Modify', '编辑测试设备', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 137, 0, 36, 'Permission.Pro.Test.TestDevice.Delete', '删除测试设备', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 138, 0, 37, 'Permission.Pro.Test.TestReport.View', '查看测试报告', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 139, 0, 37, 'Permission.Pro.Test.TestReport.Create', '创建测试报告', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 140, 0, 37, 'Permission.Pro.Test.TestReport.Modify', '编辑测试报告', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
    VALUES ( 141, 0, 37, 'Permission.Pro.Test.TestReport.Delete', '删除测试报告', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);