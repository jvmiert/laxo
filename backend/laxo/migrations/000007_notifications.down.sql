BEGIN;
ALTER TABLE notifications DROP CONSTRAINT fk_group_notification;
ALTER TABLE notifications_group DROP CONSTRAINT fk_user_notification_group;
DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS notifications_group;
COMMIT;
