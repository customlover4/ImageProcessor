create table pictures(
    id serial primary key,
    extension varchar(20) not null,
    status varchar(100) not null default ('not set')
);