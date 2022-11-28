# go-restwebsockets-platzi

Curso de Go Avanzado: REST y WebSockets: https://platzi.com/cursos/go-rest-websockets/

# Setup project

1- Run the following command at /database folder level , this will build the image for the postgres db:

```
  docker build . -t platzi-ws-rest
```

2- In order to run the db execute the following command:

```
  docker run -p 54321:5432 platzi-ws-rest
```

3- To run the project (API) execute the following command:

```
  go run main.go
```

# Re-run DB:

Inside of the `database` folder run the following command:

```
 docker build . -t platzi-ws-rest
```
