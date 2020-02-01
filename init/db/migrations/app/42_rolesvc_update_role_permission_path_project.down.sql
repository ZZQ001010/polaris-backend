update ppm_rol_permission set `path` = '/Org/{org_id}/Project' where id = 12;

update ppm_rol_permission_operation set is_delete = 2 where id in (38,41);

update ppm_rol_permission_operation set permission_id = 12 where id = 39;