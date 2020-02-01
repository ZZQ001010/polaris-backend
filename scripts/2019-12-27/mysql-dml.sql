update ppm_pri_tag t LEFT JOIN ppm_pri_issue_tag i on t.id = i.tag_id set t.project_id = i.project_id where t.id = i.tag_id;

update ppm_rol_role_permission_operation set operation_codes = "(View)|(Remove)|(Create)|(Delete)|(Modify)" where is_delete = 2 and updator = 0 and permission_id = 39

update ppm_rol_permission_operation set operation_codes = "Modify,Bind,Unbind" where operation_codes = "Modify,Create,Bind,Unbind"