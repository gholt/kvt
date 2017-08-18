package kvt_test

import (
	"fmt"
	"time"

	"github.com/gholt/kvt"
)

func Example_overview() {
	// Create store with a few items.
	store1 := kvt.Store{}
	store1.Set("A", "one")
	store1.Set("B", "two")
	store1.Set("C", "three")

	// Create another store with other items, some overlapping, some deleting,
	// some new.
	store2 := kvt.Store{}
	store2.Delete("B")
	store2.Set("C", "four")
	store2.Set("D", "five")
	store2.Set("E", "six")
	store2.Set("F", "seven")

	// For extra effort, update and delete items on the first store.
	store1.Set("E", "eight")
	store1.Delete("F")

	// Here they are prior to merging.
	fmt.Println("Store1:", store1.SimpleString())
	fmt.Println("Store2:", store2.SimpleString())

	// Now we merge store2 into store1.
	store1.Absorb(store2)
	fmt.Println()

	fmt.Println("Store1:", store1.SimpleString())

	// Output:
	// Store1: A=one,B=two,C=three,E=eight,F/deleted
	// Store2: B/deleted,C=four,D=five,E=six,F=seven
	//
	// Store1: A=one,B/deleted,C=four,D=five,E=eight,F/deleted
}

func ExampleStore() {
	// Create store with a few items.
	store1 := kvt.Store{}
	store1.Set("A", "one")
	store1.Set("B", "two")
	store1.Set("C", "three")

	// Create another store with other items, some overlapping, some deleting,
	// some new.
	store2 := kvt.Store{}
	store2.Delete("B")
	store2.Set("C", "four")
	store2.Set("D", "five")
	store2.Set("E", "six")
	store2.Set("F", "seven")

	// For extra effort, update and delete items on the first store.
	store1.Set("E", "eight")
	store1.Delete("F")

	// Here they are prior to merging.
	fmt.Println("Store1:", store1.SimpleString())
	fmt.Println("Store2:", store2.SimpleString())

	// Now we merge store2 into store1.
	store1.Absorb(store2)
	fmt.Println()

	fmt.Println("Store1:", store1.SimpleString())

	// Output:
	// Store1: A=one,B=two,C=three,E=eight,F/deleted
	// Store2: B/deleted,C=four,D=five,E=six,F=seven
	//
	// Store1: A=one,B/deleted,C=four,D=five,E=eight,F/deleted
}

func ExampleStore_Get() {
	store := kvt.Store{}
	store.Set("A", "one")
	store.Delete("B")
	for _, k := range []string{"A", "B", "C"} {
		fmt.Printf("Get(%q): %q\n", k, store.Get(k))
	}

	// Output:
	// Get("A"): "one"
	// Get("B"): ""
	// Get("C"): ""
}

func ExampleStore_Set() {
	store := kvt.Store{}
	store.Set("A", "one")
	fmt.Println(store.SimpleString())

	// Output:
	// A=one
}

func ExampleStore_SetTimestamped() {
	store := kvt.Store{}
	store.SetTimestamped("A", "one", time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC).UnixNano())
	fmt.Println(store)

	// Output:
	// {"A":["one",1483326245000000006]}
}

func ExampleStore_Delete() {
	store := kvt.Store{}
	store.Delete("A")
	fmt.Println(store.SimpleString())

	// Output:
	// A/deleted
}

func ExampleStore_DeleteTimestamped() {
	store := kvt.Store{}
	store.DeleteTimestamped("A", time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC).UnixNano())
	store.DeleteTimestamped("A", 1) // Discarded as old
	store.SetTimestamped("B", "two", 2)
	store.DeleteTimestamped("B", 1) // Discarded as old
	store.SetTimestamped("C", "three", 3)
	store.DeleteTimestamped("C", 4)
	fmt.Println(store)

	// Output:
	// {"A":[null,1483326245000000006],"B":["two",2],"C":[null,4]}
}

func ExampleStore_Purge() {
	store := kvt.Store{}
	now := time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC)
	store.DeleteTimestamped("A", now.Add(-1001*time.Hour).UnixNano())
	store.DeleteTimestamped("B", now.Add(-1000*time.Hour).UnixNano())
	store.DeleteTimestamped("C", now.Add(-999*time.Hour).UnixNano())
	fmt.Println("Before:", store)
	store.Purge(now.Add(-1000 * time.Hour).UnixNano())
	fmt.Println("After:", store)

	// Output:
	// Before: {"A":[null,1479722645000000006],"B":[null,1479726245000000006],"C":[null,1479729845000000006]}
	// After: {"B":[null,1479726245000000006],"C":[null,1479729845000000006]}
}

