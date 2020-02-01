# 介绍
北极星企业协作平台是一款更更专业，更便捷的团队协作工具：
- 适合专业产品研发团队，同样适用于其它行业项目团队。
- 多维度展现当前项目任务数据，实时提醒项目异常情况，比如延期等...
- 通过动态、任务数据折线，实时管控时间节点。
- 通过项目公告等，快速布置，决策团体任务，快速通知至每一名成员。

本项目为北极星服务端开源版，官网地址 [https://www.bjx.cloud](https://www.bjx.cloud)
# 使用
首先将代码clone到本地：
```
git clone https://github.com/galaxy-book/polaris-backend
```
修改``application.common.local.yaml``必要配置：
```
# Mysql
Mysql:
  # Host
  Host: 127.0.0.1
  # Port
  Port: 3306
  # User
  Usr: root
  # Pwd
  Pwd: root
  # 数据库
  Database: polaris_project_manage
# 邮件配置
Mail:
  # 别名：北极星
  Alias: 北极星
  # User
  Usr: xxx
  # Pwd
  Pwd: xxx
  # Host
  Host: xxx
  # Port
  Port: xxx
```
运行``bin``目录下的脚本，启动所有服务：
```
./bin/start.bat
```
或者
```
./bin/start.sh
```
通过访问[http://127.0.0.1:12000/](http://127.0.0.1:12000/)确认是否启动成功 !
## 架构介绍
```
├─app                   //app层对外暴漏接口，使用Graphql做接口路由，使用app/gen.sh生成api层
├─bin                   //启动脚本
├─common                //业务层公共模块
│  ├─core               
│  ├─extra              //三方扩展
│  ├─model              
│  │  ├─bo          
│  │  └─vo
│  └─test
├─config                //配置文件
├─facade                //外观模式对其他服务的接口做封装，提供外部调用，代码可运行bin/build.sh生成
├─init                  //数据库版本初始化
│  └─db
├─schedule              //定时服务
├─scripts               //数据库执行脚本
└─service               //各服务存放在该文件夹下
    ├─basic             //基础服务集合
    │  ├─appsvc     
    │  ├─commonsvc
    │  ├─idsvc
    │  └─msgsvc
    └─platform          //业务服务集合
        ├─orgsvc
        ├─processsvc
        ├─projectsvc
        ├─resourcesvc
        ├─rolesvc
        ├─trendssvc
        └─websitesvc
```