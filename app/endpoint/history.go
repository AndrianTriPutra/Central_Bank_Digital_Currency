package endpoint

import (
	"atp/cbdc/pkg/utils/echos/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) GetChain(c echo.Context) error {
	ctx := c.Request().Context()

	chain, err := h.ucase.GetChain(ctx)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed get chain",
			Cause:     err.Error(),
		}
	}

	response := util.WrapSuccessResponse("success", chain)
	return c.JSON(http.StatusOK, response)
}

func (h handler) GetHistory(c echo.Context) error {
	ctx := c.Request().Context()

	addr := c.FormValue("address")
	if addr == "" {
		return util.CustomError{
			ErrorType: util.ErrBadRequest,
			Message:   "The given data was invalid",
			Cause:     "address is required",
		}
	}

	history, err := h.ucase.History(ctx, addr)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed get history",
			Cause:     err.Error(),
		}
	}

	response := util.WrapSuccessResponse("success", history)
	return c.JSON(http.StatusOK, response)
}
