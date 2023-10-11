# 通用库使用手册 golang

为了推进 “全员贡献、全员享用”，提高复用程度，提升协同效率，沉淀成熟方案，巩固稳定性，设立通用库。

通用库分为 通用组件 与 通用业务 两部分，通用组件是业务无关的基础组件，包括独立小组件、小功能归类、SDK等，通用业务是对业务中的通用点抽象提炼，比如发奖队列、用户模块等。通用组件见下述内容，通用业务在general目录。

golang、php、java等语言，语法特性各异、流行开源包多样、面向业务场景不同、封装方式不一致，各语言的通用库包可能不严格一致。

### 正式包列表

经过代码评审、测试、线上检验的包。

| package               | 名称       | since   | 说明                     |
|-----------------------|----------|---------|------------------------|
| v1/algorithm          | 算法       | v1.0 | 加解密，编码解码。              |
| v1/array              | 数组       | v1.0 | go数组的补充方法。             |
| v1/config             | 配置       | v1.0 | 基于Viper，按约定方式读取本地配置文件。 |
| v1/rand               | 随机       | v1.0 | 随机方法集合。                |
| v1/sign               | 签名       | v1.0 | 签名方法集合。                |
| v1/ip                 | IP       | v1.0 | 获取、解析。                 |
| v1/encoding/json      | Json     | v1.0 | json。                  |
| v1/gosafe             | 安全运行协程   | v1.0 | 安全运行协程。                |
| v1/db                 | 数据库      | v1.0 | 基于gorm，支持多种数据库。        |
| v1/redis              | Redis    | v1.0 | 基于go-redis v8。         |
| v1/logger             | 日志       | v1.0 | 输出json日志，支持按级别发送通知。    |
| v1/notice             | 通知       | v1.0 | 邮件，飞书。                 |
| v1/datetime           | 日期时间     | v1.0 | 各种时间方法，主要围绕时区封装。       |
| v1/signal             | 信号监听     | v1.0 | 信号监听                   |
| v1/response           | 响应       | v1.0 | 响应方法集合，比如json          |
| v1/bytedance/feishu   | 字节跳动飞书   | v1.0 | 自定义群机器人。               |

**注意：正式包通过业务测试与线上校验，不代表包整体没有问题，因为线上业务可能只用到了包的一部分。新业务在使用包时，依然需要严格测试。**

### 待验收包

实验性质，仅用于测试，成熟后移至正式包列表。

| package               | 名称       | since   | 说明            |
|-----------------------|----------|---------|---------------|
|||||

### 贡献须知

要求全面遵守开发规范，包括编码格式、拆分粒度、不限特定场景、避免过度封装等等。

要求维护使用手册，便于查找使用。由简要的列表说明，与详细说明组成。

业务未使用到的多余功能，不能顺带放进通用库包，以免误导其他使用者。

要求以master分支整体发行，标签为三级版本号。package包版本体现在路径上，不兼容升级时必须创建新版本，业务类的包收录在厂商名下。

# 详细说明

## v1/config

基于Viper v1.7.1封装，单例模式，仅支持一个配置文件，yaml格式。区分环境，优先读取自定义环境设置，为空时读取-m参数。

默认环境为release。建议将默认环境设置为生产环境，以确保在任何情况下不出错，同时阻止开发测试环境连接生产环境的资源。

### 定义

**Start(envKey, path)**

使用前需要启动

envKey：自定义环境（os.Getenv），比如 MY_ENV，建议提前部署变量到所有环境。  
path：配置目录，比如 ./config。

### 实例

```go
import "github.com/lynnclub/go/v1/config"

// 修改默认环境值
config.Env="production"

// 启动
// 当MY_ENV=dev，读取./config/dev.yaml
config.Start("MY_ENV", "./config")

// 使用
config.Env //获取环境值
config.Viper.GetString("name") //读取配置
```

