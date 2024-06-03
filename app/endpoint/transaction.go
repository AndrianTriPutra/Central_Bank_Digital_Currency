package endpoint

import (
	"atp/cbdc/pkg/utils/domain"
	"atp/cbdc/pkg/utils/echos/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h handler) Transaction(c echo.Context) error {
	ctx := c.Request().Context()

	var data domain.Trans
	err := c.Bind(&data)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrBadRequest,
			Message:   "The given data was invalid",
			Cause:     "failed decode input",
		}
	}

	if data.From == data.To {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "duplicated data",
			Cause:     "data sender and receiver is same",
		}
	}

	w, existS, err := h.ucase.CheckWallet(ctx, data.From)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed CheckWallet",
			Cause:     err.Error(),
		}
	}

	if !existS {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed identify sender",
			Cause:     "sender not yet registered",
		}
	}

	err = h.ucase.Transaction(ctx, data)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed Transaction",
			Cause:     err.Error(),
		}
	}

	value, err := h.ucase.GetBalance(ctx, data.From)
	if err != nil {
		return util.CustomError{
			ErrorType: util.ErrInternalServer,
			Message:   "failed get balance",
			Cause:     err.Error(),
		}
	}

	balance := domain.Balance{
		Address: data.From,
		Balance: value,
	}

	existR := true
	for _, val := range w {
		if data.To == val {
			existR = false
		}
	}

	if existR {
		w = append(w, data.To)
		err = h.ucase.UpdateWallet(ctx, w)
		if err != nil {
			return util.CustomError{
				ErrorType: util.ErrInternalServer,
				Message:   "failed update wallet",
				Cause:     err.Error(),
			}
		}
	}

	response := util.WrapSuccessResponse("success", balance)
	return c.JSON(http.StatusOK, response)
}
