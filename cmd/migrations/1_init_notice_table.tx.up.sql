CREATE TABLE IF NOT EXISTS public.notifications
(
    id             bigserial    NOT NULL,
    type           varchar(255) NULL,
    status         text         NULL,
    status_message text         NULL,
    raw            json,
    created_at     timestamp    NULL,
    updated_at     timestamp    NULL,
    deleted_at     timestamp    NULL,
    debug          bool default false,
    CONSTRAINT notifications_pkey PRIMARY KEY (id)
)