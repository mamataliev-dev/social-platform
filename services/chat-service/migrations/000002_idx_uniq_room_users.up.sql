CREATE UNIQUE INDEX uniq_room_users ON rooms (
  LEAST(user_a_id, user_b_id),
  GREATEST(user_a_id, user_b_id)
);

CREATE INDEX idx_rooms_user_a ON rooms (user_a_id);
CREATE INDEX idx_rooms_user_b ON rooms (user_b_id);
