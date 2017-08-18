// Package kvt offers a simple Key|Value|Timestamp store.
//
// You can merge multiple Stores and the resulting Store will have the newest
// Key|Value|Timestamp triplets.
//
// This store might be useful where you have a small set of metadata and want
// to allow updates to it on multiple machines to sync up later.
//
// It offers a Hash to quickly identify if there are any changes, and offers
// JSON marshaling and unmarshaling for persistence.
package kvt

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"time"
)

// Store is a Key|Value|Timestamp simple store.
type Store map[string]*ValueTimestamp

// Get returns the value for a key; if the key does not exist or is marked
// deleted, an empty string is returned.
func (store Store) Get(key string) string {
	valueTimestamp := store[key]
	if valueTimestamp == nil || valueTimestamp.Value == nil {
		return ""
	}
	return *valueTimestamp.Value
}

// Set is equivalent to SetTimestamped(key, value, time.Now().UnixNano()).
func (store Store) Set(key string, value string) {
	store.SetTimestamped(key, value, time.Now().UnixNano())
}

// SetTimestamped stores the value for the key as long as there isn't already a
// value for that key with a newer or equal timestamp.
func (store Store) SetTimestamped(key string, value string, timestamp int64) {
	valueTimestamp := store[key]
	if valueTimestamp == nil {
		store[key] = &ValueTimestamp{&value, timestamp}
	} else if valueTimestamp.Timestamp < timestamp {
		valueTimestamp.Value = &value
		valueTimestamp.Timestamp = timestamp
	}
}

// Delete is equivalent to DeleteTimestamped(key, time.Now().UnixNano()).
func (store Store) Delete(key string) {
	store.DeleteTimestamped(key, time.Now().UnixNano())
}

// DeleteTimestamped records a deletion marker for the key as long as there
// isn't already a value for that key with a newer or equal timestamp.
func (store Store) DeleteTimestamped(key string, timestamp int64) {
	valueTimestamp := store[key]
	if valueTimestamp == nil {
		store[key] = &ValueTimestamp{nil, timestamp}
	} else if valueTimestamp.Timestamp < timestamp {
		valueTimestamp.Value = nil
		valueTimestamp.Timestamp = timestamp
	}
}

// Purge discards any deletion markers older than the cutoff timestamp given.
func (store Store) Purge(cutoff int64) {
	for key, valueTimestamp := range store {
		if valueTimestamp.Value == nil && valueTimestamp.Timestamp < cutoff {
			delete(store, key)
		}
	}
}

// Absorb will update store with any newer items from store2; after Absorb, you
// should no longer use store2.
func (store Store) Absorb(store2 Store) {
	for key, valueTimestamp2 := range store2 {
		valueTimestamp := store[key]
		if valueTimestamp == nil || valueTimestamp.Timestamp < valueTimestamp2.Timestamp {
			store[key] = valueTimestamp2
		}
	}
}

// Hash returns a computed hash string that can be used to quickly detect if
// two stores are in sync.
func (store Store) Hash() string {
	ks := make([]string, 0, len(store))
	for k := range store {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	hasher := fnv.New64a()
	for _, k := range ks {
		hasher.Write([]byte(fmt.Sprintf("%s\n%d\n", k, store[k].Timestamp)))
	}
	return fmt.Sprintf("%016x", hasher.Sum64())
}

// String returns the JSON encoded string representation of the store contents.
func (store Store) String() string {
	b, err := json.Marshal(store)
	if err != nil {
		return fmt.Sprintf("error encoding %#v: %#v", store, err)
	}
	return string(b)
}

// SimpleString returns a simple key=value[,key=value] string form of the store
// contents; useful in tests when you want to omit the timestamps.
func (store Store) SimpleString() string {
	ks := make([]string, 0, len(store))
	for k := range store {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	var msg string
	for i, k := range ks {
		if store[k].Value == nil {
			if i == 0 {
				msg += fmt.Sprintf("%s/deleted", k)
			} else {
				msg += fmt.Sprintf(",%s/deleted", k)
			}
		} else {
			if i == 0 {
				msg += fmt.Sprintf("%s=%s", k, *store[k].Value)
			} else {
				msg += fmt.Sprintf(",%s=%s", k, *store[k].Value)
			}
		}
	}
	return msg
}

// ValueTimestamp is the Value|Timestamp pair stored for each Key. If Value is
// nil, it indicates a deletion marker. These deletion markers are usually
// purged after some time using Store.Purge.
type ValueTimestamp struct {
	Value     *string
	Timestamp int64
}

// MarshalJSON returns the JSON encoded version of valueTimestamp or an error.
func (valueTimestamp *ValueTimestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{valueTimestamp.Value, valueTimestamp.Timestamp})
}

// MarshalJSON loads valueTimestamp with data from the JSON encoded b or
// returns an error.
func (valueTimestamp *ValueTimestamp) UnmarshalJSON(b []byte) error {
	jsonValueTimestamp := make([]interface{}, 0, 2)
	if err := json.Unmarshal(b, &jsonValueTimestamp); err != nil {
		return err
	}
	if len(jsonValueTimestamp) != 2 {
		return fmt.Errorf("expected [value,timestamp] from: %s", b)
	}
	if jsonValueTimestamp[0] == nil {
		valueTimestamp.Value = nil
	} else if value, ok := jsonValueTimestamp[0].(string); !ok {
		return fmt.Errorf("invalid value from: %s", b)
	} else {
		valueTimestamp.Value = &value
	}
	if t, ok := jsonValueTimestamp[1].(float64); !ok || float64(int64(t)) != t {
		return fmt.Errorf("invalid timestamp from: %s", b)
	} else {
		valueTimestamp.Timestamp = int64(t)
	}
	return nil
}

// String returns a quick string representation of valueTimestamp.
func (valueTimestamp *ValueTimestamp) String() string {
	if valueTimestamp.Value == nil {
		return fmt.Sprintf("nil,%d", valueTimestamp.Timestamp)
	}
	return fmt.Sprintf("%s,%d", *valueTimestamp.Value, valueTimestamp.Timestamp)
}
