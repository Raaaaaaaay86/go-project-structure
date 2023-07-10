drop table if exists user_roles;

alter table users add column role_id int;
alter table users add constraint users_role_id_fkey foreign key (role_id) references roles(id);
