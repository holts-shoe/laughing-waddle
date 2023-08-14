package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
  "time"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
  EndTime time.Time
  ResetTime time.Time
}

var db *gorm.DB

func InitDB() {
	dsn := "./task.db"
	var err error
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Task{})
}

func CreateTask(name, description string, resetTime time.Time) *Task {
  task := &Task{Name: name, Description: description, ResetTime: resetTime}
	db.Create(task)
	return task
}

func GetAllTasks() []Task {
	var tasks []Task
	db.Find(&tasks)
	return tasks
}

func GetTaskByID(id uint) *Task {
	var task Task
	db.First(&task, id)
	return &task
}

func UpdateTaskName(id uint, newName string) {
	var task Task
	db.First(&task, id)
	db.Model(&task).Update("Name", newName)
}
