package actions

import (
	"fmt"
	"sync"

	"github.com/tcolar/goed/core"
)

// TODO : This is kind of memory heavy ....
// TODO : group together quick succesive undos (insert a, insert b, insert c) + Flushing

var maxUndos = 500

// viewId keyed map of undo actions
var undos map[int64][]actionTuple = map[int64][]actionTuple{}

// viewId keyed map of redo actions
var redos map[int64][]actionTuple = map[int64][]actionTuple{}

var lock sync.Mutex

// a do/undo combo
type actionTuple struct {
	do   []core.Action
	undo []core.Action
}

// or group by alphanum sequence ??
func Undo(viewId int64) {
	action, err := func() ([]core.Action, error) {
		lock.Lock()
		defer lock.Unlock()
		tuples, found := undos[viewId]
		if !found || len(tuples) == 0 {
			return nil, fmt.Errorf("Nothing to undo.")
		}
		tuple := tuples[len(tuples)-1]
		undos[viewId] = undos[viewId][:len(tuples)-1]
		redos[viewId] = append(redos[viewId], tuple)
		return tuple.undo, nil
	}()
	if err != nil {
		Ar.EdSetStatusErr(err.Error())
		return
	}
	for _, a := range action {
		a.Run()
	}
}

func Redo(viewId int64) {
	action, err := func() ([]core.Action, error) {
		lock.Lock()
		defer lock.Unlock()
		tuples, found := redos[viewId]
		if !found || len(tuples) == 0 {
			return nil, fmt.Errorf("Nothing to redo.")
		}
		tuple := tuples[len(tuples)-1]
		redos[viewId] = redos[viewId][:len(tuples)-1]
		undos[viewId] = append(undos[viewId], tuple)
		return tuple.do, nil
	}()
	if err != nil {
		Ar.EdSetStatusErr(err.Error())
		return
	}
	for _, a := range action {
		a.Run()
	}
}

func UndoAdd(viewId int64, do, undo []core.Action) {
	lock.Lock()
	defer lock.Unlock()
	delete(redos, viewId)
	if len(undos[viewId]) < maxUndos {
		undos[viewId] = append(undos[viewId], actionTuple{do, undo})
	} else {
		copy(undos[viewId], undos[viewId][1:])
		undos[viewId][len(undos[viewId])-1] = actionTuple{do, undo}
	}
}

func UndoClear(viewId int64) {
	lock.Lock()
	defer lock.Unlock()
	delete(undos, viewId)
	delete(redos, viewId)
}

// Dump prints out the undo/redo stack of a view, for debugging
func Dump(viewId int64) {
	fmt.Printf("Undos:\n")
	for _, u := range undos[viewId] {
		fmt.Printf("\t %#v\n", u)
	}
	fmt.Printf("Redos:\n")
	for _, r := range redos[viewId] {
		fmt.Printf("\t %#v\n", r)
	}
}
