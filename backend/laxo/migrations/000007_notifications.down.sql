BEGIN;
ALTER TABLE notifications DROP CONSTRAINT fk_group_notification;
ALTER TABLE notifications_group DROP CONSTRAINT fk_user_notification_group;
ALTER TABLE notifications_group DROP CONSTRAINT notifications_group_platform_name;
DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS notifications_group;
COMMIT;
