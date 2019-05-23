FROM mysql/mysql-server:5.5

WORKDIR /parkbot/

COPY services/ $WORKDIR/services

ENTRYPOINT ["mysql"]