package logrusutil

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/tsaikd/KDGoLib/errutil"
)

var (
	DefaultConsoleLogger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &ConsoleLogFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
)

const (
	Llongfile  = log.Llongfile
	Lshortfile = log.Lshortfile
	Ltime      = log.Ltime
	Llevel     = log.Lshortfile << 1
	LstdFlags  = Ltime | Lshortfile | Llevel
)

var (
	isTerminal = logrus.IsTerminal()
)

type ConsoleLogFormatter struct {
	TimestampFormat string
	Flag            int
}

func addspace(text string, addspaceflag bool) (string, bool) {
	if addspaceflag {
		return " " + text, true
	} else {
		return text, true
	}
}

func (t *ConsoleLogFormatter) Format(entry *logrus.Entry) (data []byte, err error) {
	buffer := bytes.Buffer{}
	addspaceflag := false

	if t.Flag == 0 {
		t.Flag = LstdFlags
	}

	if t.TimestampFormat == "" {
		t.TimestampFormat = logrus.DefaultTimestampFormat
	}

	if t.Flag&Ltime != 0 {
		timetext := entry.Time.Format(t.TimestampFormat)
		timetext, addspaceflag = addspace(timetext, addspaceflag)
		if _, err = buffer.WriteString(timetext); err != nil {
			err = errutil.New("write timestamp to buffer failed", err)
			return
		}
	}

	if t.Flag&(Lshortfile|Llongfile) != 0 {
		_, file, line, ok := runtime.Caller(7)
		if !ok {
			file = "???"
			line = 0
		}

		if t.Flag&Lshortfile != 0 {
			file = filepath.Base(file)
		}

		filelinetext := fmt.Sprintf("%s:%d", file, line)
		filelinetext, addspaceflag = addspace(filelinetext, addspaceflag)
		if _, err = buffer.WriteString(filelinetext); err != nil {
			err = errutil.New("write fileline to buffer failed", err)
			return
		}
	}

	if t.Flag&Llevel != 0 {
		leveltext := fmt.Sprintf("[%s]", entry.Level.String())
		leveltext, addspaceflag = addspace(leveltext, addspaceflag)
		if _, err = buffer.WriteString(leveltext); err != nil {
			err = errutil.New("write level to buffer failed", err)
			return
		}
	}

	message := entry.Message
	message, addspaceflag = addspace(message, addspaceflag)
	if _, err = buffer.WriteString(message); err != nil {
		err = errutil.New("write message to buffer failed", err)
		return
	}

	if err = buffer.WriteByte('\n'); err != nil {
		err = errutil.New("write newline to buffer failed", err)
		return
	}

	data = buffer.Bytes()
	return
}
