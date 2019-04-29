#!/bin/bash

sudo docker pull mysql/mysql-server:latest

sudo docker run -dit -p 3306:3306 --name=parkbot_db -e MYSQL_DATABASE=parkbot -e MYSQL_ROOT_PASSWORD=test-pw mysql/mysql-server:latest