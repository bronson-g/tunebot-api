cd ~/tunebot-api
sudo apt update
sudo apt upgrade -y
sudo apt install -y make golang mysql-server
git pull
mysql --defaults-extra-file=/etc/mysql/debian.cnf < data/main.sql
sudo make