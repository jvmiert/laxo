// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: notification.sql

package sqlc

import (
	"context"

	null "gopkg.in/guregu/null.v4"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
  notification_group_id, read, current_main_step,
  current_sub_step, main_message, sub_message, error
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, redis_id, notification_group_id, created, read, current_main_step, current_sub_step, main_message, sub_message, error
`

type CreateNotificationParams struct {
	NotificationGroupID string      `json:"notificationGroupID"`
	Read                null.Time   `json:"read"`
	CurrentMainStep     null.Int    `json:"currentMainStep"`
	CurrentSubStep      null.Int    `json:"currentSubStep"`
	MainMessage         null.String `json:"mainMessage"`
	SubMessage          null.String `json:"subMessage"`
	Error               null.Bool   `json:"error"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRow(ctx, createNotification,
		arg.NotificationGroupID,
		arg.Read,
		arg.CurrentMainStep,
		arg.CurrentSubStep,
		arg.MainMessage,
		arg.SubMessage,
		arg.Error,
	)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.RedisID,
		&i.NotificationGroupID,
		&i.Created,
		&i.Read,
		&i.CurrentMainStep,
		&i.CurrentSubStep,
		&i.MainMessage,
		&i.SubMessage,
		&i.Error,
	)
	return i, err
}

const createNotificationsGroup = `-- name: CreateNotificationsGroup :one
INSERT INTO notifications_group (
  user_id, workflow_id, entity_id, entity_type,
  total_main_steps, total_sub_steps, platform_name
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, workflow_id, entity_id, entity_type, platform_name, total_main_steps, total_sub_steps
`

type CreateNotificationsGroupParams struct {
	UserID         string      `json:"userID"`
	WorkflowID     null.String `json:"workflowID"`
	EntityID       string      `json:"entityID"`
	EntityType     string      `json:"entityType"`
	TotalMainSteps null.Int    `json:"totalMainSteps"`
	TotalSubSteps  null.Int    `json:"totalSubSteps"`
	PlatformName   string      `json:"platformName"`
}

func (q *Queries) CreateNotificationsGroup(ctx context.Context, arg CreateNotificationsGroupParams) (NotificationsGroup, error) {
	row := q.db.QueryRow(ctx, createNotificationsGroup,
		arg.UserID,
		arg.WorkflowID,
		arg.EntityID,
		arg.EntityType,
		arg.TotalMainSteps,
		arg.TotalSubSteps,
		arg.PlatformName,
	)
	var i NotificationsGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WorkflowID,
		&i.EntityID,
		&i.EntityType,
		&i.PlatformName,
		&i.TotalMainSteps,
		&i.TotalSubSteps,
	)
	return i, err
}

const getNotificationByID = `-- name: GetNotificationByID :one
SELECT id, redis_id, notification_group_id, created, read, current_main_step, current_sub_step, main_message, sub_message, error FROM notifications
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNotificationByID(ctx context.Context, id string) (Notification, error) {
	row := q.db.QueryRow(ctx, getNotificationByID, id)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.RedisID,
		&i.NotificationGroupID,
		&i.Created,
		&i.Read,
		&i.CurrentMainStep,
		&i.CurrentSubStep,
		&i.MainMessage,
		&i.SubMessage,
		&i.Error,
	)
	return i, err
}

const getNotificationsByUserID = `-- name: GetNotificationsByUserID :many
SELECT notifications_group.id,
       notifications_group.user_id,
       notifications_group.workflow_id,
       notifications_group.entity_id,
       notifications_group.entity_type,
       notifications_group.total_main_steps,
       notifications_group.total_sub_steps,
       notifications_group.platform_name,
       n.id as notification_id,
       n.redis_id as notification_redis_id,
       n.created as notification_created,
       n.read as notification_read,
       n.current_main_step as notification_current_main_step,
       n.current_sub_step as notification_current_sub_step,
       n.main_message as notification_main_message,
       n.sub_message as notification_sub_message,
       n.error as notification_error
