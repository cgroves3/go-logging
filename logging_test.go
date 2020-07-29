package logging

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func LogTestHandlerFunc(writer http.ResponseWriter, request *http.Request) {
	expected := "Logged response."
	writer.Write([]byte(expected))
}

func NewServer(port string, handler http.Handler) http.Server {
	return http.Server{
		Addr:    port,
		Handler: handler,
	}
}

func TestLogHandlerFunc(t *testing.T) {
	expected := "Logged response."
	serverLogFileName := "server.log"
	file, _ := os.OpenFile(serverLogFileName, os.O_CREATE|os.O_WRONLY, 0666)
	Logger = log.New(file, "", 0)
	port := ":8000"
	server := NewServer(port, LogHandlerFunc(LogTestHandlerFunc))
	fmt.Printf("Starting server %v\n", server)
	go server.ListenAndServe()
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	fmt.Println("Running client")
	client.Get(fmt.Sprintf("http://localhost%s/", port))
	fmt.Println("Reading file.")
	actual := getLastNonEmptyLine(serverLogFileName)
	if actual != expected {
		t.Errorf("actual=%s, expected=%s", actual, expected)
	}
}

func getLastNonEmptyLine(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to read %s %v", fileName, err)
	}
	fileStr := string(bytes)
	lines := strings.Split(fileStr, "\n")
	last := ""
	for i := len(lines) - 1; i > 0; i-- {
		if len(lines[i]) > 0 {
			last = lines[i]
			break
		}
	}
	return last
}

type LogTestHandler struct {}
func (ltHandler LogTestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	LogTestHandlerFunc(writer, request)
}

func TestLogHandler(t *testing.T) {
	expected := "Logged response."
	serverLogFileName := "server2.log"
	file, _ := os.OpenFile(serverLogFileName, os.O_CREATE|os.O_WRONLY, 0666)
	Logger = log.New(file, "", 0)
	logTestHandler := LogTestHandler{}
	port := ":8001"
	server := NewServer(port, LogHandler(logTestHandler))
	fmt.Printf("Starting server %v\n", server)
	go server.ListenAndServe()
	go http.ListenAndServe(port, nil)
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	fmt.Println("Running client")
	client.Get(fmt.Sprintf("http://localhost%s/", port))
	fmt.Println("Reading file.")
	actual := getLastNonEmptyLine(serverLogFileName)
	if actual != expected {
		t.Errorf("actual=%s, expected=%s", actual, expected)
	}
}
