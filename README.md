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
Install postgres

Create database and table
```sql
CREATE TABLE data(
    timestamp bigint,
    generated double precision,
    consumption double precision,
    gridexport double precision,
    accuexport double precision,
    accuimport double precision,
    gridimport double precision,
    PRIMARY KEY (timestamp)
);
```

```sql
CREATE TABLE local_data
(
    stat_state integer,
    gui_bat_data_power double precision,
    gui_inverter_power double precision,
    gui_house_pow double precision,
    gui_grid_pow double precision,
    gui_bat_data_fuel_charge double precision,
    gui_charging_info integer,
    gui_boosting_info integer,
    ts bigint
);

```
```sql
CREATE TABLE weather_data (
    WeatherCode DOUBLE PRECISION,
    Temperature2MMax DOUBLE PRECISION,
    Temperature2MMin DOUBLE PRECISION,
    Sunrise TEXT,
    Sunset TEXT,
    DaylightDuration DOUBLE PRECISION,
    SunshineDuration DOUBLE PRECISION,
    UvIndexMax DOUBLE PRECISION,
    UvIndexClearSkyMax DOUBLE PRECISION,
    Date TEXT
);

```



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
docker buildx build --platform linux/arm64 -f Dockerfile -t username/test-repo:senec-monitor-latest --push  .
```



