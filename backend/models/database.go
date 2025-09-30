package models

import (
	"github.com/dakshin-2107/BillSplitterLite/backend/logger"
	"github.com/dakshin-2107/BillSplitterLite/backend/managers"
)

type IDatabase interface {
	Init(logger logger.ILogger) error
	IsConnected() bool
	ConnectToDatabase() error
	DisconnectFromDatabase() error
	ISplitModifier
	managers.IActionManager
}

type DatabaseProvider struct {
	Uri           string
	IsDBConnected bool
}
