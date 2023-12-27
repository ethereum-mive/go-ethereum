package syncx

import "github.com/ethereum/go-ethereum/internal/syncx"

type ClosableMutex = syncx.ClosableMutex

func NewClosableMutex() *ClosableMutex {
	return syncx.NewClosableMutex()
}
