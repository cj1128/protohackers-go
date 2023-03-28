package util

import (
	"os"

	"github.com/fatih/color"
)

var fatalf = color.New(color.FgRed, color.Bold).PrintfFunc()
var fatalln = color.New(color.FgRed, color.Bold).PrintlnFunc()

func Fatalf(msg string, v ...any) {
	fatalf(msg, v...)
	os.Exit(1)
}
func Fatalln(v ...any) {
	fatalln(v...)
	os.Exit(1)
}

var Infof = color.New(color.FgCyan).PrintfFunc()
var Infoln = color.New(color.FgCyan).PrintlnFunc()

var Debugf = color.New(color.FgBlack).PrintfFunc()
var Debugln = color.New(color.FgBlack).PrintlnFunc()

var Errorf = color.New(color.FgRed).PrintfFunc()
var Errorln = color.New(color.FgRed).PrintlnFunc()
