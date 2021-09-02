package log


import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	pkgerr "github.com/pkg/errors"
)
/* 
type iLogger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})

	InfoDepth(string, int, ...interface{})
	WarnDepth(string, int, ...interface{})
	ErrorDepth(string, int, ...interface{})
} */

/* //Logger ....
type Logger struct{} */

var level string
var lvls []string

//var lg *logger

//Log ...
//var Log iLogger = lg

//Init ..
func Init(Logfile, lvl string) {

	f, err := os.OpenFile(Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	lvls = []string{"INFO", "WARN", "ERROR"}
	level = lvl
}

/*
######################################################
*/

//Info ..
func Info(mes string, args ...interface{}) {
	printfMsg("INFO", 0, mes, args...)
}

//Debug ..
func Debug(mes string, args ...interface{}) {
	printfMsg("WARN", 0, mes, args...)
}

//Error ..
func Error(mes string, args ...interface{}) {
	printfMsg("ERROR", 0, mes, args...)
}

/*
######################################################
*/

//InfoDepth ...
func InfoDepth(mes string, depth int, args ...interface{}) {
	printfMsg("INFO", depth, mes, args...)
}

//DebugDepth ..
func DebugDepth(mes string, depth int, args ...interface{}) {
	printfMsg("WARN", depth, mes, args...)
}

//ErrorDepth ..
func ErrorDepth(mes string, depth int, args ...interface{}) {
	printfMsg("ERROR", depth, mes, args...)
}

/*
--
--
--
--
--
--
--
--
--
*/

func printfMsg(level string, depth int, mes string, args ...interface{}) {
	// Chek for appropriate level of logging
	if checkLevel(level) {
		argsStr := getArgsString(args...)

		if level == "ERROR" {
			trace := fmt.Sprintf("'%+v'", pkgerr.New(""))
			if argsStr == "" {
				log.Printf("[%s] - %s - %s - trace = %v \n-\n", level, caller(depth+3), mes, trace)
				return
			}
			log.Printf("[%s] - %s - %s [%s] - trace = %v \n-\n", level, caller(depth+3), mes, argsStr, trace)
			return

		}
		if argsStr == "" {
			log.Printf("[%s] - %s - %s ", level, caller(depth+3), mes)
		} else {
			log.Printf("[%s] - %s - %s [%s] ", level, caller(depth+3), mes, argsStr)
		}
	}
}

func checkLevel(lvl string) bool {

	j := 0
	var str string
	for j, str = range lvls {
		if str == level {
			break
		}
	}
	for i, v := range lvls {
		if v == lvl {
			if j <= i {
				return true
			}
		}
	}
	return false
}

// getArgsString return formated string with arguments
func getArgsString(args ...interface{}) (argsStr string) {
	for _, arg := range args {
		if arg != nil {
			argsStr = argsStr + fmt.Sprintf("'%v', ", arg)
		}
	}
	argsStr = strings.TrimRight(argsStr, ", ")
	return
}

// caller returns a Valuer that returns a file and line from a specified depth in the callstack.
func caller(depth int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(depth+1, pc)
	frame, _ := runtime.CallersFrames(pc[:n]).Next()
	idxFile := strings.LastIndexByte(frame.File, '/')
	idx := strings.LastIndexByte(frame.Function, '/')
	idxName := strings.IndexByte(frame.Function[idx+1:], '.') + idx + 1

	return frame.File[idxFile+1:] + ":[" + strconv.Itoa(frame.Line) + "] - " + frame.Function[idxName+1:] + "()"
}
