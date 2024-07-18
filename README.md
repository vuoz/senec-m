# SENEC Monitor
This is a monitor accessing the unofficial senec api.
It retrieves the data for your photovoltaic system  
and saves it to a postgres database.  

# Endpoints

All the data from the Senec-Api
```
/full 
```

Data for a specific Timestamp
```
/data?ts=1234
```

Latest data from the local-Api updated every 10 Seconds
```
/localLatest 
```

Websocket Connection that gets every update of the local-Api & weather in real time
```
/subscribe
```



## Instructions

Create database and table with docker
- [DB Setup with docker compose](https://github.com/zvup/senec-monitor-db)


Clone the repo


Add a .env file to your directory:

    SENEC_PASS=
    SENEC_USER=
    DB_HOST=
    DB_PORT=
    DB_PASS=
    DB_USER=
    DB_NAME= 
    LONGITUDE=
    LATITUDE=
    SENEC_IP=

Run via:

    go run .

Or build to exe:
    
    go build .

Visit on:
```shell
localhost:4000/
```

Building for docker arm64 requires this command:
```shell
docker buildx build --platform linux/arm64 -f Dockerfile -t your_username/test-repo:senec-monitor-latest --push  .
```

## Todos
- [ ] add ci/cd pipeline for automatic container builds


