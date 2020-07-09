# TMS lite

The project serves education purposes for learning backend-side development.

Up-to-date project version is located at `release` branch.
## Usage
API doc: https://app.swaggerhub.com/apis/iDeveloper34/y_tmslite

The project is deployed at: https://y-tmslite.herokuapp.com

In case of need to run locally, open a terminal in project root dir and use commands:

```set DB_URL=<actual_url>```

```set PORT=<desired_port>```

```go run .```

## Project design

ER diagram is located on: https://drive.google.com/drive/folders/12auxnZaop03MY_Ly653L4gV_t1Sh2J2Z?usp=sharing

## Project structure
Project is organized into several packages each taking care of its point of responsibility 

- root

Responsible for request handlers assigment and initial application setup
- middleware

Responsible for extracting data from HTTP request and passing it for further procession

- service

Takes data from middleware and does needed procession before passing to a DB

- persistence

Responsible for data storage and extraction from a permanent storage   