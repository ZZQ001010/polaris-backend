ALTER TABLE `ppm_prs_project_type` ADD PRIMARY KEY (`id`);

UPDATE ppm_prs_project_type SET `remark` = '简单易用的通用任务处理模板，适用 于诸如个人安排等活动管理。' where `id` = 1;
UPDATE ppm_prs_project_type SET `remark` = '通过内置的敏捷研发管理组件，可以轻 松实现迭代管控、需求分配、缺陷管理 等研发工作，实时掌控项目进度状况。' where `id` = 2;
