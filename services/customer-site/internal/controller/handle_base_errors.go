package controller

// import (
// 	"errors"
// 	"log/slog"
// 	"net/http"

// 	"github.com/fedotovmax/microservices-shop/customer-site/internal/openapiclient"
// 	"github.com/starfederation/datastar-go/datastar"
// )

// func (c *controller) handleErrors(w http.ResponseWriter, r *http.Request, sse *datastar.ServerSentEventGenerator, err error) {

// 	var genericOpenAPIError *openapiclient.GenericOpenAPIError

// 	switch {
// 	case errors.Is(err, ErrUnauthorized):
// 		http.Redirect(w, r, "/logout", http.StatusFound)
// 		return
// 	case errors.As(err, &genericOpenAPIError):
// 		switch t := genericOpenAPIError.Model().(type) {
// 		case openapiclient.ErrdetailsBadRequest:
// 			slog.Error("validation errors (from grpc)", slog.Any("violations", t.FieldViolations))
// 		case openapiclient.HttputilsErrorResponse:
// 			slog.Error("http error", slog.Any("error", t.GetMessage()))
// 		case openapiclient.GithubComFedotovmaxMicroservicesShopApiGatewayInternalDomainLoginErrorResponse:
// 			slog.Warn("login error response", slog.Any("current type", t))
// 		default:
// 			slog.Info("unknown http error", slog.Any("unknown error", err.Error()))
// 		}

// 	}
// }
