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

/*
   Operation operation;
   BigDecimal operand1;
   BigDecimal operand2;
*/

const (
	MUL = "MUL"
	DIV = "DIV"
	SUB = "SUB"
	SUM = "SUM"
)

type Math struct {
	Operation string  `json:"operation"`
	Operand1  float64 `json:"operand1"`
	Operand2  float64 `json:"operand2"`
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Input first operand -> ")
	operand1, _ := reader.ReadString('\n')
	operand1 = strings.TrimSpace(operand1)
	operand1Float, err := strconv.ParseFloat(operand1, 64)

	fmt.Println("Input second operand -> ")
	operand2, _ := reader.ReadString('\n')
	operand2 = strings.TrimSpace(operand2)
	operand2Float, err := strconv.ParseFloat(operand2, 64)

	fmt.Println("Input operation(SUM, SUB, MUL, DIV) -> ")
	operationWithSpaces, _ := reader.ReadString('\n')

	operation := strings.TrimSpace(operationWithSpaces)

	/*	requestBody, err := json.Marshal(map[string]string{
		"operation": operation,
		"operand1":  operand1,
		"operand2":  operand2,
	})*/
	requestBody, err := json.Marshal(Math{operation, operand1Float, operand2Float})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/calculate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}
