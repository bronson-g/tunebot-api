all: db go

db:
	"/c/Program Files/MySQL/MySQL Server 8.0/bin/mysqld"&

go:
	cd src && go run main.go&