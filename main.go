package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type Arguments map[string]string

type Users struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}
	fileName := args["fileName"]
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}
	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
		model := unmarshalFile(fileName)
		newItem := unmarshalString(item)
		newId, _ := strconv.Atoi(newItem.Id)
		if findUserbyId(model, newId) != "" {
			fmt.Fprintf(writer, "Item with id %d already exists", newId)
			return nil
		}
		model = append(model, newItem)
		marshalToFile(fileName, model)

	case "list":
		model := unmarshalFile(fileName)
		if len(model) > 0 {
			file, _ := ioutil.ReadFile(fileName)
			fmt.Fprint(writer, string(file))
		}

	case "findById":
		idF := args["id"]
		if idF == "" {
			return errors.New("-id flag has to be specified")
		}
		model := unmarshalFile(fileName)
		userId, _ := strconv.Atoi(idF)
		fmt.Fprint(writer, findUserbyId(model, userId))

	case "remove":
		idF := args["id"]
		if idF == "" {
			return errors.New("-id flag has to be specified")
		}
		model := unmarshalFile(fileName)
		userId, _ := strconv.Atoi(idF)
		result := deleteUserbyId(&model, userId)
		if result != "" {
			fmt.Fprint(writer, result)
		}
		marshalToFile(fileName, model)

	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

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

	fmt.Println(arguments)
	return arguments
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func unmarshalFile(filename string) []Users {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return []Users{}
	}
	var data []Users
	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func unmarshalString(content string) Users {
	var data Users
	_ = json.Unmarshal([]byte(content), &data)

	return data
}

func marshalToFile(filename string, list []Users) {
	j, _ := json.Marshal(list)
	ioutil.WriteFile(filename, j, 0644)
}

func deleteUserbyId(model *[]Users, id int) string {
	var index = -1
	for k, v := range *model {
		val, _ := strconv.Atoi(v.Id)
		if val == id {
			index = k
			break
		}
	}
	if index == -1 {
		return fmt.Sprintf("Item with id %d not found", id)
	}
	*model = append((*model)[0:index], (*model)[index+1:]...)
	return ""
}

func findUserbyId(model []Users, id int) string {
	for _, v := range model {
		val, _ := strconv.Atoi(v.Id)
		if val == id {
			r, _ := json.Marshal(v)
			return string(r)
		}
	}
	return ""
}
