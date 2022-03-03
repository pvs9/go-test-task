CREATE TABLE todo_items
(
    id            serial       not null unique,
    description   varchar(255) not null,
    due_date      timestamp    not null
);