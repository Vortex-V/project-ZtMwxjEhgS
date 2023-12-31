create table if not exists accounts
(
    id              bigint generated by default as identity,
    username        varchar(255)                   not null,
    password        varchar(255)                   not null,
    type            integer          default 2     not null,
    status          integer          default 1     not null,
    is_need_relogin boolean          default false not null,
    balance         double precision default 0     not null,
    created_at      timestamp        default CURRENT_TIMESTAMP,
    updated_at      timestamp        default CURRENT_TIMESTAMP
    );

alter table accounts
    add constraint accounts_pk
        primary key (id);

alter table accounts
    add constraint accounts_username_uk
        unique (username);

create table if not exists transports
(
    id            bigint generated by default as identity,
    account_id    bigint                              not null,
    can_be_rented boolean   default false             not null,
    type          varchar(255)                        not null,
    model         varchar(255)                        not null,
    color         varchar(255)                        not null,
    identifier    varchar(255)                        not null,
    description   text,
    latitude      double precision                    not null,
    longitude     double precision                    not null,
    minute_price  double precision,
    day_price     double precision,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    updated_at    timestamp default CURRENT_TIMESTAMP not null
    );

alter table transports
    add constraint transports_pk
        primary key (id);

alter table transports
    add constraint transports_accounts_id_fk
        foreign key (account_id) references accounts
            on update cascade on delete cascade;

create table if not exists rental
(
    id            bigint generated by default as identity,
    account_id    bigint                                           not null,
    type          varchar(255)                                     not null,
    transport_id  bigint                                           not null,
    time_start    timestamp    default CURRENT_TIMESTAMP,
    time_end      timestamp,
    price_of_unit double precision                                 not null,
    final_price   double precision,
    status        varchar(255) default 'active'::character varying not null,
    created_at    timestamp    default CURRENT_TIMESTAMP           not null,
    updated_at    timestamp    default CURRENT_TIMESTAMP           not null
    );

alter table rental
    add constraint rental_pk
        primary key (id);

alter table rental
    add constraint rental_accounts_id_fk
        foreign key (account_id) references accounts
            on update cascade on delete cascade;

alter table rental
    add constraint rental_transports_id_fk
        foreign key (transport_id) references transports
            on update cascade on delete cascade;

INSERT INTO accounts (username, password, type, status, is_need_relogin, balance)
VALUES ('admin', '$2a$16$0wgQ4IOHv7ZgFtcJFx.JvO6zAk04Jio.eiQPZ2wmm5iM.0VItasY6', 1, 1, false, 300000);