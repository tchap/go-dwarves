package dwarves

type TaskBag struct {
	roots []*Task
}

func NewTaskBag(tasks ...*Task) (*TaskBag, error) {
	if len(tasks) == 0 {
		return nil, ErrNoTasksSpecified
	}

	return &TaskBag{tasks}, nil
}

func (bag *TaskBag) Roots() []*Task {
	return bag.roots
}

func (bag *TaskBag) Leaves() []*Task {
	var leaves []*Task
	for _, task := range bag.roots {
		task.visit(func(t *Task) {
			if len(t.downstreamTasks) == 0 {
				leaves = append(leaves, t)
			}
		})
	}
	return leaves
}
