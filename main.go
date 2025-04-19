package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

type (
	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}
)

var (
	taskIDCounter int64 = 1
	tasks               = make(map[int64]Task)
)

func main() {
	webApp := fiber.New(fiber.Config{
		ReadBufferSize: 16 * 1024})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	taskHandler := TaskHandler{}
	webApp.Post("/tasks", taskHandler.CreateTask)
	webApp.Patch("/tasks/:id", taskHandler.UpdateTask)
	webApp.Get("/tasks/:id", taskHandler.GetTask)
	webApp.Delete("/tasks/:id", taskHandler.DeleteTask)

	logrus.Fatal(webApp.Listen(":8080"))
}

type TaskHandler struct{}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var request CreateTaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}
	task := Task{ID: taskIDCounter, Desc: request.Desc, Deadline: request.Deadline}
	tasks[task.ID] = task
	taskIDCounter++

	response := CreateTaskResponse{ID: task.ID}

	return c.JSON(response)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	var request UpdateTaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}
	idS, _ := url.QueryUnescape(c.Params("id"))
	id, _ := strconv.ParseInt(idS, 10, 64)
	task, ok := tasks[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}
	task.Desc = request.Desc
	task.Deadline = request.Deadline
	tasks[id] = task

	return c.SendString("OK")
}

func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	idS, _ := url.QueryUnescape(c.Params("id"))
	id, _ := strconv.ParseInt(idS, 10, 64)
	task, ok := tasks[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	response := GetTaskResponse{ID: task.ID, Desc: task.Desc, Deadline: task.Deadline}

	return c.JSON(response)
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	idS, _ := url.QueryUnescape(c.Params("id"))
	id, _ := strconv.ParseInt(idS, 10, 64)
	_, ok := tasks[id]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}
	delete(tasks, id)

	return c.SendString("OK")
}
