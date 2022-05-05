package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/notification"
)

type notificationHandler struct {
  service *notification.Service
}

func InitNotificationHandler(s *notification.Service, r *mux.Router, n *negroni.Negroni) {
  h := notificationHandler{
		service: s,
	}

	r.Handle("/notifications", n.With(
		negroni.WrapFunc(laxo.AssureAuth(h.GetNotifications)),
	)).Methods("GET")
}

func (h *notificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request, uID string) {
  js, err := h.service.GetNotificationsJSON(uID, 0, 50)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

