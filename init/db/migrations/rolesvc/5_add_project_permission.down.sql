delete from ppm_rol_operation where id between 18 and 24;

delete from ppm_rol_permission where id between 38 and 41;

delete from ppm_rol_permission_operation where id between 142 and 159;

Update ppm_rol_permission_operation set is_delete = 2 where id = 102;
Update ppm_rol_permission_operation set permission_id = 11 where id = 35;
Update ppm_rol_permission_operation set permission_id = 33, `remark` = '', `name` = '绑定/解绑角色用户' where id = 129;
Update ppm_rol_permission set `name` = '项目对象类型管理' where id = 15;
Update ppm_rol_permission_operation set `remark` = '', `name` = '查看看问题对象类型' where id = 50;
Update ppm_rol_permission_operation set `remark` = '', `name` = '编辑看问题对象类型' where id = 51;
Update ppm_rol_permission_operation set `remark` = '', `name` = '创建看问题对象类型' where id = 52;
Update ppm_rol_permission_operation set `remark` = '', `name` = '删除看问题对象类型' where id = 53;