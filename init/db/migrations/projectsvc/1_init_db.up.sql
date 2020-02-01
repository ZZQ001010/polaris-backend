/*==============================================================*/
/* Table: ppm_pri_issue                                         */
/*==============================================================*/
create table if not exists ppm_pri_issue
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   project_id           bigint not null default 0,
   project_object_type_id bigint not null default 0,
   title                varchar(512) not null default '',
   owner                bigint not null default 0,
   priority_id          bigint not null default 0,
   source_id            bigint not null default 0,
   issus_object_type_id bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   start_time           datetime not null default '1970-01-01 00:00:00',
   end_time             datetime not null default '1970-01-01 00:00:00',
   plan_work_hour       int not null default -1,
   iteration_id         bigint not null default 0,
   version_id           bigint not null default 0,
   module_id            bigint not null default 0,
   parent_id            bigint not null default 0,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_org_id                            */
/*==============================================================*/
create index index_ppm_pri_issue_org_id on ppm_pri_issue
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_project_id_project_objct_type_id  */
/*==============================================================*/
create index index_ppm_pri_issue_project_id_project_objct_type_id on ppm_pri_issue
(
   project_id,
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_code                              */
/*==============================================================*/
create index index_ppm_pri_issue_code on ppm_pri_issue
(
   code
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_handler                           */
/*==============================================================*/
create index index_ppm_pri_issue_handler on ppm_pri_issue
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_parent_id                         */
/*==============================================================*/
create index index_ppm_pri_issue_parent_id on ppm_pri_issue
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_status                            */
/*==============================================================*/
create index index_ppm_pri_issue_status on ppm_pri_issue
(
   status
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_iteration_id                      */
/*==============================================================*/
create index index_ppm_pri_issue_iteration_id on ppm_pri_issue
(
   iteration_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_creator                           */
/*==============================================================*/
create index index_ppm_pri_issue_creator on ppm_pri_issue
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_create_time                       */
/*==============================================================*/
create index index_ppm_pri_issue_create_time on ppm_pri_issue
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_end_time                          */
/*==============================================================*/
create index index_ppm_pri_issue_end_time on ppm_pri_issue
(
   end_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_start_time                        */
/*==============================================================*/
create index index_ppm_pri_issue_start_time on ppm_pri_issue
(
   start_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_plan_start_time                   */
/*==============================================================*/
create index index_ppm_pri_issue_plan_start_time on ppm_pri_issue
(
   plan_start_time
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_plan_end_time                     */
/*==============================================================*/
create index index_ppm_pri_issue_plan_end_time on ppm_pri_issue
(
   plan_end_time
);

/*==============================================================*/
/* Table: ppm_pri_issue_detail                                  */
/*==============================================================*/
create table if not exists ppm_pri_issue_detail
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   issue_id             bigint not null default 0,
   project_id           bigint not null default 0,
   story_point          int not null default -1,
   tags                 varchar(1024) not null default '',
   remark               text,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_org_id                     */
/*==============================================================*/
create index index_ppm_pri_issue_detail_org_id on ppm_pri_issue_detail
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_issue_id                   */
/*==============================================================*/
create index index_ppm_pri_issue_detail_issue_id on ppm_pri_issue_detail
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_project_id                 */
/*==============================================================*/
create index index_ppm_pri_issue_detail_project_id on ppm_pri_issue_detail
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_detail_create_time                */
/*==============================================================*/
create index index_ppm_pri_issue_detail_create_time on ppm_pri_issue_detail
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pri_issue_relation                                */
/*==============================================================*/
create table if not exists ppm_pri_issue_relation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   issue_id             bigint not null default 0,
   relation_id          bigint not null default 0,
   relation_type        tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_project_id               */
/*==============================================================*/
create index index_ppm_pri_issue_relation_project_id on ppm_pri_issue_relation
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_org_id                   */
/*==============================================================*/
create index index_ppm_pri_issue_relation_org_id on ppm_pri_issue_relation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_issue_id                 */
/*==============================================================*/
create index index_ppm_pri_issue_relation_issue_id on ppm_pri_issue_relation
(
   issue_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_relation_id              */
/*==============================================================*/
create index index_ppm_pri_issue_relation_relation_id on ppm_pri_issue_relation
(
   relation_id
);

/*==============================================================*/
/* Index: index_ppm_pri_issue_relation_create_time              */
/*==============================================================*/
create index index_ppm_pri_issue_relation_create_time on ppm_pri_issue_relation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pri_iteration                                     */
/*==============================================================*/
create table if not exists ppm_pri_iteration
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   owner                bigint not null default 0,
   version_id           bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   plan_work_hour       int not null default -1,
   story_point          int not null default -1,
   remark               text,
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_org_id                        */
/*==============================================================*/
create index index_ppm_pri_iteration_org_id on ppm_pri_iteration
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_project_id                    */
/*==============================================================*/
create index index_ppm_pri_iteration_project_id on ppm_pri_iteration
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_owner                         */
/*==============================================================*/
create index index_ppm_pri_iteration_owner on ppm_pri_iteration
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_version_id                    */
/*==============================================================*/
create index index_ppm_pri_iteration_version_id on ppm_pri_iteration
(
   version_id
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_creator                       */
/*==============================================================*/
create index index_ppm_pri_iteration_creator on ppm_pri_iteration
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_pri_iteration_create_time                   */
/*==============================================================*/
create index index_ppm_pri_iteration_create_time on ppm_pri_iteration
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project                                       */
/*==============================================================*/
create table if not exists ppm_pro_project
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   name                 varchar(256) not null default '',
   pre_code             varchar(16) not null default '',
   owner                bigint not null default 0,
   project_type_id      bigint not null default 1,
   priority_id          bigint not null default 0,
   plan_start_time      datetime not null default '1970-01-01 00:00:00',
   plan_end_time        datetime not null default '1970-01-01 00:00:00',
   public_status        tinyint not null default 1,
   resource_id          bigint not null default 0,
   is_filing            tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_org_id                          */
/*==============================================================*/
create index index_ppm_pro_project_org_id on ppm_pro_project
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_pre_code                        */
/*==============================================================*/
create index index_ppm_pro_project_pre_code on ppm_pro_project
(
   pre_code
);

/*==============================================================*/
/* Index: index_ppm_pro_project_code                            */
/*==============================================================*/
create index index_ppm_pro_project_code on ppm_pro_project
(
   code
);

/*==============================================================*/
/* Index: index_ppm_pro_project_owner                           */
/*==============================================================*/
create index index_ppm_pro_project_owner on ppm_pro_project
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_name                            */
/*==============================================================*/
create index index_ppm_pro_project_name on ppm_pro_project
(
   name
);

/*==============================================================*/
/* Index: index_ppm_pro_project_create_time                     */
/*==============================================================*/
create index index_ppm_pro_project_create_time on ppm_pro_project
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_pro_project_plan_start_time                 */
/*==============================================================*/
create index index_ppm_pro_project_plan_start_time on ppm_pro_project
(
   plan_start_time
);

/*==============================================================*/
/* Index: index_ppm_pro_project_plan_end_time                   */
/*==============================================================*/
create index index_ppm_pro_project_plan_end_time on ppm_pro_project
(
   plan_end_time
);

/*==============================================================*/
/* Table: ppm_pro_project_detail                                */
/*==============================================================*/
create table if not exists ppm_pro_project_detail
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   notice               varchar(4096) not null default '',
   is_enable_work_hours tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_org_id                   */
/*==============================================================*/
create index index_ppm_pro_project_detail_org_id on ppm_pro_project_detail
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_project_id               */
/*==============================================================*/
create index index_ppm_pro_project_detail_project_id on ppm_pro_project_detail
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_detail_create_time              */
/*==============================================================*/
create index index_ppm_pro_project_detail_create_time on ppm_pro_project_detail
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_module                                */
/*==============================================================*/
create table if not exists ppm_pro_project_module
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(32) not null default '',
   owner                bigint not null default 0,
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
/* Index: index_ppm_pro_project_module_org_id                   */
/*==============================================================*/
create index index_ppm_pro_project_module_org_id on ppm_pro_project_module
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_project_id               */
/*==============================================================*/
create index index_ppm_pro_project_module_project_id on ppm_pro_project_module
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_owner                    */
/*==============================================================*/
create index index_ppm_pro_project_module_owner on ppm_pro_project_module
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_module_create_time              */
/*==============================================================*/
create index index_ppm_pro_project_module_create_time on ppm_pro_project_module
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_relation                              */
/*==============================================================*/
create table if not exists ppm_pro_project_relation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   team_id              bigint not null default 0,
   relation_id          bigint not null default 0,
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
/* Index: index_ppm_pro_project_relation_org_id                 */
/*==============================================================*/
create index index_ppm_pro_project_relation_org_id on ppm_pro_project_relation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_project_id             */
/*==============================================================*/
create index index_ppm_pro_project_relation_project_id on ppm_pro_project_relation
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_relation_id            */
/*==============================================================*/
create index index_ppm_pro_project_relation_relation_id on ppm_pro_project_relation
(
   relation_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_relation_create_time            */
/*==============================================================*/
create index index_ppm_pro_project_relation_create_time on ppm_pro_project_relation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_pro_project_version                               */
/*==============================================================*/
create table if not exists ppm_pro_project_version
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(64) not null default '',
   code                 varchar(32) not null default '',
   owner                bigint not null default 0,
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
/* Index: index_ppm_pro_project_version_org_id                  */
/*==============================================================*/
create index index_ppm_pro_project_version_org_id on ppm_pro_project_version
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_project_id              */
/*==============================================================*/
create index index_ppm_pro_project_version_project_id on ppm_pro_project_version
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_owner                   */
/*==============================================================*/
create index index_ppm_pro_project_version_owner on ppm_pro_project_version
(
   owner
);

/*==============================================================*/
/* Index: index_ppm_pro_project_version_create_time             */
/*==============================================================*/
create index index_ppm_pro_project_version_create_time on ppm_pro_project_version
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_issue_object_type                             */
/*==============================================================*/
create table if not exists ppm_prs_issue_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   sort                 int not null default 0,
   project_object_type_id bigint not null default 0,
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
/* Index: index_ppm_prs_issue_object_type_org_id                */
/*==============================================================*/
create index index_ppm_prs_issue_object_type_org_id on ppm_prs_issue_object_type
(
   org_id
);

/*===============================================================*/
/* Index: index_ppm_prs_issue_object_type_project_object_type_id */
/*===============================================================*/
create index index_ppm_prs_issue_object_type_project_object_type_id on ppm_prs_issue_object_type
(
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_object_type_create_time           */
/*==============================================================*/
create index index_ppm_prs_issue_object_type_create_time on ppm_prs_issue_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_issue_source                                  */
/*==============================================================*/
create table if not exists ppm_prs_issue_source
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   sort                 int not null default 0,
   project_object_type_id bigint not null default 0,
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
/* Index: index_ppm_prs_issue_source_org_id                     */
/*==============================================================*/
create index index_ppm_prs_issue_source_org_id on ppm_prs_issue_source
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_source_project_object_type_id     */
/*==============================================================*/
create index index_ppm_prs_issue_source_project_object_type_id on ppm_prs_issue_source
(
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_issue_source_create_time                */
/*==============================================================*/
create index index_ppm_prs_issue_source_create_time on ppm_prs_issue_source
(
   create_time
);



/*==============================================================*/
/* Table: ppm_prs_priority                                      */
/*==============================================================*/
create table if not exists ppm_prs_priority
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   type                 tinyint not null default 1,
   sort                 int not null default 1,
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
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
/* Index: index_ppm_prs_priority_org_id                         */
/*==============================================================*/
create index index_ppm_prs_priority_org_id on ppm_prs_priority
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_priority_create_time                    */
/*==============================================================*/
create index index_ppm_prs_priority_create_time on ppm_prs_priority
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_object_type                           */
/*==============================================================*/
create table if not exists ppm_prs_project_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   pre_code             varchar(16) not null default '',
   name                 varchar(32) not null default '',
   icon                 varchar(8) not null default '',
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
   object_type          tinyint not null default 1,
   sort                 int not null default 1,
   remark               varchar(512) not null default '',
   is_readonly          tinyint not null default 2,
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
/* Index: index_ppm_prs_project_object_type_org_id_pre_code     */
/*==============================================================*/
create index index_ppm_prs_project_object_type_org_id_pre_code on ppm_prs_project_object_type
(
   org_id,
   pre_code
);

/*==============================================================*/
/* Index: index_ppm_prs_project_objct_type_create_time          */
/*==============================================================*/
create index index_ppm_prs_project_objct_type_create_time on ppm_prs_project_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_object_type_process                   */
/*==============================================================*/
create table if not exists ppm_prs_project_object_type_process
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   project_object_type_id bigint not null default 0,
   process_id           bigint not null default 0,
   sort                 int not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_org_id      */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_org_id on ppm_prs_project_object_type_process
(
   org_id,
   project_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_project_id  */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_project_id on ppm_prs_project_object_type_process
(
   project_id,
   project_object_type_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_object_type_process_create_time */
/*==============================================================*/
create index index_ppm_prs_project_object_type_process_create_time on ppm_prs_project_object_type_process
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_type                                  */
/*==============================================================*/
create table if not exists ppm_prs_project_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   sort                 int not null default 1,
   default_process_id   bigint not null default 0,
   category             bigint not null default 0,
   mode                 tinyint not null default 1,
   is_readonly          tinyint not null default 2,
   remark               varchar(512) not null default '',
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2
);

/*==============================================================*/
/* Index: index_ppm_prs_project_type_org_id                     */
/*==============================================================*/
create index index_ppm_prs_project_type_org_id on ppm_prs_project_type
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_project_type_create_time                */
/*==============================================================*/
create index index_ppm_prs_project_type_create_time on ppm_prs_project_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_project_type_project_object_type              */
/*==============================================================*/
create table if not exists ppm_prs_project_type_project_object_type
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_type_id      bigint not null default 0,
   project_object_type_id bigint not null default 0,
   remark               varchar(512) not null default '',
   default_process_id   bigint not null default 0,
   is_readonly          tinyint not null default 2,
   status               tinyint not null default 1,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*=========================================================================*/
/* Index: index_ppm_prs_project_type_project_object_type_org_id_project_ty */
/*=========================================================================*/
create index index_ppm_prs_project_type_project_object_type_org_id_project_ty on ppm_prs_project_type_project_object_type
(
   org_id,
   project_type_id
);

/*==================================================================*/
/* Index: index_ppm_prs_project_type_project_objct_type_create_time */
/*==================================================================*/
create index index_ppm_prs_project_type_project_objct_type_create_time on ppm_prs_project_type_project_object_type
(
   create_time
);

/*==============================================================*/
/* Table: ppm_sha_share                                         */
/*==============================================================*/
create table if not exists ppm_sha_share
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   name                 varchar(128) not null default '',
   remark               varchar(512) not null default '',
   logo                 varchar(512) not null default '',
   type                 int not null default 1,
   content              text,
   content_md5          varchar(32) not null default '',
   finish_time          datetime not null default '1970-01-01 00:00:00',
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
/* Index: index_ppm_sha_share_org_id                            */
/*==============================================================*/
create index index_ppm_sha_share_org_id on ppm_sha_share
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_sha_share_project_id                        */
/*==============================================================*/
create index index_ppm_sha_share_project_id on ppm_sha_share
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_sha_share_finish_time                       */
/*==============================================================*/
create index index_ppm_sha_share_finish_time on ppm_sha_share
(
   finish_time
);

/*==============================================================*/
/* Index: index_ppm_sha_share_content_md5                       */
/*==============================================================*/
create index index_ppm_sha_share_content_md5 on ppm_sha_share
(
   content_md5
);

/*==============================================================*/
/* Index: index_ppm_sha_share_creator                           */
/*==============================================================*/
create index index_ppm_sha_share_creator on ppm_sha_share
(
   creator
);

/*==============================================================*/
/* Index: index_ppm_sha_share_create_time                       */
/*==============================================================*/
create index index_ppm_sha_share_create_time on ppm_sha_share
(
   create_time
);

/*==============================================================*/
/* Table: ppm_sta_iteration_stat                                */
/*==============================================================*/
create table if not exists ppm_sta_iteration_stat
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   iteration_id         bigint not null default 0,
   issue_count          int not null default 0,
   issue_wait_count     int not null default 0,
   issue_running_count  int not null default 0,
   issue_end_count      int not null default 0,
   demand_count         int not null default 0,
   demand_wait_count    int not null default 0,
   demand_running_count int not null default 0,
   demand_end_count     int not null default 0,
   story_point_count    int not null default 0,
   story_point_wait_count int not null default 0,
   story_point_running_count int not null default 0,
   story_point_end_count int not null default 0,
   task_count           int not null default 0,
   task_wait_count      int not null default 0,
   task_running_count   int not null default 0,
   task_end_count       int not null default 0,
   bug_count            int not null default 0,
   bug_wait_count       int not null default 0,
   bug_running_count    int not null default 0,
   bug_end_count        int not null default 0,
   testtask_count       int not null default 0,
   testtask_wait_count  int not null default 0,
   testtask_running_count int not null default 0,
   testtask_end_count   int not null default 0,
   ext                  varchar(4096) not null default '',
   stat_date            date not null default '1970-01-01',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_org_id                   */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_org_id on ppm_sta_iteration_stat
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_iteration_id             */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_iteration_id on ppm_sta_iteration_stat
(
   iteration_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_project_id               */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_project_id on ppm_sta_iteration_stat
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_create_time              */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_create_time on ppm_sta_iteration_stat
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_sta_iteration_stat_stat_date                */
/*==============================================================*/
create index index_ppm_sta_iteration_stat_stat_date on ppm_sta_iteration_stat
(
   stat_date
);

/*==============================================================*/
/* Table: ppm_sta_project_day_stat                              */
/*==============================================================*/
create table if not exists ppm_sta_project_day_stat
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   issue_count          int not null default 0,
   issue_wait_count     int not null default 0,
   issue_running_count  int not null default 0,
   issue_overdue_count  int not null default 0,
   issue_end_count      int not null default 0,
   demand_count         int not null default 0,
   demand_wait_count    int not null default 0,
   demand_running_count int not null default 0,
   demand_overdue_count int not null default 0,
   demand_end_count     int not null default 0,
   story_point_count    int not null default 0,
   story_point_wait_count int not null default 0,
   story_point_running_count int not null default 0,
   story_point_overdue_count int not null default 0,
   story_point_end_count int not null default 0,
   task_count           int not null default 0,
   task_wait_count      int not null default 0,
   task_running_count   int not null default 0,
   task_overdue_count   int not null default 0,
   task_end_count       int not null default 0,
   bug_count            int not null default 0,
   bug_wait_count       int not null default 0,
   bug_running_count    int not null default 0,
   bug_overdue_count    int not null default 0,
   bug_end_count        int not null default 0,
   testtask_count       int not null default 0,
   testtask_wait_count  int not null default 0,
   testtask_running_count int not null default 0,
   testtask_overdue_count int not null default 0,
   testtask_end_count   int not null default 0,
   ext                  varchar(4096) not null default '',
   stat_date            date not null default '1970-01-01',
   status               bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_org_id                       */
/*==============================================================*/
create index ppm_sta_project_day_stat_org_id on ppm_sta_project_day_stat
(
   org_id
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_project_id                   */
/*==============================================================*/
create index ppm_sta_project_day_stat_project_id on ppm_sta_project_day_stat
(
   project_id
);

/*==============================================================*/
/* Index: ppm_sta_project_day_stat_create_time                  */
/*==============================================================*/
create index ppm_sta_project_day_stat_create_time on ppm_sta_project_day_stat
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_sta_project_day_stat_stat_date              */
/*==============================================================*/
create index index_ppm_sta_project_day_stat_stat_date on ppm_sta_project_day_stat
(
   stat_date
);