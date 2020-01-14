#demo demo table
DROP TABLE IF EXISTS `demo`;
 CREATE TABLE if not exists `demo` (
  `id` INT  NOT NULL auto_increment,	#主键编号
	`name` varchar(20) not null default '' COMMENT '名称',
	`gmt_create` bigint NOT NULL COMMENT '创建时间',
  `gmt_modified` bigint not null default 0 COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;