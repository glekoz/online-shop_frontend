package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/glekoz/online-shop_proto/product"
	"github.com/glekoz/online-shop_proto/user"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func (web *Weber) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ne naydeno"))
		return
	}
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

func (web *Weber) GetCreateProductHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./internal/web/static/prod_create.html"))
	tmpl.Execute(w, nil)
}

func (web *Weber) PostCreateProductHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("sorry"))
		return
	}

	name := r.PostForm.Get("name")
	pricestr := r.PostForm.Get("price")
	description := r.PostForm.Get("description")

	if name == "" || pricestr == "" || description == "" {
		w.Write([]byte("чего-то не хватает"))
		return
	}

	price, err := strconv.Atoi(pricestr)
	if err != nil {
		w.Write([]byte("в цене указана не цена"))
		return
	}

	id, err := web.ProductClient.Create(r.Context(), &product.Product{Name: name, Description: description, Price: int32(price)})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			w.Write([]byte(err.Error()))
			return
		}
		WriteStatusCode(w, st)
		w.Write([]byte(st.Message()))
		return
	}
	w.Write([]byte(id.GetId()))
}

func (web *Weber) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	sort := q.Get("sort")

	// оформить эти множественные проверки через валидатор
	lp, err := ParseIntFromQuery(q.Get("low-price"), 0)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("нормально мин цену напиши"))
		return
	}
	hp, err := ParseIntFromQuery(q.Get("high-price"), 100_000_000)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("нормально макс цену напиши"))
		return
	}
	page, err := ParseIntFromQuery(q.Get("page"), 1)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("нормально страницу напиши"))
		return
	}
	pageSize, err := ParseIntFromQuery(q.Get("page-size"), 10)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("нормально размер страницы напиши"))
		return
	}

	resp, err := web.ProductClient.GetAll(r.Context(), &product.Filter{
		Name:      name,
		LowPrice:  int32(lp) * 100,
		HighPrice: int32(hp) * 100,
		OrderBy:   sort,
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			WriteStatusCode(w, st)
			w.Write([]byte(st.Message()))
			for _, d := range st.Details() {
				det, ok := d.(*errdetails.BadRequest)
				if ok {
					for _, field := range det.FieldViolations {
						w.Write([]byte(field.Field))
					}
				} else {
					w.Write([]byte("errdetails не прочитались\n"))
				}
			}
			return
		}
		return
	}
	sb := strings.Builder{}
	for _, r := range resp.GetProducts() {
		sb.WriteString(fmt.Sprintf("ID: %s,\t Name:%s,\t Price: %d\n", r.Id, r.Name, r.Price/100))
	}
	w.Write([]byte(sb.String()))
}

func (web *Weber) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("net takogo producta"))
		return
	}
	res, err := web.ProductClient.Get(r.Context(), &product.ID{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			w.Write([]byte("owu5ka"))
			return
		}
		w.Write([]byte(st.Message()))
		return
	}
	w.Write([]byte(fmt.Sprintf("Name: %s\t Description: %s\nPrice: %d rubles %d kopecks", res.Name, res.Description, res.Price/100, res.Price%100)))
	name, ok := readNameFromCtx(r.Context())
	if ok {
		w.Write([]byte(fmt.Sprintf("\n\n\nHello, %s", name)))
	}
}
