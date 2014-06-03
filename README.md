# go-dwarves #

[![Build
Status](https://drone.io/github.com/tchap/go-dwarves/status.png)](https://drone.io/github.com/tchap/go-dwarves/latest)
[![Coverage
Status](https://coveralls.io/repos/tchap/go-dwarves/badge.png?branch=master)](https://coveralls.io/r/tchap/go-dwarves?branch=master)

Little dwarves that can be asked politely to perform various tasks.
With maximum concurrency!

## Usage ##

```go
import "github.com/tchap/go-dwarves/dwarves"
```

## Documentation ###

[GoDoc](http://godoc.org/github.com/tchap/go-dwarves/dwarves)

## Examples ##

This section shows how multiple tasks can be connected with `After`. There are
also a few container task wrappers - `TaskChain` and `TaskBag` - that can be
used to group tasks together and make them act as a single logical task.

See the documentation for more complete examples.

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
supervisor.DispatchDwarves()
supervisor.RevertChanges()
// Output:
// task A
// task A reverted
```

## License ##

MIT, see the `LICENSE` file.
