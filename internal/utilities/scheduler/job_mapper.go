package scheduler

import "sync"

type jobMapper struct {
	mu       sync.RWMutex
	jobsByID map[int64]string
}

func newJobMapper() jobMapper {
	return jobMapper{
		jobsByID: make(map[int64]string),
	}
}

func (m *jobMapper) Add(id int64, schedulerJobID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsByID[id] = schedulerJobID
}

func (m *jobMapper) GetSchedulerJobID(id int64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobID, ok := m.jobsByID[id]
	return jobID, ok
}

func (m *jobMapper) GetDBJobID(schedulerJobID string) (int64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for dbJobID, localJobID := range m.jobsByID {
		if localJobID == schedulerJobID {
			return dbJobID, true
		}
	}

	return 0, false
}

func (m *jobMapper) Remove(id int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.jobsByID, id)
}
