package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	fmt.Fprint(writer, args)

	return nil
}

func parseArgs() Arguments {
	arguments := make(Arguments)
	idFlag := flag.String("id", "", "Id")
	itemFlag := flag.String("item", "", "Item")
	operationFlag := flag.String("operation", "", "Operation")
	fileNameFlag := flag.String("fileName", "", "File name")
	flag.Parse()
	arguments["id"] = *idFlag
	arguments["item"] = *itemFlag
	arguments["operation"] = *operationFlag
	arguments["fileName"] = *fileNameFlag

	return arguments
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
