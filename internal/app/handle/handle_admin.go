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

		pass, err := process(requireOne(r, admincfg.LOGIN_PASSPHRASE_PARAM), htmlUnescapeString)
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

		executeAdminIndex(w, r, newAdminView())
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

type adminView struct {
	Error      error
	Message    string
	WriteCount int
}

func newAdminView() *adminView {
	return &adminView{
		WriteCount: -1,
	}
}

func handleAdmin() ServiceHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, sv *service.BoardService) {

		view := newAdminView()

		fp1 := ps.ByName("fp1")
		fp2 := ps.ByName("fp2")

		switch fp1 {
		case "create-board":
			switch fp2 {
			case "news4vip", "poverty":
				view.Error = sv.Admin.CreateBoard(fp2)
			default:
				view.Error = fmt.Errorf("unsupported: %v", fp2)
			}
		case "write-limit":
			switch fp2 {
			case "get":
				view.WriteCount, view.Error = sv.Admin.GetWriteCount()
			case "reset":
				view.Error = sv.Admin.ResetWriteCount()
			default:
				view.Error = fmt.Errorf("unsupported: %v", fp2)
			}
		default:
			view.Error = fmt.Errorf("unknown func %v/%v", fp1, fp2)
		}
		executeAdminIndex(w, r, view)
	}
}

func executeAdminIndex(w http.ResponseWriter, r *http.Request, view *adminView) {

	if view.Error == nil {
		view.Message = "NO ERRORS."
	} else {
		view.Message = fmt.Sprint(view.Error)
	}

	if err := adminIndexTmpl.Execute(w, view); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
