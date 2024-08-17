
package server

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func HandleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading:", err)
            return
        }

        response := HandleCommand(strings.TrimSpace(message))
        conn.Write([]byte(response + "\n"))
    }
}
