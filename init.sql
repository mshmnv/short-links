ALTER USER postgres PASSWORD 'postgres';
ALTER USER postgres WITH CREATEDB;
DROP DATABASE short_db IF EXISTS;
CREATE DATABASE short_db;