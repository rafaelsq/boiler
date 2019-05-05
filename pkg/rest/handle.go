package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rafaelsq/boiler/pkg/iface"
	er "github.com/rafaelsq/boiler/pkg/repository/email"
	ur "github.com/rafaelsq/boiler/pkg/repository/user"
	"github.com/rafaelsq/boiler/pkg/service"
)

func UserHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	st := ctx.Value("storage").(iface.Storage)

	if userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err == nil && userID > 0 {
		user, err := service.NewUser(ur.New(st)).ByID(ctx, int(userID))
		if err == nil {
			if err := json.NewEncoder(w).Encode(user); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func EmailsHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	st := ctx.Value("storage").(iface.Storage)

	if userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64); err == nil && userID > 0 {
		es, err := service.NewEmail(er.New(st)).ByUserID(ctx, int(userID))
		if err == nil {
			if err := json.NewEncoder(w).Encode(es); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func UsersHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ucase := service.NewUser(ur.New(ctx.Value("storage").(iface.Storage)))

	users, err := ucase.List(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
