# zonecopy
## 编译运行
### 编译
go build -o main main.go
### 使用说明
./main -help
```bash
Usage of ./main:
  -module string
        导入指定模块配置 
        origin: 源站组 
        domain: 域名管理 
        rule: 规则引擎 
        all: 全部模块
```
### 执行示例
// 按顺序全部导入  
./main -module all 

// 导入域名  
./main -module domain 

## 模块说明
- origin
对应控制台 源站配置-源站组 中源站相关配置
- domain
对应控制台 域名服务-域名管理 中三级域名相关配置
- rule
对应控制台 规则引擎 中所有规则配置

## 导入顺序说明
因为模块存在依赖关系(origin > domain > rule)，如域名服务依赖于源站组服务，规则引擎服务依赖源站服务和域名服务，分模块导入时，务必确保导入顺序。   
如 ./main -module rule 单独导入规则引擎配置时，务必确保domain和origin配置已经导入，否则可能会导致导入失败


