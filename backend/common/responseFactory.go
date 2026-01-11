package common

import "github.com/gin-gonic/gin"

func TallyResponse(tally Tally) gin.H {
	return gin.H{
		"success": true,
		"message": "Tally has been calculated",
		"tally":   tally,
	}
}

func BillResponse(bill Split) gin.H {
	return gin.H{
		"success": true,
		"message": "Bill has been fetched",
		"bill":    bill,
	}
}

func BillDeleteResponse() gin.H {
	return gin.H{
		"success": true,
		"message": "Bill has been deleted",
		"bill":    nil,
	}
}

// ActionExecutionSuccessResponse creates a success response for action execution
func ActionExecutionSuccessResponse(action IAction) gin.H {
	return gin.H{
		"success": true,
		"message": "Action executed successfully",
		"action":  action,
	}
}

// ActionExecutionFailureResponse creates a failure response for action execution
func ActionExecutionFailureResponse(action IAction, err error) gin.H {
	return gin.H{
		"success": false,
		"message": "action could not be executed or socket was closed",
		"action":  action,
		"err":     err,
	}
}

// SessionCreationSuccessResponse creates a success response for session creation
func SessionCreationSuccessResponse(newSplit Split) gin.H {
	return gin.H{
		"success": true,
		"message": "Session has been created",
		"bill":    newSplit,
	}
}

// SessionCreationSuccessResponse creates a success response for session creation
func SessionAlreadyExistsResponse(split Split) gin.H {
	return gin.H{
		"success": true,
		"message": "Session already exists",
		"bill":    split,
	}
}

func SessionCreationFailureResponse() gin.H {
	return gin.H{
		"success": false,
		"message": "Parsing was not successful. Session could not be created.",
	}
}

func MissingCookieResponse() gin.H {
	return gin.H{
		"success": false,
		"message": "Missing bill_id or user id cookie",
	}
}

func InvalidRequestBody() gin.H {
	return gin.H{
		"success": false,
		"message": "Invalid request body",
	}
}

// SessionUriResponse creates a success response for session URI generation
func SessionUriResponse(billID string) gin.H {
	return gin.H{
		"success":   true,
		"message":   "Session URI generated",
		"sessionId": billID,
	}
}
