# SENEC Monitor
This is a monitor accessing the unofficial senec api.
It retrieves the data for your photovoltaic system  
and saves it to a postgres database.  

## Utility
- Websocket connection can be used by embedded devices to get real time data [Check out my implementation with ESP32](https://github.com/vuoz/senec-c)
- Creates a record of yield data that can be used for analysis and prediction ( also saves weather data )
- Combines current photovoltaic system data with weather data to allow for accurate prediction / forecast
- Could also be developed into a more sophisticated web interface in combination with weather data, graphs etc

## Under the hood
- Retrieves yield data for every full day ( unofficial senec api )
- Gets current data every 10s ( your local api )
- Obtains weather data for you location on a hourly basis
- Retrieves Daily weather forecast every 24 hours and saves it to be combined with total yield values for later prediction purposes
- Provides an Api ( for specific data ) and a websocket connection for frequent updates ( weather, current power levels,  total yield )

## Endpoints

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

Copy  ```compose.yaml``` to your directory and run:
```shell
docker compose pull
```
Also copy the ```/docker-db/initFiles/init.sql``` to the directory containing the compose file. 
Please perserve the directory structure so that the init.sql file lies iside the ```/docker-db/initFiles/``` directory to make sure the db   
is intialized properly.
This is important as it initializes the database with the required tables
```shell
Createa a ```.docker.env``` with the following values
```shell
# your my-senec login password
SENEC_PASS=

# your my-senec login username
SENEC_USER= 

# db host default should be  "docker.host.internal"
DB_HOST=docker.host.internal

# db port default should be 6000
DB_PORT=6000

# db password default is yourpass
DB_PASS=yourpassj

# db user default is myuser
DB_USER=myuser

# just search the coordinates of your location and put them here
LONGITUDE=
LATITUDE=

# the ip of your local senec dashboard
SENEC_IP=
```
Then run 
```shell
docker compose up
```

**Please note that currently this docker compose file utilizes watchtower to check for container updates.   
If you do not plan on continuously working on this project yourself you can just remove everything associated with watchtower from the compose file.**

Now you should be able to visit the service on:
```shell
localhost:4000/
``` 




## Instructions (dev)

Create database and table with docker (to have a separate database)
- [DB Setup with docker compose](https://github.com/vuoz/senec-monitor-db)

Clone the repo


Add a .env file to your directory:
```shell
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
```
Run via:

    go run .

Or build to exe:
    
    go build .

Visit on:
```shell
localhost:4000/
```

Building for docker arm64 requires this command if you want to run it as a container
```shell
docker buildx build --platform linux/arm64 -f Dockerfile -t your_username/test-repo:senec-monitor-latest --push  .
```

## Todos
- [x] add ci/cd pipeline for automatic container builds


