alter table users drop column role_id;

create table user_roles  (
    user_id int,
    role_id int,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(role_id) REFERENCES roles(id)
)