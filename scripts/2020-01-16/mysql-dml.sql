update ppm_prs_process_status set bg_style = "#67D287" where id = 26;

update ppm_prs_process_status set bg_style = "#CACACA" where id = 7;

update ppm_prs_process_status set bg_style = "#FFC700" where id = 15;

update ppm_rol_role_permission_operation set is_delete = 1 where role_id = 8 and permission_id in (21,22,23,25,26,27,28,29);

update ppm_prs_process_status set name = "未完成" where id = 7;

update ppm_rol_role_permission_operation set is_delete = 2 and operation_codes = "(View)|(Comment)|(Attention)|(UnAttention)|(Create)" where role_id = 8 and permission_id in (23,25,26,27,28,29);