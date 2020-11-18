# 日志
## 级别
- 日志默认应分为6个级别，分别为：debug、info、warning、error、panic、fatal
- debug：调试用的日志级别，如果日志记录器可以分级，则在生产环境中不应该记录这部分日志
- info：日常信息级别
- warning：当出现一些次要错误，如类型不严谨，调用准备废弃的函数之类的情况下，应使用该级别日志
- error：当出现一些错误，但又不希望进程结束时应使用该级别日志，如出现网络错误，业务错误之类的
- panic：panic级别的日志在写完日志后，应调用panic方法结束进程
- fatal：fatal级别的日志在写完日志后，应调用os.Exit(1)方法结束进程

## 消息
- 每条日志必须有一条消息
- 消息必须为string类型

## 上下文
- 每条消息可以携带一个上下文数据
- 上下文为选填参数，可以不携带
- 接口中将上下文定义为可变长度参数是为了让它变为一个选填参数，而不必分别定义两组方法，所以不应传入多个context
- 上下文可以是任何类型的数据，如果是struct，调用方应保证希望记录下的字段可以被访问到

## 接口定义
日志记录器应实现此接口`log.LoggerInterface`
```go
type Logger interface {
	//Debug级别日志
	Debug(msg string, context ...interface{})
	//Info级别日志
	Info(msg string, context ...interface{})
	//Warning级别日志
	Warn(msg string, context ...interface{})
	//Error级别日志
	Error(msg string, context ...interface{})
	//Panic级别日志，应以panic(msg)方式结束进程
	Panic(msg string, context ...interface{})
	//Fatal级别日志，应以os.Exit(1)方式结束进程
	Fatal(msg string, context ...interface{})
	//自定义级别日志
	Log(level string, msg string, context ...interface{})
}
```

## 助手类
可以通过继承`log.AbstractLogger`，较快捷的实现`log.Logger`接口，它会将六个级别的日志和上下文转发到`Log`方法中。  
go的继承示例：go的继承比较麻烦，需要在子类中标记，并且要
```go

// 第一步：testLogger 继承于 AbstractLogger
type testLogger struct {
	AbstractLogger
}

func (l *testLogger) Log(level string, msg string, context ...interface{}) {}

func NewTestLogger() *testLogger {
	l := &testLogger{}
	// 第二步：将父类的指针改为自己？
	l.Logger = l
	return l
}

```

## 全局日志
可以通过`log.SetLogger`方法注册一个全局日志记录器，使用方法如下
```go
func init() {
	//注册
	log.SetLogger(NewTestLogger())
	//使用
	log.Info("test")
}
```