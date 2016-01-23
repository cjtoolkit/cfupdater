package updater

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

/*
Implements:
	loggerInterface
*/
type loggerMock struct {
	*log.Logger
	Buf *bytes.Buffer
}

func newLoggerMock() *loggerMock {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	return &loggerMock{logger, &bytes.Buffer{}}
}

func (lM loggerMock) Fatal(v ...interface{}) {
	lM.Panic(v...)
}

func (lM loggerMock) Fatalf(format string, v ...interface{}) {
	lM.Panicf(format, v...)
}

func (lM loggerMock) Fatalln(v ...interface{}) {
	lM.Panicln(v...)
}

func (lM loggerMock) Print(v ...interface{}) {
	fmt.Fprint(lM.Buf, v...)
}

func (lM loggerMock) Printf(format string, v ...interface{}) {
	fmt.Fprintf(lM.Buf, format, v...)
}

func (lM loggerMock) Println(v ...interface{}) {
	fmt.Fprintln(lM.Buf, v...)
}
