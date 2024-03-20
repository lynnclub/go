# 通用库使用手册 golang

为了推进 “全员贡献、全员享用”，提高复用程度，提升协同效率，沉淀成熟方案，巩固稳定性，设立通用库。

通用库分为 通用组件 与 通用业务 两部分，通用组件是业务无关的基础组件，包括独立小组件、小功能归类、SDK 等，通用业务是对业务中的通用点抽象提炼，比如发奖队列、用户模块等。通用组件见下述内容，通用业务在 general 目录。

golang、php、java 等语言，语法特性各异、流行开源包多样、面向业务场景不同、封装方式不一致，各语言的通用库包可能不严格一致。

### 正式包列表

经过代码评审、测试、线上检验的包。

| package             | 名称         | since | 说明                                     |
| ------------------- | ------------ | ----- | ---------------------------------------- |
| v1/algorithm        | 算法         | v1.0  | 加解密，编码解码。                       |
| v1/array            | 数组         | v1.0  | go 数组的补充方法。                      |
| v1/config           | 配置         | v1.0  | 基于 Viper，按约定方式读取本地配置文件。 |
| v1/rand             | 随机         | v1.0  | 随机方法集合。                           |
| v1/sign             | 签名         | v1.0  | 签名方法集合。                           |
| v1/ip               | IP           | v1.0  | 获取、解析。                             |
| v1/encoding/json    | Json         | v1.0  | json。                                   |
| v1/safe             | 安全运行协程 | v1.0  | 错误捕获、重试、协程重试。               |
| v1/db               | 数据库       | v1.0  | 基于 gorm，支持多种数据库。              |
| v1/mongo            | Mongodb      | v1.0  | 基于官方 mongo 包。                      |
| v1/redis            | Redis        | v1.0  | 基于 go-redis v8。                       |
| v1/logger           | 日志         | v1.0  | 输出 json 日志，支持按级别发送通知。     |
| v1/notice           | 通知         | v1.0  | 邮件，飞书。                             |
| v1/datetime         | 日期时间     | v1.0  | 各种时间方法，主要围绕时区封装。         |
| v1/signal           | 信号监听     | v1.0  | 信号监听                                 |
| v1/response         | 响应         | v1.0  | 响应方法集合，比如 json                  |
| v1/bytedance/feishu | 字节跳动飞书 | v1.0  | 自定义群机器人。                         |

**注意：正式包通过业务测试与线上校验，不代表包整体没有问题，因为线上业务可能只用到了包的一部分。新业务在使用包时，依然需要严格测试。**

### 待验收包

实验性质，仅用于测试，成熟后移至正式包列表。

| package | 名称 | since | 说明 |
| ------- | ---- | ----- | ---- |
|         |      |       |      |

### 贡献须知

要求全面遵守开发规范，包括编码格式、拆分粒度、不限特定场景、避免过度封装等等。

要求维护使用手册，便于查找使用。由简要的列表说明，与详细说明组成。

业务未使用到的多余功能，不能顺带放进通用库包，以免误导其他使用者。

要求以 master 分支整体发行，标签为三级版本号。package 包版本体现在路径上，不兼容升级时必须创建新版本，业务类的包收录在厂商名下。

# 详细说明

## v1/config

基于 Viper 封装，单例模式，仅支持一个配置文件，yaml 格式。区分环境，优先读取-e 参数，为空时读取自定义环境设置。

默认环境为 release。建议将默认环境设置为生产环境，以确保在任何情况下不出错，同时阻止开发测试环境连接生产环境的资源。

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

