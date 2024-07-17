# oss-script

## 功能描述

用于对比2个bucket中同名但内容不同的文件

## 场景说明

源bucket数据量约13G，2.6w+文件；目标bucket数据量1.2T；执行时间15S左右

通过遍历源bucket文件，并在目标bucket中查询同名文件，如果有同名文件，再对比`content-md5`是否相同，如果不同则记录该文件地址。

如果对比期间因网络原因导致任务终止，则会记录当前对比位置信息，方便之后继续执行脚本。

## 使用说明

1. 设置环境变量：
   ```sh
   export ORIGIN_REGION=<origin-region>
   export ACCESS_KEY_ID=<access-key-id>
   export ACCESS_KEY_SECRET=<access-key-secret>
   export ORIGIN_BUCKET=<origin-bucket>
   export TARGET_REGION=<target-region>
   export TARGET_BUCKET=<target-bucket>
   export TARGET_ENDPOINT=<target-endpoint>
   ```

2. 运行Go程序：
   ```sh
   go run main.go
   ```

## 变量说明

名称 | 描述
--- | ---
ACCESS_KEY_ID | RAM账号
ACCESS_KEY_SECRET | RAM密码
ORIGIN_BUCKET | 源bucket名称
ORIGIN_REGION | 源bucket区域
TARGET_BUCKET | 目标bucket名称
TARGET_REGION | 目标bucket区域
TARGET_ENDPOINT | 目标bucket域名
