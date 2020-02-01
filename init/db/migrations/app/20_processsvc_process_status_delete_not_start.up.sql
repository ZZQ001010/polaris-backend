update ppm_prs_process_status set is_delete = 1 where id = 1;

update ppm_prs_process_process_status set is_delete = 1 where id = 1;
update ppm_prs_process_process_status set `is_init_status` = 1 where id = 2;

update ppm_prs_process_step set is_delete = 1 where id in (1,2,4);
update ppm_prs_process_step set `end_status` = 2 where id = 5;