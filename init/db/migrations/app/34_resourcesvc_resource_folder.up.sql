CREATE TABLE `ppm_res_folder` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `org_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '组织id',
  `project_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '项目id',
  `name` varchar(256) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '文件夹名',
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父级文件夹id',
  `file_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '文件夹类型,0其他,1文档,2图片,3视频,4音频',
  `creator` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建人',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updator` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新人',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(11) NOT NULL DEFAULT '1' COMMENT '乐观锁',
  `is_delete` tinyint(4) NOT NULL DEFAULT '2' COMMENT '是否删除,1是,2否',
  PRIMARY KEY (`id`),
  KEY `index_ppm_res_folder_parent_id` (`parent_id`),
  KEY `index_ppm_res_folder_org_id` (`org_id`),
  KEY `index_ppm_res_folder_project_id` (`project_id`),
  KEY `index_ppm_res_folder_create_time` (`create_time`),
  KEY `index_ppm_res_folder_creator` (`creator`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `ppm_res_folder_resource` (
    `id` bigint(20) NOT NULL COMMENT '主键',
    `org_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '组织id',
    `resource_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资源id',
    `folder_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件夹id',
    `creator` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建人',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updator` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新人',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `version` int(11) NOT NULL DEFAULT '1' COMMENT '乐观锁',
    `is_delete` tinyint(4) NOT NULL DEFAULT '2' COMMENT '是否删除,1是,2否',
    PRIMARY KEY (`id`),
    KEY `index_ppm_res_folder_resource_resource_id` (`resource_id`),
    KEY `index_ppm_res_folder_resource_org_id` (`org_id`),
    KEY `index_ppm_res_folder_resource_folder_id` (`folder_id`),
    KEY `index_ppm_res_folder_resource_create_time` (`create_time`),
    KEY `index_ppm_res_folder_resource_creator` (`creator`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;