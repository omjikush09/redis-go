package data_structure

type ListStoreStructure struct {
	Storage map[string][]string
}

func (store *ListStoreStructure) Add(key string, elements ...string) int {

	list, exist := store.Storage[key]
	length := len(elements) + len(list)
	if !exist {
		store.Storage[key] = make([]string, 0)

	}
	store.Storage[key] = append(store.Storage[key], elements...)
	return length
}

func (store *ListStoreStructure) Remove(key string, elements ...string) {

}
