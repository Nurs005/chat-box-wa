package utils

import (
	"github.com/chatbox/whatsapp/internal/domain"
	"github.com/chatbox/whatsapp/pkg/errors"
	"github.com/go-chi/render"
	"net/http"
	"reflect"
)

func WriteJSONResponse(w http.ResponseWriter, r *http.Request, status int, message any, success bool) {
	if reflect.TypeOf(message).Kind() == reflect.Ptr || reflect.TypeOf(message).Kind() == reflect.Chan || reflect.TypeOf(message).Kind() == reflect.Interface {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, domain.BaseResponse{
			Success:    false,
			StatusCode: http.StatusInternalServerError,
			Message:    errors.InvalidTypeOfResponse.Error(),
		})
	}
	render.Status(r, status)
	render.JSON(w, r, domain.BaseResponse{
		Success:    success,
		StatusCode: status,
		Message:    message,
	})

}
