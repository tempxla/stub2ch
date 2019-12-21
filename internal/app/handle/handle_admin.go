package handle

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/internal/app/service"
	"log"
	"net/http"
)

func authenticate(sh ServiceHandle) ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		const invalid_user_message = "invalid user."

		c, err := r.Cookie(admincfg.LOGIN_COOKIE_NAME)
		if err != nil {
			log.Println("admin cookie is missing")
			http.Error(w, invalid_user_message, http.StatusForbidden) // 403
			return
		}
		if err != sv.Admin.VerifySession(c.Value) {
			log.Println(err)
			http.Error(w, invalid_user_message, http.StatusForbidden) // 403
			return
		}

		sh(w, r, ps, sv)
	}
}

func handleAdminLogin() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		const login_fail_message = "login failed."

		pass, err := process(requireOne(r, admincfg.LOGIN_PASSPHRASE_PARAM), HtmlUnescapeString)
		if err != nil {
			log.Print(err)
			http.Error(w, login_fail_message, http.StatusForbidden) // 403
			return
		}
		sig, err := process(requireOne(r, admincfg.LOGIN_SIGNATURE_PARAM))
		if err != nil {
			log.Print(err)
			http.Error(w, login_fail_message, http.StatusForbidden) // 403
			return
		}

		sid, err := sv.Admin.Login(pass, sig)
		if err != nil {
			log.Print(err)
			http.Error(w, login_fail_message, http.StatusForbidden) // 403
			return
		}

		cookie := &http.Cookie{
			Name:  admincfg.LOGIN_COOKIE_NAME,
			Value: sid,
			Path:  "/",
			//Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		executeAdminIndex(w, r)
	}
}

func handleAdminLogout() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {
		err := sv.Admin.Logout()
		if err != nil {
			fmt.Fprint(w, "Logout failed.")
			return
		}
		fmt.Fprint(w, "Logout success.")
	}
}

func handleAdmin() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

	}
}

func executeAdminIndex(w http.ResponseWriter, r *http.Request) {

	if err := adminIndexTmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
