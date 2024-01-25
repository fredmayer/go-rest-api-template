package dummy

import (
	"context"
	"github.com/fredmayer/go-rest-api-template/internal/controllers"
	"github.com/fredmayer/go-rest-api-template/internal/domain/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	ctx     context.Context
	service ServiceDummy
}

type ServiceDummy interface {
	Find(id int) (models.Dummy, error)
}

func NewHandler(ctx context.Context, service ServiceDummy) *Handler {
	return &Handler{
		ctx,
		service,
	}
}

func (h Handler) FindDummy(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return controllers.ErrBadRequest(c, "id", "must be number")
	}

	model, err := h.service.Find(id)
	if err != nil {
		return controllers.NotFound(c)
	}

	return c.JSON(http.StatusOK, model)
}
