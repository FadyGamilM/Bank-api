# Bank-api
Bank API for money transfer secured via JWT token.

### A command for running postgres via docker
```cmd
> docker run -d --name go_psql_container -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=fady -e POSTGRES_PASSWORD=gobankingpassword -e POSTGRES_DB=bankDB -p 2345:5432 postgres
```

### To access the spinned up container 
```cmd
> docker exec -it [containerID] psql -U [YourUserName] -d [YourDBName]
```