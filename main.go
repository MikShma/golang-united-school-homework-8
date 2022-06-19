package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	arg_fileName  = "fileName"
	arg_operation = "operation"
	arg_item      = "item"
	arg_id        = "id"
	filePerm      = 0644
	errSpecFlag   = "flag has to be specified"
)

type Arguments map[string]string

type Person struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func parseArgs() (args Arguments) {
	args = make(map[string]string)
	filePtr := flag.String("fileName", "", "file with users")
	operationPtr := flag.String("operation", "", "add remove findById")
	itemPtr := flag.String("item", "", "items to add")
	idPtr := flag.String("id", "", "person id")

	flag.Parse()

	args[arg_fileName] = *filePtr
	args[arg_item] = *itemPtr
	args[arg_operation] = *operationPtr
	args[arg_id] = *idPtr
	return args
}

func (args Arguments) ReadJsonFile(checkArg string) (people []Person, err error) {

	if args[checkArg] == "" {
		return nil, fmt.Errorf("-%s %s", checkArg, errSpecFlag)
	}

	content, err := ioutil.ReadFile(args[arg_fileName])
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, nil
	}

	if err := json.Unmarshal(content, &people); err != nil {
		return nil, err
	}
	return people, nil
}

func (args Arguments) WriteJsonFile(people []Person) error {

	file, err := os.OpenFile(args[arg_fileName], os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return err
	}

	people_out, err := json.Marshal(people)
	if err != nil {
		return err
	}
	file.Write(people_out)
	file.Close()
	return nil
}

func (args Arguments) AddOper(writer io.Writer) error {

	people, err := args.ReadJsonFile(arg_item)
	if err != nil {
		return err
	}
	var manToAdd Person
	if err := json.Unmarshal([]byte(args[arg_item]), &manToAdd); err != nil {
		return err
	}

	for _, p := range people {
		if p.Id == manToAdd.Id {
			writer.Write([]byte("Item with id " + manToAdd.Id + " already exists"))
			return nil
		}
	}

	people = append(people, Person{
		Id:    manToAdd.Id,
		Email: manToAdd.Email,
		Age:   manToAdd.Age,
	})

	if err := args.WriteJsonFile(people); err != nil {
		return err
	}

	return nil
}

func (args Arguments) RemoveOper(writer io.Writer) error {

	people, err := args.ReadJsonFile(arg_id)
	if err != nil {
		return err
	}
	notFound := true
	for i, p := range people {
		if p.Id == args[arg_id] {
			people = append(people[:i], people[i+1:]...)
			notFound = false
			break
		}
	}

	if notFound {
		writer.Write([]byte("Item with id " + args[arg_id] + " not found"))
	}

	if err := args.WriteJsonFile(people); err != nil {
		return err
	}

	return nil
}

func (args Arguments) FindOper(writer io.Writer) error {

	people, err := args.ReadJsonFile(arg_id)
	if err != nil {
		return err
	}

	for _, man := range people {
		if man.Id == args[arg_id] {
			man_out, err := json.Marshal(man)
			if err != nil {
				return err
			}
			writer.Write(man_out)
			return nil
		}
	}
	writer.Write([]byte(""))
	return nil
}

func (args Arguments) ListOper(writer io.Writer) error {
	data, err := os.ReadFile(args[arg_fileName])
	if err != nil {
		return err
	}

	writer.Write(data)
	return nil
}

func (args Arguments) validateArgsCreateFile() error {

	if args[arg_fileName] == "" {
		return fmt.Errorf("-%s %s", arg_fileName, errSpecFlag)
	}

	if args[arg_operation] == "" {
		return fmt.Errorf("-%s %s", arg_operation, errSpecFlag)
	}

	validOpers := []string{"add", "remove", "findById", "list"}
	err := fmt.Errorf("Operation %s not allowed!", args[arg_operation])
	for _, oper := range validOpers {
		if args[arg_operation] == oper {
			err = nil
		}
	}
	if err != nil {
		return err
	}

	file, err := os.OpenFile(args[arg_fileName], os.O_RDWR|os.O_CREATE, filePerm)
	file.Close()
	if err != nil {
		return err
	}

	return err
}

func Perform(args Arguments, writer io.Writer) (err error) {

	if err = args.validateArgsCreateFile(); err != nil {
		return err
	}

	switch args[arg_operation] {
	case "add":
		err = args.AddOper(writer)
	case "remove":
		err = args.RemoveOper(writer)
	case "findById":
		err = args.FindOper(writer)
	case "list":
		err = args.ListOper(writer)
	}

	return err
}

func main() {

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
