package raft

import (
	"fmt"
	"log"
	"os"
)

// Debugging
const Debug = true

var f *os.File
var logger *log.Logger

func init() {
	file := "d:/Code/go/src/新建文件夹/raft-go/src/raft/raft_logger.txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	logger = log.New(logFile, "", log.LstdFlags) // 将文件设置为loger作为输出
	return
}

func DPrintf(format string, a ...interface{}) {
	if Debug {
		logger.Println(fmt.Sprintf(format, a...))
	}
	return
}

func DPrintVerbose(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
