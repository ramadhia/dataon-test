## Deployment

To deploy this project, please run with docker compose to start the SQL Server 2022

```bash
docker compose -f docker-compose.yaml up
```

### Init Project Golang
- Please run this command in the BE project
```bash
go mod download
```

#### Migrate database 
- Please create `bosnetdb` database 
- Then run the migration first before run the server BE
```bash
make migrate
```

#### Run the server, the server will run at `:3000`
```bash
make run-api
```

- Build image
```bash
make docker
```