更多使用方法，请查阅 [《Viper文档》](https://github.com/spf13/viper/tree/v1.7.1#getting-values-from-viper)

## v1/mysql

基于gorm，依赖v1/config包，仅封装了配置与实例池。

### 定义

**Use(configName string) \*gorm.DB**

选择使用的数据库，根据配置组名称返回对应数据库实例，不存在时新建连接，留空使用default配置。

### 实例

```yaml
# 配置 yaml
mysql:
  default:
    dsn: "root:123456@tcp(0.0.0.0:3306)/test?charset=utf8mb4"
    max_open_conn: 100 #打开数据库连接的最大数量（可选）
    max_idle_conn: 10 #空闲连接池的最大数量（可选）
    log_level: 4 #Silent 1、Error 2、Warn 3、Info 4（可选）
    slow_threshold: 1 #慢SQL阈值，单位秒（可选）
```

```go
import "github.com/lynnclub/go/v1/mysql"

// 留空使用default配置
defaultDB := mysql.Use("")
// 初始化后，还可以直接通过Default访问
mysql.Default

// Use 使用
testDB := mysql.Use("test")
```

更多使用方法，请查阅 [《gorm文档》](https://gorm.io/zh_CN/docs/index.html)

## v1/db

基于gorm，仅封装了配置与实例池。

### 定义

**Add(name string, option Option)**

添加数据库配置，使用db.Option。

**AddMap(name string, setting map[string]interface{})**

便捷方法，使用map。

**AddMapBatch(batch map[string]interface{})**

便捷方法，使用map批量添加。

**Use(configName string) \*gorm.DB**

选择使用的数据库，根据配置组名称返回对应数据库实例，不存在时新建连接，留空使用default配置。

### 实例

```yaml
# 配置 yaml
db:
  default:
    dsn: "root:123456@tcp(0.0.0.0:3306)/test?charset=utf8mb4"
    max_open_conn: 100 #打开连接的最大数量，默认100
    max_idle_conn: 10 #空闲连接的最大数量，默认10
    log_level: 4 #Silent 1、Error 2、Warn 3、Info 4，默认3
    slow_threshold: 1 #慢SQL阈值，单位秒，默认1
  postgres:
    driver: "postgres" #驱动，mysql、postgres、sqlite、sqlserver、clickhouse，默认mysql
    dsn: "host=localhost user=postgres password= dbname=test port=5432"
    max_open_conn: 100 #打开连接的最大数量，默认100
    max_idle_conn: 10 #空闲连接的最大数量，默认10
    log_level: 4 #Silent 1、Error 2、Warn 3、Info 4，默认3
    slow_threshold: 1 #慢SQL阈值，单位秒，默认1
```

```go
import "github.com/lynnclub/go/v1/db"

// 添加配置
db.AddMapBatch(config.Viper.GetStringMap("db"))

// 留空使用default配置
defaultDB := db.Use("")
// 初始化后，还可以直接通过Default访问
db.Default

// Use 使用
testDB := db.Use("test")
```

更多使用方法，请查阅 [《gorm文档》](https://gorm.io/zh_CN/docs/index.html)

## v1/redis

基于go-redis v8，仅封装了配置与实例池。

### 定义

**Add(name string, option Option)**

添加数据库配置，使用redis.Option。

**AddMap(name string, setting map[string]interface{})**

便捷方法，使用map。

**AddMapBatch(batch map[string]interface{})**

便捷方法，使用map批量添加。

**Use(configName string) \*redis.Client**

选择使用的数据库，根据配置组名称返回对应数据库实例，不存在时新建连接，留空使用default配置。

**Lock(name string, expire time.Duration) bool**

加锁，expire过期时间（单位秒）。

**Unlock(name string)**

解锁。可以不解锁，等过期。

**MaxMin对象**

最大值最小值，记录、获取极值，超过极值才会覆盖。

### 实例

```yaml
# 配置 yaml
redis:
  default:
    address: "0.0.0.0:6379"
    password: "" #密码，默认空
    db: 0
    pool_size: 100 #连接池最大数量，默认100
```

```go
import "github.com/lynnclub/go/v1/redis"

// 添加配置
redis.AddMapBatch(config.Viper.GetStringMap("redis"))

// 留空使用default配置
defaultRedis := redis.Use("")
// 初始化后，还可以直接通过Default访问
redis.Default
// context.Background()
redis.Ctx

// Use 使用
testRedis := redis.Use("test")

// 锁
redis.Default.Lock("login", 1)
redis.Default.UnLock("login")

// 最大值最小值
maxMin := redis.MaxMin{CacheKey: "cache", Name: "test"}
maxId := maxMin.Get()
newMaxId := 10
if err = maxMin.SetMax(newMaxId); err == nil {
    // 设置成功
}
```

更多使用方法，请查阅 [《go-redis文档》](https://github.com/go-redis/redis/)

## v1/logger

基于官方log包，支持函数或对象两种封装，支持按级别发送通知。日志格式遵守Json规范。

### 定义

**SetLevel(level int)**

设置等级，小于该等级的日志将过滤。

**SetPrefix(prefix string)**

设置前缀，拼接在title字段之前。

**Debug(v ...interface{})**  
**Info(v ...interface{})**  
**Warn(v ...interface{})**  
**Error(v ...interface{})**  
**Panic(v ...interface{})**  
**Fatal(v ...interface{})**

调试、信息、警告、错误、恐慌、致命错误

### 实例

```yaml
# yaml配置-飞书通知
feishu:
  group:
    alert:
      env: "dev"
      webhook: "https://open.feishu.cn/xxx"
      sign_key: ""
      user_id: ""
```

```go
import "github.com/lynnclub/go/v1/logger"

// 自定义启动
// 默认为 New("", 0, "asia/shanghai", nil)
logger.Logger = logger.New(
  "",
  logger.DEBUG,
  "asia/shanghai",
  config.Viper.GetStringMapString("feishu.group.alert"),
)

// 记日志
logger.Info("[通用队列发奖]启动")
logger.Error("错误", err)
```

## v1/rand

### 定义

**Range(min, max int) int**

区间随机，使用nanoseconds做种子，全开区间[min, max]

### 实例

```go
import "github.com/lynnclub/go/v1/rand"

// Range 区间随机
num := rand.Range(6, 8)
```

## v1/sign

### 定义

**MD5(params map[string]string, secret string) string**

常规md5 get拼接

**FeiShu(secret string, timestamp int64) (string, error)**

飞书签名

### 实例

```go
import "github.com/lynnclub/go/v1/sign"

// MD5 常规md5
params := map[string]string{"test": "123"}
result := sign.MD5(params, "123")

// FeiShu 飞书
result, err := sign.FeiShu("123", 1667820457)
```

## v1/ip

### 定义

**Local(ipv4 bool) []string**

获取本地IP

**GetClientIP(c \*gin.Context) (string, error)**

获取Header client-ip 的内容。

### 实例

```go
import "github.com/lynnclub/go/v1/ip"

// Local 本地IP
ips := ip.Local(true)

// Resolve 解析IP
response, errs := ip.Resolve("14.1.44.228")
```

## v1/datetime

官方内置的time包实现简洁，无需封装，建议简单场景下直接使用，比如获取当前时间戳。

datetime包主要围绕时区封装，将字符串时间、时间戳等形式，统一解析为带时区的官包time.Time对象，然后以此为基础，提供更多快捷方法。提供函数与对象两种封装方式，效果相同。如需支持多个时区，建议使用函数形式；单一时区，建议使用对象形式。

**注意：如果timezone值无效，将自动使用系统时区**。

### 定义

**const LayoutDate**

常量，日期格式，2006-01-02

**const LayoutDateTime**

常量，日期时间格式，2006-01-02 15:04:05

**ParseDateTime(datetime string, timezone string) time.Time**

解析日期时间，基础函数。

**ToUnix(datetime string, timezone string) int64**

转成时间戳，基于ParseDateTime封装

**ParseTimestamp(timestamp int64, timezone string) time.Time**

解析时间戳，基础函数。

**ToAny(timestamp int64, timezone string, layout string) string**

转成任意格式，基于ParseTimestamp封装

**ToDate(timestamp int64, timezone string) string**

转成日期，基于ToAny封装

**ToDateTime(timestamp int64, timezone string) string**

转成日期时间，基于ToAny封装

**ToISOWeek(timestamp int64, timezone string) string**

转成年周 例如2020_5，基于ParseTimestamp封装

**ToISOWeekByDate(datetime string, timezone string) string**

转成年周 例如2020_5，基于ParseDateTime封装

**Unix(timezone string) int64**

当前时间戳，基于ToUnix封装

**Date(timezone string) string**

当前日期，基于ToAny封装

**DateTime(timezone string) string**

当前日期时间，基于ToAny封装

**ISOWeek(timezone string) string**

当前年周 例如2020_5，基于ToISOWeek封装

**Any(timezone, layout string) string**

当前任意格式

**CheckTime(timestamp int64, start string, end string, timezone string) int**

检查时间 0未开始、1正常、2已结束，半闭合区间[start, end)

**CheckTimeNow(start string, end string, timezone string) int**

检查当前时间 0未开始、1正常、2已结束，半闭合区间[start, end)，基于CheckTime封装

### 实例

```go
import "time"
import "github.com/lynnclub/go/v1/datetime"


// 默认时区，建议直接使用time包


// 当前时间戳，无时区之分
time.Now().Unix()

// 时间戳转日期时间
time.Unix(timestamp, 0).Format(datetime.LayoutDateTime)

// 日期时间转时间戳
time.Parse(datetime.LayoutDateTime, datetime).Unix()


// 多个时区，建议使用函数形式


// ParseDateTime 解析日期时间
goTime := datetime.ParseDateTime("2022", timezone)
goTime2 := datetime.ParseDateTime("2022-11-08 15:30:02", timezone)

// ToUnix 转成时间戳
timestamp := datetime.ToUnix("2022-11-08 15:30:02", timezone)

// ParseTimestamp 解析时间戳 秒级
goTime := datetime.ParseTimestamp(1667895429, timezone)

// ToAny 转成任意格式
datetimeStr := datetime.ToAny(1667895429, timezone, datetime.LayoutDateTime)

// ToDate 转成日期
dateStr := datetime.ToDate(1667895429, timezone)

// ToDateTime 转成日期时间
datetimeStr := datetime.ToDateTime(1667895429, timezone)

// ToISOWeek 转成年周
week := datetime.ToISOWeek(1667895429, timezone)

// ToISOWeekByDate 转成年周
week := datetime.ToISOWeekByDate("2022-11-08", timezone)

// Date 当前日期
dateStr := datetime.Date(timezone)

// DateTime 当前日期时间
datetimeStr := datetime.DateTime(timezone)

// ISOWeek 当前年周
week := datetime.ISOWeek(timezone)

// CheckTime 检查时间
datetime.CheckTime(1667895429, "2022-11-01", "2023-01-01 00:00:00", timezone)


// 单一时区，建议使用对象形式


// 设置时区
datetime.Single.SetTimeZone(timezone)

// ToUnix 转成时间戳
timestamp := datetime.Single.ToUnix("2022-11-08 15:30:02")

// ToAny 转成任意格式
datetimeStr := datetime.Single.ToAny(1667895429, datetime.LayoutDateTime)
```

## v1/notice

### 定义

**NewFeiShuGroup(webhook, signKey, env string) \*FeiShuGroup**

飞书群实例化

**Send(title string, content map[string]interface{}, userId string)**

飞书群发送

### 实例

```go
import "github.com/lynnclub/go/v1/notice"

// NewFeiShuGroup 飞书群实例化
group := notice.NewFeiShuGroup(url, signKey, "dev")

// Send 飞书群发送
content := map[string]interface{}{
  "tag":  "text",
  "text": json.Encode(full),
}
group.Send("title", content, userId)
```

## v1/encoding/json

### 定义

**Encode(v interface{}) string**

Json编码

**Decode(str string, v interface{}) error**

Json解码

### 实例

```go
import "github.com/lynnclub/go/v1/encoding/json"

// Encode Json编码
result := json.Encode(map[string]interface{}{"adc": "123", "No": 1234})

// Decode Json解码
err := json.Decode(jsonStr, &data)
```

## v1/signal

### 定义

**Listen(signals ...os.Signal)**

监听系统信号，可以用于脚本平滑退出。调用时会启动一个监听信号的协程，将其阻塞直至收到信号。

相关知识：

1. 如果主程不循环也不阻塞，即使协程循环或阻塞，整个进程还是会退出。也就是说，常驻程序的主程应循环或阻塞，并且在平滑退出时，主程需要等协程先退出；
2. channel作为协程的形参时，优先级很高，因为从go1.14开始，协程调度器是 基于信号的真抢占式调度；
3. channel可以被多处接收，包括主程和协程；
4. select需要关闭channel才会退出。
5. 建议结合sync/atomic原子计数器，实现高并发安全的协程平滑退出。

### 实例

```go
import "github.com/lynnclub/go/v1/signal"
import "sync/atomic"

// Listen 监听
signal.Listen()

// 模拟业务协程
var counter int64 = 0
go func() {
  atomic.AddInt64(&counter, 1)

    for loop := 0; loop < 100; loop++ {
        // 收到停机信号，主动退出业务
        if signal.Now != nil {
      atomic.AddInt64(&counter, -1)
            fmt.Println("business stop signal:", Now)
            break
        }

        // do something...
        
        time.Sleep(100 * time.Millisecond)
    }

    if signal.Now == nil {
        panic("no signal")
    }
}()

// 发送信号
err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
if err != nil {
    panic("signal send failed")
}

// 主程形式一：循环
for {
  // 主程需要等待协程停止
  if Now != nil && counter <= 0 {
    fmt.Println("main stop signal:", Now)
    break
  }
  
  // do something...

  time.Sleep(100 * time.Millisecond)
}

// 主程形式二：阻塞
//select {
//case <-ChannelOS:
//  // 主程需要等待协程停止
//  for {
//    if counter <= 0 {
//      fmt.Println("main stop signal:", Now)
//      break
//    }
//
//    time.Sleep(100 * time.Millisecond)
//  }
//}
```

## v1/array

### 定义

**In(array []string, find string) bool**

是否存在

**NotIn(array []string, find string) bool**

是否不存在

### 实例

```go
import "github.com/lynnclub/go/v1/array"

testArray := []string{
    "adc",
    "mon",
    "测试",
}

array.In(testArray, "测试abc")
array.NotIn(testArray, "测试")
```

## v1/response

响应方法集合，目前支持json。响应码建议遵守http code规范，使用官方http包的编码定义。

```go
import "http"
import "github.com/lynnclub/go/v1/response"

response.Json(http.StatusUnauthorized, "请登录")
```

## v1/bytedance/feishu

### 定义

**NewGroupRobot(webhook, signKey string) \*GroupRobot**

群机器人实例化

**Send(request interface{}) (entity.GroupRobotResponse, error)**

发送消息

### 实例

```go
import (
  "github.com/lynnclub/go/v1/bytedance/feishu"
  "github.com/lynnclub/go/v1/bytedance/feishu/entity"
)

// NewGroupRobot 群机器人实例化
robot := feishu.NewGroupRobot(webhook, signKey)

// Send 发送

var data entity.PostData
data.Title = "title"

// 艾特用户
if userId == "" {
  data.Content = [][]map[string]interface{}{{content}}
} else {
  data.Content = [][]map[string]interface{}{{content, map[string]interface{}{
    "tag":     "at",
    "user_id": userId,
  }}}
}

var richText entity.MsgTypePost
richText.Post = map[string]entity.PostData{"zh_cn": data}

_, _ = robot.Send(richText)
```
