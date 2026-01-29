package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/omjikush09/redis-go/app/data_structure"
	// "github.com/omjikush09/redis-go/main"
	
)

func HandleCommands(input []string, config ConnectionConfig) string {
	if config.inTransaction && input[0] != "EXEC" && input[0] != "DISCARD" {
		if input[0] == "MULTI" {
			return encodeAsError("ERR MULTI calls can not be nested")
		}
		config.commands = append(config.commands, input)
		return encodeAsSimpleString("QUEUED")
	}
	switch input[0] {
	case "GET":
		return handleGetCommand(input[1:])
	case "SET":
		return handleSetCommand(input[1:])
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		input := input[1]
		output := fmt.Sprintf("$%d\r\n%s\r\n", len(input), input)
		return output
	case "RPUSH":
		return handleRPUSHCommand(input[1:])
	case "INCR":
		return handleINCRCommand(input[1:])
	case "MULTI":
		return handleMULTICommand(input[1:], config)
	case "EXEC":
		return handleExecuteCommand(config)
	case "DISCARD":
		return handleDiscardCommand(config)
	default:
		return "unknown command"
	}
}

var mapStore *data_structure.MapStoreStructure
var listStore *data_structure.ListStoreStructure

func InitializeStore() {
	mapStore = &data_structure.MapStoreStructure{
		Storage: make(map[string]data_structure.Data),
	}
	listStore = &data_structure.ListStoreStructure{
		Storage: make(map[string][]string),
	}
}

func handleSetCommand(input []string) string {
	key := input[0]
	value := input[1]
	var duration time.Duration
	tll := false
	if len(input) > 2 {
		if input[2] == "PX" {

			ms, _ := strconv.ParseInt(input[3], 10, 64)
			duration = time.Duration(ms) * time.Millisecond
			tll = true
		}
	}
	mapStore.Add(key, value, duration, tll)
	return encodeAsSimpleString("OK")
}

func handleGetCommand(input []string) string {
	key := input[0]
	value, exist := mapStore.Get(key)
	fmt.Printf("value is %s and exist is %t", value, exist)
	if !exist {

		return "$-1\r\n"
	}
	return encodeAsBulkString(value)
}

func handleRPUSHCommand(input []string) string {
	key := input[0]
	value, exist := listStore.Storage[key]
	if !exist {
		listStore.Storage[key] = make([]string, 0)
	}
	length := len(value) + len(input[1:])
	listStore.Storage[key] = append(listStore.Storage[key], input[1:]...)
	return encodeAsInteger(length)
}

func handleINCRCommand(input []string) string {
	key := input[0]
	value, err := mapStore.Increment(key)
	if err != nil {
		return encodeAsError(err.Error())
	}
	valueInt, _ := strconv.Atoi(value)
	return encodeAsInteger(valueInt)
}

func handleMULTICommand(input []string, config *ConnectionConfig) string {
	config.inTransaction = true
	config.commands = make([][]string, 0)
	return encodeAsSimpleString("OK")
}

func handleExecuteCommand(config *connectionConfig) string {
	fmt.Println(config.inTransaction)
	if !config.inTransaction {
		return encodeAsError("ERR EXEC without MULTI")
	}
	config.inTransaction = false
	responses := make([]string, 0)
	for _, command := range config.commands {
		response := HandleCommands(command, config)
		responses = append(responses, response)
	}

	config.commands = make([][]string, 0)
	fmt.Println(config.commands)
	return encodeAsArray(responses)
}

func handleDiscardCommand(config *connectionConfig) string {
	if !config.inTransaction {
		return encodeAsError("ERR DISCARD without MULTI")
	}
	config.inTransaction = false
	config.commands = make([][]string, 0)
	return encodeAsSimpleString("OK")
}
