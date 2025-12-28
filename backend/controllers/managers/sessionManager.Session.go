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

func (s *SessionManager) GetAllSessionConnections(sessionId string) (map[string]*websocket.Conn, error) {

	if session, isPresent := s.SessionDict[sessionId]; isPresent {
		return session.ClientConnections, nil
	}

	return nil, fmt.Errorf("no session found")
}

func (s *SessionManager) GetUserSessionConnection(sessionId string, userID string, ctx *gin.Context) (*websocket.Conn, error) {

	if sess, isSessionPresent := s.SessionDict[sessionId]; isSessionPresent {
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
			ActionCount:       0,
			RequiresNewTally:  true,
			LastUsed:          time.Now(),
			BroadcastChannel:  make(chan gin.H, 100),
			SignalChannel:     make(chan Action, 100),
		}

		sess.ClientConnections[userID] = wsConn

		go sess.RunPublisherService()
		go sess.RunTallyService()

		s.SessionDict[sessionId] = sess

		return wsConn, nil
	}
}

func (s *SessionManager) AddConnection(sessionId string, newConnection *websocket.Conn) error {

	return nil
}

func (s *SessionManager) DeleteSession(sessionId string) error {
	session, isPresent := s.SessionDict[sessionId]
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
	delete(s.SessionDict, sessionId)

	return nil
}

func (s *SessionManager) CleanSessions() error {
	return nil
}
