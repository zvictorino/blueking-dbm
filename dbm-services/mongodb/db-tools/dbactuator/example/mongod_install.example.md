### mongod_install
初始化新机器:

```json
./dbactuator_redis  --uid={{uid}} --root_id={{root_id}} --node_id={{node_id}} --version_id={{version_id}} --atom-job-list="mongod_install" --data_dir=/path/to/data  --backup_dir=/path/to/backup  --payload='{{payload_base64}}'
```


原始payload

## shardsvr
```json
{
  "mediapkg":{
    "pkg":"mongodb-linux-x86_64-3.4.20.tar.gz",
    "pkg_md5":"e68d998d75df81b219e99795dec43ffb"
  },
  "ip":"1.1.1.1",
  "port":27001,
  "dbVersion":"3.4.20",
  "instanceType":"mongod",
  "setId":"s1",
  "keyFile": "xxx",
  "auth": true,
  "clusterRole":"shardsvr",
  "dbConfig":{
    "slowOpThresholdMs":200,
    "cacheSizeGB":1,
    "oplogSizeMB":500,
    "destination":"file"
  }
}
```
部署复制集时"clusterRole"字段为空

## configsvr
```json
{
  "mediapkg":{
    "pkg":"mongodb-linux-x86_64-3.4.20.tar.gz",
    "pkg_md5":"e68d998d75df81b219e99795dec43ffb"
  },
  "ip":"1.1.1.1",
  "port":27002,
  "dbVersion":"3.4.20",
  "instanceType":"mongod",
  "setId":"conf",
  "keyFile": "xxx",
  "auth": true,
  "clusterRole":"configsvr",
  "dbConfig":{
    "slowOpThresholdMs":200,
    "cacheSizeGB":1,
    "oplogSizeMB":500,
    "destination":"file"
  }
}
```