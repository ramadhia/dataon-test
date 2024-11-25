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
---

### 3. Backend

#### 3.1 Deploy the Project Using Docker Compose

To deploy the project and start all services (Backend Server, Worker Server, PostgreSQL, RabbitMQ), run the following command:

```bash
docker compose -f docker-compose.yaml up
```

This command will pull the necessary images and start the services defined in the `docker-compose.yaml` file, including:

- **Backend API**: Running at port `40001`
- **PostgreSQL**: Database service

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
- This proses will produces the data seed for testing 

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
- - You can access this API `localhost:40001/organizations` or `localhost:15000/organizations` to fetch the structure organization data
- **PostgreSQL**: Accessible on `localhost:45432` (username: `root`, password: `password`)

---

### Example Response
- you can download the example response file on this project ``example-response.json``
