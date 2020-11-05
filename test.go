package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"github.com/Dasem/trrp3client/pb"
)

var operations = map[string]bool{
	"MUL": true,
	"DIV": true,
	"SUB": true,
	"SUM": true,
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

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:5300", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewCalculatorClient(conn)
	request := &pb.Request{
		Operation: op,
		Operand1:  op1,
		Operand2:  op2,
	}

	ctx := context.Background()
	resp, err := client.Calculate(ctx, request)
	if err != nil {
		log.Fatalf("failed to invoke calculate: %v", err)
	}

	fmt.Println(resp.String())
}