func ExampleStore_Absorb() {
	// Create store with a few items.
	store1 := kvt.Store{}
	store1.Set("A", "one")
	store1.Set("B", "two")
	store1.Set("C", "three")

	// Create another store with other items, some overlapping, some deleting,
	// some new.
	store2 := kvt.Store{}
	store2.Delete("B")
	store2.Set("C", "four")
	store2.Set("D", "five")
	store2.Set("E", "six")
	store2.Set("F", "seven")

	// For extra effort, update and delete items on the first store.
	store1.Set("E", "eight")
	store1.Delete("F")

	// Here they are prior to merging.
	fmt.Println("Store1:", store1.SimpleString())
	fmt.Println("Store2:", store2.SimpleString())

	// Now we merge store2 into store1.
	store1.Absorb(store2)
	fmt.Println()

	fmt.Println("Store1:", store1.SimpleString())

	// Output:
	// Store1: A=one,B=two,C=three,E=eight,F/deleted
	// Store2: B/deleted,C=four,D=five,E=six,F=seven
	//
	// Store1: A=one,B/deleted,C=four,D=five,E=eight,F/deleted
}

func ExampleStore_Hash() {
	store1 := kvt.Store{}
	now := time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC).UnixNano()
	store1.SetTimestamped("A", "one", now)
	store1.SetTimestamped("B", "two", now)
	store1.SetTimestamped("C", "three", now)
	fmt.Println("store1 has hash", store1.Hash())
	store2 := kvt.Store{}
	store2.SetTimestamped("A", "one", now)
	store2.SetTimestamped("B", "two", now)
	store2.SetTimestamped("C", "three", now)
	fmt.Println("store2 has hash", store2.Hash())
	store2.SetTimestamped("C", "three", now+1)
	fmt.Println("store2 now has hash", store2.Hash())

	// Output:
	// store1 has hash 3d6a6976bcf5dfb9
	// store2 has hash 3d6a6976bcf5dfb9
	// store2 now has hash 3d670f76bcf310f4
}

func ExampleStore_String() {
	store := kvt.Store{}
	now := time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC).UnixNano()
	store.SetTimestamped("A", "one", now)
	store.DeleteTimestamped("B", now)
	fmt.Println(store.String())
	fmt.Println(store)

	// Output:
	// {"A":["one",1483326245000000006],"B":[null,1483326245000000006]}
	// {"A":["one",1483326245000000006],"B":[null,1483326245000000006]}
}

func ExampleStore_SimpleString() {
	store := kvt.Store{}
	now := time.Date(2017, 1, 2, 3, 4, 5, 6, time.UTC).UnixNano()
	store.SetTimestamped("A", "one", now)
	store.DeleteTimestamped("B", now)
	fmt.Println(store.SimpleString())

	// Output:
	// A=one,B/deleted
}

func ExampleValueTimestamp() {
	one := "one"
	vtA := &kvt.ValueTimestamp{Value: &one, Timestamp: 1}
	vtB := &kvt.ValueTimestamp{Value: nil, Timestamp: 2}
	store := kvt.Store{"A": vtA, "B": vtB}
	fmt.Println(store)
	// A bit simpler:
	store = kvt.Store{"A": {&one, 1}, "B": {nil, 2}}
	fmt.Println(store)

	// Output:
	// {"A":["one",1],"B":[null,2]}
	// {"A":["one",1],"B":[null,2]}
}

func ExampleValueTimestamp_MarshalJSON() {
	one := "one"
	b, err := (&kvt.ValueTimestamp{Value: &one, Timestamp: 1}).MarshalJSON()
	fmt.Println(string(b), err)
	b, err = (&kvt.ValueTimestamp{Value: nil, Timestamp: 2}).MarshalJSON()
	fmt.Println(string(b), err)

	// Output:
	// ["one",1] <nil>
	// [null,2] <nil>
}

func ExampleValueTimestamp_UnmarshalJSON() {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`["one",1]`))
	fmt.Println(vt, err)
	vt = &kvt.ValueTimestamp{}
	err = vt.UnmarshalJSON([]byte(`[null,2]`))
	fmt.Println(vt, err)

	// Output:
	// one,1 <nil>
	// nil,2 <nil>
}

func ExampleValueTimestamp_String() {
	one := "one"
	vt1 := &kvt.ValueTimestamp{Value: &one, Timestamp: 1}
	vt2 := &kvt.ValueTimestamp{Value: nil, Timestamp: 2}
	fmt.Println(vt1.String(), vt2.String())
	fmt.Println(vt1, vt2)

	// Output:
	// one,1 nil,2
	// one,1 nil,2
}
