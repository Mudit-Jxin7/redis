package server

import (
	"redis/internal/storage"
	"strconv"
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
	case "DECR":
		if len(parts) != 2 {
			return "ERROR: Invalid DECR command"
		}
		value := storage.Decr(parts[1])
		return value
	case "LPUSH":
		if len(parts) < 3 {
			return "ERROR: Invalid LPUSH command"
		}
		value := storage.LPush(parts[1], parts[2:]...)
		return value
	case "RPUSH":
		if len(parts) < 3 {
			return "ERROR: Invalid RPUSH command"
		}
		value := storage.RPush(parts[1], parts[2:]...)
		return value
	case "LPOP":
		if len(parts) != 2 {
			return "ERROR: Invalid LPOP command"
		}
		value := storage.LPop(parts[1])
		return value
	case "RPOP":
		if len(parts) != 2 {
			return "ERROR: Invalid RPOP command"
		}
		value := storage.RPop(parts[1])
		return value
	case "LRANGE":
		if len(parts) != 4 {
			return "ERROR: LRANGE requires 3 arguments"
		}
		start, _ := strconv.Atoi(parts[2])
		stop, _ := strconv.Atoi(parts[3])
		value := storage.LRange(parts[1], start, stop)
		return value
	case "LLEN":
		if len(parts) != 2 {
			return "ERROR: Invalid LLEN command"
		}
		value := storage.LLen(parts[1])
		return value
	default:
		return "ERROR: Unknown command"
	}
}
