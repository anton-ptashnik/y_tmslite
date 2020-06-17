# TMS lite

_It is an initial document ver._

## Project design

ER diagram is located on: https://drive.google.com/drive/folders/12auxnZaop03MY_Ly653L4gV_t1Sh2J2Z?usp=sharing

## Project structure
Project is organized into several packages each taking care of its point of responsibility 

- root

Responsible for request handlers assigment and initial application setup
- middleware

Responsible for extracting data from HTTP request and passing it for further procession

- persistence

Responsible for data storage and extraction from a permanent storage   