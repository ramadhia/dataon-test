# DataOn Test
## Prerequisites

Before you begin, ensure you have the following tools installed on your machine:

- **Make**: Make sure you have `Make` installed. If you don’t have it, you can follow the instructions [here](https://www.gnu.org/software/make/) or use a package manager to install it.
- **Docker**: This project uses Docker to run services in isolated containers. If you haven’t installed Docker yet, you can get it from [here](https://www.docker.com/get-started).
- **Docker Compose**: Ensure that you have Docker Compose installed to run multiple containers. You can find installation instructions [here](https://docs.docker.com/compose/install/).

## Getting Started

Follow these steps to get the project up and running on your local machine.

### 1. Initialize the Golang Project

- Rename the `env.example` file to `.env` and configure your environment variables as needed.

- Run the following command in the project directory to download Go dependencies:

```bash
go mod download
```

### 2. Algo Test

#### Test 1
Run the following command for **Test 1**:

```bash
go run github.com/ramadhia/dataon-test/cmd/dataon algo-test1
```

#### Test 2
Run the following command for **Test 2**, where you can specify the total amount and the amount paid:

```bash
go run github.com/ramadhia/dataon-test/cmd/dataon algo-test2 --total-belanja=52312 --uang-dibayar=100000
```

#### Test 3
Run the following command for **Test 3**:

```bash
go run github.com/ramadhia/dataon-test/cmd/dataon algo-test3
```

#### Test 4
Run the following command for **Test 4**, where you specify the number of vacation days, leave duration, and join date:

```bash
go run github.com/ramadhia/dataon-test/cmd/dataon algo-test4 --cuti-bersama=7 --cuti-durasi=1 --join-date=2021-05-01 --join-date=2021-07-05
```

---

### 3. Backend

#### 3.1 Deploy the Project Using Docker Compose

To deploy the project and start all services (Backend Server, Worker Server, PostgreSQL, RabbitMQ), run the following command:

```bash
docker compose -f docker-compose.yaml up
```

This command will pull the necessary images and start the services defined in the `docker-compose.yaml` file, including:

- **Backend API**: Running at port `40001`
- **Worker Server**: Responsible for handling background tasks
- **PostgreSQL**: Database service
- **RabbitMQ**: Message queue service, Dashboard queue running at `localhost:15672`

#### 3.2 Initialize the Golang Project

- Rename the `env.example` file to `.env` and configure your environment variables as needed.

- Run the following command in the project directory to download Go dependencies:

```bash
go mod download
```

#### 3.3 Migrate the Database

- Create a PostgreSQL database. For this project, we’ll use `dbrafli` as the database name. You can use a PostgreSQL client to create it manually or automate the process.
- Once the database is created, run the database migrations to set up the schema:
- Or you can check the path of migrations file at ``database/migration``

```bash
make migrate
```


#### 3.4 **Optional: Run the Backend Server Manually**

If you prefer to run the Backend API server manually (instead of using Docker Compose), you can start the server using the following command. It will run on `localhost:40001`:

```bash
make run-api
```

#### 3.5 **Optional: Run the Worker Server Manually**

To start the Worker server manually (instead of using Docker Compose), run the following command. The Worker server is responsible for processing tasks from RabbitMQ:

```bash
make run-worker
```

#### 3.6 Accessing Services

- **Backend API**: Accessible on `http://localhost:40001`
- **RabbitMQ Management UI**: Accessible on `http://localhost:15672` (username: `guest`, password: `guest`)
- **PostgreSQL**: Accessible on `localhost:45432` (username: `root`, password: `password`)

---

### Postman Collection
- you can download the Postman collection file on this project ``MNC.postman_collection.json``