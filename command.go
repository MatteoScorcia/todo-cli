package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdFlags struct {
	Add          string
	Del          int
	Edit         string
	UpdateStatus string
	List         bool
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo title")
	flag.StringVar(&cf.Edit, "edit", "", "Edit a todo by id & specify a new title. id:new_title")
	flag.StringVar(&cf.UpdateStatus, "updateStatus", "", "Edit a todo by id & specify a new status. id:new_status")
	flag.IntVar(&cf.Del, "del", -1, "Specify a todo by id to delete")
	flag.BoolVar(&cf.List, "list", false, "List all todos")

	flag.Parse()

	return &cf
}

func (cf *CmdFlags) Execute(todos *Todos) {
	switch {
	case cf.List:
		todos.Print()

	case cf.Add != "":
		todos.Add(cf.Add, "")

	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error, invalid format for edit. Please use id:new_title")
			os.Exit(1)
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error: invalid id for edit")
			os.Exit(1)
		}

		todos.EditTitle(id, parts[1])

	case cf.UpdateStatus != "":
		parts := strings.SplitN(cf.UpdateStatus, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error, invalid format for status update. Please use id:new_status")
			os.Exit(1)
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error: invalid id for updateStatus")
			os.Exit(1)
		}

		parsedStatus, err := ParseStatus(parts[1])
		if err != nil {
			fmt.Println("Error: invalid status for updateStatus")
			os.Exit(1)
		}

		todos.UpdateStatus(id, parsedStatus)

	case cf.Del != -1:
		todos.Delete(cf.Del)

	default:
		fmt.Println("Invalid statement")
	}
}
