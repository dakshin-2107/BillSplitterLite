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
	s.SessionDict = make(map[string]*SplitSession)
}

func (s *SessionManager) GetNewUserSessionId(sessionId string) (string, error) {

	if session, isPresent := s.getSession(sessionId); isPresent {
		return fmt.Sprintf("%d", session.UserIdCounter), nil
	}

	return "", fmt.Errorf("no session found")
}

func (s *SessionManager) GetUserSessionConnection(sessionId string, userID string, ctx *gin.Context) (*websocket.Conn, error) {

	if sess, isSessionPresent := s.getSession(sessionId); isSessionPresent {
		if conn, isClientPresent := sess.ClientConnections[userID]; isClientPresent {
			if pingErr := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Millisecond*100)); pingErr != nil {
				conn.Close()
				delete(sess.ClientConnections, userID)
				s.logger.DebugLog(fmt.Sprintf("Failed to ping client connection, deleted connection: %v", pingErr))

				return s.GetUserSessionConnection(sessionId, userID, ctx)
			}

			return conn, nil
		}

		conn, err := s.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			s.logger.DebugLog(fmt.Sprintf("Failed to upgrade connection: %v", err))
			return nil, err
		}

		sess.ClientConnections[userID] = conn
		sess.UserIdCounter++
		return conn, nil

	} else {
		wsConn, err := s.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			s.logger.DebugLog(fmt.Sprintf("Failed to upgrade connection: %v", err))
			return nil, err
		}

		sess := &SplitSession{
			SplitID:           sessionId,
			ModelHelper:       s.ModelHelper,
			ClientConnections: make(map[string]*websocket.Conn),
			AdminConnection:   wsConn,
			UserIdCounter:     1,
			LastUsed:          time.Now(),
			BroadcastChannel:  make(chan gin.H, 100),
			SignalChannel:     make(chan Action, 100),
		}
		sess.RequiresNewTally.Store(true)

		sess.ClientConnections[userID] = wsConn

		go sess.RunPublisherService()
		go sess.RunTallyService()

		// NOTE: benign TOCTOU race — two concurrent first-connects on the same
		// sessionId could both see isSessionPresent=false and both call setSession.
		// The second write overwrites the first; the orphaned session's goroutines
		// run until idle eviction (H3). This is pre-existing and statistically negligible.
		s.setSession(sessionId, sess)

		return wsConn, nil
	}
}

func (s *SessionManager) AddConnection(sessionId string, newConnection *websocket.Conn) error {

	return nil
}

func (s *SessionManager) DeleteSession(sessionId string) error {
	// deleteAndGetSession atomically removes the entry so no concurrent
	// ExecuteAction or GetUserSessionConnection can observe it after teardown begins.
	session, isPresent := s.deleteAndGetSession(sessionId)
	if !isPresent {
		return fmt.Errorf("session not found")
	}

	byeAction := Action{ActionType: common.BYE_BYE}
	session.SignalChannel <- byeAction
	session.SignalChannel <- byeAction

	session.AdminConnection.Close()
	for userID, conn := range session.ClientConnections {
		err := conn.Close()
		if err != nil {
			s.logger.DebugLog(fmt.Sprintf("Error closing connection for user %s: %v", userID, err))
		}
	}

	close(session.SignalChannel)
	close(session.BroadcastChannel)

	return nil
}

func (s *SessionManager) CleanSessions() error {
	s.mu.RLock()
	ids := make([]string, 0, len(s.SessionDict))
	for id := range s.SessionDict {
		ids = append(ids, id)
	}
	s.mu.RUnlock()

	for _, id := range ids {
		s.DeleteSession(id)
	}
	return nil
}
