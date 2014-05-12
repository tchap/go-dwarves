package dwarves

import "fmt"

func ExampleTask_After() {
	taskA := newTask("A")
	taskB := newTask("B").After(taskA)
	taskC := newTask("C").After(taskB)
	taskD := newTask("D").After(taskC)
	newTask("E").After(taskD)

	supervisor := NewSupervisor(taskA)
	if err := supervisor.DispatchDwarves(nil); err != nil {
		panic(err)
	}
	supervisor.WaitFinished()
	// Output:
	// A
	// B
	// C
	// D
	// E
}

func ExampleTask_RevertChanges() {
	taskA := newTask("A")
	taskB := newTask("B").After(taskA)
	taskC := newTask("C").After(taskB)
	taskD := newTask("D").After(taskC)
	newTask("E").After(taskD)

	supervisor := NewSupervisor(taskA)
	if err := supervisor.DispatchDwarves(nil); err != nil {
		panic(err)
	}
	if err := supervisor.RevertChanges(nil); err != nil {
		panic(err)
	}
	supervisor.WaitReverted()
	// Output:
	// A
	// B
	// C
	// D
	// E
	// -E
	// -D
	// -C
	// -B
	// -A
}

func ExampleTask_Uses() {
	repository := NewResource("repository")

	// The following tasks will be run one at a time, in random order,
	// since every task requires the repository resource to run.
	taskA := newTask("A").Uses(repository)
	taskB := newTask("B").Uses(repository)
	taskC := newTask("C").Uses(repository)

	supervisor := NewSupervisor(taskA, taskB, taskC)
	if err := supervisor.DispatchDwarves(nil); err != nil {
		panic(err)
	}
	supervisor.WaitFinished()
}

func ExampleTask_MonitorChannel() {
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
