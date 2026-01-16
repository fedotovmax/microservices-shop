package customercontroller

import (
	"net/http"

	"github.com/fedotovmax/httputils"
	"github.com/fedotovmax/microservices-shop-protos/gen/go/sessionspb"
	controolerPkg "github.com/fedotovmax/microservices-shop/api-gateway/internal/controller"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/keys"
	"google.golang.org/grpc/metadata"
)

// @Summary      Refresh session
// @Description  Refresh session by refresh token from cookie
// @Router       /customers/session/refresh-session [post]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param X-Request-Locale header string false "Locale"
// @Success      201  {object}  sessionspb.CreateSessionResponse
// @Failure      400  {object}  errdetails.BadRequest
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) refreshSession(w http.ResponseWriter, r *http.Request) {
	const op = "controller.customer.refreshSession"

	//l := c.log.With(slog.String("op", op))

	locale := r.Header.Get(keys.HeaderLocale)

	if locale == "" {
		locale = keys.FallbackLocale
	}

	userAgent := r.UserAgent()

	ip := controolerPkg.GetRealIP(r)

	//TODO: get from cookie
	refreshToken := "123"

	userSessionActionReq := &sessionspb.RefreshSessionRequest{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		Ip:           ip,
		Issuer:       c.issuer,
	}

	md := metadata.Pairs(
		keys.MetadataLocaleKey, locale,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)

	response, err := c.sessions.RefreshSession(ctx, userSessionActionReq)

	if err != nil {
		httputils.HandleErrorFromGrpc(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, response)

}
