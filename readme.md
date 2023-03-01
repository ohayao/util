## UTIL
> golang library

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