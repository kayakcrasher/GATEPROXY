package main

import (
"log"
"net"
)

func main() {
listener, err := net.Listen("tcp", ":8080")
if err != nil {
log.Fatalf("failed to listen: %v", err)
}
defer listener.Close()
log.Println("GATEPROXY listening on :8080")

for {
conn, err := listener.Accept()
if err != nil {
log.Printf("accept error: %v", err)
continue
}
go handleClient(conn)
}
}
