package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var logFile *os.File

func SetLog(p string) {

	if !strings.HasSuffix(p, ".txt") {
		panic(".txt is not suffixed at logName Env")
	} else {
		if f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			if os.IsNotExist(err) {
				if f, err = os.Create(p); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			logFile = f
			log.SetOutput(f)
		}
	}

}

func GetLogFile() *os.File {
	return logFile
}

func InfoLog(w ...string) {
	msg := makeMsg("[INFO] ", w...)
	fmt.Println(msg)
	log.Printf(msg)
}

func ErrLog(w ...string) {
	msg := makeMsg("[ERR] ", w...)
	fmt.Println(msg)
	log.Printf(msg)
}

func CritLog(w ...string) {
	msg := makeMsg("[Crit] ", w...)
	fmt.Println(msg)
	log.Printf(msg)
	panic(w)
}

func makeMsg(baseMsg string, msg ...string) string {
	combinedMsg := strings.Join(msg, " ")
	// 기본 메시지와 결합된 메시지를 결합
	finalMsg := baseMsg + combinedMsg
	return finalMsg
}
