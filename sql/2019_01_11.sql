CREATE TABLE `consumer_backup` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `consumer_id` varchar(64) NOT NULL COMMENT '消费协程唯一id',
  `topic` varchar(64) NOT NULL COMMENT 'topic名称',
  `consumer_group` varchar(64) NOT NULL COMMENT '消费者组名称',
  `callback_url` varchar(255) NOT NULL COMMENT '回调地址',
  `running` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否正在运行 true : 运行中 false: 已宕机',
  `ip` varchar(64) DEFAULT NULL,
  `created_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '消费者创建时间',
  `restarted_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '重启时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8 COMMENT='消费者协程备份表';

CREATE TABLE `topic_callback_template` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `topic` varchar(64) NOT NULL COMMENT 'kafka topic的名字',
  `callback_url` varchar(255) NOT NULL COMMENT 'kafka消费数据后转发url',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '模板创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8 COMMENT='消费者模板表';