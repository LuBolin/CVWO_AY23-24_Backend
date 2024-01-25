# CVWO 23/24 Project Golang Backend & PostgreSQL Database
Lu Bolin


## Hosted demo:
The entire experience is accessible at [LocaleLookout](https://localelookout.onrender.com/).



## Setup Instructions

Remember to setup the frontend too. It can be found at [LuBolin/CVWO_AY23-24_Frontend (github.com)](https://github.com/LuBolin/CVWO_AY23-24_Frontend).

#### Prerequisites:
- NodeJS 20.10.0
- npm 10.2.3

Forward compatability is expected as there were no major changes around these versions.
However, they have not been tested.



**Clone this repository**
```bash
git clone https://github.com/LuBolin/CVWO_AY23-24_Backend.git
```


### PostgreSQL hosting

Set up a PostgreSQL server of version >=15.

Take down the Hostname, Port, Database name, and a set of Username and Password. They will be used later by the backend service.

In the backend GitHub, there is a PostgreSQL/db_schema.sql. It is the schema of the database's tables. Execute it in PostgreSQL. (Note: in the case of hosting on render.com, a shell was not provided. In such situations, connect to the psql server yourself from bash, and then create the tables within your database).



### Backend hosting
Clone this repo.
Make sure you have golang >= 1.21.4.
Backward compatibility has not been tested.

#### Environment Variables:
These can be set in the .env file.
When hosting online, for example in the case of render.com, you can set them from the hosting dashboard.
Set these in the .env or in your environment variables, before building and running the golang server.

- DB_HOST: database host name. If local, this will be localhost.
- DB_PORT: database port. The default port of PostgreSQL is 5432.
- DB_NAME: name of the actual database.
- DB_USER: username.
- DB_PASSWORD: password.
- FRONTEND_IP: address of your frontend hosting server. This is provided for the sake of CORS.
  For the onrender.com example, this is https://localelookout.onrender.com.
- HMAC_SECRET: arbitrary secret string. This is used in hashing.

If you are on windows, you can add .exe incase your cmd gives errors.
The running would fail if the database has not been setup.
#### Build & Run command: 
```bash
go install
go build -o app
./app
```