更多使用方法，请查阅 [《Viper 文档》](https://github.com/spf13/viper/tree/v1.7.1#getting-values-from-viper)

## v1/db

基于 gorm，仅封装了配置与实例池，支持在并发下使用。

### 定义

**Add(name string, option Option)**

添加数据库配置，使用 db.Option。

**AddMap(name string, setting map[string]interface{})**

便捷方法，使用 map。

**AddMapBatch(batch map[string]interface{})**

便捷方法，使用 map 批量添加。

**Use(configName string) \*gorm.DB**

选择使用的数据库，根据配置组名称返回对应数据库实例，不存在时新建连接，留空使用 default 配置。

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

更多使用方法，请查阅 [《gorm 文档》](https://gorm.io/zh_CN/docs/index.html)

## v1/redis

基于 go-redis v8，仅封装了配置与实例池。

### 定义

**Add(name string, option Option)**

添加数据库配置，使用 redis.Option。

**AddMap(name string, setting map[string]interface{})**

便捷方法，使用 map。

**AddMapBatch(batch map[string]interface{})**

便捷方法，使用 map 批量添加。

**Use(name string) \*redis.Client**

使用数据库，根据配置组名称返回对应数据库实例，不存在时新建连接，留空使用 default 配置。

**Cluster(name string) \*redis.ClusterClient**

使用集群数据库，根据配置组名称返回对应数据库实例，不存在时新建连接。

**Lock(name string, expire time.Duration) bool**

加锁，expire 过期时间（单位秒）。

**Unlock(name string)**

解锁。可以不解锁，等过期。

**MaxMin 对象**

最大值最小值，记录、获取极值，超过极值才会覆盖。

### 实例

```yaml
# 配置 yaml
redis:
  default:
    address:
      - "0.0.0.0:6379"
    password: "" #密码，默认空
    db: 0
    pool_size: 100 #连接池最大数量，默认100
  cluster:
    address:
      - "0.0.0.0:6379"
      - "0.0.0.0:6380"
    password: "" #密码，默认空
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

更多使用方法，请查阅 [《go-redis 文档》](https://github.com/go-redis/redis/)

## v1/logger

基于官方 log 包，支持函数或对象两种封装，支持按级别发送通知。日志格式遵守 Json 规范。

### 定义

**SetLevel(level int)**

设置等级，小于该等级的日志将过滤。

**SetPrefix(prefix string)**

设置前缀，拼接在 title 字段之前。

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
import "github.com/lynnclub/go/v1/datetime"
import "gopkg.in/natefinch/lumberjack.v2"

// 自定义启动
// 默认为 New(log.New(os.Stderr, "", log.Lmsgprefix), DEBUG, "local", "asia/shanghai", nil)
lumberjack := &lumberjack.Logger{
	Filename:   "foo.log",
	MaxSize:    500, // megabytes
	MaxBackups: 3,
	MaxAge:     14,   //days
	Compress:   true, // disabled by default
}
logger.Logger = logger.New(
  log.New(lumberjack, "", log.Lmsgprefix),
  logger.DEBUG,
  "local",
  "asia/shanghai",
  datetime.LayoutDateTimeZoneT,
  config.Viper.GetStringMapString("feishu.group.alert"),
)

// 记日志
logger.Info("[通用队列发奖]启动")
logger.Error("错误", err)
```

## v1/rand

### 定义

**Range(min, max int) int**

区间随机，使用 nanoseconds 做种子，全开区间[min, max]

### 实例

```go
import "github.com/lynnclub/go/v1/rand"

// Range 区间随机
num := rand.Range(6, 8)
```

## v1/sign

### 定义

**MD5(params map[string]string, secret string) string**

常规 md5 get 拼接

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

获取本地 IP

**GetClientIP(c \*gin.Context) (string, error)**

获取 Header client-ip 的内容。

### 实例

```go
import "github.com/lynnclub/go/v1/ip"

// Local 本地IP
ips := ip.Local(true)

// Resolve 解析IP
response, errs := ip.Resolve("14.1.44.228")
```

## v1/datetime

官方内置的 time 包实现简洁，无需封装，建议简单场景下直接使用，比如获取当前时间戳。

datetime 包主要围绕时区封装，将字符串时间、时间戳等形式，统一解析为带时区的官包 time.Time 对象，然后以此为基础，提供更多快捷方法。提供函数与对象两种封装方式，效果相同。如需支持多个时区，建议使用函数形式；单一时区，建议使用对象形式。

**注意：如果 timezone 值无效，将自动使用系统时区**。

### 定义

**const LayoutDate**

常量，日期格式，2006-01-02

**const LayoutDateTime**

常量，日期时间格式，2006-01-02 15:04:05

**ParseDateTime(datetime string, timezone string) time.Time**

解析日期时间，基础函数。

**ToUnix(datetime string, timezone string) int64**

转成时间戳，基于 ParseDateTime 封装

**ParseTimestamp(timestamp int64, timezone string) time.Time**

解析时间戳，基础函数。

**ToAny(timestamp int64, timezone string, layout string) string**

转成任意格式，基于 ParseTimestamp 封装

**ToDate(timestamp int64, timezone string) string**

转成日期，基于 ToAny 封装

**ToDateTime(timestamp int64, timezone string) string**

转成日期时间，基于 ToAny 封装

**ToISOWeek(timestamp int64, timezone string) string**

转成年周 例如 2020_5，基于 ParseTimestamp 封装

**ToISOWeekByDate(datetime string, timezone string) string**

转成年周 例如 2020_5，基于 ParseDateTime 封装

**Unix(timezone string) int64**

当前时间戳，基于 ToUnix 封装

**Date(timezone string) string**

当前日期，基于 ToAny 封装

**DateTime(timezone string) string**

当前日期时间，基于 ToAny 封装

**ISOWeek(timezone string) string**

当前年周 例如 2020_5，基于 ToISOWeek 封装

**Any(timezone, layout string) string**

当前任意格式

**CheckTime(timestamp int64, start string, end string, timezone string) int**

检查时间 0 未开始、1 正常、2 已结束，半闭合区间[start, end)

**CheckTimeNow(start string, end string, timezone string) int**

检查当前时间 0 未开始、1 正常、2 已结束，半闭合区间[start, end)，基于 CheckTime 封装

**ClickhouseDatatimeRange() (time.Time, time.Time)**

clickhouse datetime 类型底层以时间戳存储，不包含时区，时间范围 [1970-01-01 00:00:00, 2106-02-07 06:28:15]

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

Json 编码

**Decode(str string, v interface{}) error**

Json 解码

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

SIGHUP 挂起（hangup），当终端关闭或者连接的会话结束时，由内核发送给进程  
SIGINT 中断（interrupt），通常由用户按下 Ctrl+C 产生，进程接收到信号后应立即停止当前的工作  
SIGQUIT 退出（quit），通常由用户按下 Ctrl+\ 产生，进程接收到信号后应立即退出，并清理自己占用的资源  
SIGTERM 终止（terminate），这是一个通用信号，通常用于要求进程正常终止  
SIGFPE 在发生致命的算术运算错误时发出，如除零操作、数据溢出等  
SIGKILL 立即结束程序的运行  
SIGALRM 时钟定时信号  
SIGBUS SIGSEGV 进程访问非法地址

相关知识：

1. 如果主程不循环也不阻塞，即使协程循环或阻塞，整个进程还是会退出。也就是说，常驻程序的主程应循环或阻塞，并且在平滑退出时，主程需要等协程先退出；
2. channel 作为协程的形参时，优先级很高，因为从 go1.14 开始，协程调度器是 基于信号的真抢占式调度；
3. channel 可以被多处接收，包括主程和协程；
4. select 需要关闭 channel 才会退出。
5. 建议结合 sync/atomic 原子计数器，实现高并发安全的协程平滑退出。

### 实例

```go
import "github.com/lynnclub/go/v1/signal"
import "sync"

// Listen 监听
signal.Listen()

// 模拟业务协程
var wg sync.WaitGroup
wg.Add(1)
go func() {
    for loop := 0; loop < 100; loop++ {
        // 收到停机信号，主动退出业务
        if signal.Now != nil {
            wg.Done()
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
  if Now != nil {
    wg.Wait()
    fmt.Println("main stop signal:", Now)
    break
  }

  // do something...

  time.Sleep(100 * time.Millisecond)
}

// 主程形式二：阻塞
// select {
// case <-ChannelOS:
// 	// 主程需要等待协程停止
// 	wg.Wait()
// 	fmt.Println("main stop signal:", Now)
// 	break
// }
```

## v1/array

### 定义

**In[T comparable](array []T, find T) bool**

是否存在

**NotIn[T comparable](array []T, find T) bool**

是否不存在

**Chunk[T any](array []T, chunkSize int) [][]T**

分组

**ChunkMap[K comparable, V any](array map[K]V, chunkSize int) []map[K]V**

Map 分组

**Keys[V comparable](array []V, find V) []int**

获取切片的 key

**KeysMap[K comparable](array map[K]V) []K**

获取 Map 的 key

**KeysFind[K comparable, V comparable](array map[K]V, find V) []K**

获取 Map 指定值的 key

**Values[K comparable, V any](array map[K]V) []V**

获取 Map 的 value

**Column[T any, N comparable, K comparable](array map[N]map[K]T, columnKey K) []T**

获取 Map 指定 column

**Diff[T comparable](a []T, b ...[]T) []T**

获取切片 a 中不存在于 b 的元素

**Intersect[T comparable](a, b []T) []T**

获取切片a、b的交集

**ToLower(array []string) []string**

转小写

**ToUpper(array []string) []string**

转大写

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

响应方法集合，目前支持 json。响应码建议遵守 http code 规范，使用官方 http 包的编码定义。

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
