package dwarves

import "fmt"

type Resource struct {
	name   string
	locked bool
}

func NewResource(name string) *Resource {
	return &Resource{name: name}
}

func (res *Resource) String() string {
	return res.name
}

func (res *Resource) isAvailable() bool {
	return res.locked
}

func (res *Resource) lock() {
	if res.locked {
		panic(fmt.Errorf("resource %q is already locked", res))
	}
	res.locked = true
}

func (res *Resource) unlock() {
	if !res.locked {
		panic(fmt.Errorf("resource %q is already unlocked", res))
	}
	res.locked = false
}
