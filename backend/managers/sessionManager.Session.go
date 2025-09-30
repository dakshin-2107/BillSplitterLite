package managers

import (
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
)

func testSessionManager() ISessionManager {
	return &SessionManager{}
}

func (s *SessionManager) Init(logger logger.ILogger, connUpgrader *ConnectionUpgrader) {
	s.logger = logger
	s.upgrader = connUpgrader
	s.ConnectionDict = make(map[string]*SplitSession)
}

func (s *SessionManager) GetSessionConnnections(sessionId string) ([]*MyWebSocketConnection, error) {
	session, exists := s.ConnectionDict[sessionId]
	if exists {
		s.logger.DebugLog("Finding connections for split ID : " + sessionId)
		return session.ClientConnections, nil
	}

	s.logger.DebugLog("No entry found for split ID : " + sessionId)
	return nil, fmt.Errorf("no split ID found")
}

func (s *SessionManager) CreateNewSession(sessionId string) error {
	_, exists := s.ConnectionDict[sessionId]
	if !exists {
		s.logger.DebugLog("Session with ID : " + sessionId + " already exists")
		s.ConnectionDict[sessionId] = &SplitSession{
			LastUsed:        time.Now(),
			IsSessionActive: true,
		}

		s.logger.DebugLog("Created new session with ID : " + sessionId)
		return nil
	}

	s.logger.DebugLog("Could not create new session with ID : " + sessionId)
	return fmt.Errorf("could not create new session with ID : %s", sessionId)
}

func (s *SessionManager) AddConnection(sessionId string, newConnection *MyWebSocketConnection) error {
	session, exists := s.ConnectionDict[sessionId]
	if exists {
		s.logger.DebugLog("Adding new connection to session : " + sessionId)
		session.ClientConnections = append(session.ClientConnections, newConnection)

		if len(session.ClientConnections) == 1 {
			session.AdminConnection = newConnection // first connection = admin connection
		}

		s.logger.DebugLog("Added new connection to session ID : " + sessionId)
		return nil
	}

	s.logger.DebugLog("No session found with ID : " + sessionId)
	return fmt.Errorf("no session found")
}

func (s *SessionManager) DeleteSession(sessionId string) error {
	return nil
}

func (s *SessionManager) CleanSessions() error {
	return nil
}
