-- name: GetNotificationByID :one
SELECT * FROM notifications
WHERE id = $1
LIMIT 1;

-- name: GetNotificationsGroupByID :one
SELECT * FROM notifications_group
WHERE id = $1
LIMIT 1;

-- name: GetNotificationsGroupByWorkflowID :one
SELECT * FROM notifications_group
WHERE workflow_id = $1 AND user_id = $2
LIMIT 1;

-- name: CreateNotificationsGroup :one
INSERT INTO notifications_group (
  user_id, workflow_id, entity_id, entity_type,
  total_main_steps, total_sub_steps
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: CreateNotification :one
INSERT INTO notifications (
  notification_group_id, read, current_main_step,
  current_sub_step, main_message, sub_message
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateRedisIDByNotificationID :exec
UPDATE notifications SET
  redis_id = $1
WHERE id = $2;
