MYSQL = service mysql start

ifeq ($(OS),Windows_NT)
	MYSQL = "/c/Program Files/MySQL/MySQL Server 8.0/bin/mysqld"
endif

all: db go

db:
	$(MYSQL)&

go:
	cd src && go run main.go&