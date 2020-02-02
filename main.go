package main

import (
	"fmt"
	"os"

	"github.com/geek1011/dictutil/dictgen"
)

func main() {
	df, err := dictgen.ParseDictFile(os.Stdin)
	if err != nil {
		panic(err)
	}
	for _, dfe := range df {
		fmt.Printf("%#v\n", dfe)
	}
	df.WriteDictFile(os.Stdout)
}
