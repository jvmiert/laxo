-- name: GetNotificationsByUserID :many
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
  SELECT DISTINCT ON (notification_group_id) *
	FROM notifications
	ORDER BY notification_group_id, id desc
) n ON notifications_group.id = n.notification_group_id
WHERE notifications_group.user_id = $1 AND n.redis_id IS NOT NULL
ORDER BY notifications_group.id desc
LIMIT $2 OFFSET $3;

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
  total_main_steps, total_sub_steps, platform_name
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
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
    THEN @total_sub_steps::BIGINT ELSE total_sub_steps END,

  platform_name = CASE WHEN @platform_name_do_update::boolean
    THEN @platform_name::VARCHAR(300) ELSE platform_name END

WHERE id = @id
RETURNING *;

-- name: CreateNotification :one
INSERT INTO notifications (
  notification_group_id, read, current_main_step,
  current_sub_step, main_message, sub_message, error
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateRedisIDByNotificationID :exec
UPDATE notifications SET
  redis_id = $1
WHERE id = $2;
