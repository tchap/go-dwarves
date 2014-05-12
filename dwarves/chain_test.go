package dwarves

func ExampleTaskChain() {
	chain, _ := NewTaskChain(
		newTask("A"),
		newTask("B"),
		newTask("C"),
		newTask("D"),
		newTask("E"),
		newTask("F"),
		newTask("G"),
		newTask("H"),
		newTask("I"),
		newTask("J"),
	)
	supervisor := NewSupervisor(chain)
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
	// F
	// G
	// H
	// I
	// J
}
