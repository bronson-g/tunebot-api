MYSQL = service mysql start

ifeq ($(OS),Windows_NT)
	MYSQL = "/c/Program Files/MySQL/MySQL Server 8.0/bin/mysqld"
endif

all:
	cd src/; $(MYSQL); go run main.go
