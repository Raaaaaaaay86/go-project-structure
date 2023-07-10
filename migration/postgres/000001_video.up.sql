create table roles (
    id int primary key,
    name varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);

create table users (
    id serial primary key,
    role_id int not null,
    username varchar(255) not null,
    password varchar(255) not null,
    email varchar(255) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp,
    foreign key (role_id) references roles(id)
);