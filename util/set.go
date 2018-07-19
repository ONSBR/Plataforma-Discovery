package util

import (
	"sync"
)

type StringSet struct {
	hashmap map[string]interface{}
	mux     sync.Mutex
}

//NewStringSet creates a new Set
func NewStringSet() *StringSet {
	set := new(StringSet)
	set.hashmap = make(map[string]interface{})
	return set
}

//Add new string to the Set
func (set *StringSet) Add(value string, store ...interface{}) {
	set.mux.Lock()
	defer set.mux.Unlock()
	if !set.Exist(value) {
		if len(store) == 0 {
			set.hashmap[value] = true
		} else {
			set.hashmap[value] = store[0]
		}
	}
}

func (set *StringSet) Len() int {
	return len(set.hashmap)
}

//Exist check if some value exist in Set
func (set *StringSet) Exist(value string) bool {
	_, exist := set.hashmap[value]
	return exist
}

func (set *StringSet) Get(value string) interface{} {
	obj, exist := set.hashmap[value]
	if exist {
		return obj
	}
	return nil
}

//List all values in set
func (set *StringSet) List() []string {
	set.mux.Lock()
	defer set.mux.Unlock()
	keys := make([]string, len(set.hashmap))
	i := 0
	for k := range set.hashmap {
		keys[i] = k
		i++
	}
	return keys
}

//ListValues return all vallues
func (set *StringSet) ListValues() interface{} {
	set.mux.Lock()
	defer set.mux.Unlock()
	keys := make([]interface{}, len(set.hashmap))
	i := 0
	for k := range set.hashmap {
		keys[i] = k
		i++
	}
	return keys
}
