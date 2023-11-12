package models

import "net/http"

type errorResponse struct {
	StatusCode  int    `json:"status_code"`
	Error       string `json:"error"`
	Information any    `json:"information,omitempty"`
}

func NewBadRequestErrorResponse(information any) errorResponse {
	return errorResponse{
		StatusCode:  http.StatusBadRequest,
		Error:       "Bad request",
		Information: information,
	}
}

func NewInternalServerErrorResponse() errorResponse {
	return errorResponse{
		StatusCode:  http.StatusInternalServerError,
		Error:       "Internal server error",
		Information: "An error occurred that could not be processed, please try again later",
	}
}

func NewConflictErrorResponse(information any) errorResponse {
	return errorResponse{
		StatusCode:  http.StatusConflict,
		Error:       "Conflict",
		Information: information,
	}
}

func NewUnprocessableEntityErrorResponse(information any) errorResponse {
	return errorResponse{
		StatusCode:  http.StatusUnprocessableEntity,
		Error:       "Unprocessable entity",
		Information: information,
	}
}

func NewValidationErrorResponse(information any) errorResponse {
	return errorResponse{
		StatusCode:  http.StatusBadRequest,
		Error:       "Validation failed",
		Information: information,
	}
}

func NewAuthorizationErrorResponse(information any) errorResponse {
	return errorResponse{
		StatusCode:  http.StatusUnauthorized,
		Error:       "Authorization failed",
		Information: information,
	}
}
