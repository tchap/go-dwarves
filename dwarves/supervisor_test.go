package dwarves

import "fmt"

func ExampleSupervisor() {
	taskA := newTask("A")
	taskB := newTask("B").After(taskA)
	taskC := newTask("C").After(taskB)
	taskD := newTask("D").After(taskC)
	newFailingTask("E").After(taskD)

	supervisor := NewSupervisor(taskA)

	monitorCh := make(chan *TaskFinishedEvent)
	if err := supervisor.DispatchDwarves(monitorCh); err != nil {
		panic(err)
	}
	for {
		event, ok := <-monitorCh
		if !ok {
			return
		}
		if event.Error != nil {
			fmt.Println(event.Error)
			break
		}
	}

	monitorCh = make(chan *TaskFinishedEvent)
	if err := supervisor.RevertChanges(monitorCh); err != nil {
		panic(err)
	}
	var events []*TaskFinishedEvent
	for {
		event, ok := <-monitorCh
		if !ok {
			break
		}
		events = append(events, event)
	}
	for _, event := range events {
		fmt.Println(event.Error)
	}
	// Output:
	// A
	// B
	// C
	// D
	// E
	// task E failed
	// -E
	// -D
	// -C
	// -B
	// -A
	// task E rollback failed
}
