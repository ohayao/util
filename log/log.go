package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Log struct {
	console *os.File
	file    *os.File
	mode    uint8
	lvs     int
	sync.Mutex
}

// Level 日志级别
type Level int

const (
	LvDebug Level = 1 << 0
	LvInfo  Level = 1 << 1
	LvWarn  Level = 1 << 2
	LvError Level = 1 << 3
	LvStack Level = 1 << 4
)

var lvs = map[Level]string{
	LvDebug: "DEBUG",
	LvInfo:  "INFO",
	LvWarn:  "WARN",
	LvError: "ERROR",
	LvStack: "STACK",
}

func New() *Log {
	return &Log{
		console: os.Stdout,
		mode:    0,
	}
}

// 设置日志级别
// lv_flag ==>> LvDebug|LvInfo
func (l *Log) SetLevel(lv_flag int) {
	l.lvs = lv_flag
}

// Close 关闭释放资源
func (l *Log) Close() {
	l.Lock()
	defer l.Unlock()
	if l.file != nil {
		l.file.Close()
	}
	l.console.Close()
}

// 设置日志文件
func (l *Log) SetFile(filepath string) error {
	l.Lock()
	defer l.Unlock()
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return err
		}
	}
	if file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644); err == nil {
		l.file = file
		return nil
	} else {
		return err
	}
}

// UseConsole 使用控制台
func (l *Log) UseConsole() {
	l.Lock()
	defer l.Unlock()
	l.mode = 0
}

// UseFile 使用文件
func (l *Log) UseFile() error {
	l.Lock()
	defer l.Unlock()
	if l.file == nil {
		return errors.New("no file set")
	}
	l.mode = 1
	return nil
}

func (l *Log) write(data string) {
	l.Lock()
	defer l.Unlock()
	if l.mode == 0 {
		l.console.WriteString(data)
	} else {
		l.file.WriteString(data)
	}
}

// dnf date.now.format
func dnf() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// 格式化日志信息
func fmtstring(lv Level, msg string) string {
	return fmt.Sprintf("[%s | % 5s] %s\n", dnf(), lvs[lv], msg)
}

// 格式化debug日志信息
func fmtdebugstring(lv Level, start, end int, msg string) string {
	var sb strings.Builder
	var last_file, last_line = "", 0
	for x := start; x <= end; x++ {
		if _, file, line, ok := runtime.Caller(x); ok {
			if last_file == file && last_line == line {
				break
			} else {
				if x > start {
					sb.WriteString(" ===>> ")
				}
				sb.WriteString(fmt.Sprintf("%s:%d", file, line))
			}
			last_file = file
			last_line = line
		} else {
			break
		}
	}
	return fmt.Sprintf("[%s | % 5s] [%s] %s\n", dnf(), lvs[lv], sb.String(), msg)
}

func (l *Log) Info(args ...any) {
	if l.lvs&int(LvInfo) == 0 {
		return
	}
	l.write(fmtstring(LvInfo, fmt.Sprint(args...)))
}
func (l *Log) Infof(format string, args ...any) {
	if l.lvs&int(LvInfo) == 0 {
		return
	}
	l.write(fmtstring(LvInfo, fmt.Sprintf(format, args...)))
}

func (l *Log) Warn(args ...any) {
	if l.lvs&int(LvWarn) == 0 {
		return
	}
	l.write(fmtstring(LvWarn, fmt.Sprint(args...)))
}
func (l *Log) Warnf(format string, args ...any) {
	if l.lvs&int(LvWarn) == 0 {
		return
	}
	l.write(fmtstring(LvWarn, fmt.Sprintf(format, args...)))
}

func (l *Log) Error(args ...any) {
	if l.lvs&int(LvError) == 0 {
		return
	}
	l.write(fmtstring(LvError, fmt.Sprint(args...)))
}
func (l *Log) Errorf(format string, args ...any) {
	if l.lvs&int(LvError) == 0 {
		return
	}
	l.write(fmtstring(LvError, fmt.Sprintf(format, args...)))
}

func (l *Log) Debug(args ...any) {
	if l.lvs&int(LvDebug) == 0 {
		return
	}
	l.write(fmtdebugstring(LvDebug, 2, 2, fmt.Sprint(args...)))
}
func (l *Log) Debugf(format string, args ...any) {
	if l.lvs&int(LvDebug) == 0 {
		return
	}
	l.write(fmtdebugstring(LvDebug, 2, 2, fmt.Sprintf(format, args...)))
}

// Json 使用格式化的Json格式日志信息
// data是Json数据会换行打印; args参数会在json格式前打印
func (l *Log) Json(lv Level, data any, args ...any) {
	if l.lvs&int(lv) == 0 {
		return
	}
	var sb strings.Builder
	bs, _ := json.MarshalIndent(data, "", "  ")
	if lv == LvDebug {
		if len(args) > 0 {
			sb.WriteString(fmtdebugstring(lv, 2, 2, fmt.Sprint(args...)))
		} else {
			sb.WriteString(fmtdebugstring(lv, 2, 2, ""))
		}
	} else {
		if len(args) > 0 {
			sb.WriteString(fmtstring(lv, fmt.Sprint(args...)))
		} else {
			sb.WriteString(fmtstring(lv, ""))
		}
	}
	sb.Write(bs)
	sb.WriteByte('\n')
	l.write(sb.String())
}

// Stack 调用堆栈信息
// depth 文件深度，最小=2
func (l *Log) Stack(depth int, args ...any) {
	if l.lvs&int(LvStack) == 0 {
		return
	}
	if depth < 2 {
		depth = 2
	}
	l.write(fmtdebugstring(LvStack, 2, depth, fmt.Sprint(args...)))
}

// Stackf 调用堆栈信息
// depth 文件深度，最小=2
func (l *Log) Stackf(depth int, format string, args ...any) {
	if l.lvs&int(LvStack) == 0 {
		return
	}
	if depth < 2 {
		depth = 2
	}
	l.write(fmtdebugstring(LvStack, 2, depth, fmt.Sprintf(format, args...)))
}
