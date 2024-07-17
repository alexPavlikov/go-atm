package router

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexPavlikov/go-atm/internal/server/locations"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	ACCOUNT_ID string = "Account_id"
)

type RouterBuilder struct {
	LocationsHandler *locations.Handler
}

func New(locationsHandler *locations.Handler) *RouterBuilder {
	return &RouterBuilder{
		LocationsHandler: locationsHandler,
	}
}

func (r *RouterBuilder) Build() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/accounts", mware(handlerWrapper(r.LocationsHandler.AccountsHandler)))
	router.Post("/accounts/{id}/deposit", mware(handlerWrapper(r.LocationsHandler.DepositAccountsHandler)))
	router.Post("/accounts/{id}/withdraw", mware(handlerWrapper(r.LocationsHandler.WithdrawAccountsHandler)))
	router.Get("/accounts/{id}/balance", mware(handlerWrapper(r.LocationsHandler.BalanceAccountsHandler)))

	return router
}

type wrappedFunc[Input, Output any] func(r *http.Request, data Input) (Output, error)

func handlerWrapper[Input, Output any](fn wrappedFunc[Input, Output]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Input

		if r.Method != "GET" {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&data); err != nil {
				slog.ErrorContext(r.Context(), "can't decode data", "error", err)
				http.Error(w, "Bad request"+err.Error(), http.StatusBadRequest)
				return
			}
		}

		response, err := fn(r, data)
		if err != nil {
			slog.ErrorContext(r.Context(), "can't handle request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(&response); err != nil {
			slog.ErrorContext(r.Context(), "can't encode request", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func mware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		str := chi.URLParam(r, "id")
		var id int
		var err error
		if str != "" {
			id, err = strconv.Atoi(str)
			if err != nil {
				slog.ErrorContext(r.Context(), "failed to convert dricer_id", "error", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), ACCOUNT_ID, id))
		h.ServeHTTP(w, r)
	}
}
