package endpoint

import (
	"atp/cbdc/pkg/utils/domain"
	"atp/cbdc/pkg/utils/echos/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) GetBalance(c echo.Context) error {
	ctx := c.Request().Context()

	addr := c.FormValue("address")
	if addr == "" {
		return util.CustomError{
			ErrorType: util.ErrBadRequest,
			Message:   "The given data was invalid",
			Cause:     "address is required",
		}
	}

	value, err := h.ucase.GetBalance(ctx, addr)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed get balance",
			Cause:     err.Error(),
		}
	}

	balance := domain.Balance{
		Address: addr,
		Balance: value,
	}

	response := util.WrapSuccessResponse("success", balance)
	return c.JSON(http.StatusOK, response)
}

func (h handler) GetWallet(c echo.Context) error {
	ctx := c.Request().Context()

	w, err := h.ucase.GetWallet(ctx)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed get wallet",
			Cause:     err.Error(),
		}
	}
	response := util.WrapSuccessResponse("success", w)
	return c.JSON(http.StatusOK, response)
}
