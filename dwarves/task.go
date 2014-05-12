package dwarves

type (
	TaskFunc   func(interruptCh <-chan struct{}) error
	RevertFunc func() error
)

type TaskForest interface {
	Roots() []*Task
	Leaves() []*Task
}

type Task struct {
	upstreamTasks     map[*Task]struct{}
	downstreamTasks   map[*Task]struct{}
	requiredResources map[*Resource]struct{}
	taskFunc          TaskFunc
	revertFunc        RevertFunc
	interruptCh       chan struct{}
}

func NewTask(taskFunc TaskFunc) *Task {
	return &Task{
		upstreamTasks:     make(map[*Task]struct{}),
		downstreamTasks:   make(map[*Task]struct{}),
		requiredResources: make(map[*Resource]struct{}),
		taskFunc:          taskFunc,
		interruptCh:       make(chan struct{}),
	}
}

func (task *Task) After(forests ...TaskForest) (this *Task) {
	for _, forest := range forests {
		for _, t := range forest.Leaves() {
			task.upstreamTasks[t] = struct{}{}
			t.downstreamTasks[task] = struct{}{}
		}
	}
	return task
}

func (task *Task) Uses(resources ...*Resource) (this *Task) {
	for _, r := range resources {
		task.requiredResources[r] = struct{}{}
	}
	return task
}

func (task *Task) RevertChangesWith(revertFunc RevertFunc) (this *Task) {
	task.revertFunc = revertFunc
	return task
}

func (task *Task) Roots() []*Task {
	return []*Task{task}
}

func (task *Task) Leaves() []*Task {
	return []*Task{task}
}

func (task *Task) run() error {
	return task.taskFunc(task.interruptCh)
}

func (task *Task) interrupt() {
	close(task.interruptCh)
}

func (task *Task) isRevertible() bool {
	return task.revertFunc != nil
}

func (task *Task) revert() error {
	return task.revertFunc()
}

func (task *Task) visit(visitor func(*Task)) {
	visitor(task)
	for t := range task.downstreamTasks {
		t.visit(visitor)
	}
}
