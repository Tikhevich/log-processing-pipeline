# 🧪 Log Processing Pipeline — Pet Project in Go
This pet project is built as a way to learn the Go programming language and to practice designing a microservice architecture using Kafka, MySQL, Redis, and Docker.

## Project Architecture
![Project Architecture](./assets/diagrams/architecture.svg)

The project implements a simple log processing pipeline consisting of the following components:

### log-generator-app
Generates random HTTP logs and sends them via HTTP to the log-receiver service.

### log-receiver-api-app
Accepts logs through a REST API, processes them, and publishes messages to Kafka.

### log-processor-app
Consumes messages from Kafka, filters out logs with status code 200, and saves only 404 and 500 logs into a MySQL database.

### ui-api-app
Exposes a REST API that allows users to:

Retrieve logs for specific time ranges (hour, day, week)

Get statistics on errors, traffic, and latency
These results are cached in Redis for a limited time. If not found in cache, data is fetched from MySQL.

## 🗂 Project Structure

```lua
├── LICENSE
├── README.md
├── apps
│   ├── log-generator-app
│   ├── log-processor-app
│   ├── log-receiver-api-app
│   └── ui-api-app
├── assets
│   └── diagrams
└── docker-compose.yml
```
Each app is an independent Go service with its own go.mod, go.sum, internal logic and a main file located under cmd/.

## Interaction Diagram
![Interaction Diagram](./assets/diagrams/interaction.svg)

### Frame 1 - Receiving logs
`log-generator-app` generates logs, which are sent via API request to `log-receiver-api-app`. They are then processed and sent to the `Kafka` broker.

### Frame 2 - Processing logs
`log-processor-app` receives messages from the `Kafka` broker, processes them, discarding logs with a status code of 200 OK, and stores everything else in a `MySQL` database.

### Frame 3 - Get stats from cache
`ui-api-app` receives an API request to get statistics, then checks if these statistics are cached in `Redis`. If so, we immediately give them to the user without delay

### Frame 4 - Get stats from DB and add to cache
`ui-api-app` receives an API request to get statistics, then checks if these statistics are cached in `Redis`. If not, a request is sent to the `MySQL` database to get statistics, after that, it is additionally cached in `Redis` and returned to the user

### Frame 5 - Get logs
`ui-api-app` receives an API request to retrieve logs, collects the appropriate filter if specified, then a request is made to `MySQL` database and the prepared data with the possibility of pagination is returned to the user.


## API description for ui-api-app 

#### `GET /logs?...`

Returns all available logs. It is possible to filter by the fields specified below. Pagination is also available

**Parametrs:**

| Parametr | Type | Description |
| - | - | - |
| `ip` | string  | IP-addres |
| `status` | string | Status Code |
| `method` | string | Method |
| `path` | string | Path |
| `from` | string | Date from |
| `to` | string | Date to |
| `limit` | string | Limit |
| `offset` | string | Offset |

**Request example:**
```
GET http://localhost:8081/logs?status=500&limit=1&method=post
```

**Response example:**
```json
{
    "data": [
        {
            "id": 59432,
            "Timestamp": "2025-06-10T10:41:40Z",
            "IP": "184.116.4.244",
            "Method": "POST",
            "Path": "/products",
            "Status": 500,
            "LatencyMs": 971
        }
    ],
    "meta": {
        "total": 4208,
        "limit": 1,
        "offset": 0
    }
}
```

---

#### `GET /stats?...`

Returns statistics for the specified period and type

**Parametrs:**

| Parametr | Type | Description |
| - | - | - |
| `type` | string  | Type of stats, one of  `latency`, `errors`, `traffic` |
| `range` | string | One of `hour`, `day`, `week` |


**Request example:**
```
GET http://localhost:8081/stats?range=day&type=latency
```

**Response example:**
```json
// Response for type "latency"
{
    "type": "latency",
    "data": {
        "avg_latency": 503.3166,
        "max_latency": 999
    }
}
// Response for type "traffic"
{
    "type": "traffic",
    "data": {
        "total_request": 20384,
        "unique_ips": 20384
    }
}
// Response for type "errors"
{
    "type": "errors",
    "data": {
        "404": 13484,
        "500": 6900
    }
}
```

## ⚠️ Development Status
***The project is under active development. Structure, code, and documentation may evolve.***

The goal is to build a working prototype with a clear architecture and potential for scaling.