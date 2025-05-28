# 🧪 Log Processing Pipeline — Pet Project in Go
This pet project is built as a way to learn the Go programming language and to practice designing a microservice architecture using Kafka, MySQL, Redis, and Docker.

## 🔧 Idea & Structure
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
└── docker-compose.yml
```

Each app is an independent Go service with its own go.mod, go.sum, internal logic and a main file located under cmd/.


## ⚠️ Development Status
***The project is under active development. Structure, code, and documentation may evolve.***

The goal is to build a working prototype with a clear architecture and potential for scaling.