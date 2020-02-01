

INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (1, 'View', '查看', '查看对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (2, 'Modify', '更新', '更新对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (3, 'Delete', '删除', '删除对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (4, 'Create', '创建', '创建对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (5, 'Check', '审核', '审核对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (6, 'Invite', '邀请', '邀请对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (7, 'Bind', '绑定', '绑定对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (8, 'Unbind', '解绑', '解绑对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (9, 'Attention', '关注', '关注对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (10, 'UnAttention', '解除关注', '解除关注对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (11, 'ModifyStatus', '更新状态', '更新对象状态（流程流转）', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (12, 'Comment', '评论', '评论对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (13, 'Transfer', '转让', '转让对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (14, 'Init', '初始化', '初始化对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (15, 'Drop', '注销', '注销组织', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (16, 'Filing', '归档', '归档对象', 1	, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_operation (`id`, `code`, `name`, `remark`, `status`	, `creator`, `create_time`, `updator`, `update_time`, `version`	, `is_delete`)
    VALUES (17, 'UnFiling', '解除归档', '解除归档对象', 1	, 1, now(), 1, now(), 1, 2);



INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 1, 'Permission.Sys.Sys', 'Sys', '系统管理', 0, 1, '/Sys', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 1, 'Permission.Sys.Dic', 'Dic', '数据字典管理', 1, 1, '/Sys/Dic', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 3, 1, 'Permission.Sys.Source', 'Source', '来源渠道管理', 1, 1, '/Sys/Source', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 4, 1, 'Permission.Sys.PayLevel', 'PayLevel', '订购级别管理', 1, 1, '/Sys/PayLevel', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 5, 0, 'Permission.Org.Org', 'Org', '组织相关权限', 0, 2, '/Org/{org_id}', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 6, 0, 'Permission.Org.Config', 'OrgConfig', '组织设置管理', 5, 2, '/Org/{org_id}/OrgConfig', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 7, 0, 'Permission.Org.MessageConfig', 'MessageConfig', '组织系统消息设置', 5, 2, '/Org/{org_id}/MessageConfig', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 8, 0, 'Permission.Org.User', 'User', '组织用户管理', 5, 2, '/Org/{org_id}/User', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 9, 0, 'Permission.Org.Team', 'Team', '组织团队管理', 5, 2, '/Org/{org_id}/Team', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 10, 0, 'Permission.Org.RoleGroup', 'RoleGroup', '组织角色组管理', 5, 2, '/Org/{org_id}/RoleGroup', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 11, 0, 'Permission.Org.Role', 'Role', '组织角色管理', 5, 2, '/Org/{org_id}/Role', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 12, 0, 'Permission.Org.Project', 'Project', '组织项目管理', 5, 2, '/Org/{org_id}/Project', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 13, 0, 'Permission.Org.ProjectType', 'ProjectType', '项目类型管理', 5, 3, '/Org/{org_id}/ProjectType', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 14, 0, 'Permission.Org.IssueSource', 'IssueSource', '问题来源管理', 5, 3, '/Org/{org_id}/IssueSource', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 15, 0, 'Permission.Org.ProjectObjectType', 'ProjectObjectType', '项目对象类型管理', 5, 3, '/Org/{org_id}/ProjectObjectType', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 16, 0, 'Permission.Org.Priority', 'Priority', '优先级管理', 5, 3, '/Org/{org_id}/Priority', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 17, 0, 'Permission.Org.ProcessStatus', 'ProcessStatus', '流程状态管理', 5, 3, '/Org/{org_id}/ProcessStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 18, 0, 'Permission.Org.Process', 'Process', '流程管理', 5, 3, '/Org/{org_id}/Process', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 19, 0, 'Permission.Org.ProcessStep', 'ProcessStep', '流程步骤管理', 5, 3, '/Org/{org_id}/ProcessStep', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 20, 0, 'Permission.Pro.Pro', 'Pro', '项目相关权限', 0, 3, '/Org/{org_id}/Pro/{pro_id}', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 21, 0, 'Permission.Pro.Config', 'ProConfig', '项目设置管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/ProConfig', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 22, 0, 'Permission.Pro.Ban', 'Ban', '项目面板管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Ban', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 23, 0, 'Permission.Pro.Iteration', 'Iteration', '项目迭代管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Iteration', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 24, 0, 'Permission.Pro.Issue', 'Issue', '项目问题管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Issue', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 25, 0, 'Permission.Pro.Issue.2', '2', '项目特性管理', 24, 3, '/Org/{org_id}/Pro/{pro_id}/Issue/2', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 26, 0, 'Permission.Pro.Issue.3', '3', '项目需求管理', 24, 3, '/Org/{org_id}/Pro/{pro_id}/Issue/3', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 27, 0, 'Permission.Pro.Issue.4', '4', '项目任务管理', 24, 3, '/Org/{org_id}/Pro/{pro_id}/Issue/4', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 28, 0, 'Permission.Pro.Issue.5', '5', '项目缺陷管理', 24, 3, '/Org/{org_id}/Pro/{pro_id}/Issue/5', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 29, 0, 'Permission.Pro.Issue.6', '6', '项目测试任务管理', 24, 3, '/Org/{org_id}/Pro/{pro_id}/Issue/6', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 30, 0, 'Permission.Pro.Comment', 'Comment', '评论管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 31, 0, 'Permission.Pro.ProjectVersion', 'ProjectVersion', '项目版本管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/ProjectVersion', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 32, 0, 'Permission.Pro.ProjectModule', 'ProjectModule', '项目模块管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/ProjectModule', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission( `id`, `org_id`, `lang_code`, `code`, `name`, `parent_id`, `type`, `path`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 33, 0, 'Permission.Pro.Role', 'Role', '项目角色管理', 20, 3, '/Org/{org_id}/Pro/{pro_id}/Role', '', 1, 1, now(), 1, now(), 1, 2);



INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 1, 2, 'PermissionOperation.Sys.Dic.View', '查看数据字典', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 1, 2, 'PermissionOperation.Sys.Dic.Create', '创建数据字典', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 3, 1, 2, 'PermissionOperation.Sys.Dic.Modify', '编辑数据字典', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 4, 1, 2, 'PermissionOperation.Sys.Dic.Delete', '删除数据字典', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 5, 1, 3, 'PermissionOperation.Sys.Source.View', '查看来源渠道', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 6, 1, 3, 'PermissionOperation.Sys.Source.Create', '创建来源渠道', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 7, 1, 3, 'PermissionOperation.Sys.Source.Modify', '编辑来源渠道', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 8, 1, 3, 'PermissionOperation.Sys.Source.Delete', '删除来源渠道', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 9, 1, 4, 'PermissionOperation.Sys.PayLevel.View', '查看订购级别', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 10, 1, 4, 'PermissionOperation.Sys.PayLevel.Create', '创建订购级别', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 11, 1, 4, 'PermissionOperation.Sys.PayLevel.Modify', '编辑订购级别', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 12, 1, 4, 'PermissionOperation.Sys.PayLevel.Delete', '删除订购级别', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 13, 0, 6, 'Permission.Org.Config.View', '查看组织设置', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 14, 0, 6, 'Permission.Org.Config.Modify', '编辑组织设置', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 15, 0, 6, 'Permission.Org.Config.Transfer', '转让组织', 'Transfer', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 16, 0, 7, 'Permission.Org.MessageConfig.View', '查看组织系统消息设置', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 17, 0, 7, 'Permission.Org.MessageConfig.Modify', '编辑系统消息设置', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 18, 0, 8, 'Permission.Org.User.View', '查看成员信息', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 19, 0, 8, 'Permission.Org.User.ModifyStatus', '编辑/审核成员状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 20, 0, 8, 'Permission.Org.User.Invite', '邀请成员', 'Invite', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 21, 0, 9, 'Permission.Org.Team.View', '查看团队信息', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 22, 0, 9, 'Permission.Org.Team.Create', '创建团队', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 23, 0, 9, 'Permission.Org.Team.Modify', '编辑团队信息', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 24, 0, 9, 'Permission.Org.Team.Delete', '删除团队', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 25, 0, 9, 'Permission.Org.Team.ModifyStatus', '编辑团队状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 26, 0, 9, 'Permission.Org.Team.Bind', '加入/解除团队成员', 'Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 27, 0, 10, 'Permission.Org.RoleGroup.View', '查看角色组', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 28, 0, 10, 'Permission.Org.RoleGroup.Create', '创建角色组', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 29, 0, 10, 'Permission.Org.RoleGroup.Modify', '编辑角色组', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 30, 0, 10, 'Permission.Org.RoleGroup.Delete', '删除角色组', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 31, 0, 11, 'Permission.Org.Role.View', '查看角色', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 32, 0, 11, 'Permission.Org.Role.Create', '创建角色', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 33, 0, 11, 'Permission.Org.Role.Modify', '编辑角色', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 34, 0, 11, 'Permission.Org.Role.Delete', '删除角色', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 35, 0, 11, 'Permission.Org.Role.Bind', '加入/解除角色成员', 'Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 36, 0, 12, 'Permission.Org.Project.View', '查看项目信息', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 37, 0, 12, 'Permission.Org.Project.Create', '创建项目', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 38, 0, 12, 'Permission.Org.Project.Modify', '编辑项目', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 39, 0, 12, 'Permission.Org.Project.Delete', '删除项目', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 40, 0, 12, 'Permission.Org.Project.Attention', '关注项目', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 41, 0, 12, 'Permission.Org.Project.Filing', '归档项目', 'Filing,UnFiling', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 42, 0, 13, 'Permission.Org.ProjectType.View', '查看项目类型', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 43, 0, 13, 'Permission.Org.ProjectType.Modify', '编辑项目类型', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 44, 0, 13, 'Permission.Org.ProjectType.Create', '创建项目类型', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 45, 0, 13, 'Permission.Org.ProjectType.Delete', '删除项目类型', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 46, 0, 14, 'Permission.Org.IssueSource.View', '查看问题来源', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 47, 0, 14, 'Permission.Org.IssueSource.Modify', '编辑问题来源', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 48, 0, 14, 'Permission.Org.IssueSource.Create', '创建问题来源', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 49, 0, 14, 'Permission.Org.IssueSource.Delete', '删除问题来源', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 50, 0, 15, 'Permission.Org.ProjectObjectType.View', '查看问题对象类型', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 51, 0, 15, 'Permission.Org.ProjectObjectType.Modify', '编辑问题对象类型', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 52, 0, 15, 'Permission.Org.ProjectObjectType.Create', '创建问题对象类型', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 53, 0, 15, 'Permission.Org.ProjectObjectType.Delete', '删除问题对象类型', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 54, 0, 16, 'Permission.Org.Priority.View', '查看优先级', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 55, 0, 16, 'Permission.Org.Priority.Modify', '编辑优先级', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 56, 0, 16, 'Permission.Org.Priority.Create', '创建优先级', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 57, 0, 16, 'Permission.Org.Priority.Delete', '删除优先级', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 58, 0, 17, 'Permission.Org.ProcessStatus.View', '查看流程状态', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 59, 0, 17, 'Permission.Org.ProcessStatus.Modify', '编辑流程状态', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 60, 0, 17, 'Permission.Org.ProcessStatus.Create', '创建流程状态', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 61, 0, 17, 'Permission.Org.ProcessStatus.Delete', '删除流程状态', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 62, 0, 18, 'Permission.Org.Process.View', '查看流程', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 63, 0, 18, 'Permission.Org.Process.Modify', '编辑流程', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 64, 0, 18, 'Permission.Org.Process.Create', '创建流程', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 65, 0, 18, 'Permission.Org.Process.Delete', '删除流程', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 66, 0, 19, 'Permission.Org.ProcessStep.View', '查看流程步骤', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 67, 0, 19, 'Permission.Org.ProcessStep.Modify', '编辑流程步骤', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 68, 0, 19, 'Permission.Org.ProcessStep.Create', '创建流程步骤', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 69, 0, 19, 'Permission.Org.ProcessStep.Delete', '删除流程步骤', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 70, 0, 21, 'Permission.Pro.Config.View', '查看项目设置', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 71, 0, 21, 'Permission.Pro.Config.Modify', '编辑项目设置', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 72, 0, 21, 'Permission.Pro.Config.Filing', '归档项目', 'Filing,UnFiling', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 73, 0, 21, 'Permission.Pro.Config.ModifyStatus', '变更项目状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 74, 0, 22, 'Permission.Pro.Ban.View', '查看面板', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 75, 0, 23, 'Permission.Pro.Iteration.View', '查看迭代', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 76, 0, 23, 'Permission.Pro.Iteration.Modify', '编辑迭代', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 77, 0, 23, 'Permission.Pro.Iteration.Create', '创建迭代', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 78, 0, 23, 'Permission.Pro.Iteration.Delete', '删除迭代', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 79, 0, 23, 'Permission.Pro.Iteration.ModifyStatus', '变更迭代状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 80, 0, 23, 'Permission.Pro.Iteration.Bind', '规划迭代', 'Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 81, 0, 23, 'Permission.Pro.Iteration.Attention', '关注迭代', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 82, 0, 25, 'Permission.Pro.Issue.2.View', '查看特性', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 83, 0, 25, 'Permission.Pro.Issue.2.Modify', '编辑特性', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 84, 0, 25, 'Permission.Pro.Issue.2.Create', '创建特性', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 85, 0, 25, 'Permission.Pro.Issue.2.Delete', '删除特性', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 86, 0, 25, 'Permission.Pro.Issue.2.ModifyStatus', '变更特性状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 87, 0, 25, 'Permission.Pro.Issue.2.Comment', '评论特性', 'Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 88, 0, 25, 'Permission.Pro.Issue.2.Attention', '关注特性', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 89, 0, 26, 'Permission.Pro.Issue.3.View', '查看需求', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 90, 0, 26, 'Permission.Pro.Issue.3.Modify', '编辑需求', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 91, 0, 26, 'Permission.Pro.Issue.3.Create', '创建需求', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 92, 0, 26, 'Permission.Pro.Issue.3.Delete', '删除需求', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 93, 0, 26, 'Permission.Pro.Issue.3.ModifyStatus', '变更需求状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 94, 0, 26, 'Permission.Pro.Issue.3.Comment', '评论需求', 'Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 95, 0, 26, 'Permission.Pro.Issue.3.Attention', '关注需求', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 96, 0, 27, 'Permission.Pro.Issue.4.View', '查看任务', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 97, 0, 27, 'Permission.Pro.Issue.4.Modify', '编辑任务', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 98, 0, 27, 'Permission.Pro.Issue.4.Create', '创建任务', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 99, 0, 27, 'Permission.Pro.Issue.4.Delete', '删除任务', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 100, 0, 27, 'Permission.Pro.Issue.4.ModifyStatus', '变更任务状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 101, 0, 27, 'Permission.Pro.Issue.4.Comment', '评论任务', 'Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 102, 0, 27, 'Permission.Pro.Issue.4.Attention', '关注任务', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 103, 0, 28, 'Permission.Pro.Issue.5.View', '查看缺陷', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 104, 0, 28, 'Permission.Pro.Issue.5.Modify', '编辑缺陷', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 105, 0, 28, 'Permission.Pro.Issue.5.Create', '创建缺陷', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 106, 0, 28, 'Permission.Pro.Issue.5.Delete', '删除缺陷', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 107, 0, 28, 'Permission.Pro.Issue.5.ModifyStatus', '变更缺陷状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 108, 0, 28, 'Permission.Pro.Issue.5.Comment', '评论缺陷', 'Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 109, 0, 28, 'Permission.Pro.Issue.5.Attention', '关注缺陷', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 110, 0, 29, 'Permission.Pro.Issue.5.View', '查看测试任务', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 111, 0, 29, 'Permission.Pro.Issue.5.Modify', '编辑测试任务', 'Modify,Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 112, 0, 29, 'Permission.Pro.Issue.5.Create', '创建测试任务', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 113, 0, 29, 'Permission.Pro.Issue.5.Delete', '删除测试任务', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 114, 0, 29, 'Permission.Pro.Issue.5.ModifyStatus', '变更测试任务状态', 'ModifyStatus', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 115, 0, 29, 'Permission.Pro.Issue.5.Comment', '评论测试任务', 'Comment', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 116, 0, 29, 'Permission.Pro.Issue.5.Attention', '关注测试任务', 'Attention,UnAttention', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 117, 0, 30, 'Permission.Pro.Comment.Modify', '编辑评论', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 118, 0, 30, 'Permission.Pro.Comment.Delete', '删除评论', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 119, 0, 31, 'Permission.Pro.ProjectVersion.View', '查看项目版本', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 120, 0, 31, 'Permission.Pro.ProjectVersion.Modify', '编辑项目版本', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 121, 0, 31, 'Permission.Pro.ProjectVersion.Create', '创建项目版本', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 122, 0, 31, 'Permission.Pro.ProjectVersion.Delete', '删除项目版本', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 123, 0, 32, 'Permission.Pro.ProjectModule.View', '查看项目模块', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 124, 0, 32, 'Permission.Pro.ProjectModule.Modify', '编辑项目模块', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 125, 0, 32, 'Permission.Pro.ProjectModule.Create', '创建项目模块', 'Create', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 126, 0, 32, 'Permission.Pro.ProjectModule.Delete', '删除项目模块', 'Delete', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 127, 0, 33, 'Permission.Pro.Role.View', '查看项目角色', 'View', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 128, 0, 33, 'Permission.Pro.Role.Modify', '编辑项目角色', 'Modify', '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_permission_operation( `id`, `org_id`, `permission_id`, `lang_code`, `name`, `operation_codes`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 129, 0, 33, 'Permission.Pro.Role.Bind', '绑定/解绑角色用户', 'Bind,Unbind', '', 1, 1, now(), 1, now(), 1, 2);



INSERT INTO ppm_rol_role_group( `id`, `org_id`, `lang_code`, `name`, `remark`, `type`, `is_readonly`, `is_show`, `is_default`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 0, 'RoleGroup.Sys', '系统角色组', '', 1, 1, 1, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_group( `id`, `org_id`, `lang_code`, `name`, `remark`, `type`, `is_readonly`, `is_show`, `is_default`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 0, 'RoleGroup.Special', '特殊角色组', '', 1, 1, 1, 2, 1, 1, now(), 1, now(), 1, 2);



INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 1, 1, 'Role.Sys.Admin', '系统超级管理员', '', 1, 2, 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 2, 1, 'Role.Sys.Manager', '系统管理员', '', 1, 1, 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 3, 1, 'Role.Sys.Member', '系统成员', '', 1, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 4, 0, 'RoleGroup.Special.Creator', 'Creator', '创建人', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 5, 0, 'RoleGroup.Special.Owner', 'Owner', '负责人', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 6, 0, 'RoleGroup.Special.Worker', 'Worker', '参与人', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 7, 0, 'RoleGroup.Special.Attention', 'Attention', '关注人', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 8, 0, 'RoleGroup.Special.Member', 'Member', '组织成员', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role( `id`, `org_id`, `lang_code`, `name`, `remark`, `is_readonly`, `is_modify_permission`, `is_default`, `role_group_id`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
VALUES ( 9, 0, 'RoleGroup.Special.Visitor', 'Visitor', '访客', 1, 1, 2, 2, 1, 1, now(), 1, now(), 1, 2);


INSERT INTO ppm_rol_role_user( `id`, `org_id`, `project_id`, `role_id`, `user_id`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` ) VALUES ( 1, 1, 0, 1, 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 1, 1, 0, 1, '/Sys', '*', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 1, 2, 0, 2, '/Sys/Dic', '(View)|(Create)|(Modify)', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 1, 2, 0, 3, '/Sys/Source', '(View)|(Create)|(Modify)', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 1, 2, 0, 4, '/Sys/PayLevel', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 1, 3, 0, 2, '/Sys/Dic', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 1, 3, 0, 3, '/Sys/Source', 'View', 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_rol_role_permission_operation( `id`, `org_id`, `role_id`, `project_id`, `permission_id`, `permission_path`, `operation_codes`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 1, 3, 0, 4, '/Sys/PayLevel', 'View', 1, now(), 1, now(), 1, 2);