CREATE TABLE IF NOT EXISTS public.sms
(
    id             bigserial NOT NULL,
    notificationId bigserial NOT NULL,
    smsId          text      NOT NULL UNIQUE,
    code           int8      NOT NULL,
    status         text      NULL,
    phrase         text      null,
    cost           float     null,
    number         int8      null,
    created_at     timestamp NULL,
    updated_at     timestamp NULL,
    deleted_at     timestamp NULL,
    CONSTRAINT sms_pkey PRIMARY KEY (id)
)