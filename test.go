package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var operations = map[string]bool{
	"MUL": true,
	"DIV": true,
	"SUB": true,
	"SUM": true,
}

type Math struct {
	Operation string  `json:"operation"`
	Operand1  float64 `json:"operand1"`
	Operand2  float64 `json:"operand2"`
}

type REPLReader struct {
	r *bufio.Reader
}

func NewREPLReader() *REPLReader {
	return &REPLReader{r: bufio.NewReader(os.Stdin)}
}

func (r *REPLReader) printAndRead(msg string) (string, error) {
	fmt.Println(msg)
	v, err := r.r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read string from stdin: %w", err)
	}
	return strings.TrimSpace(v), nil
}

func (r *REPLReader) ReadString(msg string) (string, error) {
	return r.printAndRead(msg)
}

func (r *REPLReader) ReadFloat(msg string) (float64, error) {
	str, err := r.printAndRead(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to print and read: %w", err)
	}
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse data: %w", err)
	}

	return res, nil
}

func main() {
	r := NewREPLReader()

	op1, err := r.ReadFloat("Input first operand -> ")
	if err != nil {
		log.Fatalf("failed to read op1: %v", err)
	}

	op2, err := r.ReadFloat("Input second operand -> ")
	if err != nil {
		log.Fatalf("failed to read op2: %v", err)
	}

	op, err := r.ReadString("Input operation(SUM, SUB, MUL, DIV) -> ")
	if err != nil {
		log.Fatalf("failed to read operation: %v", err)
	}

	if _, ok := operations[op]; !ok {
		log.Fatalf("invalid operation: %v", op)
	}

	requestBody, err := json.Marshal(Math{op, op1, op2})
	if err != nil {
		log.Fatalf("failed to marshal body: %v", err)
	}

	b := bytes.NewReader(requestBody)
	resp, err := http.Post("http://127.0.0.1:8080/calculate", "application/json", b)
	if err != nil {
		log.Fatalf("failed to do request: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read resp body: %v", err)
	}

	log.Printf("%s", body)
}
