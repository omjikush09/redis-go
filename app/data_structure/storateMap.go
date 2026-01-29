package data_structure

import (
	"fmt"
	"strconv"
	"time"
)

type Data struct {
	value  string
	tll    bool
	expiry time.Time
}
type MapStoreStructure struct {
	Storage map[string]Data
}

func (store *MapStoreStructure) Add(key string, value string, expiry time.Duration, tll bool) {
	store.Storage[key] = Data{value: value, expiry: time.Now().Add(expiry), tll: tll}
}

func (store *MapStoreStructure) Get(key string) (string, bool) {
	output, exist := store.Storage[key]

	expired := time.Now().After(output.expiry)

	if (expired || !exist) && output.tll {
		if exist {
			delete(store.Storage, key)
		}
		return "", false
	}
	return output.value, exist
}

func (store *MapStoreStructure) Increment(key string) (string, error) {
	output, exist := store.Storage[key]
	if !exist {
		store.Storage[key] = Data{value: "1", expiry: output.expiry, tll: output.tll}
		return "1", nil
	}

	intValue, err := strconv.Atoi(output.value)
	if err != nil {
		return "", fmt.Errorf("ERR value is not an integer or out of range")
	}
	resultValue := strconv.Itoa(intValue + 1)
	store.Storage[key] = Data{value: resultValue, expiry: output.expiry, tll: output.tll}
	return resultValue, nil
}
