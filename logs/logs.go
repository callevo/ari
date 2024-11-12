package logs

import (
	"os"

	"github.com/rs/zerolog"
)

var TLogger zerolog.Logger

func init() {
	// remote syslog over unencrypted tcp
	//var zsyslog zerolog.SyslogWriter

	//fmt.Println("trying to connect t remote logs")
	//zsyslog, err := syslog.Dial("tcp", "172.30.0.216:514", syslog.LOG_KERN|syslog.LOG_EMERG|syslog.LOG_ERR|syslog.LOG_INFO|syslog.LOG_CRIT|syslog.LOG_WARNING|syslog.LOG_NOTICE|syslog.LOG_DEBUG, "queuelog")
	//if err != nil {
	//	panic(err)
	//}

	w := zerolog.MultiLevelWriter(os.Stdout)

	TLogger = zerolog.New(w).With().Timestamp().Caller().Logger().Level(zerolog.DebugLevel)

}
