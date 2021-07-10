package main

import (
	"context"
	"fmt"
	"strconv"
)

var ctx = context.Background()

func main() {
	s:=fmt.Sprintf("%07d", 2)
	fmt.Println(s)
	id, _ := strconv.Atoi(s)
	print(id)

}
