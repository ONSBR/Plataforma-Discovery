package util

import (
	"sync"
)

type StringSet struct {
	hashmap map[string]bool
	mux     sync.Mutex
}

//NewStringSet creates a new Set
func NewStringSet() *StringSet {
	set := new(StringSet)
	set.hashmap = make(map[string]bool)
	return set
}

//Add new string to the Set
func (set *StringSet) Add(value string) {
	set.mux.Lock()
	defer set.mux.Unlock()
	set.hashmap[value] = true
}

func (set *StringSet) Len() int {
	return len(set.hashmap)
}

//Exist check if some value exist in Set
func (set *StringSet) Exist(value string) bool {
	_, exist := set.hashmap[value]
	return exist
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
