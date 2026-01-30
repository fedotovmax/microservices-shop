package customercontroller

import (
	"net/http"

	"github.com/fedotovmax/httputils"
	_ "github.com/fedotovmax/microservices-shop/api-gateway/internal/domain"
	"github.com/fedotovmax/microservices-shop/api-gateway/internal/utils"
)

// @Summary      Check session auth
// @Description  Checking session auth for protected client routing
// @Router       /customers/session/check [get]
// @Tags         customers
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param X-Request-Locale header string false "Locale"
// @Success      200  {object}  domain.LocalSession
// @Failure      401  {object}  httputils.ErrorResponse
// @Failure      500  {object}  httputils.ErrorResponse
func (c *controller) checkSession(w http.ResponseWriter, r *http.Request) {
	session, err := utils.GetLocalSession(r)

	if err != nil {
		httputils.WriteJSON(w, http.StatusUnauthorized, httputils.NewError(err.Error()))
		return
	}

	httputils.WriteJSON(w, http.StatusUnauthorized, session)

}
