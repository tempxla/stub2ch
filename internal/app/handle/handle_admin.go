package handle

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/service"
	"log"
	"net/http"
)

func authenticate(sh ServiceHandle) ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		const invalid_user_message = "invalid user."

		c, err := r.Cookie(config.ADMIN_COOKIE_NAME)
		if err != nil || len(c.Value) != 1 {
			log.Println("admin cookie is missing")
			http.Error(w, invalid_user_message, http.StatusForbidden) // 403
			return
		}
		if err != sv.Admin.VerifySessionId(c.Value) {
			log.Println(err)
			http.Error(w, invalid_user_message, http.StatusForbidden) // 403
			return
		}

		sh(w, r, ps, sv)
	}
}

func handleAdminLogin() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

	}
}

func handleAdmin() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

	}
}
