![](https://img.shields.io/badge/version-v1.x-green.svg) &nbsp; ![](https://img.shields.io/badge/version-go1.21-green.svg) &nbsp;  ![](https://img.shields.io/badge/builder-success-green.svg) &nbsp;

> Zap 日志框架助手，旨在简化配置及使用，开发者无需过多关心日志配置。

## 一、前言

### 1、安装

* Get

```bash
go get github.com/archine/zaplogger@v1.0.3
```

* Mod

```bash
# go.mod文件加入下面的一条
github.com/archine/zaplogger v1.0.3

# 命令行在该项目目录下执行
go mod tidy
```

## 二、使用说明

### 1、初始化

下面的代码展示了如何初始化 zaplogger，只需调用 Init 方法即可，无需关心其他配置。在实际使用中，可以根据需要自定义配置。<br/>
具体配置项请参考下面的配置项说明。

```go
package main

import "github.com/archine/zaplogger"

func main() {
	// 全部使用默认配置
	var conf zaplogger.Config
	// 调用初始化方法
	err := zaplogger.Init(&conf)
	if err != nil {
		panic(err)
	}
}
```

* 配置项说明：

| 属性               | 描述                                                                 | 默认值     |
|------------------|--------------------------------------------------------------------|---------|
| Level            | 日志级别，可选项： info、warn、error、panic、fatal                              | debug   |
| LevelColor       | 日志级别色彩，只在 format 等于 console 时生效                                    | False   |
| Format           | 日志格式，可选项：json。console 会将日志按照指定分隔符打印，<br/>json 则按键值对形式打印            | console |
| ConsoleSeparator | 日志格式为 console 时的 分隔符                                               | \t      |
| PrintStacktrace  | 日志级别为 Error 及以上时是否打印堆栈                                             | False   |
| Options          | 可选属性，默认追加 zap.AddCaller 和 zap.AddStaceTrace                        |         |
| ApplyFields      | 当打印日志通过 ``zaplogger.WithContext(ctx)``调用时，会触发该函数，<br/>可通过该函数来动态设置值 |         |
| ApplyEncoder     | 编码器，默认根据 Format 来选择                                                |         |
| ApplyCores       | 指定zap核心，默认打印到控制台，可通过该属性来配置如写入到文件等                                  |         |

### 2、基本使用

```go
package main

import "github.com/archine/zaplogger"

func main() {
	// 基于默认配置初始化...
	// ...

	// 使用
	zaplogger.Info("我是Info日志")
	zaplogger.Warn("我是Warn日志")
	zaplogger.Debug("我是Debug日志")
	zaplogger.Error("我是错误日志")
}
```

控制台输出

```bash
2024-06-02 18:28:18   INFO   我是Info日志
2024-06-02 18:28:18   WARN   我是Warn日志
2024-06-02 18:28:18   DEBUG  我是Debug日志
2024-06-02 18:28:18   ERROR  我是错误日志
```

### 3、设置分隔符

```go
package main

import "github.com/archine/zaplogger"

func main() {
	conf := &zaplogger.Config{
		BasicConfig: zaplogger.BasicConfig{
			ConsoleSeparator: " | ",
		},
	}
	_ = zaplogger.Init(conf)
	zaplogger.Info("我是Info日志")
}
```

控制台输出

```bash
2024-06-02 18:33:39 | INFO | 我是Info日志
```

### 4、开启日志级别色彩

```go
package main

import "github.com/archine/zaplogger"

func main() {
	conf := &zaplogger.Config{
		BasicConfig: zaplogger.BasicConfig{
			LevelColor: true, // 开启日志级别色彩
		},
	}
	_ = zaplogger.Init(conf)
	zaplogger.Info("我是Info日志")
	zaplogger.Warn("我是Warn日志")
	zaplogger.Debug("我是Debug日志")
	zaplogger.Error("我是错误日志")
}
```

控制台输出

![img](https://github.com/archine/zaplogger/assets/35919643/e4d1e95e-92c5-4512-94dd-71ca664d7516)


### 5、打印堆栈

```go
package main

import "github.com/archine/zaplogger"

func main() {
	conf := &zaplogger.Config{
		BasicConfig: zaplogger.BasicConfig{
			PrintStacktrace: true, // 打印堆栈
		},
	}
	_ = zaplogger.Init(conf)
	zaplogger.Error("我是错误日志")
}
```

控制台输出

```bash
2024-06-02 18:38:18  ERROR  我是错误日志
goroutine 1 [running]:
main.main()
    /Users/archine/go/src/github.com/archine/zaplogger/example/main.go:13 +0x1b4
exit status 1
```

### 6、写入日志到文件

我们可以通过 ApplyCores 方法来指定写入日志到文件，下面的示例展示了如何将日志写入到文件中。这里需要注意的是<br/>
`文件所在目录要存在，否则会报错`。`zap.Open()`支持传入多个参数，如需同时打印控制台，那么可以传入`"stdout"`<br/>
或者
`"stderr"`。

```go
package main

import (
	"github.com/archine/zaplogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	conf := &zaplogger.Config{
		ApplyCores: func(enc zapcore.Encoder, level zapcore.LevelEnabler, conf zaplogger.BasicConfig) (zapcore.Core, error) {
			mySyncer, closeAll, err := zap.Open("./app.log")
			// mySyncer, closeAll, err := zap.Open("./app.log", "stdout")
			if err != nil {
				closeAll()
				return nil, err
			}
			return zapcore.NewCore(enc, mySyncer, level), nil
		},
	}
	if err := zaplogger.Init(conf); err != nil {
		panic(err)
	}
	zaplogger.Info("我是Info日志")
}

```

### 7、日志文件切割

我们可以通过 `lumberjack` 库来实现日志文件的切割，下面的示例展示了如何实现日志文件的切割。

```go
package main

import (
	"github.com/archine/zaplogger"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	conf := &zaplogger.Config{
            ApplyCores: func(enc zapcore.Encoder, level zapcore.LevelEnabler, conf zaplogger.BasicConfig) (zapcore.Core, error) {
                splitSyncer := zapcore.AddSync(&lumberjack.Logger{
                    Filename:   "./app.log",
                    MaxSize:    1,    // 每个日志文件的最大尺寸，单位是 MB
                    MaxBackups: 3,    // 保留旧文件的最大个数
                    MaxAge:     7,    // 保留旧文件的最大天数
                    LocalTime:  true, // 使用本地时间
                    Compress:   true, // 是否压缩旧文件
                })
                return zapcore.NewCore(enc, splitSyncer, level), nil
            },
	}
	if err := zaplogger.Init(conf); err != nil {
		panic(err)
	}
	zaplogger.Info("我是Info日志")
}
```

### 8、动态设置字段

我们可以通过 `ApplyFields` 方法来动态设置字段，下面的示例展示了如何动态设置字段。

```go
package main

import (
	"context"
	"github.com/archine/zaplogger"
	"go.uber.org/zap"
)

func main() {
    conf := &zaplogger.Config{
        ApplyFields: func(ctx context.Context) []zap.Field {
            var fields []zap.Field
            // 读取ctx中的trace_id，然后设置到日志中
            traceId := ctx.Value("trace_id")
            if traceId != nil {
                fields = append(fields, zap.String("trace_id", traceId.(string)))
            }
            return fields
        },
    }
    if err := zaplogger.Init(conf); err != nil {
        panic(err)
    }
    
    // 通过context传递trace_id
    ctx := context.WithValue(context.Background(), "trace_id", "123456")
    
    zaplogger.WithContext(ctx).Info("我是Info日志")
}
```
控制台输出
```bash
2024-06-02 21:13:51  INFO  我是Info日志    {"trace_id": "123456"}
```

### 9、替换Gin-Plus 默认的logger
通过监听器去替换Gin-Plus默认的logger，下面的示例展示了如何替换Gin-Plus默认的logger。``以下基于Gin-Plus v3.3.0版本``，低于该版本请自行调整。
```go
package main

import (
	"github.com/archine/gin-plus/v3/application"
	"github.com/archine/zaplogger"
	"github.com/archine/zaplogger/ginplus"
	"github.com/spf13/viper"
)

type AppConfigListener struct {}

func (a *AppConfigListener) Read(v *viper.Viper) error {
    return v.ReadInConfig()
}

func (a *AppConfigListener) After(v *viper.Viper) {
    var conf zaplogger.Config
    if err := zaplogger.Init(&conf); err != nil {
        panic(err)
    }
    application.ChangeLogger(&zaplogger.GinPlusLoggerImpl{})
}

//go:generate gp-ast
func main() {
    application.Default(&AppConfigListener{}).Run()
}
```
控制台输出
```bash
2024-06-02 21:22:13   INFO   Application start success on Ports:[4006]
```

---
以上是一些基本的使用案例，更多的方式大家可以根据自己的需求来进行配置。

