# E-Wallet

An e-wallet built using Golang and provides a comprehensive solution for managing digital transactions with a focus on the Indonesian market. Leveraging Midtrans for payment processing, users can effortlessly top up their wallets using QRIS, Gopay, and Bank Transfer, all in Indonesian currency.
## Features

- **Currency Top-Up**: Top up wallet using QRIS, Gopay, or Bank Transfer via Midtrans.
- **Centralized Logging**: Utilize Filebeat, Kibana, and Elasticsearch to centralize and monitor logs from the main e-wallet system, e-wallet-queue, and e-wallet-scheduler.
- **Email Notifications**: Receive OTP codes for secure transactions and weekly summaries of spending and earning directly into email inbox.
- **Live Notifications**: Receive live notification when top-ups and transfers are successful using Server-Sent Events (SSE).
- **Performance Optimization**: Leverage Redis for efficient caching and queuing.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Docker

## How to Run the Application
1. First, build the docker image of [E-Wallet-Queue](https://github.com/reyhanyogs/e-wallet-queue) and [E-Wallet-Scheduler](https://github.com/reyhanyogs/e-wallet-scheduler).
2. Insert your Midtrans key in ```main.env```
3. Then, build ```docker-compose.yaml``` using:
   ```bash
   make composeup
   ```
4. After all services are running, use http://localhost:8080 to access the API.

## Logging and Monitoring
To access the centralized logging system:
1. Navigate to Kibana at http://localhost:5601.
2. Log in to Kibana using the following credentials
   ```bash
   kibana
   Nqy0fHREsStfASF
   ```
4. Configure the index pattern for Elasticsearch.
5. Explore the logs for the main e-wallet system, e-wallet-queue, and e-wallet-scheduler.
