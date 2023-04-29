## UTIL
> golang library

### log
``` go
// 新建
var logger=log.New()
// 在控制台打印
logger.UseConsole()
// 在文件中打印
// 需要先传输一个文件
logger.SetFile('log_file_path')
logger.UseFile()
// 控制台和文件可以切换打印，按最后一个Use为准
// 打印有5种标记 Debug|Info|Warn|Error|Stack
// 每种标记都可以自己指定格式化方式
// 其中Debug会打印所在的文件名以及行号
// Stack指定深度来打印调用的堆栈
// 需要指定打印的标记，使用或‘｜’操作符添加多种
logger.SetLevel(int(log.LvInfo|log.LvDebug))
logger.Info("hello","abc")
logger.Infof(`{"name":"%s","age":%d}`,"James",25)
logger.Stack(5,"从调用的地方可是向上查找5个调用堆栈","如果深度大，则到顶为止")
// 支持把数据打印成Json结构的形式，方便查看
logger.Json(log.LvInfo,map[string]any{"name":"zhangsan","age":14})
```

### 测试用例
``` shell
cd http
# 测试所有
go test -v
# 测试指定方法
go test -test.run TestUserCookie -v
# 基准测试 -bench=通配符 （.所有）-run=不存在的函数用来排除测试方法
go test -bench=BenchmarkGetMenus -run=123 -v
```