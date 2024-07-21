package save

import (
	"Rest1/internal/lib/api/response"
	"Rest1/internal/lib/logger/sl"
	"Rest1/internal/lib/random"
	"Rest1/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		var req Request
		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(writer, request, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			validateErr := err.(validator.ValidationErrors)
			render.JSON(writer, request, response.ValidationError(validateErr))
			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.RandomString(aliasLength)
		}
		id, err := saver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(writer, request, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(writer, request, response.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOk(writer, request, alias)
	}
}

func responseOk(writer http.ResponseWriter, request *http.Request, alias string) Response {
	return Response{
		Response: response.OK(),
		Alias:    alias,
	}
}
