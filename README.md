# Golang - Shop

## Documentation.
Documentation was done using Postman you can check it through the [Link](https://documenter.getpostman.com/view/30922024/2sAY4sj4JW#ee889508-9a0f-4bab-b862-6422c7b154ea), it contains examples with it.

Unfortunately Postman does not support publishing Websocket requests and examples on free plan.

## Running project.
### Local.
**Note**: for this project you are required to have make functional on your pc so you can use Makefile commands.
use the given commands to run the project or you could use the commands they do run inside the [Makefile](./Makefile).
```make
make run
```

for running the tests 
```make
make test
```

### Backup & Restore.
- **Backup.**

In order to be able to create a backup you must be connection already to your database and use the following command:
```make
make backup ARGS="--create"
```
backup data will be created in a `.json` file in your main directory under folder `backups`.
every backup file is created with a unique string which indicate to the exact time that has been created at.

- **Restore.**

You must restore your data from a `.json` file made from the command made specifically for creating backup or with any json file that satisfies the schemas of the database.
in order to make a restore you must choose the backup file name and must be inside `backups` folder.
example:
you have this file:
- `backups`
    - `data_backup_20241025_015150.json`

to restore this file run this command:
```make
make backup ARGS="--restore --dbname=golang_shop_test --password=123456 --backupFile=data_backup_20241025_015150.json"
```
you may need to specify more flags to your database connection some flags are set to default.

**Flags**:
- port: database port **default is 3306**.
- user: database user **default is root**.
- host: database host **default is 127.0.0.1**.
- dbname: the database name you want to move your data to.
- password: database user password.
- backupFile: the backup file name as it is exactly inside `backups` file.


### Running Tests

Use the following command to run the tests:

```make
make test
```

### Local with Docker.
you can run your project in the 3 following environments.
make sure to check [variables environment example](.env.example) with your keys (especially for cloudinary)
1. **Production.**
```make
make up-prod
```
you can also pass --build flag to the command if you wish to build the image before run:

```make
make up-prod ARGS="--build"
```

to stop and remove all the running containers run:
```make
make down-prod
```

2. **Development.**

```make
make up-dev
```
to stop and remove all the running containers run:
```make
make down-dev
```

3. **Testing.**
There is a database setup for testing it has the original database name + "_test", you must ensure to migrate all your data to the test database using the entry point that meant for backup and restore.

>Methodology:
- Integration test:
    - we are using integration test to test the returned response from the endpoint, it must satisfies the returned type as well as the returned status code to pass.
- Unit testing:
    - we are using the validation functionality to the expected payloads (DTO's).

```make
make up-test
```

to stop and remove all the running containers run:
```make
make down-test
```

You can also use the file [curl requests](./curlreqs.sh) to test the behavior of [Nginx config](./nginx/default.conf) if it works as expected or not, just run the following command in bash:
```bash
./curlreqs.sh
```
In the current set nginx can hold 10 burst so you would need to increase the amount of requests in the file or run the command twice to see the error 429 is returned.

### Running with Docker-swarm.
We have made swarm only for the production environment.
Before running with docker-swarm, you need to type the database name inside [docker-compose](docker-compose.yml) file and set it implicitly there without importing it from `.env` file or you can use external secrets as work around.
If you wish to try running swarm with the other environments you need to specify in their `docker-compose.(env name).yml` file.
To create a stack and run the services files for production use the following command:
```make
make swarm-up
```

To remove the stack and stop the services use the following command:
```make
make swarm-down
```

### Creating fake data.
To create a fake data, so far there was need to create fake for **users**, **products** and **reviews**, you can simply go to the seed main file [main.go](./cmd/seed/main.go) and set the flag to true manually if u wish to create to them fake data.
There was no need to add flags to it, as it only needed to be run once.

## Database Models.
![Database models](https://res.cloudinary.com/doxhxgz2g/image/upload/v1730256130/Others/golang_shop_yquhuz.jpg)

## Improvements

- Separate the Websocket to an independent server so it can be scaled independently from the rest of the project.
- We have messaging system in our project therefore we must add a mechanism that allow us to restrict the user from sending messages.
- Add More tests to have full coverage.
- Add an endpoint to black list the signed-out users, can be either by using redis or using a map.