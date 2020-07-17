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

func LogHandler(writer http.ResponseWriter, request *http.Request) {
	expected := "Logged response."
	writer.Write([]byte(expected))
}

func TestLogHandler(t *testing.T) {
	expected := "Logged response."
	fileName := "server.log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	Logger = log.New(file, "", 0)
	http.HandleFunc("/", LogHandlerFunc(LogHandler))
	port := ":8000"
	fmt.Println("Starting server")
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
	actual := last
	fmt.Printf("last=%s", last)
	if actual != expected {
		t.Errorf("actual=%s, expected=%s", actual, expected)
	}
}
