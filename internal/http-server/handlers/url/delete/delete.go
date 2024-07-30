package delete

import (
	"Rest1/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		alias := chi.URLParam(request, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(writer, request, response.Error("invalid request"))
			return
		}

		err := deleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url")
			render.JSON(writer, request, response.Error("failed to delete request"))
			return
		}
		log.Info("url deleted", slog.String("alias", alias))
	}
}
