# near-indexer

Data indexer and API service for Near protocol networks

## Requirements

- PostgreSQL 10.x+
- Go 1.14

## Installation

*Not available yet*

## Configuration

You can configure the service using either a config file or environment variables.

### Config File

Example:

```json
{
  "app_env": "production",
  "rpc_endpoint": "http://YOUR_NODE_RPC_IP:PORT",
  "server_addr": "127.0.0.1",
  "server_port": 8081,
  "database_url": "postgres://user:pass@host/dbname",
  "sync_interval": "500ms",
  "cleanup_interval": "10m"
}
```

### Environment Variables

| Name                | Description             | Default Value
|---------------------|-------------------------|-----------------
| `APP_ENV`           | Application environment | `development`
| `DATABASE_URL`      | PostgreSQL database URL | REQUIRED
| `NEAR_RPC_ENDPOINT` | Near RPC endpoint       | REQUIRED
| `START_HEIGHT`      | Initial start height    | optional, will use genesis if blank
| `SERVER_ADDR`       | Server listen addr      | `0.0.0.0`
| `SERVER_PORT`       | Server listen port      | `8081`
| `SYNC_INTERVAL`     | Data sync interval      | `500ms`
| `CLEANUP_INTERVAL`  | Data cleanup interval   | `10m`
| `DEBUG`             | Turn on debugging mode  | `false`

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

If previous steps did not produce any error you can start the indexer worker:

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
| GET    | /health                         | Healthcheck endpoint
| GET    | /status                         | App version info and sync status
| GET    | /height                         | Current indexed blockchain height
| GET    | /blocks                         | Blocks search
| GET    | /blocks/:hash                   | Block details by ID or Hash
| GET    | /block_times                    | Block times stats
| GET    | /block_times_interval           | Block creation stats
| GET    | /validators                     | List of chain validators
| GET    | /validator_times_interval       | Active validator stats
| GET    | /accounts/:id                   | Account details by ID or Key

## License

Apache License v2.0
