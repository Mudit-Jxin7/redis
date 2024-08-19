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
	case "HSET":
		if len(parts) != 4 {
			return "ERROR: HSET requires 3 arguments"
		}
		value := storage.HSet(parts[1], parts[2], parts[3])
		return value
	case "HGET":
		if len(parts) != 3 {
			return "ERROR: HGET requires 2 arguments"
		}
		value := storage.HGet(parts[1], parts[2])
		return value
	case "HMSET":
		if len(parts) < 4 || len(parts)%2 != 0 {
			return "ERROR: HMSET requires an even number of arguments"
		}
		fieldValues := make(map[string]string)
		for i := 2; i < len(parts); i += 2 {
			fieldValues[parts[i]] = parts[i+1]
		}
		value := storage.HMSet(parts[1], fieldValues)
		return value
	case "HMGET":
		if len(parts) < 3 {
			return "ERROR: HMGET requires at least 2 arguments"
		}
		values := storage.HMGet(parts[1], parts[2:]...)
		return strings.Join(values, " ")
	case "HGETALL":
		if len(parts) != 2 {
			return "ERROR: HGETALL requires 1 argument"
		}
		values := storage.HGetAll(parts[1])
		var result []string
		for field, value := range values {
			result = append(result, field, value)
		}
		return strings.Join(result, " ")
	case "HDEL":
		if len(parts) < 3 {
			return "ERROR: HDEL requires at least 2 arguments"
		}
		count := storage.HDel(parts[1], parts[2:]...)
		return strconv.Itoa(count)
	default:
		return "ERROR: Unknown command"
	}
}
