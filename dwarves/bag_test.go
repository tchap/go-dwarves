package dwarves

import "testing"

func TestTaskBag(t *testing.T) {
	tasks := []*Task{
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
		newEmptyTask(),
	}

	bag, err := NewTaskBag(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	monitorCh := make(chan *TaskFinishedEvent)
	if err := NewSupervisor(bag).DispatchDwarves(monitorCh); err != nil {
		t.Fatal(err)
	}

	var counter int
	for {
		event, ok := <-monitorCh
		if !ok {
			break
		}
		if event.Error != nil {
			t.Error(err)
		}
		counter++
	}
	if counter != len(tasks) {
		t.Errorf("Expected %v tasks to finish, got %v", len(tasks), counter)
	}
}
