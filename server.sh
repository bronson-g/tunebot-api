cd ~/tunebot-api
apt update
apt upgrade -y
apt install -y make golang mysql-server
git pull
mysql --defaults-extra-file=/etc/mysql/debian.cnf < data/main.sql
make