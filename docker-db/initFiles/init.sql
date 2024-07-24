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