FROM notifications_group
JOIN (
  SELECT DISTINCT ON (notification_group_id) id, redis_id, notification_group_id, created, read, current_main_step, current_sub_step, main_message, sub_message, error
	FROM notifications
	ORDER BY notification_group_id, id desc
) n ON notifications_group.id = n.notification_group_id
WHERE notifications_group.user_id = $1 AND n.redis_id IS NOT NULL
ORDER BY notifications_group.id desc
LIMIT $2 OFFSET $3
`

type GetNotificationsByUserIDParams struct {
	UserID string `json:"userID"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type GetNotificationsByUserIDRow struct {
	ID                          string      `json:"id"`
	UserID                      string      `json:"userID"`
	WorkflowID                  null.String `json:"workflowID"`
	EntityID                    string      `json:"entityID"`
	EntityType                  string      `json:"entityType"`
	TotalMainSteps              null.Int    `json:"totalMainSteps"`
	TotalSubSteps               null.Int    `json:"totalSubSteps"`
	PlatformName                string      `json:"platformName"`
	NotificationID              string      `json:"notificationID"`
	NotificationRedisID         null.String `json:"notificationRedisID"`
	NotificationCreated         null.Time   `json:"notificationCreated"`
	NotificationRead            null.Time   `json:"notificationRead"`
	NotificationCurrentMainStep null.Int    `json:"notificationCurrentMainStep"`
	NotificationCurrentSubStep  null.Int    `json:"notificationCurrentSubStep"`
	NotificationMainMessage     null.String `json:"notificationMainMessage"`
	NotificationSubMessage      null.String `json:"notificationSubMessage"`
	NotificationError           null.Bool   `json:"notificationError"`
}

