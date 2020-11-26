create role mailowl with password 'mailowl' NOSUPERUSER NOCREATEDB NOCREATEROLE INHERIT LOGIN;

create database mailowl owner mailowl;
