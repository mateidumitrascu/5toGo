package server

import (
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type ServerData struct {
	startTime       time.Time
	activeConnCount atomic.Int64
}

func (s *ServerData) serverHealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(s.startTime)

	responseString := fmt.Sprintf(">> SERVER HEALTH\n\tUPTIME: %s\n\tACTIVE CONNECTIONS: %d\n", uptime.String(), s.activeConnCount.Load())

	if _, err := fmt.Fprintf(w, "%s", responseString); err != nil {
		fmt.Printf("Error responding with uptime: %v", err)
	}
}

func (s *ServerData) connState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		s.activeConnCount.Add(1)
	case http.StateClosed:
		s.activeConnCount.Add(-1)
	}
}
