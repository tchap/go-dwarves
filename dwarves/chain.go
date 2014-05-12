package dwarves

type TaskChain struct {
	first *Task
	last  *Task
}

func NewTaskChain(tasks ...*Task) (*TaskChain, error) {
	if len(tasks) == 0 {
		return nil, ErrNoTasksSpecified
	}
	chain := &TaskChain{
		first: tasks[0],
		last:  tasks[0],
	}
	chain.Append(tasks[1:]...)
	return chain, nil
}

func (chain *TaskChain) Append(tasks ...*Task) error {
	if len(tasks) == 0 {
		return ErrNoTasksSpecified
	}
	for _, t := range tasks {
		t.After(chain.last)
		chain.last = t
	}
	return nil
}

func (chain *TaskChain) Roots() []*Task {
	return []*Task{chain.first}
}

func (chain *TaskChain) Leaves() []*Task {
	return []*Task{chain.last}
}
