/*==============================================================*/
/* Table: ppm_rol_operation                                     */
/*==============================================================*/
create table if not exists ppm_rol_operation
(
   id                   bigint not null,
   code                 varchar(16) not null default '',
   name                 varchar(128) not null default '',
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
/* Index: index_ppm_rol_operation_code                          */
/*==============================================================*/
create index index_ppm_rol_operation_code on ppm_rol_operation
(
   code
);

/*==============================================================*/
/* Index: index_ppm_rol_operation_create_time                   */
/*==============================================================*/
create index index_ppm_rol_operation_create_time on ppm_rol_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_permission                                    */
/*==============================================================*/
create table if not exists ppm_rol_permission
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   code                 varchar(32) not null default '',
   name                 varchar(128) not null default '',
   parent_id            bigint not null default 0,
   type                 tinyint not null default 1,
   path                 varchar(512) not null default '',
   is_show              tinyint not null default 1,
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
/* Index: index_ppm_rol_permission_org_id                       */
/*==============================================================*/
create index index_ppm_rol_permission_org_id on ppm_rol_permission
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_parent_id                    */
/*==============================================================*/
create index index_ppm_rol_permission_parent_id on ppm_rol_permission
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_path                         */
/*==============================================================*/
create index index_ppm_rol_permission_path on ppm_rol_permission
(
   path
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_create_time                  */
/*==============================================================*/
create index index_ppm_rol_permission_create_time on ppm_rol_permission
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_permission_operation                          */
/*==============================================================*/
create table if not exists ppm_rol_permission_operation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   permission_id        bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(128) not null default '',
   operation_codes      varchar(128) not null default '',
   remark               varchar(512) not null default '',
   is_show              tinyint not null default 1,
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
/* Index: index_ppm_rol_permission_operation_org_id             */
/*==============================================================*/
create index index_ppm_rol_permission_operation_org_id on ppm_rol_permission_operation
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_operation_permission_id      */
/*==============================================================*/
create index index_ppm_rol_permission_operation_permission_id on ppm_rol_permission_operation
(
   permission_id
);

/*==============================================================*/
/* Index: index_ppm_rol_permission_operation_create_time        */
/*==============================================================*/
create index index_ppm_rol_permission_operation_create_time on ppm_rol_permission_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role                                          */
/*==============================================================*/
create table if not exists ppm_rol_role
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   is_readonly          tinyint not null default 2,
   is_modify_permission tinyint not null default 1,
   is_default           tinyint not null default 2,
   role_group_id        bigint not null default 0,
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
/* Index: index_ppm_rol_role_org_id                             */
/*==============================================================*/
create index index_ppm_rol_role_org_id on ppm_rol_role
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_role_group_id                      */
/*==============================================================*/
create index index_ppm_rol_role_role_group_id on ppm_rol_role
(
   role_group_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_create_time                        */
/*==============================================================*/
create index index_ppm_rol_role_create_time on ppm_rol_role
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_group                                    */
/*==============================================================*/
create table if not exists ppm_rol_role_group
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   remark               varchar(512) not null default '',
   type                 tinyint not null default 1,
   is_readonly          tinyint not null default 2,
   is_show              tinyint not null default 1,
   is_default           tinyint not null default 2,
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
/* Index: index_ppm_rol_role_group_org_id                       */
/*==============================================================*/
create index index_ppm_rol_role_group_org_id on ppm_rol_role_group
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_group_create_time                  */
/*==============================================================*/
create index index_ppm_rol_role_group_create_time on ppm_rol_role_group
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_permission_operation                     */
/*==============================================================*/
create table if not exists ppm_rol_role_permission_operation
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   role_id              bigint not null default 0,
   project_id           bigint not null default 0,
   permission_id        bigint not null default 0,
   permission_path      varchar(512) not null default '',
   operation_codes      varchar(128) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_id        */
/*==============================================================*/
create index index_ppm_rol_role_permission_operation_org_id on ppm_rol_role_permission_operation
(
   org_id,
   project_id
);

/*==================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_permission_id */
/*==================================================================*/
create index index_ppm_rol_role_permission_operation_org_permission_id on ppm_rol_role_permission_operation
(
   permission_id
);

/*====================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_permission_path */
/*====================================================================*/
create index index_ppm_rol_role_permission_operation_org_permission_path on ppm_rol_role_permission_operation
(
   permission_path
);

/*================================================================*/
/* Index: index_ppm_rol_role_permission_operation_org_create_time */
/*================================================================*/
create index index_ppm_rol_role_permission_operation_org_create_time on ppm_rol_role_permission_operation
(
   create_time
);

/*==============================================================*/
/* Table: ppm_rol_role_user                                     */
/*==============================================================*/
create table if not exists ppm_rol_role_user
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   role_id              bigint not null default 0,
   user_id              bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_org_id                        */
/*==============================================================*/
create index index_ppm_rol_role_user_org_id on ppm_rol_role_user
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_role_id                       */
/*==============================================================*/
create index index_ppm_rol_role_user_role_id on ppm_rol_role_user
(
   role_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_user_id_project_id            */
/*==============================================================*/
create index index_ppm_rol_role_user_user_id_project_id on ppm_rol_role_user
(
   user_id,
   project_id
);

/*==============================================================*/
/* Index: index_ppm_rol_role_user_create_time                   */
/*==============================================================*/
create index index_ppm_rol_role_user_create_time on ppm_rol_role_user
(
   create_time
);


