FROM postgres:latest
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=mydatabase
ENV POSTGRES_USER=postgres
COPY init.sql /docker-entrypoint-initdb.d/
