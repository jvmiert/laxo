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

-- name: UpdateNotificationGroup :one
UPDATE notifications_group SET
  user_id = CASE WHEN @user_id_do_update::boolean
    THEN @user_id::CHAR(26) ELSE user_id END,

  workflow_id = CASE WHEN @workflow_id_do_update::boolean
    THEN @workflow_id::CHAR(64) ELSE workflow_id END,

  entity_id = CASE WHEN @entity_id_do_update::boolean
    THEN @entity_id::CHAR(64) ELSE entity_id END,

  entity_type = CASE WHEN @entity_type_do_update::boolean
    THEN @entity_type::CHAR(64) ELSE entity_type END,

  total_main_steps = CASE WHEN @total_main_steps_do_update::boolean
    THEN @total_main_steps::BIGINT ELSE total_main_steps END,

  total_sub_steps = CASE WHEN @total_sub_steps_do_update::boolean
    THEN @total_sub_steps::BIGINT ELSE total_sub_steps END

WHERE id = @id
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
