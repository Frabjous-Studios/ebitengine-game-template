package debug

import "log"

var Enabled = true

func Println(v ...any) {
	if Enabled {
		log.Println(v...)
	}
}
func Printf(format string, v ...any) {
	if Enabled {
		log.Printf(format, v...)
	}
}
