**Docker Image**

`docker pull yunuskilicdev/suss`

The Easiest way to use my system is using docker container.

Sample Requests and Responses

I used 3 different location to demonstrate results.

London 213.123.58.0

NY Brooklyn 4.14.4.0

NJ Atlantic City 12.151.181.192

**Request#1**:

`curl --location --request POST 'localhost:5000' \
--header 'Content-Type: application/json' \
--data-raw '{
"username": "yunus",
"unix_timestamp": 1619437515,
"event_uuid": "1",
"ip_address": "4.14.4.0"
}'`

**Response#1**:

`{
"currentGeo": {
"lat": 40.7252,
"lon": -73.944,
"radius": 100 },
"travelToCurrentGeoSuspicious": false,
"travelFromCurrentGeoSuspicious": false,
"precedingIpAccess": {
"lat": 0,
"lon": 0,
"radius": 0,
"speed": 0,
"ip": "",
"timestamp": 0 },
"subsequentIpAccess": {
"lat": 0,
"lon": 0,
"radius": 0,
"speed": 0,
"ip": "",
"timestamp": 0 } }
`

**Request#2**:

`curl --location --request POST 'localhost:5000' \
--header 'Content-Type: application/json' \
--data-raw '{
"username": "yunus",
"unix_timestamp": 1619438355,
"event_uuid": "3",
"ip_address": "213.123.58.0"
}'
`

**Response#2**:

`{
"currentGeo": {
"lat": 51.5096,
"lon": -0.0972,
"radius": 200 },
"travelToCurrentGeoSuspicious": true,
"travelFromCurrentGeoSuspicious": false,
"precedingIpAccess": {
"lat": 40.7252,
"lon": -73.944,
"radius": 100,
"speed": 14023.874946643242,
"ip": "4.14.4.0",
"timestamp": 1619437515 },
"subsequentIpAccess": {
"lat": 0,
"lon": 0,
"radius": 0,
"speed": 0,
"ip": "",
"timestamp": 0 } }
`

**Request#3**:

`curl --location --request POST 'localhost:5000' \
--header 'Content-Type: application/json' \
--data-raw '{
"username": "yunus",
"unix_timestamp": 1619439315,
"event_uuid": "2",
"ip_address": "12.151.181.192"
}'`

**Response#3**:

`{
"currentGeo": {
"lat": 39.363,
"lon": -74.4324,
"radius": 50 },
"travelToCurrentGeoSuspicious": true,
"travelFromCurrentGeoSuspicious": false,
"precedingIpAccess": {
"lat": 51.5096,
"lon": -0.0972,
"radius": 200,
"speed": 12685.743130072375,
"ip": "213.123.58.0",
"timestamp": 1619438355 },
"subsequentIpAccess": {
"lat": 0,
"lon": 0,
"radius": 0,
"speed": 0,
"ip": "",
"timestamp": 0 } }`

**Project Structure**

This project calculates suspicious activity with respect its IP address.

In order to find coordinates of the IP address, I used  https://dev.maxmind.com/geoip/geoip2/geolite2

MaxMind City database has CIDR blocks and coordinates of locations.

Using the geoip2-csv-converter tool provided by MaxMind, we will add the integer range representation of each network to our CSV.

Now searching becomes easier with start and end ip address.

If you want to create your own local structure, follow steps below:

https://blog.maxmind.com/2020/02/13/the-easy-way-to-use-maxmind-geoip-with-redshift/

Now our csv is ready.

I used sqlite for database operations.

Create a **suss.sqlite**

Below scripts will be used to create table and import csv into it.

create table city_blocks
(
network TEXT,
network_start_integer INTEGER,
network_last_integer INTEGER,
geoname_id INTEGER,
registered_country_geoname_id INTEGER,
represented_country_geoname_id TEXT,
is_anonymous_proxy INTEGER,
is_satellite_provider INTEGER,
postal_code TEXT,
latitude REAL,
longitude REAL,
accuracy_radius INTEGER
);

.mode csv

.import --skip 1 GeoLite2-City-Blocks-IPv4-with-ranges.csv CityBlock

create index city_blocks_network_last_integer_index
on city_blocks (network_last_integer);

create index city_blocks_network_start_integer_index
on city_blocks (network_start_integer);

Our working environment is ready for usage.

You can run your application.

**Dependencies**

go-sqlite3 for sqlite

gorm for orm. With orm our codebase becomes understandable and clear. 

**Assumed Cases**

If a user is different locations at the same second, I assumed that speed is equal to distance.

If a request has subsequent request, I assumed there is a network delay that previous request comes later.

If distance between two location is bigger than sum of accuracy radius of two request, I subtracted sum of accuracy radius to handle uncertainty


