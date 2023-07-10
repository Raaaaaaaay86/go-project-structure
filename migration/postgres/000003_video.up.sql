create table video_posts
(
    id          serial primary key,
    title       varchar(255) not null,
    description varchar(255),
    uuid  varchar(255) not null,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now(),
    deleted_at  timestamp
);