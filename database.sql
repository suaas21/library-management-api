create database library_management;

create table users (
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL,
    image varchar(100) NOT NULL,
    mail varchar(50) NOT NULL,
    password varchar(300) NOT NULL,
    phone_no varchar(15) NOT NULL,
    user_type varchar(20) NOT NULL
);

create table books (
    id SERIAL PRIMARY KEY,
    book_name varchar(100) NOT NULL,
    author varchar(50) NOT NULL,
    available BOOLEAN DEFAULT TRUE
);

create table book_loan_histories (
    id SERIAL PRIMARY KEY,
    book_id integer NOT NULL,
    user_id integer NOT NULL ,
    purchased_date varchar(50) NOT NULL,
    return_date varchar(50) NOT NULL,
    returned BOOLEAN DEFAULT False
);


create table book_loan_requests (
    id SERIAL PRIMARY KEY,
    book_id integer NOT NULL,
    user_id integer NOT NULL ,
    status varchar(20) DEFAULT 'Pending'
);


