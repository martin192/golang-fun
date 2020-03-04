package main

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestThatWhenNewWorklogEntryIsPersisted_ThenItIsStoredCorrectly(t *testing.T) {
	worklogManager := NewWorklogManager(1)
	entriesToPersist := []WorklogEntry{"apple", "orange", "lemon"}

	for _, entry := range entriesToPersist {
		worklogManager.Persist(entry)
	}

	if found := worklogManager.GetAll(); !reflect.DeepEqual(found, entriesToPersist) {
		t.Errorf("Result was incorrect, got: %s, wanted: %s.", found, entriesToPersist)
	}
}

func TestThatMultipleThreadsPersistWorklogEntries_ThenCorrectEntriesAreStored(t *testing.T) {
	numberOfThreads := 20
	numberOfRecords := 100
	waitGroup := sync.WaitGroup{}
	worklogManager := NewWorklogManager(1)

	waitGroup.Add(numberOfThreads)

	for index := 0; index < numberOfThreads; index++ {
		go func(threadNumber int) {
			defer waitGroup.Done()

			for index := 0; index < numberOfRecords; index++ {
				worklogManager.Persist(fmt.Sprintf("Thread%d", threadNumber))
			}
		}(index)
	}

	waitGroup.Wait()

	if expectedResultsCount, totalRecordsCount := (numberOfThreads * numberOfRecords), len(worklogManager.GetAll()); expectedResultsCount != totalRecordsCount {
		t.Errorf("Result was incorrect, expected %d entries, found %d", expectedResultsCount, totalRecordsCount)
	}

	for index := 0; index < numberOfThreads; index++ {
		expectedCount := numberOfRecords
		expectedEntry := fmt.Sprintf("Thread%d", index)

		for _, entry := range worklogManager.GetAll() {
			if entry == expectedEntry {
				expectedCount--
			}
		}

		if expectedCount != 0 {
			t.Errorf("Result was incorrect, number of entries '%s' differs by %d.", expectedEntry, expectedCount)
		}
	}
}
