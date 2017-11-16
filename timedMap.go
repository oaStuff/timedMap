package timedMap

import (
	"time"
	"sync"
	"container/list"
)

type TimedMap struct {
	data sync.Map
	timerList *list.List
	storageTime time.Duration
}

var emptyCond *sync.Cond

func init()  {
	emptyCond = sync.NewCond(&sync.Mutex{})
}

// Creates and new TimedMap with the desired storage time.
// The storageTime determines how long a value will be valid in the TimedMap
func NewTimeMap(storageTime time.Duration)  *TimedMap {
	l := list.New()
	l.Init()
	tm := &TimedMap{data:sync.Map{}, timerList: l, storageTime:storageTime}
	go cleaner(tm)

	return tm
}


// Adds a key value pair to the timed map.
func (tm *TimedMap) Add(key, value interface{})  {
	expiry := time.Now().Add(tm.storageTime)
	ele := tm.timerList.PushBack(newTimerListData(key,expiry))
	tm.data.Store(key,newTimedMapData(value,expiry,ele))
	emptyCond.Signal()
}


// Removes the value with the associated key.
func (tm *TimedMap) Remove(key interface{})  {
	if v, ok := tm.data.Load(key); ok {
		tm.data.Delete(key)
		tm.timerList.Remove(v.(*timedMapData).element)
	}
}

// Retrieves the value with the associated key.
// Returns nil if key those not exist in the map
func (tm *TimedMap) Get(key interface{}) interface{}  {
	v, ok := tm.data.Load(key)
	if !ok {
		return nil
	}

	value := v.(*timedMapData)
	if time.Now().After(value.expiryTime) {
		tm.data.Delete(key)
		return nil
	}

	return value.itemData
}

// Checks if the given key exist in the map.
func (tm *TimedMap) Contains(key interface{}) bool {
	if value, ok := tm.data.Load(key); ok {
		if time.Now().After(value.(*timedMapData).expiryTime) {
			tm.data.Delete(key)
			return false
		}
		return true
	}

	return false
}


// Cleaner function that runs in its own goroutine.
// Removes all expired keys with associated values from the map
func cleaner(tm *TimedMap)  {
	for {
		el := tm.timerList.Front()
		for  ; nil == el ; el = tm.timerList.Front() {
			emptyCond.L.Lock()
			emptyCond.Wait()
			emptyCond.L.Unlock()
		}

		time.Sleep(el.Value.(*timerListData).expiryTime.Sub(time.Now()))
		tm.data.Delete(el.Value.(*timerListData).key)
		tm.timerList.Remove(el)
	}
}


