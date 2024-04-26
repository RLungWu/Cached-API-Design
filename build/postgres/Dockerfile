FROM postgres:latest
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=mydatabase
COPY init.sql /docker-entrypoint-initdb.d/
