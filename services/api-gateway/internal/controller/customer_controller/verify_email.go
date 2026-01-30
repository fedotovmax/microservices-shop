package customercontroller

import (
	"log/slog"
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/userspb"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Verify email
// @Description  Verify email by verification link from user email address
// @Router       /customers/users/verify-email/{link} [get]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param X-Request-Locale header string false "Locale"
// @Param link   path string true "Verification link"
// @Success      200  {object}  userspb.EmailVerifiedSuccess
// @Failure      403  {object}  userspb.VerifyEmailLinkExpired
// @Failure      404  {object}  userspb.VerifyEmailLinkNotFound
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) verifyEmail(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.refreshSession"

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	link := r.PathValue("link")

	if link == "" {
		httputils.WriteJSON(w, http.StatusBadRequest, httputils.NewError("path parameter {link} is empty"))
		return
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	res, err := c.users.VerifyEmail(ctx, &userspb.VerifyEmailRequest{
		Link: link,
	})

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	c.handleVerifyEmailResponse(w, res)

}

func (c *controller) handleVerifyEmailResponse(w http.ResponseWriter, res *userspb.VerifyEmailResponse) {
	const op = "controller.customers.handleVerifyEmailResponse"

	l := c.log.With(slog.String("op", op))

	switch t := res.Payload.(type) {
	case *userspb.VerifyEmailResponse_LinkExpired:
		httputils.WriteJSON(w, http.StatusForbidden, t.LinkExpired)
		return
	case *userspb.VerifyEmailResponse_NotFound:
		httputils.WriteJSON(w, http.StatusNotFound, t.NotFound)
		return
	case *userspb.VerifyEmailResponse_Ok:
		httputils.WriteJSON(w, http.StatusOK, t.Ok)
		return
	default:
		l.Error(
			"unknown response from grpc.users-service.VerifyEmail",
			slog.Any("response.payload", res.Payload),
		)
		httputils.WriteJSON(
			w,
			http.StatusInternalServerError,
			httputils.NewError("unknown verify email response"),
		)
	}

}
