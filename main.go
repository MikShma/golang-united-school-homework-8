package main

import (
	"encoding/json"
	"errors"
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

func (args Arguments) AddOper(writer io.Writer) error {

	if args[arg_item] == "" {
		return errors.New("-item flag has to be specified")
	}

	content, err := ioutil.ReadFile(args[arg_fileName])
	if err != nil {
		return err
	}
	var people []Person
	if len(content) != 0 {

		if err := json.Unmarshal(content, &people); err != nil {
			return err
		}
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

	file, err := os.OpenFile(args[arg_fileName], os.O_RDWR|os.O_CREATE, filePerm)
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

func (args Arguments) RemoveOper(writer io.Writer) error {

	if args[arg_id] == "" {
		return errors.New("-id flag has to be specified")
	}

	content, err := ioutil.ReadFile(args[arg_fileName])
	if err != nil {
		return err
	}

	var people []Person
	if err := json.Unmarshal(content, &people); err != nil {
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

	file, err := os.OpenFile(args[arg_fileName], os.O_RDWR|os.O_TRUNC, filePerm)
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

func (args Arguments) FindOper(writer io.Writer) error {

	if args[arg_id] == "" {
		return errors.New("-id flag has to be specified")
	}

	content, err := ioutil.ReadFile(args[arg_fileName])
	if err != nil {
		return err
	}
	var people []Person
	if err := json.Unmarshal(content, &people); err != nil {
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

func (args Arguments) validateArgs() error {

	err := errors.New("flag has to be specified")

	if args[arg_fileName] == "" {
		return fmt.Errorf("-fileName %w", err)
	}

	if args[arg_operation] == "" {
		return fmt.Errorf("-operation %w", err)
	}

	validOpers := []string{"add", "remove", "findById", "list"}
	err = errors.New("Operation " + args[arg_operation] + " not allowed!")
	for _, oper := range validOpers {
		if args[arg_operation] == oper {
			return nil
		}
	}
	return err
}

func (args Arguments) openFile() error {
	file, err := os.OpenFile(args[arg_fileName], os.O_RDWR|os.O_CREATE, filePerm)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) (err error) {

	if err = args.validateArgs(); err != nil {
		return err
	}

	if err = args.openFile(); err != nil {
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
