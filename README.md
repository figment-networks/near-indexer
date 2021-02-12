# near-indexer ![CI](https://github.com/figment-networks/near-indexer/workflows/CI/badge.svg)

Data indexer and API service for Near protocol networks

*Project is under active development*

## Requirements

- PostgreSQL 10.x+
- Go 1.14+

## Installation

Please see the sections below for all available methods of installation.

### Binary Releases

See [Github Releases](https://github.com/figment-networks/near-indexer/releases) page for details.

### Docker

Pull the official Docker image:

```bash
docker pull figmentnetworks/near-indexer
```

### Golang

```bash
go get -u github.com/figment-networks/near-indexer
```

## Usage

```bash
$ ./near-indexer --help

Usage of ./near-indexer:
  -cmd string
    	Command to run
  -config string
    	Path to config
  -v	Show application version
```

Executing commands:

```bash
near-indexer -c path/to/config.json -cmd=COMMAND
```

Available commands:

| Name      | Description
|-----------|-----------------------------------------------------
| `status`  | Print out current indexer and node status
| `migrate` | Perform database migration
| `sync`    | Run a one-time indexer sync (for testing purposes)
| `worker`  | Start the indexer sync worker
| `server`  | Start the indexer API server
| `reset`   | Reset the database

## Configuration

You can configure the service using a config file or environment variables.

### Config File

Example:

```json
{
  "app_env": "production",
  "rpc_endpoint": "http://YOUR_NODE_RPC_IP:PORT",
  "server_addr": "127.0.0.1",
  "server_port": 8081,
  "database_url": "postgres://user:pass@host/dbname?sslmode=mode",
  "sync_interval": "500ms",
  "cleanup_interval": "10m",
  "cleanup_threshold": 3600,
  "start_height": 0,
  "rollbar_token": "rollbar access token",
  "rollbar_namespace": "rollbar app name"
}
```

### Environment Variables

| Name                 | Description             | Default Value
|----------------------|-------------------------|-----------------
| `APP_ENV`            | Application environment | `development`
| `DATABASE_URL`       | PostgreSQL database URL | REQUIRED
| `NEAR_RPC_ENDPOINT`  | Near RPC endpoint       | REQUIRED
| `START_HEIGHT`       | Initial start height    | optional, will use genesis if 0
| `SERVER_ADDR`        | Server listen addr      | `0.0.0.0`
| `SERVER_PORT`        | Server listen port      | `8081`
| `SYNC_INTERVAL`      | Data sync interval      | `500ms`
| `CLEANUP_INTERVAL`   | Data cleanup interval   | `10m`
| `CLEANUP_THRESHOLD`  | Max number of heights   | `3600`
| `DEBUG`              | Turn on debugging mode  | `false`
| `ROLLBAR_TOKEN`      | Rollbar access token    |
| `ROLLBACK_NAMESPACE` | Rollbar app name        |

## Running Application

Once you have created a database and specified all configuration options, you
need to migrate the database. You can do that by running the command below:

```bash
near-indexer -config path/to/config.json -cmd=migrate
```

Perform the indexer check:

```bash
near-indexer -config path/to/config.josn -cmd=status
```

Perform the initial sync:

```bash
near-indexer -config path/to/config.josn -cmd=sync
```

If previous steps did not produce any errors you can start the indexer worker:

```bash
near-indexer -config path/to/config.json -cmd=worker
```

Start the API server:

```bash
near-indexer -config path/to/config.json -cmd=server
```

## API Reference

| Method | Path                            | Description
|--------|---------------------------------|------------------------------------
| GET    | /                               | See all available endpoints
| GET    | /health                         | Healthcheck endpoint
| GET    | /status                         | App version info and sync status
| GET    | /height                         | Current indexed blockchain height
| GET    | /block                          | Get latest block
| GET    | /blocks                         | Blocks search
| GET    | /blocks/:hash                   | Block details by ID or Hash
| GET    | /block_stats                    | Block times stats for a time bucket
| GET    | /block_times                    | Block average times
| GET    | /block_times_interval           | Block creation stats
| GET    | /epochs                         | Get list of epochs
| GET    | /epochs/:id                     | Epoch details by ID
| GET    | /validators                     | List of chain validators
| GET    | /validators/:id/epochs          | Validator Epochs performance by ID
| GET    | /validators/:id/events          | Validator Events by ID
| GET    | /transactions                   | List of transactions
| GET    | /transactions/:id               | Get transaction details
| GET    | /accounts/:id                   | Account details by ID or Key
| GET    | /delegations/:id                | Account delegations by ID
| GET    | /events                         | List of Events

## License

Apache License v2.0
