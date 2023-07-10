insert into roles (id, name) values (1, 'SUPER_ADMIN'), (2, 'ADMIN'), (3, 'USER');

alter table users add constraint users_username_unique unique (username);
