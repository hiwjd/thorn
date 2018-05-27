package main

import "fmt"

func main() {
	buf := make([]byte, 10)
	s := buf[0:9]
	//s = buf[0:10]
	fmt.Println(s)
}