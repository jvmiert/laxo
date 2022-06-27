package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/user"
)

type userHandler struct {
	server  *laxo.Server
	service *user.Service
}

func InitUserHandler(server *laxo.Server, u *user.Service, r *mux.Router, n *negroni.Negroni) {
	h := userHandler{
		server:  server,
		service: u,
	}

	r.Handle("/user", n.With(
		negroni.WrapFunc(h.server.Middleware.AssureAuth(h.HandleGetUser)),
	)).Methods("GET")

	r.Handle("/logout", n.With(
		negroni.WrapFunc(h.HandleLogout),
	)).Methods("POST")

	r.Handle("/login", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.HandleLogin),
	)).Methods("POST")

	r.Handle("/user", n.With(
		negroni.HandlerFunc(h.server.Middleware.AssureJSON),
		negroni.WrapFunc(h.HandleCreateUser),
	)).Methods("POST")
}

func (h *userHandler) HandleGetUser(w http.ResponseWriter, r *http.Request, uID string) {
	fmt.Fprintf(w, "Hello, your uID is: %s\n", uID)
}

func (h *userHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(h.server.Config.AuthCookieName)

	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		h.server.Logger.Errorw("Error in auth handler function (cookie parsing)",
			"error", err,
		)
		return
	}

	err = h.service.RemoveUserSession(c.Value)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.service.RemoveUserCookie(w)
	fmt.Fprintf(w, "ok")
}

func (h *userHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest user.LoginRequest

	if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &loginRequest); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	printer := laxo.GetLocalePrinter(r)

	user, err := h.service.LoginUser(loginRequest.Email, loginRequest.Password, printer)

	if err != nil {
		laxo.ErrorJSONEncode(w, err, http.StatusUnauthorized)
		return
	}

	expireT, sessionKey, err := h.service.SetUserSession(user)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	user.SessionKey = sessionKey

	h.service.SetUserCookie(user.SessionKey, w, expireT)

	js, err := user.JSON()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (h *userHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User

	if err := laxo.DecodeJSONBody(h.server.Logger, w, r, &u.Model); err != nil {
		var mr *laxo.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	printer := laxo.GetLocalePrinter(r)
	if err := h.service.ValidateNew(&u, printer); err != nil {
		laxo.ErrorJSONEncode(w, err, http.StatusUnprocessableEntity)
		return
	}

	if _, err := h.service.SaveNewUserToDB(&u); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	expireT, sessionKey, err := h.service.SetUserSession(&u)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.SessionKey = sessionKey
	h.service.SetUserCookie(u.SessionKey, w, expireT)

	js, err := u.JSON()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}
