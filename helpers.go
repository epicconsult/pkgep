package pkgep

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AppResponseStatus int

type ApiResponse struct {
	StatusCode AppResponseStatus `json:"status"`
	Message    string            `json:"message"`
	Data       interface{}       `json:"data,omitempty"`
}

const (
	Success                         AppResponseStatus = 20000
	InfoNotFound                    AppResponseStatus = 20001
	Created                         AppResponseStatus = 20100
	Accepted                        AppResponseStatus = 20200
	NoContent                       AppResponseStatus = 20400
	BadRequest                      AppResponseStatus = 40000
	EmpytyPostBody                  AppResponseStatus = 40001
	InvalidJsonBody                 AppResponseStatus = 40002
	MandatoryRequired               AppResponseStatus = 40003
	InvalidRequestField             AppResponseStatus = 40004
	IncorrectUserOrPwd              AppResponseStatus = 40005
	InvalidToken                    AppResponseStatus = 40006
	LockedAccount                   AppResponseStatus = 40007
	Unauthorized                    AppResponseStatus = 40100
	Forbidden                       AppResponseStatus = 40300
	NotFound                        AppResponseStatus = 40400
	MethodNotAllowed                AppResponseStatus = 40500
	EntityTooLarge                  AppResponseStatus = 41300
	Conflict                        AppResponseStatus = 40900
	Internal                        AppResponseStatus = 50000
	HaveErrorOccurredData           AppResponseStatus = 50001
	HaveErrorOccurredWhileQueryData AppResponseStatus = 50002
	NotImplemented                  AppResponseStatus = 50100
	BadGateway                      AppResponseStatus = 50200
	ServiceUnavailable              AppResponseStatus = 50300
	GatewayTimeout                  AppResponseStatus = 50400
	CustomStatus                    AppResponseStatus = 00000
)

var StatusCodeMap = map[AppResponseStatus]string{
	Success:        "Success",
	InfoNotFound:   "Information not found",
	Created:        "Created",
	Accepted:       "Accepted",
	NoContent:      "No Content",
	BadRequest:     "Bad Request",
	Internal:       "Internal Server Error",
	InvalidToken:   "The token is incorrect",
	Unauthorized:   "Unauthorized",
	NotFound:       "Not Found",
	EntityTooLarge: "Request body size limit exceeded (max 30MB)",
	Conflict:       "Conflict",
	Forbidden:      "Forbidden",
	CustomStatus:   "Custom Status",
}

func ErrorResponse(c *fiber.Ctx, appStatus AppResponseStatus, messages ...string) error {

	var message string
	if len(messages) > 0 {
		message = messages[0]
	}

	var resError ApiResponse
	httpStatus := http.StatusInternalServerError

	logNew.Path = c.Path()

	switch appStatus {
	case InfoNotFound:
		if message == "" {
			message = StatusCodeMap[InfoNotFound]
		}
		resError = ApiResponse{
			StatusCode: InfoNotFound,
			Message:    message,
		}
		httpStatus = http.StatusNotFound
	case BadRequest:
		msg := StatusCodeMap[BadRequest]
		if strings.TrimSpace(message) != "" {
			msg = fmt.Sprintf("%s, %s", msg, message)
		}
		resError = ApiResponse{
			StatusCode: BadRequest,
			Message:    msg,
		}
		httpStatus = http.StatusBadRequest
	case NotFound:
		if message == "" {
			message = StatusCodeMap[NotFound]
		}
		resError = ApiResponse{
			StatusCode: NotFound,
			Message:    message,
		}
		httpStatus = http.StatusNotFound
	case Conflict:
		if message == "" {
			message = StatusCodeMap[Conflict]
		}
		resError = ApiResponse{
			StatusCode: Conflict,
			Message:    message,
		}
		httpStatus = http.StatusConflict
	case InvalidToken:
		if message == "" {
			message = StatusCodeMap[InvalidToken]
		}
		resError = ApiResponse{
			StatusCode: InvalidToken,
			Message:    message,
		}
		httpStatus = http.StatusBadRequest
	case Unauthorized:
		if message == "" {
			message = StatusCodeMap[Unauthorized]
		}
		resError = ApiResponse{
			StatusCode: Unauthorized,
			Message:    message,
		}
		httpStatus = http.StatusUnauthorized
	case Forbidden:
		if message == "" {
			message = StatusCodeMap[Forbidden]
		}
		resError = ApiResponse{
			StatusCode: Forbidden,
			Message:    message,
		}
		httpStatus = http.StatusForbidden
	default:
		resError = ApiResponse{
			StatusCode: Internal,
			Message:    StatusCodeMap[Internal],
		}
		httpStatus = http.StatusInternalServerError
	}

	jsonStrRes, _ := json.Marshal(resError)
	logNew.LogInformation(HTTPRESPONSE, string(jsonStrRes))

	return c.Status(httpStatus).JSON(resError)
}

func SuccessResponse(c *fiber.Ctx, data interface{}, messages ...string) error {

	msg := StatusCodeMap[Success]
	if len(messages) > 0 {
		msg = messages[0]
	}

	var resError ApiResponse

	logNew.Path = c.Path()

	resError = ApiResponse{
		StatusCode: Success,
		Message:    msg,
		Data:       data,
	}

	jsonStrRes, _ := json.Marshal(msg)
	logNew.LogInformation(HTTPRESPONSE, string(jsonStrRes))

	return c.Status(http.StatusOK).JSON(resError)
}

var logNew *Logger

func NewHelpers(log Logger) {
	logNew = &log
}
