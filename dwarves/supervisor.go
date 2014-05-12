package dwarves

import (
	"errors"
	//"sync/atomic"

	"launchpad.net/tomb"
)

var (
	ErrNotStarted      = errors.New("the dwarf were not started yet")
	ErrNotRunning      = errors.New("the dwarfs are not running")
	ErrAlreadyStarted  = errors.New("the dwarfs were already started")
	ErrAlreadyReverted = errors.New("the dwarfs were already reverted")
)

const (
	ctxStateInit = iota
	ctxStateRunning
	ctxStateStopped
	ctxStateReverted
)

const (
	taskUnknown = iota
	taskStarted
	taskInterrupted
	taskFinished
	taskReverted
)

var ErrNoTasksSpecified = errors.New("no tasks were specified")

type TaskFinishedEvent struct {
	Task  *Task
	Error error
}

type Supervisor struct {
	state             int32
	tomb              tomb.Tomb
	initialTasks      []*Task
	pendingTasks      map[*Task]struct{}
	startedTasks      []*Task
	taskStates        map[*Task]int
	numWorkingDwarves int
	returnedCh        chan *TaskFinishedEvent
	revertedCh        chan struct{}
}

func NewSupervisor(forests ...TaskForest) *Supervisor {
	roots := make([]*Task, 0, len(forests))
	for _, forest := range forests {
		roots = append(roots, forest.Roots()...)
	}
	return &Supervisor{
		initialTasks: roots,
		pendingTasks: make(map[*Task]struct{}),
		taskStates:   make(map[*Task]int),
		returnedCh:   make(chan *TaskFinishedEvent),
		revertedCh:   make(chan struct{}),
	}
}

// Public methods --------------------------------------------------------------

func (s *Supervisor) DispatchDwarves(monitorCh chan<- *TaskFinishedEvent) error {
	go s.work(monitorCh)
	return nil
}

func (s *Supervisor) WithdrawDwarves() error {
	s.tomb.Kill(nil)
	return nil
}

func (s *Supervisor) DwarvesFinished() <-chan struct{} {
	return s.tomb.Dead()
}

func (s *Supervisor) WaitFinished() {
	<-s.DwarvesFinished()
}

func (s *Supervisor) RevertChanges(monitorCh chan<- *TaskFinishedEvent) error {
	s.WaitFinished()
	go func() {
		defer func() {
			close(s.revertedCh)
			if monitorCh != nil {
				close(monitorCh)
			}
		}()
		for i := len(s.startedTasks) - 1; i >= 0; i-- {
			task := s.startedTasks[i]
			if s.taskStates[task] < taskFinished {
				continue
			}
			if task.isRevertible() {
				if err := task.revert(); err != nil && monitorCh != nil {
					monitorCh <- &TaskFinishedEvent{task, err}
				}
			}
		}
	}()
	return nil
}

func (s *Supervisor) ChangesReverted() <-chan struct{} {
	return s.revertedCh
}

func (s *Supervisor) WaitReverted() {
	<-s.ChangesReverted()
}

// Private methods -------------------------------------------------------------

func (s *Supervisor) work(monitorCh chan<- *TaskFinishedEvent) {
	for _, task := range s.initialTasks {
		s.tryDispatchDwarf(task)
	}

	dyingCh := s.tomb.Dying()
	for {
		select {
		case finishedEvent := <-s.returnedCh:
			var (
				task = finishedEvent.Task
				err  = finishedEvent.Error
			)
			if monitorCh != nil {
				monitorCh <- finishedEvent
			}

			// Unlock all relevant resources.
			for r := range task.requiredResources {
				r.unlock()
			}

			// Mark the task as finished.
			s.numWorkingDwarves--
			s.taskStates[task] = taskFinished

			// Try to dispatch more dwarves,
			// unless the supervisor is dying.
			if dyingCh != nil {
				for t := range s.pendingTasks {
					s.tryDispatchDwarf(t)
				}
				if err == nil {
					for t := range task.downstreamTasks {
						s.tryDispatchDwarf(t)
					}
				}
			}

			// Exit if there are no more dwarves that are working.
			if s.numWorkingDwarves == 0 {
				s.tomb.Done()
				if monitorCh != nil {
					close(monitorCh)
				}
				return
			}

		case <-dyingCh:
			for task := range s.taskStates {
				task.interrupt()
			}
			dyingCh = nil
		}
	}
}

func (s *Supervisor) tryDispatchDwarf(task *Task) {
	// Make sure the task has not been started yet.
	if s.taskStates[task] != taskUnknown {
		return
	}

	// Make sure all the upstream tasks have been finished.
	for t := range task.upstreamTasks {
		if s.taskStates[t] != taskFinished {
			s.pendingTasks[t] = struct{}{}
			return
		}
	}

	// Make sure all required resources are available.
	for r := range task.requiredResources {
		if !r.isAvailable() {
			s.pendingTasks[task] = struct{}{}
			return
		}
	}
	// Lock all required resources.
	for r := range task.requiredResources {
		r.lock()
	}

	// Start the task.
	s.numWorkingDwarves++
	delete(s.pendingTasks, task)
	s.taskStates[task] = taskStarted
	s.startedTasks = append(s.startedTasks, task)
	go func() {
		var err error
		defer func() {
			s.returnedCh <- &TaskFinishedEvent{task, err}
		}()
		err = task.run()
	}()
}
