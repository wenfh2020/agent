package mysql

import "time"

/*
CREATE DATABASE IF NOT EXISTS lhl_product DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE `device_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `device_mac` varchar(64) NOT NULL COMMENT '设备网卡 mac 地址',
  `device_version` varchar(64) NOT NULL COMMENT '请求流水 id',
  `activation` varchar(128) COMMENT '设备激活码',
  `client_type` varchar(64) NOT NULL COMMENT '用户类型',
  `client_level` varchar(64) NOT NULL COMMENT '用户等级',
  `status` tinyint(1) unsigned DEFAULT '1' COMMENT '数据状态，默认 1 有效，0 无效',
  `active_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '激活时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `mac` (`device_mac`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
*/

type DeviceInfo struct {
	Id            int       `gorm:"id" json:"id"`
	DeviceMac     string    `gorm:"device_mac" json:"device_mac"`
	DeviceVersion string    `gorm:"device_version" json:"device_version"`
	Activation    string    `gorm:"activation" json:"activation"`
	ClientType    string    `gorm:"client_type" json:"client_type"`
	ClientLevel   string    `gorm:"client_level" json:"client_level"`
	Status        int       `gorm:"status" json:"status"`
	ActiveTime    time.Time `gorm:"active_time" json:"active_time"`
	CreateTime    time.Time `gorm:"create_time" json:"create_time"`
	UpdateTime    time.Time `gorm:"update_time" json:"update_time"`
}

func (DeviceInfo) TableName() string {
	return "device_info"
}
