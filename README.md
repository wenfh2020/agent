# agent

simple golang http demo: client/server for communication's verification.

---

## 1. protocol

| protocol | method | data  |                 url                 |
| :------: | :----: | :---: | :---------------------------------: |
|   http   |  post  | json  | http://xxxx.com/product/agent/check |

* request

```json
{
    "client": {
        "level": "23424343",
        "type": "fdsfdsfa",
    },
    "device" : {
        "mac": "dhsjfhjasfhjadfhjdasf",
        "version": "12.34.4354"
    },
    "time": "2010-03-13 10:00:11",
    "sign": "fhkdsahfjashfjkdshfjkdafda"
}
```

* response

```json
{
    "errno": 0,
    "errstr": "xxx",
    "device" : {
        "mac": "dhsjfhjasfhjadfhjdasf",
        "version": "12.34.4354"
    },
    "activation": "ewruhfdjdsahfjkhfjsirewure",
    "time": "2010-03-14 10:00:11",
    "sign": "fdhsjfhdasjfhjasdfhjka"
}
```

---

## 2. database

```sql
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
```

---

## 3. run

```shell
cd agent
# client
go run client/client.go
# server
go run main.go
```
