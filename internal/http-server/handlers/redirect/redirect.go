package redirect

import (
	"Rest1/internal/lib/api/response"
	"Rest1/internal/lib/logger/sl"
	"Rest1/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		alias := chi.URLParam(request, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(writer, request, response.Error("invalid request"))
		}

		resURL, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(writer, request, response.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(writer, request, response.Error("internal error"))
			return
		}
		log.Info("url found", slog.String("url", resURL))
		http.Redirect(writer, request, resURL, http.StatusFound)
	}
}
