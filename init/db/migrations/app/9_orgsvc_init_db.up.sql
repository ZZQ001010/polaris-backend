
/*==============================================================*/
/* Table: ppm_bas_source_channel                                */
/*==============================================================*/
create table if not exists ppm_bas_source_channel
(
   id                   bigint not null,
   code                 varchar(32) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_source_channel_create_time              */
/*==============================================================*/
create index index_ppm_bas_source_channel_create_time on ppm_bas_source_channel
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_source_channel_code                     */
/*==============================================================*/
create index index_ppm_bas_source_channel_code on ppm_bas_source_channel
(
   code
);

/*==============================================================*/
/* Table: ppm_bas_pay_level                                     */
/*==============================================================*/
create table if not exists ppm_bas_pay_level
(
   id                   bigint not null,
   lang_code            varchar(64) not null default 'zh-CN',
   name                 varchar(32) not null default '',
   storage              bigint not null default 0,
   member_count         int not null default 10,
   price                bigint not null default 0,
   member_price         bigint not null default 0,
   duration             bigint not null default 0,
   is_show              tinyint not null default 1,
   sort                 int not null default 0,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_pay_level_create_time                   */
/*==============================================================*/
create index index_ppm_bas_pay_level_create_time on ppm_bas_pay_level
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_pay_level_sort                          */
/*==============================================================*/
create index index_ppm_bas_pay_level_sort on ppm_bas_pay_level
(
   sort
);


/*==============================================================*/
/* Table: ppm_orc_config                                        */
/*==============================================================*/
create table if not exists ppm_orc_config
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   time_zone            varchar(32) not null default 'Asia/Shanghai',
   time_difference      varchar(8) not null default '+08:00',
   pay_level            smallint not null default 1,
   pay_start_time       datetime not null default CURRENT_TIMESTAMP,
   pay_end_time         datetime not null default '2038-01-01 00:00:00',
   web_site             varchar(256) not null default '',
   language             varchar(8) not null default 'zh-CN',
   datetime_format      varchar(32) not null default 'yyyy-MM-dd HH:mm:ss',
   password_length      tinyint not null default 6,
   password_rule        tinyint not null default 1,
   max_login_fail_count int not null default 0,
   remind_send_time     varchar(8) not null default '09:00',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_orc_config_org_id                           */
/*==============================================================*/
create index index_ppm_orc_config_org_id on ppm_orc_config
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_orc_config_create_time                      */
/*==============================================================*/
create index index_ppm_orc_config_create_time on ppm_orc_config
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_department                                    */
/*==============================================================*/
create table if not exists ppm_org_department
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(64) not null default '',
   parent_id            bigint not null default 0,
   sort                 int not null default 0,
   is_hide              tinyint not null default 2,
   source_channel       varchar(16) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_department_org_id                       */
/*==============================================================*/
create index index_ppm_org_department_org_id on ppm_org_department
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_code                         */
/*==============================================================*/
create index index_ppm_org_department_code on ppm_org_department
(
   code
);

/*==============================================================*/
/* Index: index_ppm_org_department_parent_id                    */
/*==============================================================*/
create index index_ppm_org_department_parent_id on ppm_org_department
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_create_time                  */
/*==============================================================*/
create index index_ppm_org_department_create_time on ppm_org_department
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_department_out_info                           */
/*==============================================================*/
create table if not exists ppm_org_department_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   department_id        bigint not null default 0,
   source_channel       varchar(16) not null default '',
   out_org_department_id varchar(64) not null default '',
   out_org_department_code varchar(64) not null default '',
   name                 varchar(64) not null default '',
   out_org_department_parent_id varchar(64) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_org_id              */
/*==============================================================*/
create index index_ppm_org_department_out_info_org_id on ppm_org_department_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_department_id       */
/*==============================================================*/
create index index_ppm_org_department_out_info_department_id on ppm_org_department_out_info
(
   department_id
);

/*================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_id */
/*================================================================*/
create index index_ppm_org_department_out_info_out_org_department_id on ppm_org_department_out_info
(
   out_org_department_id
);

/*==================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_code */
/*==================================================================*/
create index index_ppm_org_department_out_info_out_org_department_code on ppm_org_department_out_info
(
   out_org_department_code
);

/*==============================================================*/
/* Index: index_ppm_org_department_out_info_create_time         */
/*==============================================================*/
create index index_ppm_org_department_out_info_create_time on ppm_org_department_out_info
(
   create_time
);

/*=======================================================================*/
/* Index: index_ppm_org_department_out_info_out_org_department_parent_id */
/*=======================================================================*/
create index index_ppm_org_department_out_info_out_org_department_parent_id on ppm_org_department_out_info
(
   out_org_department_parent_id
);

/*==============================================================*/
/* Table: ppm_org_organization                                  */
/*==============================================================*/
create table if not exists ppm_org_organization
(
   id                   bigint not null,
   name                 varchar(256) not null default '',
   web_site             varchar(512) not null default '',
   industry_id          bigint not null default 0,
   scale                varchar(32) not null default '',
   source_channel       varchar(16) not null default '',
   country_id           bigint not null default 0,
   province_id          bigint not null default 0,
   city_id              bigint not null default 0,
   address              varchar(256) not null default '',
   logo_url             varchar(512) not null default '',
   resource_id          bigint not null default 0,
   owner                bigint not null default 0,
   is_authenticated     tinyint not null default 1,
   remark               varchar(512) not null default '',
   init_status          tinyint not null default 1,
   init_version         int not null default 1,
   status               tinyint not null default 1,
   is_show              tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_organization_create_time                */
/*==============================================================*/
create index index_ppm_org_organization_create_time on ppm_org_organization
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_organization_city_id                    */
/*==============================================================*/
create index index_ppm_org_organization_city_id on ppm_org_organization
(
   city_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_province_id                */
/*==============================================================*/
create index index_ppm_org_organization_province_id on ppm_org_organization
(
   province_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_country_id                 */
/*==============================================================*/
create index index_ppm_org_organization_country_id on ppm_org_organization
(
   country_id
);

/*==============================================================*/
/* Table: ppm_org_organization_out_info                         */
/*==============================================================*/
create table if not exists ppm_org_organization_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   out_org_id           varchar(64) not null default '',
   source_channel       varchar(16) not null default '',
   name                 varchar(64) not null default '',
   industry             varchar(64) not null default '',
   is_authenticated     tinyint not null default 1,
   auth_ticket          varchar(256) not null default '',
   auth_level           varchar(32) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_org_id            */
/*==============================================================*/
create index index_ppm_org_organization_out_info_org_id on ppm_org_organization_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_out_org_id        */
/*==============================================================*/
create index index_ppm_org_organization_out_info_out_org_id on ppm_org_organization_out_info
(
   out_org_id
);

/*==============================================================*/
/* Index: index_ppm_org_organization_out_info_create_time       */
/*==============================================================*/
create index index_ppm_org_organization_out_info_create_time on ppm_org_organization_out_info
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_user                                          */
/*==============================================================*/
CREATE TABLE `ppm_org_user` (
  `id` bigint(20) NOT NULL,
  `org_id` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `name_pinyin` varchar(128) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `login_name` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `login_name_edit_count` int(11) NOT NULL DEFAULT '0',
  `email` varchar(128) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `mobile_region` varchar(8) COLLATE utf8mb4_bin NOT NULL DEFAULT '+86',
  `mobile` varchar(16) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `avatar` varchar(1024) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `birthday` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `sex` tinyint(4) NOT NULL DEFAULT '99',
  `password` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `password_salt` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `source_platform` varchar(16) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `source_channel` varchar(16) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `source_obj_id` varchar(32) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `language` varchar(8) COLLATE utf8mb4_bin NOT NULL DEFAULT 'zh-CN',
  `motto` varchar(512) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `last_login_ip` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `last_login_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `login_fail_count` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `creator` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updator` bigint(20) NOT NULL DEFAULT '0' ,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
  `version` int(11) NOT NULL DEFAULT '1' ,
  `is_delete` tinyint(4) NOT NULL DEFAULT '2' ,
  PRIMARY KEY (`id`),
  KEY `index_ppm_org_user_org_id` (`org_id`),
  KEY `index_ppm_org_user_mobile` (`mobile`),
  KEY `index_ppm_org_user_name_pinyin` (`name_pinyin`),
  KEY `index_ppm_org_user_name` (`name`),
  KEY `index_ppm_org_user_login_name` (`login_name`),
  KEY `index_ppm_org_user_email` (`email`),
  KEY `index_ppm_org_user_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*==============================================================*/
/* Table: ppm_org_user_config                                   */
/*==============================================================*/
create table if not exists ppm_org_user_config
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   daily_report_message_status tinyint not null default 2,
   owner_range_status   tinyint not null default 1,
   participant_range_status tinyint not null default 1,
   attention_range_status tinyint not null default 1,
   create_range_status  tinyint not null default 2,
   remind_message_status tinyint not null default 2,
   comment_at_message_status tinyint not null default 1,
   modify_message_status tinyint not null default 1,
   relation_message_status tinyint not null default 2,
   ext                  varchar(4096) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_user_id                     */
/*==============================================================*/
create index index_ppm_org_user_config_user_id on ppm_org_user_config
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_org_id                      */
/*==============================================================*/
create index index_ppm_org_user_config_org_id on ppm_org_user_config
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_config_create_time                 */
/*==============================================================*/
create index index_ppm_org_user_config_create_time on ppm_org_user_config
(
   create_time
);

/*==============================================================*/
/* Table: ppm_org_user_department                               */
/*==============================================================*/
create table if not exists ppm_org_user_department
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   department_id        bigint not null default 0,
   is_leader            tinyint not null default 2,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_org_id                  */
/*==============================================================*/
create index index_ppm_org_user_department_org_id on ppm_org_user_department
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_create_time             */
/*==============================================================*/
create index index_ppm_org_user_department_create_time on ppm_org_user_department
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_user_id                 */
/*==============================================================*/
create index index_ppm_org_user_department_user_id on ppm_org_user_department
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_department_department_id           */
/*==============================================================*/
create index index_ppm_org_user_department_department_id on ppm_org_user_department
(
   department_id
);

/*==============================================================*/
/* Table: ppm_org_user_organization                             */
/*==============================================================*/
create table if not exists ppm_org_user_organization
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   check_status         tinyint not null default 1,
   use_status           tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_org_id                */
/*==============================================================*/
create index index_ppm_org_user_organization_org_id on ppm_org_user_organization
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_create_time           */
/*==============================================================*/
create index index_ppm_org_user_organization_create_time on ppm_org_user_organization
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_org_user_organization_user_id               */
/*==============================================================*/
create index index_ppm_org_user_organization_user_id on ppm_org_user_organization
(
   user_id
);

/*==============================================================*/
/* Table: ppm_org_user_out_info                                 */
/*==============================================================*/
create table if not exists ppm_org_user_out_info
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   user_id              bigint not null default 0,
   source_channel       varchar(16) not null default '',
   out_org_user_id      varchar(64) not null default '',
   out_user_id          varchar(64) not null default '',
   name                 varchar(64) not null default '',
   avatar               varchar(1024) not null default '',
   is_active            tinyint not null default 1,
   job_number           varchar(32) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_org_id                    */
/*==============================================================*/
create index index_ppm_org_user_out_info_org_id on ppm_org_user_out_info
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_user_id                   */
/*==============================================================*/
create index index_ppm_org_user_out_info_user_id on ppm_org_user_out_info
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_out_org_user_id           */
/*==============================================================*/
create index index_ppm_org_user_out_info_out_org_user_id on ppm_org_user_out_info
(
   out_org_user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_out_user_id               */
/*==============================================================*/
create index index_ppm_org_user_out_info_out_user_id on ppm_org_user_out_info
(
   out_user_id
);

/*==============================================================*/
/* Index: index_ppm_org_user_out_info_create_time               */
/*==============================================================*/
create index index_ppm_org_user_out_info_create_time on ppm_org_user_out_info
(
   create_time
);

/*==============================================================*/
/* Table: ppm_tem_team                                          */
/*==============================================================*/
create table if not exists ppm_tem_team
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   name                 varchar(128) not null default '',
   nick_name            varchar(128) not null default '',
   owner                bigint not null default 0,
   department_id        bigint not null default 0,
   is_default           tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tem_team_org_id                             */
/*==============================================================*/
create index index_ppm_tem_team_org_id on ppm_tem_team
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tem_team_create_time                        */
/*==============================================================*/
create index index_ppm_tem_team_create_time on ppm_tem_team
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tem_team_department_id                      */
/*==============================================================*/
create index index_ppm_tem_team_department_id on ppm_tem_team
(
   department_id
);

/*==============================================================*/
/* Table: ppm_tem_user_team                                     */
/*==============================================================*/
create table if not exists ppm_tem_user_team
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   team_id              bigint not null default 0,
   user_id              bigint not null default 0,
   relation_type        tinyint not null default 1,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_org_id                        */
/*==============================================================*/
create index index_ppm_tem_user_team_org_id on ppm_tem_user_team
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_team_id                       */
/*==============================================================*/
create index index_ppm_tem_user_team_team_id on ppm_tem_user_team
(
   team_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_user_id                       */
/*==============================================================*/
create index index_ppm_tem_user_team_user_id on ppm_tem_user_team
(
   user_id
);

/*==============================================================*/
/* Index: index_ppm_tem_user_team_create_time                   */
/*==============================================================*/
create index index_ppm_tem_user_team_create_time on ppm_tem_user_team
(
   create_time
);

