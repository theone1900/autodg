# AutoDG oracle dataguard 自动搭建使用手册
-------
#### 使用说明
1. check 参数 主库环境检查
    1. 主要检查 数据库版本 是否大于11201，归档模式&force_loging 是否开启，密码文件orapwd$SID是否存在
    2. 不兼容性对象
    8. 注意事项
       1. 

2. 表
   1. 表
      1. 若
   2. 注意事项
      

    
#### 使用事项

```
1、下载 oracle client，参考官网下载地址 https://www.oracle.com/database/technologies/instant-client/linux-x86-64-downloads.html

2、上传 oracle client 至程序运行服务器，并解压到指定目录，比如：/data1/soft/client/instantclient_19_8

3、配置程序运行环境变量 LD_LIBRARY_PATH
export LD_LIBRARY_PATH=/data1/soft/client/instantclient_19_8
echo $LD_LIBRARY_PATH

4、配置 autodg 参数文件，config.toml 相关参数配置说明见 conf/config.toml

5、主库环境检查
$ ./autodg --config config.toml --mode check

6、自动配置dataguard 环境编辑，数据库同步备份等
$ ./autodg --config config.toml --mode prepare
```
#### prepare 模式
```sql

/* 数据库开启归档以及补充日志 */
-- 开启归档【必须选项】
alter database archivelog;
-- 强制日志【必须选项】
ALTER DATABASE force log ;




/* 查看附加日志 */
-- 数据库级别附加日志查看

```

若直接在命令行中用 `nohup` 启动程序，可能会因为 SIGHUP 信号而退出，建议把 `nohup` 放到脚本里面且不建议用 kill -9，如：

```shell
#!/bin/bash
nohup ./autodg -config config.toml --mode check > nohup.out &
```
