CREATE TABLE messages (
  id bigserial NOT NULL PRIMARY KEY,
  sender_id int NOT NULL REFERENCES users(id),
  sender_name character varying(100) NOT NULL,
  receiver_id int NOT NULL REFERENCES users(id),
  receiver_name character varying(100) NOT NULL,
  user_message character varying(100) NOT NULL,
  created_at timestamp with time zone DEFAULT now() NOT NULL
);
