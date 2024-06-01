FROM mysql:latest

COPY init.sql /docker-entrypoint-initdb.d/

RUN echo "Hey ! I'm the Dockerfile !" && \
    chown -R mysql:mysql /docker-entrypoint-initdb.d/ && \
    chmod 777 /docker-entrypoint-initdb.d/