package server

import (
	"redis/internal/storage"
	"strings"
)

func HandleCommand(command string) string {
	parts := strings.Split(command, " ")
	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "SET":
		if len(parts) != 3 {
			return "ERROR: Invalid SET command"
		}
		storage.Set(parts[1], parts[2])
		return "OK"
	case "GET":
		if len(parts) != 2 {
			return "ERROR: Invalid GET command"
		}
		value := storage.Get(parts[1])
		if value == "" {
			return "(nil)"
		}
		return value
	case "DEL":
		if len(parts) != 2 {
			return "ERROR: Invalid DEL command"
		}
		storage.Del(parts[1])
		return "OK"
	case "INCR":
		if len(parts) != 2 {
			return "ERROR: Invalid INCR command"
		}
		value := storage.Incr(parts[1])
		return value
	case "DEC":
		if len(parts) != 2 {
			return "ERROR: Invalid INCR command"
		}
		value := storage.Dec(parts[1])
		return value
	default:
		return "ERROR: Unknown command"
	}
}
