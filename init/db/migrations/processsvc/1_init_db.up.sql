/*==============================================================*/
/* Table: ppm_prs_process                                       */
/*==============================================================*/
create table if not exists ppm_prs_process
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(64) not null default '',
   is_default           tinyint not null default 0,
   type                 tinyint not null default 1,
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
/* Index: index_ppm_prs_process_org_id                          */
/*==============================================================*/
create index index_ppm_prs_process_org_id on ppm_prs_process
(
   org_id,
   name
);

/*==============================================================*/
/* Index: index_ppm_prs_process_create_time                     */
/*==============================================================*/
create index index_ppm_prs_process_create_time on ppm_prs_process
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_process_status                        */
/*==============================================================*/
create table if not exists ppm_prs_process_process_status
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   process_id           bigint not null default 0,
   process_status_id    bigint not null default 0,
   is_init_status       tinyint not null default 2,
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
/* Index: index_ppm_prs_process_process_status_org_id           */
/*==============================================================*/
create index index_ppm_prs_process_process_status_org_id on ppm_prs_process_process_status
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_process_status_process_id       */
/*==============================================================*/
create index index_ppm_prs_process_process_status_process_id on ppm_prs_process_process_status
(
   process_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_process_status_create_time      */
/*==============================================================*/
create index index_ppm_prs_process_process_status_create_time on ppm_prs_process_process_status
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_status                                */
/*==============================================================*/
create table if not exists ppm_prs_process_status
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   sort                 int not null default 1,
   bg_style             varchar(8) not null default '',
   font_style           varchar(8) not null default '',
   type                 tinyint not null default 1,
   category             tinyint not null default 1,
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
/* Index: index_ppm_prs_process_status_org_id                   */
/*==============================================================*/
create index index_ppm_prs_process_status_org_id on ppm_prs_process_status
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_status_create_time              */
/*==============================================================*/
create index index_ppm_prs_process_status_create_time on ppm_prs_process_status
(
   create_time
);

/*==============================================================*/
/* Table: ppm_prs_process_step                                  */
/*==============================================================*/
create table if not exists ppm_prs_process_step
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   process_id           bigint not null default 0,
   lang_code            varchar(64) not null default '',
   name                 varchar(32) not null default '',
   start_status         bigint not null default 0,
   end_status           bigint not null default 0,
   sort                 int not null default 1,
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
/* Index: index_ppm_prs_process_step_org_id                     */
/*==============================================================*/
create index index_ppm_prs_process_step_org_id on ppm_prs_process_step
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_process_id                 */
/*==============================================================*/
create index index_ppm_prs_process_step_process_id on ppm_prs_process_step
(
   process_id
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_start_status               */
/*==============================================================*/
create index index_ppm_prs_process_step_start_status on ppm_prs_process_step
(
   start_status
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_end_status                 */
/*==============================================================*/
create index index_ppm_prs_process_step_end_status on ppm_prs_process_step
(
   end_status
);

/*==============================================================*/
/* Index: index_ppm_prs_process_step_create_time                */
/*==============================================================*/
create index index_ppm_prs_process_step_create_time on ppm_prs_process_step
(
   create_time
);