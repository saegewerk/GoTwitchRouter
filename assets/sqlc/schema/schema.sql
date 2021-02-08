create extension if not exists "uuid-ossp";
CREATE TABLE twitch_user
(
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name varchar(255) ,
    nick varchar(255) NOT NULL UNIQUE ,
    accessLevel int NOT NULL,
    created_at timestamp    NOT NULL,
    updated_at timestamp    NOT NULL
);
CREATE TABLE msg
(
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    msg varchar(500) NOT NULL ,
    msg_id varchar(500) NOT NULL ,
    event varchar(100) NOT NULL,
    fk_user uuid REFERENCES twitch_user(id) ON DELETE CASCADE,
    created_at timestamp    NOT NULL
)