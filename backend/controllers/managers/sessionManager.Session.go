package managers

import (
	"fmt"
	"time"

	"github.com/dakshin-2107/BillSplitterLite/backend/common"
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func testSessionManager() common.ISessionManager {
	return &SessionManager{}
}

func (s *SessionManager) Init(logger logger.ILogger, connUpgrader *websocket.Upgrader, modelHelper common.ISplitModelHelper) {
	s.logger = logger
	s.Upgrader = connUpgrader
	s.ModelHelper = modelHelper
	s.ConnectionDict = make(map[string]*SplitSession)
}

func (s *SessionManager) GetSessionConnnections(sessionId string) (map[string]*websocket.Conn, error) {

	if session, isPresent := s.ConnectionDict[sessionId]; isPresent {
		return session.ClientConnections, nil
	}

	return nil, fmt.Errorf("no session found")
}

func (s *SessionManager) CreateNewSession(sessionId string, userID string, ctx *gin.Context) (*websocket.Conn, error) {

	if sess, isPresent := s.ConnectionDict[sessionId]; isPresent {
		return sess.AdminConnection, nil
	} else {
		wsConn, err := s.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			s.logger.DebugLog(fmt.Sprintf("Failed to upgrade connection: %v", err))
			return nil, err
		}

		s.ConnectionDict[sessionId] = &SplitSession{
			ClientConnections: make(map[string]*websocket.Conn),
			AdminConnection:   wsConn,
			ActionCount:       0,
			IsSessionActive:   true,
			LastUsed:          time.Now(),
		}
		return wsConn, nil
	}
}

func (s *SessionManager) AddConnection(sessionId string, newConnection *websocket.Conn) error {
	return nil
}

func (s *SessionManager) DeleteSession(sessionId string) error {
	return nil
}

func (s *SessionManager) CleanSessions() error {
	return nil
}
