drop DATABASE if EXISTS redapricot;
create database redapricot
  DEFAULT CHARACTER SET utf8
  DEFAULT COLLATE utf8_general_ci;
USE redapricot;
SET NAMES utf8;

DROP TABLE if exists data_log;
CREATE TABLE `data_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tag` varchar(100) NOT NULL, /* 服务端初始入口生成该标记，例如 xxxxxxxx_muip_1.1*/
  `maintag` varchar(50) not null default '', /* 上述标记的前置值*/
  `client_tag` VARCHAR(100) NOT NULL DEFAULT '', /*用于按压测业务标记*/
  `level` VARCHAR(10) NOT NULL DEFAULT '', /* 调用层次*/
  `funcname` varchar(50) not null default '', /* 调用功能函数*/
  `appname` varchar(50) not null DEFAULT '', /* 应用名称*/
  `exectime` VARCHAR(14) NOT NULL  DEFAULT '', /* 向该收集程序发送时的时间点*/
  `usetime` int(10) UNSIGNED NOT NULL  DEFAULT '0',/* 埋点位置调用耗时*/
  `content` BLOB,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*900150983cd24fb0d6963f7d28e17f72*/
/*insert into `user`(`name`, `passwd`, `role`, `created`) values('root', md5('abc'), 2, now());*/
/* mysql -h 127.0.0.1 -u root -p < log_schema.sql */
