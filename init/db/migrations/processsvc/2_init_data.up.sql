-- 流程状态

-- INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
--         VALUES ( 1, 0, 'ProcessStatus.Project.NotStart', '未开始', 1, '#008FFF', '#008FFF', 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'ProcessStatus.Project.Running', '进行中', 2, '#F1A102', '#F1A102', 2, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'ProcessStatus.Project.Complete', '已完成', 3, '#25B47E', '#25B47E', 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'ProcessStatus.Iteration.NotStart', '未开始', 1, '#DBDBDB', '#FFFFFF', 1, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'ProcessStatus.Iteration.Running', '进行中', 2, '#FFCD1C', '#FFFFFF', 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'ProcessStatus.Iteration.Complete', '已完成', 3, '#69A922', '#FFFFFF', 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 'ProcessStatus.Issue.NotStart', '未完成', 1, '#CACACA', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 'ProcessStatus.Issue.WaitEvaluate', '待评估', 2, '#DBDBDB', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 'ProcessStatus.Issue.WaitConfirmBug', '待确认缺陷', 3, '#DBDBDB', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 'ProcessStatus.Issue.ConfirmedBug', '已确认缺陷', 4, '#DBDBDB', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 'ProcessStatus.Issue.ReOpen', '重新打开', 5, '#DBDBDB', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 'ProcessStatus.Issue.Evaluated', '已评估', 6, '#FFCD1C', '#FFFFFF', 1, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 'ProcessStatus.Issue.Planning', '计划中', 7, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 'ProcessStatus.Issue.Design', '设计中', 8, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 'ProcessStatus.Issue.Processing', '处理中', 9, '#FFC700', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 'ProcessStatus.Issue.Development', '研发中', 10, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 'ProcessStatus.Issue.WaitTest', '待测试', 11, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 'ProcessStatus.Issue.Testing', '测试中', 12, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 'ProcessStatus.Issue.WaitRelease', '待发布', 13, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 'ProcessStatus.Issue.Wait', '挂起', 14, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 'ProcessStatus.Issue.Repair', '修复中', 15, '#FFCD1C', '#FFFFFF', 2, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 'ProcessStatus.Issue.Success', '成功', 16, '#69A922', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 'ProcessStatus.Issue.Fail', '失败', 17, '#FD1F00', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 'ProcessStatus.Issue.Released', '已发布', 18, '#69A922', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 'ProcessStatus.Issue.Closed', '已关闭', 19, '#69A922', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 'ProcessStatus.Issue.Complete', '已完成', 20, '#67D287', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_status( `id`, `org_id`, `lang_code`, `name`, `sort`, `bg_style`, `font_style`, `type`, `category`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 'ProcessStatus.Issue.Confirmed', '已确认', 21, '#69A922', '#FFFFFF', 3, 3, '', 1, 1, now(), 1, now(), 1, 2);


-- 流程

INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 1, 0, 'Process.DefaultProject', '默认项目流程', 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 'Process.DefaultIteration', '默认迭代流程', 1, 2, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 'Process.Issue.DefaultFeature', '默认特性流程', 2, 3, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 'Process.Issue.DefaultDemand', '默认需求流程', 2, 3, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 'Process.Issue.DefaultTask', '默认任务流程', 1, 3, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 'Process.Issue.DefaultAgileTask', '默认敏捷项目任务流程', 2, 3, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 'Process.Issue.DefaultBug', '默认缺陷流程', 2, 3, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process( `id`, `org_id`, `lang_code`, `name`, `is_default`, `type`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 'Process.Issue.DefaultTestTask', '默认测试任务流程', 2, 3, 6, 1, now(), 1, now(), 1, 2);


-- 流程与流程状态关联

-- INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
--         VALUES ( 1, 0, 1, 1, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 2, 0, 1, 2, 1, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 1, 3, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 4, 0, 2, 4, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 2, 5, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 6, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 3, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 3, 13, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 3, 12, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 3, 16, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 3, 26, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 3, 25, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 4, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 4, 13, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 4, 14, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 4, 16, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 4, 17, 2, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 4, 18, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 4, 19, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 4, 24, 2, 8, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 4, 25, 2, 9, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 5, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 5, 15, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 5, 26, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 6, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 6, 15, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 6, 20, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0, 6, 26, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0, 7, 9, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0, 7, 10, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 31, 0, 7, 21, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 32, 0, 7, 17, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 33, 0, 7, 11, 2, 5, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 34, 0, 7, 19, 2, 6, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 35, 0, 7, 24, 2, 7, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 36, 0, 7, 25, 2, 8, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 37, 0, 8, 7, 1, 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 38, 0, 8, 11, 2, 2, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 39, 0, 8, 20, 2, 3, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 40, 0, 8, 22, 2, 4, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_process_status( `id`, `org_id`, `process_id`, `process_status_id`, `is_init_status`, `sort`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 41, 0, 8, 23, 2, 5, 1, now(), 1, now(), 1, 2);


-- 流程步骤

-- INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
--         VALUES ( 1, 0, 1, 'ProcessStep.Project.Start', '项目开始', 1, 2, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
-- INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
--         VALUES ( 2, 0, 1, 'ProcessStep.Project.Complete', '项目完成', 1, 3, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 3, 0, 1, 'ProcessStep.Project.Complete', '项目完成', 2, 3, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
-- INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
--         VALUES ( 4, 0, 1, 'ProcessStep.Project.NotStart', '项目未开始', 2, 1, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 5, 0, 1, 'ProcessStep.Project.ReStart', '重启项目', 3, 2, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 6, 0, 2, 'ProcessStep.Iteration.Start', '迭代开始', 4, 5, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 7, 0, 2, 'ProcessStep.Iteration.Complete', '迭代完成', 4, 6, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 8, 0, 2, 'ProcessStep.Iteration.Complete', '迭代完成', 5, 6, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 9, 0, 2, 'ProcessStep.Iteration.NotStart', '迭代未开始', 5, 4, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 10, 0, 2, 'ProcessStep.Iteration.ReStart', '重启迭代', 6, 4, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 11, 0, 3, 'ProcessStep.Feature.StartPlan', '开始计划', 7, 13, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 12, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 7, 12, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 13, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 7, 16, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 14, 0, 3, 'ProcessStep.Feature.Released', '发布', 7, 26, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 15, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 7, 25, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 16, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 13, 12, 6, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 17, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 13, 16, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 18, 0, 3, 'ProcessStep.Feature.Released', '发布', 13, 26, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 19, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 13, 25, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 20, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 13, 7, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 21, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 12, 16, 11, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 22, 0, 3, 'ProcessStep.Feature.Released', '发布', 12, 26, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 23, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 12, 25, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 24, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 12, 7, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 25, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 12, 13, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 26, 0, 3, 'ProcessStep.Feature.Released', '发布', 16, 26, 16, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 27, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 16, 25, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 28, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 16, 7, 18, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 29, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 16, 13, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 30, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 16, 12, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 31, 0, 3, 'ProcessStep.Feature.Closed', '关闭', 26, 25, 21, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 32, 0, 3, 'ProcessStep.Feature.NotStart', '未开始', 26, 7, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 33, 0, 3, 'ProcessStep.Feature.ReStartPlan', '重新计划', 26, 13, 23, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 34, 0, 3, 'ProcessStep.Feature.Evaluated', '已评估', 26, 12, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 35, 0, 3, 'ProcessStep.Feature.StartDevelopment', '进入研发', 26, 16, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 36, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 7, 13, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 37, 0, 4, 'ProcessStep.Demand.Evaluated', '已评估', 7, 14, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 38, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 7, 16, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 39, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 7, 25, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 40, 0, 4, 'ProcessStep.Demand.Evaluated', '开始设计', 13, 14, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 41, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 13, 16, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 42, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 13, 25, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 43, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 13, 7, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 44, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 14, 16, 9, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 45, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 14, 25, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 46, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 14, 7, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 47, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 14, 13, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 48, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 16, 17, 13, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 49, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 16, 18, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 50, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 16, 19, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 51, 0, 4, 'ProcessStep.Demand.Released', '发布', 16, 24, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 52, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 16, 25, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 53, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 16, 7, 18, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 54, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 16, 13, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 55, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 16, 14, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 56, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 17, 18, 21, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 57, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 17, 19, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 58, 0, 4, 'ProcessStep.Demand.Released', '发布', 17, 24, 23, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 59, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 17, 25, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 60, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 17, 7, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 61, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 17, 13, 26, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 62, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 17, 14, 27, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 63, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 17, 16, 28, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 64, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 18, 19, 29, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 65, 0, 4, 'ProcessStep.Demand.Released', '发布', 18, 24, 30, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 66, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 18, 25, 31, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 67, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 18, 7, 32, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 68, 0, 4, 'ProcessStep.Demand.StartPlan', '开始计划', 18, 13, 33, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 69, 0, 4, 'ProcessStep.Demand.StartDesign', '开始设计', 18, 14, 34, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 70, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 18, 16, 35, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 71, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 18, 17, 36, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 72, 0, 4, 'ProcessStep.Demand.Released', '发布', 19, 24, 37, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 73, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 19, 25, 38, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 74, 0, 4, 'ProcessStep.Demand.StartDevelopment', '进入研发', 19, 16, 39, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 75, 0, 4, 'ProcessStep.Demand.WaitTest', '等待测试', 19, 17, 40, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 76, 0, 4, 'ProcessStep.Demand.StartTest', '进入测试', 19, 18, 41, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 77, 0, 4, 'ProcessStep.Demand.Closed', '关闭', 24, 25, 42, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 78, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 24, 7, 43, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 79, 0, 4, 'ProcessStep.Demand.WaitReleased', '待发布', 24, 19, 44, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 80, 0, 4, 'ProcessStep.Demand.NotStart', '未开始', 25, 7, 45, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 81, 0, 5, 'ProcessStep.Task.Start', '任务开始', 7, 15, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 82, 0, 5, 'ProcessStep.Task.Complete', '任务完成', 7, 26, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 83, 0, 5, 'ProcessStep.Task.Complete', '任务完成', 15, 26, 3, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 84, 0, 5, 'ProcessStep.Task.NotStart', '任务未开始', 15, 7, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 85, 0, 5, 'ProcessStep.Task.ReStart', '重启任务', 26, 7, 5, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 86, 0, 6, 'ProcessStep.AgileTask.Start', '任务开始', 7, 15, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 87, 0, 6, 'ProcessStep.AgileTask.HangUp', '任务挂起', 7, 20, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 88, 0, 6, 'ProcessStep.AgileTask.Complete', '任务完成', 7, 26, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 89, 0, 6, 'ProcessStep.AgileTask.Complete', '任务完成', 15, 26, 4, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 90, 0, 6, 'ProcessStep.AgileTask.HangUp', '任务挂起', 15, 20, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 91, 0, 6, 'ProcessStep.AgileTask.NotStart', '任务未开始', 15, 7, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 92, 0, 6, 'ProcessStep.AgileTask.ReStart', '重启任务', 20, 7, 7, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 93, 0, 6, 'ProcessStep.AgileTask.ReStart', '重启任务', 26, 7, 8, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 94, 0, 7, 'ProcessStep.Bug.Confirm', '确认', 9, 10, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 95, 0, 7, 'ProcessStep.Bug.Repair', '修复', 9, 21, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 96, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 9, 17, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 97, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 9, 19, 4, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 98, 0, 7, 'ProcessStep.Bug.Released', '已发布', 9, 24, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 99, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 9, 25, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 100, 0, 7, 'ProcessStep.Bug.Repair', '修复', 10, 21, 7, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 101, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 10, 17, 8, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 102, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 10, 19, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 103, 0, 7, 'ProcessStep.Bug.Released', '已发布', 10, 24, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 104, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 10, 25, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 105, 0, 7, 'ProcessStep.Bug.WaitConfirm', '待确认', 10, 9, 12, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 106, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 21, 17, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 107, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 21, 19, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 108, 0, 7, 'ProcessStep.Bug.Released', '已发布', 21, 24, 15, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 109, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 21, 25, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 110, 0, 7, 'ProcessStep.Bug.Confirm', '确认', 21, 10, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 111, 0, 7, 'ProcessStep.Bug.ReOpen', '待发布', 17, 19, 18, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 112, 0, 7, 'ProcessStep.Bug.WaitReleased', '重新打开', 17, 11, 19, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 113, 0, 7, 'ProcessStep.Bug.Released', '已发布', 17, 24, 20, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 114, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 17, 25, 21, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 115, 0, 7, 'ProcessStep.Bug.Repair', '修复', 17, 21, 22, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 116, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 11, 21, 23, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 117, 0, 7, 'ProcessStep.Bug.Released', '待测试', 11, 17, 24, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 118, 0, 7, 'ProcessStep.Bug.Closed', '已发布', 11, 24, 25, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 119, 0, 7, 'ProcessStep.Bug.WaitConfirm', '关闭', 11, 25, 26, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 120, 0, 7, 'ProcessStep.Bug.Repair', '确认', 11, 10, 27, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 121, 0, 7, 'ProcessStep.Bug.Released', '已发布', 19, 24, 28, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 122, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 19, 25, 29, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 123, 0, 7, 'ProcessStep.Bug.Repair', '修复', 19, 21, 30, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 124, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 19, 17, 31, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 125, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 19, 11, 32, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 126, 0, 7, 'ProcessStep.Bug.Closed', '关闭', 24, 25, 33, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 127, 0, 7, 'ProcessStep.Bug.Repair', '修复', 24, 21, 34, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 128, 0, 7, 'ProcessStep.Bug.WaitTest', '待测试', 24, 17, 35, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 129, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 24, 11, 36, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 130, 0, 7, 'ProcessStep.Bug.WaitReleased', '待发布', 24, 19, 37, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 131, 0, 7, 'ProcessStep.Bug.ReOpen', '重新打开', 25, 11, 38, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 132, 0, 8, 'ProcessStep.TestTask.Success', '成功', 7, 22, 1, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 133, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 7, 23, 2, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 134, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 7, 20, 3, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 135, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 22, 23, 4, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 136, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 22, 20, 5, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 137, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 22, 11, 6, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 138, 0, 8, 'ProcessStep.TestTask.NotStart', '未开始', 22, 7, 7, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 139, 0, 8, 'ProcessStep.TestTask.Success', '成功', 23, 22, 8, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 140, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 23, 20, 9, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 141, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 23, 11, 10, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 142, 0, 8, 'ProcessStep.TestTask.NotStart', '未开始', 23, 7, 11, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 143, 0, 8, 'ProcessStep.TestTask.ReOpen', '重新打开', 20, 11, 12, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 144, 0, 8, 'ProcessStep.TestTask.Success', '成功', 20, 22, 13, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 145, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 20, 23, 14, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 146, 0, 8, 'ProcessStep.TestTask.Success', '成功', 11, 22, 15, 1, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 147, 0, 8, 'ProcessStep.TestTask.Fail', '失败', 11, 23, 16, 2, '', 1, 1, now(), 1, now(), 1, 2);
INSERT INTO ppm_prs_process_step( `id`, `org_id`, `process_id`, `lang_code`, `name`, `start_status`, `end_status`, `sort`, `is_default`, `remark`, `status`, `creator`, `create_time`, `updator`, `update_time`, `version`, `is_delete` )
        VALUES ( 148, 0, 8, 'ProcessStep.TestTask.Wait', '挂起', 11, 20, 17, 2, '', 1, 1, now(), 1, now(), 1, 2);