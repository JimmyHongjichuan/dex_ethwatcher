package ethwatcher

import (
	"bytes"
	"fmt"
	"os"

	log "github.com/inconshreveable/log15"
)

type logLevel log.Lvl

const (
	DEBUG    = logLevel(log.LvlDebug)
	INFO     = logLevel(log.LvlInfo)
	WARN     = logLevel(log.LvlWarn)
	ERROR    = logLevel(log.LvlError)
	CRITICAL = logLevel(log.LvlCrit)

	timeFormat  = "2006-01-02 15:04:05.000"
	termMsgJust = 40
)

func (lv logLevel) String() string {
	switch lv {
	case DEBUG:
		return "DBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "EROR"
	case CRITICAL:
		return "CRIT"
	default:
		panic("bad level")
	}
}

func toLogLevel(lvl string) logLevel {
	switch lvl {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "critical":
		return CRITICAL
	default:
		return DEBUG
	}
}

func printLogCtx(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		k, ok := ctx[i].(string)
		v := fmt.Sprintf("%+v", ctx[i+1])
		if !ok {
			k, v = "ERR", fmt.Sprintf("%+v", k)
		}
		fmt.Fprintf(buf, " %s=%s", k, v)
	}
	buf.WriteByte('\n')
}

func terminalFormat() log.Format {
	return log.FormatFunc(func(r *log.Record) []byte {
		b := &bytes.Buffer{}
		level := logLevel(r.Lvl)

		fmt.Fprintf(b, "[%v][%s] %s", level, r.Time.Format(timeFormat), r.Msg)
		// try to justify the log output for short messages
		if len(r.Ctx) > 0 && len(r.Msg) < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-len(r.Msg)))
		}

		printLogCtx(b, r.Ctx)
		return b.Bytes()
	})
}

// New 返回一个logger对象
func NewLogger(lvl string, module string) log.Logger {
	handler := log.LvlFilterHandler(
		log.Lvl(toLogLevel(lvl)),
		log.CallerFileHandler(
			log.StreamHandler(os.Stdout, terminalFormat())))
	logger := log.New("module", module)
	logger.SetHandler(handler)
	return logger
}
