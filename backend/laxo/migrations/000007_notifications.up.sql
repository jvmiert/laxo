BEGIN;
CREATE TABLE IF NOT EXISTS notifications_group(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  user_id CHAR(26) NOT NULL,
  workflow_id CHAR(64),
  entity_id CHAR(26) NOT NULL,
  entity_type CHAR(32) NOT NULL,
  total_main_steps BIGINT,
  total_sub_steps BIGINT,
  CONSTRAINT fk_user_notification_group FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS notifications(
  id CHAR(26) DEFAULT ulid_create() NOT NULL PRIMARY KEY,
  redis_id CHAR(18),
  notification_group_id CHAR(26) NOT NULL,
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  read TIMESTAMP WITH TIME ZONE,
  current_main_step BIGINT,
  current_sub_step BIGINT,
  main_message CHAR(64),
  sub_message CHAR(64),
  CONSTRAINT fk_group_notification FOREIGN KEY(notification_group_id) REFERENCES notifications_group(id)
);
COMMIT;
