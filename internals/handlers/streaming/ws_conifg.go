package streaming

import (
	"sync"

	"github.com/gofrs/uuid"
)

type Notification struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

var (
	Mu          sync.RWMutex
	Connections = make(map[uuid.UUID]chan Notification)
)
