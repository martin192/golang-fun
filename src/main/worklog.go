package main

import (
	"fmt"
	"log"
	"sync"
)

type WorklogEntry interface {
}

type WorklogManager interface {
	Persist(WorklogEntry)
	GetAll() []WorklogEntry
}

type inMemoryWorklogManager struct {
	writeMutex sync.Mutex
	records    []WorklogEntry
}

func (manager *inMemoryWorklogManager) Persist(entry WorklogEntry) {
	defer manager.writeMutex.Unlock()

	log.Println(fmt.Sprintf("Storing worklog entry '%v'...", entry))
	manager.writeMutex.Lock()
	manager.records = append(manager.records, entry)
}

func (manager *inMemoryWorklogManager) GetAll() []WorklogEntry {
	log.Println("Listing all worklog entries...")
	return append([]WorklogEntry(nil), manager.records...)
}

func NewWorklogManager(initialCapacity int) WorklogManager {
	return &inMemoryWorklogManager{
		writeMutex: sync.Mutex{},
		records:    make([]WorklogEntry, 0, initialCapacity),
	}
}
