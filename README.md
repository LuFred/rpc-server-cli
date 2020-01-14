根据sql表结构脚本 自动创建grpc服务

# 准备  

准备一份sql脚本，命名为table.sql放在根目录下  
每一张表结构必须包含id,gmt_create,gmt_modified 3个字段  
例如：  
```
#demo demo table 
DROP TABLE IF EXISTS `demo`;
 CREATE TABLE if not exists `demo` (
  `id` INT  NOT NULL auto_increment,	#主键编号	
	`name` varchar(20) not null default '' COMMENT '名称',	
	`gmt_create` bigint NOT NULL COMMENT '创建时间',
  `gmt_modified` bigint not null default 0 COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```  

# 参数

-listen 服务监听地址  
-port 服务监听端口号  
-service 服务名称  

# 使用  
## 环境  
  已安装protobuf，且protoc命令可识别(最后一步需要生成proto代码)
## 步骤  
* 第一步:go run main.go -service demo -listen 0.0.0.0 -port 4000  
* 第二步:将util文件夹拷贝到生成目录下
* 第三步:执行命令 go mod init "github.com/lufred/{{serviceName}}_service"
