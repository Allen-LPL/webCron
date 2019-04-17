/*
Navicat MySQL Data Transfer

Source Server         : 192.168.11.60
Source Server Version : 50619
Source Host           : 192.168.11.60:3306
Source Database       : webcron

Target Server Type    : MYSQL
Target Server Version : 50619
File Encoding         : 65001

Date: 2019-04-17 17:00:14
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for t_resource
-- ----------------------------
DROP TABLE IF EXISTS `t_resource`;
CREATE TABLE `t_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(30) DEFAULT NULL,
  `url` char(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_role_resource
-- ----------------------------
DROP TABLE IF EXISTS `t_role_resource`;
CREATE TABLE `t_role_resource` (
  `role_id` int(11) NOT NULL DEFAULT '0',
  `resource_id` char(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`role_id`,`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_role_task_group
-- ----------------------------
DROP TABLE IF EXISTS `t_role_task_group`;
CREATE TABLE `t_role_task_group` (
  `role_id` int(5) NOT NULL AUTO_INCREMENT,
  `task_group_id` int(5) DEFAULT '0',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_roles
-- ----------------------------
DROP TABLE IF EXISTS `t_roles`;
CREATE TABLE `t_roles` (
  `role_name` varchar(20) NOT NULL,
  `user_id` varchar(10) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT NULL,
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_task
-- ----------------------------
DROP TABLE IF EXISTS `t_task`;
CREATE TABLE `t_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组ID',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT '时间表达式',
  `concurrent` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否只允许一个实例',
  `command` text NOT NULL COMMENT '命令详情',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0停用 1启用',
  `notify` tinyint(4) NOT NULL DEFAULT '0' COMMENT '通知设置',
  `notify_email` text NOT NULL COMMENT '通知人列表',
  `timeout` smallint(6) NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int(11) NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `die_time` int(10) NOT NULL DEFAULT '0' COMMENT '中止时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_task_group
-- ----------------------------
DROP TABLE IF EXISTS `t_task_group`;
CREATE TABLE `t_task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_name` varchar(50) NOT NULL DEFAULT '' COMMENT '组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_task_log
-- ----------------------------
DROP TABLE IF EXISTS `t_task_log`;
CREATE TABLE `t_task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=722 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_task_log_copy1
-- ----------------------------
DROP TABLE IF EXISTS `t_task_log_copy1`;
CREATE TABLE `t_task_log_copy1` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_user_group
-- ----------------------------
DROP TABLE IF EXISTS `t_user_group`;
CREATE TABLE `t_user_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `grunp_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户组名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_user_role
-- ----------------------------
DROP TABLE IF EXISTS `t_user_role`;
CREATE TABLE `t_user_role` (
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records
-- ----------------------------
--初始化用户
INSERT INTO `t_user` (`id`, `user_name`, `email`, `password`, `salt`, `last_login`, `last_ip`, `status`)
VALUES (1,'admin','admin@example.com','7fef6171469e80d32c0559f88b377245','',0,'',0);

INSERT INTO `t_resource` VALUES ('1', '首页', 'main/index');
INSERT INTO `t_resource` VALUES ('2', '登录', 'main/login');
INSERT INTO `t_resource` VALUES ('3', '退出', 'main/logout');
INSERT INTO `t_resource` VALUES ('4', '帮助', 'help/index');
INSERT INTO `t_resource` VALUES ('5', '任务列表管理', 'task/list');
INSERT INTO `t_resource` VALUES ('6', '任务分组列表管理', 'group/list');
INSERT INTO `t_resource` VALUES ('7', '角色列表管理', 'role/list');
INSERT INTO `t_resource` VALUES ('8', '角色写入管理', 'role/add');
INSERT INTO `t_resource` VALUES ('9', '角色修改管理', 'role/edit');
INSERT INTO `t_resource` VALUES ('10', '角色规则列表管理', 'resource/list');
INSERT INTO `t_resource` VALUES ('11', '角色规则写入管理', 'resource/add');
INSERT INTO `t_resource` VALUES ('12', '角色规则修改管理', 'resource/edit');
INSERT INTO `t_resource` VALUES ('13', '用户列表管理', 'user/list');
INSERT INTO `t_resource` VALUES ('14', '用户写入管理', 'user/add');
INSERT INTO `t_resource` VALUES ('15', '用户修改管理', 'user/edit');
INSERT INTO `t_resource` VALUES ('16', '任务激活', 'task/start');
INSERT INTO `t_resource` VALUES ('17', '任务暂停', 'task/pause');
INSERT INTO `t_resource` VALUES ('18', '任务编辑', 'task/edit');
INSERT INTO `t_resource` VALUES ('19', '任务执行', 'task/run');
INSERT INTO `t_resource` VALUES ('20', '任务日志', 'task/logs');
INSERT INTO `t_resource` VALUES ('21', '任务分组添加', 'group/add');
INSERT INTO `t_resource` VALUES ('22', '任务分组修改', 'group/edit');
INSERT INTO `t_resource` VALUES ('23', '任务日志详情', 'task/viewlog');

INSERT INTO `t_roles` VALUES ('管理员', '1', 1, '超级管理员权限, 不做权限管理', 0);
INSERT INTO `t_roles` VALUES ('棋牌', '2', 2, '棋牌任务', 0);
INSERT INTO `t_role_resource` VALUES ('1', '1001');
INSERT INTO `t_role_resource` VALUES ('1', '2001');
INSERT INTO `t_role_resource` VALUES ('1', '3001');
INSERT INTO `t_role_resource` VALUES ('1', '4001');
INSERT INTO `t_role_resource` VALUES ('2', '1001');
INSERT INTO `t_role_resource` VALUES ('2', '3001');

INSERT INTO `t_user_role` VALUES ('1', '1');
INSERT INTO `t_user_role` VALUES ('2', '2');

