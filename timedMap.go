package timedMap

import (
	"time"
	"sync"
	"container/list"
)

type EvictionNotify func (key, value interface{})

type TimedMap struct {
	data 			sync.Map
	timerList 		*list.List
	storageTime 	time.Duration
	emptyCond 		*sync.Cond
	callback 		EvictionNotify
}


// Creates and new TimedMap with the desired storage time.
// The storageTime determines how long a value will be valid in the TimedMap
func NewTimeMap(storageTime time.Duration, callback EvictionNotify)  *TimedMap {
	l := list.New()
	l.Init()
	tm := &TimedMap{data:sync.Map{}, timerList: l, storageTime:storageTime, emptyCond:sync.NewCond(&sync.Mutex{}),
					callback:callback}
	go cleaner(tm)

	return tm
}


// Adds a key value pair to the timed map.
func (tm *TimedMap) Put(key, value interface{})  {
	expiry := time.Now().Add(tm.storageTime)
	ele := tm.timerList.PushBack(newTimerListData(key,expiry))
	tm.data.Store(key,newTimedMapData(value,expiry,ele))
	tm.emptyCond.Signal()
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
			tm.emptyCond.L.Lock()
			tm.emptyCond.Wait()
			tm.emptyCond.L.Unlock()
		}

		expiry  := el.Value.(*timerListData).expiryTime
		if !time.Now().After(expiry) {
			time.Sleep(expiry.Sub(time.Now()))
		}
		key := el.Value.(*timerListData).key
		val,_ := tm.data.Load(key)
		tm.data.Delete(key)
		tm.timerList.Remove(el)
		if tm.callback != nil {
			go tm.callback(key, val.(*timedMapData).itemData)
		}
	}
}


