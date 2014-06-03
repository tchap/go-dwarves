package dwarves

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
