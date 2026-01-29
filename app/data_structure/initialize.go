package data_structure

var MapStore *MapStoreStructure
var ListStore *ListStoreStructure

func InitializeStore() {
	MapStore = &MapStoreStructure{
		Storage: make(map[string]Data),
	}
	ListStore = &ListStoreStructure{
		Storage: make(map[string][]string),
	}
}
