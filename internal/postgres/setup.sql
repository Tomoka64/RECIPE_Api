drop table ratings;
drop table recipes;
drop table users;

create table users (
  id         serial primary key,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null default CURRENT_TIMESTAMP
);

create table recipes (
  id         serial primary key,
  topic      varchar(255),
  description text,
  user_id    integer references users(id),
  prep_time integer,
  difficulty integer,
  vegetarian boolean,
  is_deleted boolean default false,
  created_at timestamp not null default CURRENT_TIMESTAMP
);

create table ratings (
  rate integer,
  user_id    integer references users(id),
  recipe_id integer references recipes(id),
  primary key (user_id, recipe_id)
);

