CREATE UNIQUE INDEX uniq_room_users ON rooms (
  LEAST(initiator_id, participant_id),
  GREATEST(initiator_id, participant_id)
);

CREATE INDEX idx_rooms_initiator_id ON rooms (initiator_id);
CREATE INDEX idx_rooms_participant_id ON rooms (participant_id);