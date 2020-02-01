/*==============================================================*/
/* Table: ppm_tre_comment                                       */
/*==============================================================*/
create table if not exists ppm_tre_comment
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   trends_id            bigint not null default 0,
   object_id            bigint not null default 0,
   object_type          varchar(32) not null default 'issue',
   content              varchar(4096) not null default '',
   parent_id            bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_org_id                          */
/*==============================================================*/
create index index_ppm_tre_comment_org_id on ppm_tre_comment
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_create_time                     */
/*==============================================================*/
create index index_ppm_tre_comment_create_time on ppm_tre_comment
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_project_id                      */
/*==============================================================*/
create index index_ppm_tre_comment_project_id on ppm_tre_comment
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_trends_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_trends_id on ppm_tre_comment
(
   trends_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_parent_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_parent_id on ppm_tre_comment
(
   parent_id
);

/*==============================================================*/
/* Index: index_ppm_tre_comment_object_id                       */
/*==============================================================*/
create index index_ppm_tre_comment_object_id on ppm_tre_comment
(
   object_id
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_create_time                 */
/*==============================================================*/
create index index_ppm_tre_trends_oper_create_time on ppm_tre_comment
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_creator                     */
/*==============================================================*/
create index index_ppm_tre_trends_oper_creator on ppm_tre_comment
(
   creator
);

/*==============================================================*/
/* Table: ppm_tre_trends                                        */
/*==============================================================*/
create table if not exists ppm_tre_trends
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   uuid                 varchar(32) not null default '',
   module1              varchar(16) not null default '',
   module2_id           bigint not null,
   module2              varchar(16) not null default '',
   module3_id           bigint not null default 0,
   module3              varchar(16) not null default '',
   oper_code            varchar(16) not null default '',
   oper_obj_id          bigint not null default 0,
   oper_obj_type        varchar(32) not null default '',
   oper_obj_property    varchar(32) not null default '',
   relation_obj_id      bigint not null default 0,
   relation_obj_type    varchar(32) not null default '',
   relation_type        varchar(32) not null default '',
   new_value            text,
   old_value            text,
   ext                  varchar(4096) not null default '',
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_org_id                           */
/*==============================================================*/
create index index_ppm_tre_trends_org_id on ppm_tre_trends
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_uuid                             */
/*==============================================================*/
create index index_ppm_tre_trends_uuid on ppm_tre_trends
(
   uuid
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_module2_id_module2               */
/*==============================================================*/
create index index_ppm_tre_trends_module2_id_module2 on ppm_tre_trends
(
   module2_id,
   module2
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_module3_id_module3               */
/*==============================================================*/
create index index_ppm_tre_trends_module3_id_module3 on ppm_tre_trends
(
   module3_id,
   module3
);

/*====================================================================*/
/* Index: index_ppm_tre_trends_oper_relation_obj_id_relation_obj_type */
/*====================================================================*/
create index index_ppm_tre_trends_oper_relation_obj_id_relation_obj_type on ppm_tre_trends
(
   relation_obj_id,
   relation_obj_type
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_obj_id                      */
/*==============================================================*/
create index index_ppm_tre_trends_oper_obj_id on ppm_tre_trends
(
   oper_obj_id,
   oper_obj_type
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_create_time                 */
/*==============================================================*/
create index index_ppm_tre_trends_oper_create_time on ppm_tre_trends
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_tre_trends_oper_creator                     */
/*==============================================================*/
create index index_ppm_tre_trends_oper_creator on ppm_tre_trends
(
   creator
);