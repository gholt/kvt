# KVT - A Simple Key|Value|Timestamp Store

Package kvt offers a simple Key|Value|Timestamp store.

You can merge multiple Stores and the resulting Store will have the newest
Key|Value|Timestamp triplets.

This store might be useful where you have a small set of metadata and want
to allow updates to it on multiple machines to sync up later.

It offers a Hash to quickly identify if there are any changes, and offers
JSON marshaling and unmarshaling for persistence.

[API Documentation](http://godoc.org/github.com/gholt/kvt)  

> Copyright See AUTHORS. All rights reserved.  
> Use of this source code is governed by a BSD-style  
> license that can be found in the LICENSE file.

# Example Usage

```go
package main

import (
    "fmt"

    "github.com/gholt/kvt"
)

func main() {
    // Create store with a few items.
    store1 := kvt.Store{}
    store1.Set("A", "one")
    store1.Set("B", "two")
    store1.Set("C", "three")

    // Create another store with other items, some overlapping, some
    // deleting, some new.
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
```
