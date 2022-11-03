FROM postgres:10.3

# copy up.sql and it is the first thing to be run
COPY up.sql /docker-entrypoint-initdb.d/1.sql

# executes command "postgres"
CMD ["postgres"]