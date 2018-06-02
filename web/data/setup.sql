-- drop table comments;
drop table recipes;
drop table sessions;
drop table users;

create table users (
  id         serial primary key,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null
);

create table sessions (
  id         serial primary key,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null
);

create table recipes (
  id         serial primary key,
  topic      text,
  user_id    integer references users(id),
  prep_time integer,
  difficulty integer,
  vegetarian boolean,
  created_at timestamp not null
);

create table rating (
  rate integer,
  user_id    integer references users(id),
  recipe_id integer references recipes(id),
  primary key (user_id, recipe_id),
);
--
-- create table comments (
--   id         serial primary key,
--   uuid       varchar(64) not null unique,
--   body       text,
--   user_id    integer references users(id),
--   thread_id  integer references threads(id),
--   created_at timestamp not null
-- );
