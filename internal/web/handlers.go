package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/glekoz/online-shop_proto/user"
	"google.golang.org/grpc/status"
)

func (web *Weber) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var guest string = "guest"
	name, ok := readNameFromCtx(r.Context())
	if ok {
		guest = name
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s", guest)))
}

func (web *Weber) GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./internal/web/static/reg.html"))
	tmpl.Execute(w, nil)
}

func (web *Weber) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("sorry"))
		return
	}
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	passwword := r.PostForm.Get("password")

	resp, err := web.UserClient.Register(r.Context(), &user.RegisterUserRequest{
		Username: name,
		Email:    email,
		Password: passwword,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			w.Write([]byte("plaki-plaki"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Code: %s, \nMessage: %s", st.Code(), st.Message())))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Weber) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./internal/web/static/log.html"))
	tmpl.Execute(w, nil)
}

func (web *Weber) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("sorry"))
		return
	}
	email := r.PostForm.Get("email")
	passwword := r.PostForm.Get("password")

	resp, err := web.UserClient.Login(r.Context(), &user.LoginUserRequest{
		Email:    email,
		Password: passwword,
	})
	if err != nil {
		w.Write([]byte("plaki-plaki"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (web *Weber) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/", http.StatusOK)
}
