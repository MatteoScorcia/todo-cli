package main

func main() {
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	storage.Load(&todos)
	cmdFlags := NewCmdFlags()
	cmdFlags.Execute(&todos)
	storage.Save(todos)
	// todos := Todos{}
	// todos.Add("Hello", "matteoscorcia")
	// todos.UpdateStatus(1, StatusCompleted)
	// todos.Print()
}
