update ppm_rol_role_permission_operation set `permission_id` = 7 where `permission_path` like '%/MessageConfig' and `permission_id` = 11;
update ppm_rol_role_permission_operation set `permission_id` = 13 where `permission_path` like '%/ProjectType' and `permission_id` = 12;
update ppm_rol_role_permission_operation set `permission_id` = 14 where `permission_path` like '%/IssueSource' and `permission_id` = 11;
update ppm_rol_role_permission_operation set `permission_id` = 15 where `permission_path` like '%/ProjectObjectType' and `permission_id` = 12;
update ppm_rol_role_permission_operation set `permission_id` = 16 where `permission_path` like '%/Priority' and `permission_id` = 11;
update ppm_rol_role_permission_operation set `permission_id` = 17 where `permission_path` like '%/ProcessStatus' and `permission_id` = 12;
update ppm_rol_role_permission_operation set `permission_id` = 18 where `permission_path` like '%/Process' and `permission_id` = 11;
update ppm_rol_role_permission_operation set `permission_id` = 19 where `permission_path` like '%/ProcessStep' and `permission_id` = 12;

update ppm_rol_role_permission_operation set `operation_codes` = 'View' where permission_id in (21, 33) and role_id in (
select id from ppm_rol_role where `lang_code` in ('RoleGroup.Pro.Member', 'RoleGroup.Special.Worker', 'RoleGroup.Special.Attention');
)

update ppm_rol_role_permission_operation set `operation_codes` = '(View)|(Create)|(Attention)|(UnAttention)' where permission_id = 12;