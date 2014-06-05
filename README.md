# go-dwarves #

[![Build
Status](https://drone.io/github.com/tchap/go-dwarves/status.png)](https://drone.io/github.com/tchap/go-dwarves/latest)
[![Coverage
Status](https://coveralls.io/repos/tchap/go-dwarves/badge.png?branch=master)](https://coveralls.io/r/tchap/go-dwarves?branch=master)

Little dwarves that can be asked politely to perform various tasks.
With maximum concurrency!

## Motivation ##

I started writing a CLI utility that needed to perform multiple independent
tasks, and I wanted to speed up the utility by running the tasks concurrently
where possible.

After writing some code and synchronizing the tasks manually by using
channels where necessary, I thought that it would be handy to come up with a
general-purpose library that would do the synchronization for me, according to
the conditions I specify. And this is the result.

## Usage ##

```go
import "github.com/tchap/go-dwarves/dwarves"
```

go-dwarves is a tiny task scheduler that is to be used in the following way:

1. Initialise desired `Task` objects.
2. Optionally specify certain (partial) ordering using `Task.After`. Calling
   `taskB.After(taskA)` means that `taskB` will run only after `taskA` has
   finished executing.
3. Optionally specify what tasks use what resources by calling `Task.Uses`.
   No two tasks using the same resource will be run at the same time.
4. Initialise a `Supervisor` with the list of tasks to be run. The downstream
   tasks will be triggered automatically, i.e. when `taskB.After(taskA)` is
   called, it is enough to list `taskA` since `taskB` will be added
   automatically.
5. Start the whole thing by calling `Supervisor.DispatchDwarves`.

Once the supervisor is started, it can be interrupted by calling
`Supervisor.WithdrawDwarves`. This means that the supervisor will wait until the
currently running tasks are finished, but it will not start any new tasks. A
channel is passed to the task functions that is closed when `WithdrawDwarves` is
called. The tasks can try to exit as soon as possible.

There is also one extra feature that is orthogonal to the rest. The supervisor
can be asked to perform a rollback, which is supposed to revert all
the changes that happened during the task execution. The rollback function can
be specified for every task by using `Task.RevertChangesWith`. The supervisor
then simply calls these function in the inverted order compared to the order the
tasks were started.

### Error Handling ###

The errors from the task functions and the task rollback functions are returned
over a channel that can be passed to `Supervisor.DispatchDwarves` and
`Supervisor.RevertChanges`. For this reason these methods are a bit against the
good practices since they do not block until the tasks are finished. Use
`Supervisor.DwarvesFinished`, `Supervisor.WaitFinished`,
`Supervisor.ChangesReverted` or `Supervisor.WaitReverted`.

## Documentation ###

[GoDoc](http://godoc.org/github.com/tchap/go-dwarves/dwarves)

## Examples ##

This section shows how multiple tasks can be connected with `After`. There are
also a few container task wrappers - `TaskChain` and `TaskBag` - that can be
used to group tasks together and make them act as a single logical task.

See the documentation for more complete
[examples](http://godoc.org/github.com/tchap/go-dwarves/dwarves#pkg-examples).

### Task ###

#### After ####

`Task.After` can be used to enforce certain task order of execution.

```go
taskA := dwarves.NewTask(...)
taskB := dwarves.NewTask(...).After(taskA)
```

#### Uses ####

`Task.Uses` can be used to tell the supervisor that the task requires certain
resource as created with `NewResource`. The supervisor will ensure that no two
tasks requiring the same resource will run at the same time.

```go
store := dwarves.NewResource("store")
taskA := dwarves.NewTask(...).Uses(store)
```

### TaskBag ###

`TaskBag` can be used to group tasks together so that they look like a single
logical task, even though no ordering is specified.

The tasks specified to run after the bag will be started once all the tasks in
the bag are finished executing.

```go
taskA := dwarves.NewTask(...)
taskB, _ := dwarves.NewTaskBag(
    subtaskB1,
	subtaskB2,
	subtaskB3,
)
taskB.After(taskA)
taskC.After(taskB)
```

### TaskChain ###

`TaskChain` can be used to group tasks together so that they look like a single
logical task, and it also ensures that the tasks forming the chain are executed
one by one. `TaskChain.Append` can be used to extend the chain, although it has
its caveats.

```go
taskA := dwarves.NewTask(...)
taskB, _ := dwarves.NewTaskChain(
    subtaskA,
    subtaskB,
    subtaskC,
)
taskB.After(taskA)
taskC.After(taskB)
```

### Supervisor ###

#### RevertChanges ####

`Supervisor.RevertChanges` can be used to revert changes implied by the tasks
that has already run. The revert function can be specified for every task with
`Task.RevertChangesWith`.

```go
task := dwarves.NewTask(func(interruptCh <-chan struct{}) error {
    fmt.Println("task A")
}).RevertChangesWith(func() error {
    fmt.Println("task A reverted")
})

supervisor := dwarves.NewSupervisor(task)
supervisor.DispatchDwarves(nil)
supervisor.RevertChanges(nil)
supervusor.WaitChangesReverted()
// Output:
// task A
// task A reverted
```

## License ##

MIT, see the `LICENSE` file.
