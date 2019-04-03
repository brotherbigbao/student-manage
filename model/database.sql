CREATE TABLE `student` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `no` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '学号',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '学生姓名',
  `c_score` float NOT NULL DEFAULT '0' COMMENT 'C语言成绩',
  `math_score` float NOT NULL DEFAULT '0' COMMENT '数学成绩',
  `english_score` float NOT NULL DEFAULT '0' COMMENT '英语成绩',
  `total_score` float NOT NULL DEFAULT '0' COMMENT '总成绩',
  `average_score` float NOT NULL DEFAULT '0' COMMENT '平均成绩',
  `ranking` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '排名',
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_no` (`no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;