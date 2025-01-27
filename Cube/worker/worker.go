package worker

import (
	"cube/stats"
	"cube/store"
	"cube/task"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-collections/collections/queue"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        store.Store
	Stats     *stats.Stats
	TaskCount int
}

func New() *Worker {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	w := Worker{
		Name:  hostname,
		Queue: *queue.New(),
	}

	var s store.Store

	filename := fmt.Sprintf("%s_tasks.db", hostname)
	s, err = store.NewTaskStore(filename, 0600, "tasks")

	if err != nil {
		log.Fatalf("Failed to create task store: %v", err)
	}
	w.Db = s
	return &w

}

func (w *Worker) CollectStats() {

	for {

		log.Println("Collecting Stats")
		w.Stats = stats.GetStats()
		w.TaskCount = w.Stats.TaskCount
		time.Sleep(5 * time.Second)

	}

}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) RunTasks() {

	for {

		if w.Queue.Len() != 0 {
			result := w.runTask()

			if result.Error != nil {
				log.Printf("Error running task: %v \n", result.Error)

			}

		} else {
			log.Printf("No tasks to process currently. \n")
		}
		log.Println("Sleeping for 10 seconds")

		time.Sleep(5 * time.Second)
	}
}

func (w *Worker) runTask() task.Result {

	t := w.Queue.Dequeue()

	if t == nil {

		log.Println("[worker] No tasks in the queue")

		return task.Result{Error: nil}
	}
	taskQueued := t.(task.Task)

	fmt.Printf("[worker] Found task in queue :%v:\n", taskQueued)

	err := w.Db.Put(taskQueued.ID.String(), &taskQueued)

	if err != nil {

		msg := fmt.Errorf("error strong task %s: %v", taskQueued.ID.String(), err)

		log.Println(msg)

		return task.Result{Error: msg}
	}

	result, err := w.Db.Get(taskQueued.ID.String())

	if err != nil {
		msg := fmt.Errorf("eror getting task %s from database :%v", taskQueued.ID.String(), err)
		log.Println(msg)
		return task.Result{Error: msg}
	}

	taskPersisted := *result.(*task.Task)

	if taskPersisted.State == task.Completed {

		return w.StopTask(taskPersisted)
	}

}
