CREATE TABLE users (
  id bigserial NOT NULL PRIMARY KEY,
  user_name character varying(50) NOT NULL,
  public_key character varying(50) NOT NULL,
  user_password character varying(50) NOT NULL,
  created_at timestamp with time zone DEFAULT now() NOT NULL
);