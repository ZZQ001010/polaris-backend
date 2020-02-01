CREATE TABLE if not exists `ppm_tre_notice` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `org_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '组织id',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '任务类型, 1项目通知,2组织通知,',
  `project_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '项目id',
  `issue_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'issueId',
  `trends_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联动态id',
  `content` varchar(4096) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '通知内容',
  `noticer` bigint(20) NOT NULL DEFAULT '0' COMMENT '被通知人',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态, 1未读,2已读',
  `creator` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updator` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新人',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(11) NOT NULL DEFAULT '1' COMMENT '乐观锁',
  `is_delete` tinyint(4) NOT NULL DEFAULT '2' COMMENT '是否删除,1是,2否',
  PRIMARY KEY (`id`),
  KEY `index_ppm_tre_notice_issue_id` (`issue_id`),
  KEY `index_ppm_tre_notice_project_id` (`project_id`),
  KEY `index_ppm_tre_notice_noticer` (`noticer`),
  KEY `index_ppm_tre_notice_org_id` (`org_id`),
  KEY `index_ppm_tre_notice_create_time` (`create_time`),
  KEY `index_ppm_tre_notice_trends_id` (`trends_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


