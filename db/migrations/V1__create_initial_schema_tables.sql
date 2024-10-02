-- This migration script creates the initial schema for the UATZ API database.
-- The schema includes tables for devices, device handlers, device webhooks, and webhook messages.
-- Version: 1
-- Author: @caiofabiodearaujo
-- Date: 2024-10-01

-- Schema 'whatsmeow' is used for the whatsmeow library.
-- Schema 'uatzapi' is used for the UATZ API.
create schema whatsmeow;
create schema uatzapi;

-- Create the sequence for the 'device' table if it does not exist
create sequence if not exists uatzapi.device_id_seq;

-- Table 'device'
-- This table stores information about each device. The device has a unique ID, WhatsApp ID, 
-- push name, business name, active status, and a timestamp for when it was created.
create table if not exists uatzapi.device (
  id bigint primary key not null default nextval('uatzapi.device_id_seq'::regclass),
  whatsapp_id character varying not null, -- WhatsApp ID associated with the device
  push_name character varying not null, -- Name pushed to the device
  business_name character varying, -- Optional business name associated with the device
  active boolean not null, -- Indicates if the device is currently active
  created_at timestamp with time zone not null default now() -- Timestamp when the record was created
);

-- Creating a unique index on 'whatsapp_id' to ensure no two devices have the same WhatsApp ID
create unique index if not exists device_whatsapp_id_key 
on uatzapi.device using btree (whatsapp_id);

-- Create the sequence for the 'device' table if it does not exist
create sequence if not exists uatzapi.device_handler_id_seq;

-- Table 'device_handler'
-- This table stores handlers associated with each device. It contains a reference to the device,
-- a boolean indicating whether the handler is active, and timestamps for activation and deactivation.
create table if not exists uatzapi.device_handler (
  id bigint primary key not null default nextval('uatzapi.device_handler_id_seq'::regclass),
  device_id bigint not null, -- Foreign key referencing the 'device' table
  active boolean not null, -- Indicates if the handler is active
  active_at timestamp with time zone not null, -- Timestamp when the handler was activated
  inactive_at timestamp with time zone, -- Optional timestamp when the handler was deactivated
  constraint fk_device_handler_device_id foreign key (device_id) references uatzapi.device (id) -- Foreign key constraint
);

-- Creating an index on 'device_id' in the 'device_handler' table to optimize queries by device reference
create index if not exists device_handler_device_id_idx 
on uatzapi.device_handler (device_id);

-- Create the sequence for the 'device' table if it does not exist
create sequence if not exists uatzapi.device_webhook_id_seq;

-- Table 'device_webhook'
-- This table stores webhook information for each device, including the URL, status, and timestamp.
create table if not exists uatzapi.device_webhook (
  id bigint primary key not null default nextval('uatzapi.device_webhook_id_seq'::regclass),
  device_id bigint not null, -- Foreign key referencing the 'device' table
  webhook_url character varying not null, -- Webhook URL associated with the device
  active boolean not null, -- Indicates if the webhook is currently active
  timestamp timestamp with time zone not null default now(), -- Timestamp when the webhook was created
  constraint fk_device_webhook_device_id foreign key (device_id) references uatzapi.device (id) -- Foreign key constraint
);

-- Creating an index on 'device_id' in the 'device_webhook' table to optimize queries by device reference
create index if not exists device_webhook_device_id_idx 
on uatzapi.device_webhook (device_id);

-- Create the sequence for the 'device' table if it does not exist
create sequence if not exists uatzapi.webhook_message_id_seq;

-- Table 'webhook_message'
-- This table stores information about messages sent to the webhook, including the message content,
-- the response, the webhook URL used, the response code, and a timestamp.
create table if not exists uatzapi.webhook_message (
  id bigint primary key not null default nextval('uatzapi.webhook_message_id_seq'::regclass),
  device_id bigint not null, -- Foreign key referencing the 'device' table
  message character varying not null, -- The message sent to the webhook
  response character varying not null, -- The response received from the webhook
  webhook_url character varying not null, -- The webhook URL to which the message was sent
  code_response bigint not null, -- HTTP status code or custom response code from the webhook
  timestamp timestamp with time zone not null default now(), -- Timestamp when the message was logged
  constraint fk_webhook_message_device_id foreign key (device_id) references uatzapi.device (id) -- Foreign key constraint
);

-- Creating an index on 'device_id' in the 'webhook_message' table to optimize queries by device reference
create index if not exists webhook_message_device_id_idx 
on uatzapi.webhook_message (device_id);

-- Creating an index on 'webhook_url' in the 'webhook_message' table to optimize queries by webhook URL
create index if not exists webhook_message_webhook_url_idx 
on uatzapi.webhook_message (webhook_url);

-- Creating an index on 'code_response' in the 'webhook_message' table to optimize queries by response code
create index if not exists webhook_message_code_response_idx 
on uatzapi.webhook_message (code_response);
