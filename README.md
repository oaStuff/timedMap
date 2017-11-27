# timedMap

timedMap is simple go 1.9 sync.map with timing added. This implies that timedMap just
stores items in a map and expires them automatically based on the time duration set.
## usage

Import the package:

```go
import (
	"github.com/oaStuff/timedMap"
)

```

```bash
go get "github.com/oaStuff/timedMap"
```


## example

```go

        //create the map
	    tm := timedMap.NewTimeMap(time.Second * 3, nil) // we want items to live in the map for only 3 seconds
    	tm.Add("one","1")
    	tm.Add("two","2")
    	tm.Add("three", "3")

    	val := tm.Get("one")
    	fmt.Println("value is " + val)

```

## another example

```go
    type user struct {
        name string
        age int
    }
    
    tm := timedMap.NewTimeMap(time.Second * 3, nil)
	tm.Add("john", &user{"john mark", 30})
	tm.Add("mary", &user{"mary jane", 26})
	tm.Add("paul", &user{"paul frank", 19})

	val = tm.Get("mary")
	if val != nil {
	    fmt.Println(val.(*user).name)
	    fmt.Println(val.(*user).age)
	}
```

The library also supports eviction callback when an item in the map expires.
```go

    type user struct {
        name string
        age int
    }
    tm := timedMap.NewTimeMap(time.Second * 3, func(key, value interface{}) {
        fmt.Println("expried callback:")
        fmt.Printf("%+v\n", key)
        fmt.Printf("%+v\n", value)
    })
    
	tm.Add("john", &user{"john mark", 30})
	tm.Add("mary", &user{"mary jane", 26})
	tm.Add("paul", &user{"paul frank", 19})

    time.Sleep(time.Second * 5)
	val = tm.Get("mary")
	if val != nil {
	    fmt.Println(val.(*user).name)
	    fmt.Println(val.(*user).age)
	}
```

## license
MIT (see [LICENSE](https://github.com/orcaman/concurrent-map/blob/master/LICENSE) file)