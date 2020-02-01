/*==============================================================*/
/* Table: ppm_bas_object_id                                     */
/*==============================================================*/
create table if not exists ppm_bas_object_id
(
   id                   bigint not null auto_increment,
   org_id               bigint not null default 0,
   code                 varchar(64) not null default '',
   max_id               bigint not null default 0,
   step                 int not null default 100,
   creator              bigint not null default 0,
   create_time          datetime not null default CURRENT_TIMESTAMP,
   updator              bigint not null default 0,
   update_time          datetime not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   version              int not null default 1,
   is_delete            tinyint not null default 2,
   primary key (id)
);

/*==============================================================*/
/* Index: index_ppm_bas_object_id_create_time                   */
/*==============================================================*/
create index index_ppm_bas_object_id_create_time on ppm_bas_object_id
(
   create_time
);

/*==============================================================*/
/* Index: index_ppm_bas_object_id_org_id                        */
/*==============================================================*/
create index index_ppm_bas_object_id_org_id on ppm_bas_object_id
(
   org_id
);