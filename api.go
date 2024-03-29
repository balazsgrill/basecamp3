package basecamp3

import "fmt"

func (bc *Basecamp) TodoSet_Lists(ctx ContextWithTokenPersistence, account int, project int, todoset int64) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todosets/%d/todolists.json", BasecampApiRootURL, account, project, todoset)
	var result []Todo
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) TodoSet(ctx ContextWithTokenPersistence, account int, project int, todoset int64) (TodoSet, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todosets/%d.json", BasecampApiRootURL, account, project, todoset)
	var result TodoSet
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) Project(ctx ContextWithTokenPersistence, account int, project int) (Project, error) {
	url := fmt.Sprintf("%s/%d/projects/%d.json", BasecampApiRootURL, account, project)
	var result Project
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) Projects(ctx ContextWithTokenPersistence, account int) ([]Project, error) {
	url := fmt.Sprintf("%s/%d/projects.json", BasecampApiRootURL, account)
	var result []Project
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) TodoList(ctx ContextWithTokenPersistence, account int, project int, todolist int) (Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d.json", BasecampApiRootURL, account, project, todolist)
	var result Todo
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) TodoList_Groups(ctx ContextWithTokenPersistence, account int, project int, todolist int) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d/groups.json", BasecampApiRootURL, account, project, todolist)
	var result []Todo
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) Todos(ctx ContextWithTokenPersistence, account int, project int, todolist int) ([]Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todolists/%d/todos.json", BasecampApiRootURL, account, project, todolist)
	var result []Todo
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}

func (bc *Basecamp) Todo(ctx ContextWithTokenPersistence, account int, project int, todo int) (Todo, error) {
	url := fmt.Sprintf("%s/%d/buckets/%d/todos/%d.json", BasecampApiRootURL, account, project, todo)
	var result Todo
	err := bc.jsonGet(ctx, url, &result)
	return result, err
}
