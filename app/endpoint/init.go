package endpoint

import (
	"atp/cbdc/app/usecase/utransaction"
	"atp/cbdc/pkg/utils/echos/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	ucase utransaction.UsecaseI
}

func NewHandler(e *echo.Echo, endpoint string, ucase utransaction.UsecaseI) {
	handler := handler{
		ucase: ucase,
	}

	e.POST(endpoint+"transaction", middleware.ErrorMiddleware(handler.Transaction))
	e.GET(endpoint+"balance", middleware.ErrorMiddleware(handler.GetBalance))
	e.GET(endpoint+"wallet", middleware.ErrorMiddleware(handler.GetWallet))
	e.GET(endpoint+"chain", middleware.ErrorMiddleware(handler.GetChain))
	e.GET(endpoint+"history", middleware.ErrorMiddleware(handler.GetHistory))
}
