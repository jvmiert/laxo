package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"laxo.vn/laxo/laxo"
	"laxo.vn/laxo/laxo/notification"
)

type notificationHandler struct {
  server *laxo.Server
  service *notification.Service
}

func InitNotificationHandler(server *laxo.Server, s *notification.Service, r *mux.Router, n *negroni.Negroni) {
  h := notificationHandler{
    server: server,
		service: s,
	}

	r.Handle("/notifications", n.With(
		negroni.WrapFunc(server.Middleware.AssureAuth(h.GetNotifications)),
	)).Methods("GET")
}

func (h *notificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request, uID string) {
  js, err := h.service.GetNotificationsJSON(uID, 0, 50)
  if err != nil {
    h.server.Logger.Errorw("GetNotifications handler error",
      "error", err,
    )
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Write(js)
}

