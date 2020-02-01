/*==============================================================*/
/* Table: ppm_res_resource                                      */
/*==============================================================*/
create table if not exists ppm_res_resource
(
   id                   bigint not null,
   org_id               bigint not null default 0,
   project_id           bigint not null default 0,
   type                 tinyint(2) not null default 1,
   bucket               varchar(32) not null default '',
   path                 varchar(512) not null default '1',
   name                 varchar(256) not null default '',
   suffix               varchar(32) not null default 'txt',
   md5                  varchar(32) not null default '',
   size                 bigint not null default 0,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_res_resource_md5                            */
/*==============================================================*/
create index index_ppm_res_resource_md5 on ppm_res_resource
(
   md5
);

/*==============================================================*/
/* Index: index_ppm_res_resource_org_id                         */
/*==============================================================*/
create index index_ppm_res_resource_org_id on ppm_res_resource
(
   org_id
);

/*==============================================================*/
/* Index: index_ppm_res_resource_project_id                     */
/*==============================================================*/
create index index_ppm_res_resource_project_id on ppm_res_resource
(
   project_id
);

/*==============================================================*/
/* Index: index_ppm_res_resource_create_time                    */
/*==============================================================*/
create index index_ppm_res_resource_create_time on ppm_res_resource
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_res_resource_creator                        */
/*==============================================================*/
create index index_ppm_res_resource_creator on ppm_res_resource
(
   creator
);