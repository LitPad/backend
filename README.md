# LITPAD V1 (WORK IN PROGRESS)

![alt text](https://github.com/LitPad/backend/blob/main/display/fiber.png?raw=true)


#### FIBER DOCS: [Documentation](https://docs.gofiber.io/)
#### GORM DOCS: [Documentation](https://gorm.io/docs/index.html)
#### PG ADMIN: [Documentation](https://pgadmin.org) 


## How to run locally

* Download this repo or run: 
```bash
    $ git clone git@github.com:LitPad/backend.git
```

#### In the root directory:
- Install all dependencies
```bash
    $ go install github.com/cosmtrek/air@latest 
    $ go mod download
```
- Create an `.env` file and copy the contents from the `.env.example` to the file and set the respective values. A postgres database can be created with PG ADMIN or psql

- Run Locally
```bash
    $ air
```

- Run With Docker
```bash
    $ docker-compose up --build -d --remove-orphans
```
OR
```bash
    $ make build
```

- Test Coverage
```bash
    $ go test ./tests -v -count=1
```
OR
```bash
    $ make test
```