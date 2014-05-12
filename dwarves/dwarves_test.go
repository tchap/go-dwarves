package dwarves

import "fmt"

func newTask(name string) *Task {
	return NewTask(func(interruptCh <-chan struct{}) error {
		fmt.Println(name)
		return nil
	}).RevertChangesWith(func() error {
		fmt.Println("-" + name)
		return nil
	})
}

func newFailingTask(name string) *Task {
	return NewTask(func(interruptCh <-chan struct{}) error {
		fmt.Println(name)
		return fmt.Errorf("task %v failed", name)
	}).RevertChangesWith(func() error {
		fmt.Println("-" + name)
		return fmt.Errorf("task %v rollback failed", name)
	})
}

func newEmptyTask() *Task {
	return NewTask(func(interruptCh <-chan struct{}) error {
		return nil
	})
}
