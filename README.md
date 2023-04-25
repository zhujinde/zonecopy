# zonecopy

站点配置拷贝，支持拷贝配置项如下：

- 源站配置
- 域名管理
- 站点加速
- 规则引擎

## 注意事项

- 拷贝方式为在线拷贝，故需保证模板站点配置在EdgeOne控制台已正确配置；
- 由于站点(二级域名)导入无法自动完成，需先手动添加站点，保证站点已生效；
- 当前仅限同一账号下不同站点间的配置拷贝；
- 配置拷贝时，目标站点如已存在相关配置时默认跳过，不会重复导入/覆盖。

## 配置文件准备

请参照配置模版 cmd/cp.yaml.template 填写配置信息
```
# cp.yaml
log_path: ./cp.log  # 运行日志保存路径
account:
  secret_id: xxx    # 账号密钥信息
  secret_key: xxx
  end_point: teo.tencentcloudapi.com  # 固定配置
  region: ap-guangzhou # 固定配置
template_zone: zjd.asia   # 模板站点名
template_zone_id: zone-2dqo3q94x9ks  # 模板站点Id
target_zone: example.com  # 目标站点名
target_zone_id: zone-2ginev8u1owi # 目标站点Id
```

## 编译运行

### 编译

go build -o zcp main.go

### 使用说明

因为模块存在依赖关系(origin > domain > rule)，如域名服务依赖于源站组服务，规则引擎服务依赖源站服务和域名服务，分模块导入时，务必确保导入顺序。

如 ./zcp -module rule 单独导入规则引擎配置时，务必确保依赖于domain和origin的相关配置已经导入，否则可能会导致导入失败。

### 示例

1. 查看使用说明

```bash
./zcp -help
  -config string
        配置文件路径 (default "./config/cp.yaml")
  -module string
        导入指定模块配置 
        origin: 源站组 
        domain: 域名管理 
        zonesetting: 站点加速配置 
        rule: 规则引擎 
        all: 全部模块
```

2. 按顺序全部拷贝

```bash
./zcp -config ./cp.yaml -module all  
```

3. 拷贝域名配置

```bash
./zcp -config ./cp.yaml -module domain
```

## 模块说明

- origin 对应控制台 源站配置-源站组 中源站相关配置
- domain 对应控制台 域名服务-域名管理 中三级域名相关配置
- zonesetting 对应控制台 站点加速 相关配置
- rule 对应控制台 规则引擎 中所有规则配置

