# oss-script

## 功能描述

用于对比2个bucket中同名但内容不同的文件

## 场景说明

源bucket数据量约13G，2.6w+文件；目标bucket数据量1.2T；执行时间15S左右

通过遍历源bucket文件，并在目标bucket中查询同名文件，如果有同名文件，再对比`content-md5`是否相同，如果不同则记录该文件地址。

如果对比期间因网络原因导致任务终止，则会记录当前对比位置信息，方便之后继续执行脚本。

## 变量说明

名称 | 描述
--- | ---
accessKeyId | RAM账号
accessKeySecret | RAM密码
originBucket | 源bucket名称
originRegion | 源bucket区域
targetBucket | 目标bucket名称
targetRegion | 目标bucket区域
targetEndpoint | 目标bucket域名
marker | 分页标记
max-keys | 返回文件最大数量 0-1000
