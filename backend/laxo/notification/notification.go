package notification

import (
	"encoding/json"

	"laxo.vn/laxo/laxo/sqlc"
)

type Notification struct {
  Model       *sqlc.Notification       `json:"notification"`
  GroupModel  *sqlc.NotificationsGroup `json:"notificationGroup"`
}

func (n *Notification) JSON() ([]byte, error) {
  bytes, err := json.Marshal(n)

  if err != nil {
    return bytes, err
  }

  return bytes, nil
}
