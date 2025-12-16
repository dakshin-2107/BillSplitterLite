package common

import "github.com/gin-gonic/gin"

// ActionExecutionResponse creates a standardized response for action execution
func ActionExecutionResponse(success bool, message string, action IAction) gin.H {
	return gin.H{
		"success": success,
		"message": message,
		"action":  action,
	}
}

// ActionExecutionSuccessResponse creates a success response for action execution
func ActionExecutionSuccessResponse(action IAction) gin.H {
	return ActionExecutionResponse(true, "action was executed", action)
}

// ActionExecutionFailureResponse creates a failure response for action execution
func ActionExecutionFailureResponse(action IAction, err error) gin.H {
	response := ActionExecutionResponse(false, "action could not be executed or socket was closed", action)
	response["err"] = err
	return response
}

// SessionCreationSuccessResponse creates a success response for session creation
func SessionCreationSuccessResponse(newSplit Split) gin.H {
	return gin.H{
		"success": true,
		"message": "Session has been created",
		"bill":    newSplit,
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
