package logging

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var Logger *log.Logger

func NewFileLogger(fileName string) *log.Logger {
	var file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func toString(request *http.Request) string {
	requestStr := fmt.Sprintf("%v %v %v\n", request.Method, request.URL, request.Proto)
	requestStr += fmt.Sprintf("Host: %v\n", request.Host)
	for name, headers := range request.Header {
		requestStr += fmt.Sprintf("%s: %v\n", name, headers)
	}
	requestStr += fmt.Sprintf("%v\n", request.Form.Encode())
	return requestStr
}

func LogHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func (writer http.ResponseWriter, request *http.Request) {
		if Logger != nil {
			Logger.Printf("Request: %v", toString(request))
		} else {
			fmt.Println("Logger was nil. Unable to log.")
		}
		logWriter := LogResponseWriter{
			writer,
		}
		handlerFunc(logWriter, request)
	}
}

type LogResponseWriter struct {
	http.ResponseWriter
}

func (logWriter LogResponseWriter) Write(b []byte) (n int, err error) {
	if Logger != nil {
		headerStr := "Headers"
		for name, headers := range logWriter.ResponseWriter.Header() {
			headerStr += fmt.Sprintf("%s: %v\n", name, headers)
		}
		Logger.Printf("Response: %v \n%s", headerStr, string(b))
	} else {
		fmt.Println("Logger was nil. Unable to log.")
	}
	return logWriter.ResponseWriter.Write(b)
}