package timedMap

import (
	"time"
	"container/list"
)

type timedMapData struct {
	itemData interface{}
	expiryTime time.Time
	element *list.Element
}

type timerListData struct {
	key interface{}
	expiryTime time.Time
}

func newTimedMapData(value interface{}, expiry time.Time, ele *list.Element) *timedMapData  {
	return &timedMapData{itemData:value,expiryTime:expiry, element:ele}
}

func newTimerListData(key interface{}, expiry time.Time) *timerListData  {
	return &timerListData{key,expiry}
}
