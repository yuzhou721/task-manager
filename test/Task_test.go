package test

import (
	"task/app/models"
	"testing"
)

// func TestInitTable(t *testing.T) {
// 	models.InitTable()
// 	t.Log("initTable")
// }

// func TestSave(t *testing.T) {
// 	task := &models.Task{
// 		Assigner: "test",
// 	}
// 	err := task.Save()
// 	if err != nil {
// 		t.Error("save fail")
// 	}
// }

func TestUpdate(t *testing.T) {

	task := new(models.Task)
	task.ID = 9
	task2, err := task.FindOne()
	if err != nil {
		t.Error("query fail")
	}
	task2.Status = models.TaskStatusDone
	task2.Update()

}

// func TestDelete(t *testing.T) {
// 	task := &models.Task{
// 		Assigner: "test",
// 	}
// 	err := task.Delete()
// 	if err != nil {
// 		t.Error("delete fail")
// 	}
// }

// func TestSaveTasks(t *testing.T) {
// 	var tasks []models.Task
// 	for i := 0; i < 5; i++ {
// 		elem := &models.Task{
// 			Assigner: "test",
// 			Title:    fmt.Sprintf("测试%d", i),
// 		}
// 		tasks = append(tasks, *elem)
// 	}
// 	if err := (&models.Task{}).SaveList(&tasks); err != nil {
// 		t.Error(err)
// 	}
// }

func TestGetTask(t *testing.T) {
	var task models.Task
	task.FindByID("2")
}

func TestGetTasks(t *testing.T) {
	tasks, count, err := (&models.Task{}).FindList(2, "5df05b6bd08e4390b7f7306f", "273413ce-1a56-11ea-9751-0050569293b2", "", 1, 10)
	if err != nil {
		t.Error("error :", err)
	}
	t.Logf("task=%v count=%v err=%v", tasks, count, err)

}
