# AutoDG oracle dataguard 自动搭建使用手册
-------
#### 功能参数使用说明
1. check 参数 主库环境预检查
    1. 主要检查数据库版本是否大于11201
    2. 主库归档模式 & force_logging 是否开启
    3. 主库密码文件orapwd$SID是否存在
    4. 主库是否使用spfile 

2. prepare 参数 主库环境检查,主备环境初始化，自动搭建oracle dataguard
   1. 主库更新tnsnames.ora 
   2. 下载主库tnsnames.ora，orapwd 密码文件
   3. 本地备库tnsnames.ora orapwd 文件同步
   4. 本地备库adump 等目录配置
   5. 本地备库listener.ora 初始化
   6. 本地备库pfile 初始化
   7. 本地备库启动到 nomount 状态
   8. 本地备库执行rman duplicate 命令    

    
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


若直接在命令行中用 `nohup` 启动程序，可能会因为 SIGHUP 信号而退出，建议把 `nohup` 放到脚本里面且不建议用 kill -9，如：

```shell
#!/bin/bash
nohup ./autodg -config config.toml --mode check > nohup.out &
```
