package test

import (
	"testing"
	"github.com/oaStuff/timedMap"
	"time"
	"fmt"
)

type user struct {
	name string
	age int
}

func TestInsert(t *testing.T)  {

	tm := timedMap.NewTimeMap(time.Second * 3, nil)
	tm.Put("one","1")
	tm.Put("two","2")
	tm.Put("three", "3")

	val := tm.Get("one")
	if val != "1" {
		t.Error("value ought to be 1")
	}

	val = tm.Get("two")
	if "2" != val {
		t.Error("value ought to be 2")
	}

	val = tm.Get("three")
	if "3" != val {
		t.Error("value ought to be 3")
	}

	val = tm.Get("four")
	if nil != val {
		t.Error("value ought to be nil")
	}
}


func TestStruct(t *testing.T)  {

	tm := timedMap.NewTimeMap(time.Second * 3, nil)
	tm.Put("john", &user{"john mark", 30})
	tm.Put("mary", &user{"mary jane", 26})
	tm.Put("paul", &user{"paul frank", 19})

	val := tm.Get("unkown")
	if nil != val {
		t.Error("value should be nil since the key those not exist")
	}

	val = tm.Get("mary")
	if val.(*user).name != "mary jane" {
		t.Error("value returned should be what was inserted")
	}
}

func TestRemove(t *testing.T)  {

	tm := timedMap.NewTimeMap(time.Minute * 3, nil)
	tm.Put("1", "one")
	tm.Put("2", "two")
	tm.Put("3","three")

	val := tm.Get("2")
	if "two" != val {
		t.Error("two ought to be returned")
	}

	tm.Remove("2")
	val = tm.Get("2")
	if nil != val {
		t.Error("the value ought to be nil since the data was removed")
	}
}

func TestContains(t *testing.T)  {

	tm := timedMap.NewTimeMap(time.Minute * 3, nil)
	tm.Put("1", "one")
	tm.Put("2", "two")
	tm.Put("3","three")

	if !tm.Contains("1") {
		t.Error("the map is suppose to contain one")
	}

	if tm.Contains("unkown") {
		t.Error("the map is not suppose to have a value for an unkown key")
	}
}


func TestExpiry(t *testing.T)  {

	tm := timedMap.NewTimeMap(time.Second * 3, func(key, value interface{}) {
		fmt.Println("expried callback:")
		fmt.Printf("%+v\n", key)
		fmt.Printf("%+v\n", value)
	})

	tm.Put("one","1")
	tm.Put("two","2")

	time.Sleep(time.Second * 2)
	tm.Put("three", "3")

	time.Sleep(time.Second * 2)

	val := tm.Get("three")
	if "3" != val {
		t.Error("value ought to be 3")
	}

	val = tm.Get("one")
	if nil != val {
		t.Error("value ought to be nil because it should have expired")
	}

}



