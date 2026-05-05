package managers

// getSession returns the session for the given id under a read lock.
// After this call returns, the caller may use the pointer freely — the Go GC
// keeps the object alive and session deletion is expected to be rare.
func (s *SessionManager) getSession(id string) (*SplitSession, bool) {
	s.mu.RLock()
	sess, ok := s.SessionDict[id]
	s.mu.RUnlock()
	return sess, ok
}

// setSession writes a new session entry under an exclusive lock.
func (s *SessionManager) setSession(id string, sess *SplitSession) {
	s.mu.Lock()
	s.SessionDict[id] = sess
	s.mu.Unlock()
}

// deleteAndGetSession atomically removes the entry for id and returns the value
// that was there (if any). The check and delete are a single critical section,
// which prevents concurrent DeleteSession calls from double-closing channels.
func (s *SessionManager) deleteAndGetSession(id string) (*SplitSession, bool) {
	s.mu.Lock()
	sess, ok := s.SessionDict[id]
	if ok {
		delete(s.SessionDict, id)
	}
	s.mu.Unlock()
	return sess, ok
}