func (q *Queries) GetNotificationsByUserID(ctx context.Context, arg GetNotificationsByUserIDParams) ([]GetNotificationsByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getNotificationsByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotificationsByUserIDRow
	for rows.Next() {
		var i GetNotificationsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.WorkflowID,
			&i.EntityID,
			&i.EntityType,
			&i.TotalMainSteps,
			&i.TotalSubSteps,
			&i.PlatformName,
			&i.NotificationID,
			&i.NotificationRedisID,
			&i.NotificationCreated,
			&i.NotificationRead,
			&i.NotificationCurrentMainStep,
			&i.NotificationCurrentSubStep,
			&i.NotificationMainMessage,
			&i.NotificationSubMessage,
			&i.NotificationError,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotificationsGroupByID = `-- name: GetNotificationsGroupByID :one
SELECT id, user_id, workflow_id, entity_id, entity_type, platform_name, total_main_steps, total_sub_steps FROM notifications_group
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNotificationsGroupByID(ctx context.Context, id string) (NotificationsGroup, error) {
	row := q.db.QueryRow(ctx, getNotificationsGroupByID, id)
	var i NotificationsGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WorkflowID,
		&i.EntityID,
		&i.EntityType,
		&i.PlatformName,
		&i.TotalMainSteps,
		&i.TotalSubSteps,
	)
	return i, err
}

const getNotificationsGroupByWorkflowID = `-- name: GetNotificationsGroupByWorkflowID :one
SELECT id, user_id, workflow_id, entity_id, entity_type, platform_name, total_main_steps, total_sub_steps FROM notifications_group
WHERE workflow_id = $1 AND user_id = $2
LIMIT 1
`

type GetNotificationsGroupByWorkflowIDParams struct {
	WorkflowID null.String `json:"workflowID"`
	UserID     string      `json:"userID"`
}

func (q *Queries) GetNotificationsGroupByWorkflowID(ctx context.Context, arg GetNotificationsGroupByWorkflowIDParams) (NotificationsGroup, error) {
	row := q.db.QueryRow(ctx, getNotificationsGroupByWorkflowID, arg.WorkflowID, arg.UserID)
	var i NotificationsGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WorkflowID,
		&i.EntityID,
		&i.EntityType,
		&i.PlatformName,
		&i.TotalMainSteps,
		&i.TotalSubSteps,
	)
	return i, err
}

const updateNotificationGroup = `-- name: UpdateNotificationGroup :one
UPDATE notifications_group SET
  user_id = CASE WHEN $1::boolean
    THEN $2::CHAR(26) ELSE user_id END,

  workflow_id = CASE WHEN $3::boolean
    THEN $4::CHAR(64) ELSE workflow_id END,

  entity_id = CASE WHEN $5::boolean
    THEN $6::CHAR(64) ELSE entity_id END,

  entity_type = CASE WHEN $7::boolean
    THEN $8::CHAR(64) ELSE entity_type END,

  total_main_steps = CASE WHEN $9::boolean
    THEN $10::BIGINT ELSE total_main_steps END,

  total_sub_steps = CASE WHEN $11::boolean
    THEN $12::BIGINT ELSE total_sub_steps END,

  platform_name = CASE WHEN $13::boolean
    THEN $14::VARCHAR(300) ELSE platform_name END

WHERE id = $15
RETURNING id, user_id, workflow_id, entity_id, entity_type, platform_name, total_main_steps, total_sub_steps
`

type UpdateNotificationGroupParams struct {
	UserIDDoUpdate         bool   `json:"userIDDoUpdate"`
	UserID                 string `json:"userID"`
	WorkflowIDDoUpdate     bool   `json:"workflowIDDoUpdate"`
	WorkflowID             string `json:"workflowID"`
	EntityIDDoUpdate       bool   `json:"entityIDDoUpdate"`
	EntityID               string `json:"entityID"`
	EntityTypeDoUpdate     bool   `json:"entityTypeDoUpdate"`
	EntityType             string `json:"entityType"`
	TotalMainStepsDoUpdate bool   `json:"totalMainStepsDoUpdate"`
	TotalMainSteps         int64  `json:"totalMainSteps"`
	TotalSubStepsDoUpdate  bool   `json:"totalSubStepsDoUpdate"`
	TotalSubSteps          int64  `json:"totalSubSteps"`
	PlatformNameDoUpdate   bool   `json:"platformNameDoUpdate"`
	PlatformName           string `json:"platformName"`
	ID                     string `json:"id"`
}

func (q *Queries) UpdateNotificationGroup(ctx context.Context, arg UpdateNotificationGroupParams) (NotificationsGroup, error) {
	row := q.db.QueryRow(ctx, updateNotificationGroup,
		arg.UserIDDoUpdate,
		arg.UserID,
		arg.WorkflowIDDoUpdate,
		arg.WorkflowID,
		arg.EntityIDDoUpdate,
		arg.EntityID,
		arg.EntityTypeDoUpdate,
		arg.EntityType,
		arg.TotalMainStepsDoUpdate,
		arg.TotalMainSteps,
		arg.TotalSubStepsDoUpdate,
		arg.TotalSubSteps,
		arg.PlatformNameDoUpdate,
		arg.PlatformName,
		arg.ID,
	)
	var i NotificationsGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WorkflowID,
		&i.EntityID,
		&i.EntityType,
		&i.PlatformName,
		&i.TotalMainSteps,
		&i.TotalSubSteps,
	)
	return i, err
}

const updateRedisIDByNotificationID = `-- name: UpdateRedisIDByNotificationID :exec
UPDATE notifications SET
  redis_id = $1
WHERE id = $2
`

type UpdateRedisIDByNotificationIDParams struct {
	RedisID null.String `json:"redisID"`
	ID      string      `json:"id"`
}

func (q *Queries) UpdateRedisIDByNotificationID(ctx context.Context, arg UpdateRedisIDByNotificationIDParams) error {
	_, err := q.db.Exec(ctx, updateRedisIDByNotificationID, arg.RedisID, arg.ID)
	return err
}
