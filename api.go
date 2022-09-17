package basecamp3

import "fmt"

func (bc *Basecamp) TodoSet_Lists(account int, project int, todoset int) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todosets/%d/todolists.json", BasecampApiRootURL, account, project, todoset)
	var result []Todo
	err := bc.jsonGet(url, result)
	return result, err
}

func (bc *Basecamp) TodoSet(account int, project int, todoset int) (TodoSet, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todosets/%d.json", BasecampApiRootURL, account, project, todoset)
	var result TodoSet
	err := bc.jsonGet(url, &result)
	return result, err
}

func (bc *Basecamp) Project(account int, project int) (Project, error) {
	url := fmt.Sprintf("%s/%d/projects/%d.json", BasecampApiRootURL, account, project)
	var result Project
	err := bc.jsonGet(url, &result)
	return result, err
}

func (bc *Basecamp) TodoList(account int, project int, todolist int) (Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d.json", BasecampApiRootURL, account, project, todolist)
	var result Todo
	err := bc.jsonGet(url, &result)
	return result, err
}

func (bc *Basecamp) TodoList_Groups(account int, project int, todolist int) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d/groups.json", BasecampApiRootURL, account, project, todolist)
	var result []Todo
	err := bc.jsonGet(url, result)
	return result, err
}

func (bc *Basecamp) Todos(account int, project int, todolist int) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d/todos.json", BasecampApiRootURL, account, project, todolist)
	var result []Todo
	err := bc.jsonGet(url, result)
	return result, err
}
