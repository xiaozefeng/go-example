package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits     = uint64(5) // 5bit workerID out of 10bit worker machine ID
	dataCenterIDBits = uint64(5) // 5bit workerID out of 10bit worker dataCenterID
	sequenceBits     = uint64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits) // The maximum value of the node ID used to prevent overflow
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // Timestamp offset left
	dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	workLeft = uint8(12) // workLeft = sequenceBits // Node IDx offset to the left

	twepoch = int64(1659674040000) // constant timestamp (milliseconds)
)

type Worker struct {
	mu           sync.Mutex
	LastStamp    int64
	WorkerID     int64
	DataCenterID int64
	Sequence     int64
}

func NewWorker(workerID, dataCenterID int64) *Worker {
	return &Worker{
		LastStamp:    workerID,
		WorkerID:     0,
		DataCenterID: 0,
		Sequence:     dataCenterID,
	}
}

func (w *Worker) getMillSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Worker) NextID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *Worker) nextID() (uint64, error) {
	timeStamp := w.getMillSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards, waiting until")
	}
	if w.LastStamp == timeStamp {
		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMillSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) | (w.DataCenterID << dataLeft) | (w.WorkerID << workLeft) | w.Sequence
	return uint64(id), nil
}
