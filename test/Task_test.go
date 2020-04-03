package test

import (
	"fmt"
	"task/app/models"
	"testing"
)

func TestInitTable(t *testing.T) {
	models.InitTable()
	t.Log("initTable")
}

func TestSave(t *testing.T) {
	task := &models.Task{
		Create: "test",
	}
	err := task.Save()
	if err != nil {
		t.Error("save fail")
	}
}

func TestUpdate(t *testing.T) {
	task := &models.Task{
		Create: "test",
	}
	task2, err := task.FindOne()
	if err != nil {
		t.Error("query fail")
	}
	task2.Create = "test2"
	task2.Update()

}

func TestDelete(t *testing.T) {
	task := &models.Task{
		Create: "test",
	}
	err := task.Delete()
	if err != nil {
		t.Error("delete fail")
	}
}

func TestSaveTasks(t *testing.T) {
	var tasks []models.Task
	for i := 0; i < 5; i++ {
		elem := &models.Task{
			Create: "test",
			Title:  fmt.Sprintf("测试%d", i),
		}
		tasks = append(tasks, *elem)
	}
	if err := (&models.Task{}).SaveList(&tasks); err != nil {
		t.Error(err)
	}
}
