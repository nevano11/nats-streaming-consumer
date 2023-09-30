create table delivery
(
    id      serial primary key,
    name    text,
    phone   text,
    zip     text,
    city    text,
    address text,
    region  text,
    email   text
);

create table payment
(
    id            serial primary key,
    transaction   text,
    request_id    text,
    currency      text,
    provider      text,
    amount        int,
    payment_dt    int,
    bank          text,
    delivery_cost int,
    goods_total   int,
    custom_fee    int
);

create table model
(
    id                 serial primary key,
    order_uid          text,
    track_number       text,
    entry              text,
    delivery_id        int,
    payment_id         int,
    locale             text,
    internal_signature text,
    customer_id        text,
    delivery_service   text,
    shardkey           text,
    sm_id              int,
    date_created       date,
    oof_shard          text,

    CONSTRAINT fk_delivery
        FOREIGN KEY (delivery_id)
            REFERENCES delivery (id),
    CONSTRAINT fk_payment
        FOREIGN KEY (payment_id)
            REFERENCES payment (id)
);

create table item
(
    id           serial primary key,
    model_id     int,
    chrt_id      int,
    track_number text,
    price        int,
    rid          text,
    name         text,
    sale         int,
    size         text,
    total_price  int,
    nm_id        int,
    brand        text,
    status       int,

    CONSTRAINT fk_model
        FOREIGN KEY (model_id)
            REFERENCES model (id)
);
