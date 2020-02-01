

INSERT INTO ppm_prs_project_type( `id`, `org_id`, `lang_code`, `name`, `sort`, `default_process_id`, `category`, `mode`, `is_readonly`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'ProjectType.NormalTask', '普通任务项目', 1, 1, 0, 1, 1, '简单易用的通用任务处理模板，适用 于诸如个人安排等活动管理。', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type( `id`, `org_id`, `lang_code`, `name`, `sort`, `default_process_id`, `category`, `mode`, `is_readonly`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'ProjectType.Agile', '敏捷研发项目', 2, 1, 0, 2, 1, '通过内置的敏捷研发管理组件，可以轻 松实现迭代管控、需求分配、缺陷管理 等研发工作，实时掌控项目进度状况。', 1, 1, now(), 1, now(), 1, 2);


INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'Project.ObjectType.Iteration', 'I', '迭代', 1, '', '', '', 1, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'Project.ObjectType.Feature', 'F', '特性', 2, '', '', '', 2, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'Project.ObjectType.Demand', 'D', '需求', 2, '', '', '', 3, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'Project.ObjectType.Task', 'T', '任务', 2, '', '', '', 4, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'Project.ObjectType.Bug', 'B', '缺陷', 2, '', '', '', 5, '', 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_object_type( `id`, `org_id`, `lang_code`, `pre_code`, `name`, `object_type`, `bg_style`, `font_style`, `icon`, `sort`, `remark`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'Project.ObjectType.TestTask', 'TT', '测试任务', 2, '', '', '', 6, '', 1, 1, 1, now(), 1, now(), 1, 2);

INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 1, 0, '', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 1, 4, '', 5, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 2, 0, '', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 2, 1, '', 2, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 2, 2, '', 3, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 3, '', 4, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 2, 4, '', 6, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 2, 5, '', 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_project_type_project_object_type( `id`, `org_id`, `project_type_id`, `project_object_type_id`, `remark`, `default_process_id`, `is_readonly`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 2, 6, '', 8, 1, 1, 1, now(), 1, now(), 1, 2);

