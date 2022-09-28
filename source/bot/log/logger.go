package log

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger = &log.Logger{}
	errLogger   *log.Logger = &log.Logger{}
	debugLogger *log.Logger = &log.Logger{}
	Logger      *log.Logger = &log.Logger{}
)

func init() {

	Logger.SetOutput(os.Stdout)
	infoLogger.SetOutput(os.Stdout)
	errLogger.SetOutput(os.Stdout)
	debugLogger.SetOutput(os.Stdout)

	infoLogger.SetPrefix("[info] ")
	errLogger.SetPrefix("[error] ")
	debugLogger.SetPrefix("[debug] ")

	infoLogger.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	errLogger.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	debugLogger.SetFlags(log.Ldate | log.Ltime | log.LUTC)

}

func Info(info string) {

	infoLogger.Println(info)

}

func Error(err string) {

	errLogger.Println(err)

}

func Debug(debug string) {

	debugLogger.Println(debug)

}
