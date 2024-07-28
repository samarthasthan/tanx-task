# Price Alert Application

Welcome to the Price Alert Application project! This application allows users to set price alerts for cryptocurrencies. When the price of a cryptocurrency reaches the user's target price, an alert is triggered and an email is sent to the user.

## Technologies Used

- **Backend**: Go
- **Database**: MySQL
- **Caching**: Redis
- **Messaging**: RabbitMQ
- **Real-Time Data**: Binance WebSocket for real-time price updates
- **Email Service**: Brevo Emails
- **Containerization**: Docker and Docker Compose

## Features

- **Create Alert**: API endpoint to set price alerts.
- **Delete Alert**: API endpoint to remove price alerts.
- **Fetch Alerts**: API endpoint to fetch all alerts with pagination and filtering options.
- **User Authentication**: JWT tokens for secure access to endpoints.
- **Real-Time Price Updates**: Uses Binance WebSocket to get real-time price updates.
- **Email Notifications**: Sends an email when the price alert is triggered.
- **Caching**: Utilizes Redis for caching alert data.

## High-Level Design

![Sign Up OTP](/samples/Untitled-2024-07-20-2114.png)

## Postman Documentation

For detailed API documentation, you can view the Postman collection [here](https://documenter.getpostman.com/view/19782195/2sA3kaBeRD).

## Sample Screenshots

Here are some screenshots:

| Sign Up OTP Email                     | Trigger Email                               |
| ------------------------------------- | ------------------------------------------- |
| ![Sign Up OTP](/samples/IMG_0740.PNG) | ![Trigger Email](/samples/IMG_0742%202.PNG) |

## Getting Started

To get started with the project, follow these steps:

1. **Clone the repository**:

   ```bash
   git clone https://github.com/samarthasthan/tanx-task.git
   cd tanx-task
   ```

2. **Create a `.env` file**:
   Create a `.env` file in the build/compose directory of the project with the necessary environment variables. Example:

   ```env
   MYSQL_ROOT_PASSWORD=<your_mysql_root_password>
   MYSQL_PORT=3306
   REDIS_PORT=6379
   RABBITMQ_DEFAULT_USER=<your_rabbitmq_user>
   RABBITMQ_DEFAULT_PASS=<your_rabbitmq_password>
   RABBITMQ_DEFAULT_PORT=5672
   SMTP_SERVER=<your_smtp_server>
   SMTP_PORT=<your_smtp_port>
   SMTP_LOGIN=<your_smtp_login>
   SMTP_PASSWORD=<your_smtp_password>
   JWT_SECRET=<your_jwt_secret>
   REST_API_PORT=8000
   ```

3. **Install Docker and Docker Compose**:
   Ensure Docker and Docker Compose are installed on your system.

4. **Start the services**:

   ```bash
   docker-compose up -d
   ```

   or use Makefile

   ```bash
   Make up
   ```

5. **Migrate Database**;

   ```bash
   make migrate-up
   ```

6. **Access the API**:
   - **Create Alert**: `POST /alerts/create/`
   - **Delete Alert**: `DELETE /alerts/delete/`
   - **Fetch Alerts**: `GET /alerts/all`

## API Endpoints

### Create Alert

- **Endpoint**: `POST /alerts/create/`
- **Description**: Create a new price alert.
- **Request Body**:
  ```json
  {
    "currency": "BTCUSDT",
    "price": 67341.99
  }
  ```

### Delete Alert

- **Endpoint**: `DELETE /alerts/delete/`
- **Description**: Delete an existing price alert.
- **Request Body**:
  ```json
  {
    "alert_id": "<alert_id>"
  }
  ```

### Fetch Alerts

- **Endpoint**: `GET /alerts/all`
- **Description**: Fetch all alerts

## Docker Configuration

The project uses Docker Compose for containerization. Below is the `docker-compose.yml` configuration for the services:

```yaml
version: "3.8"

services:
  # Databases
  mysql:
    image: mysql:latest
    container_name: mysql
    networks:
      - tanx
    ports:
      - "${MYSQL_PORT}:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=tanx

  redis:
    image: redis:latest
    container_name: redis
    networks:
      - tanx
    ports:
      - "${REDIS_PORT}:6379"

  # Message Brokers
  rabbitmq:
    image: rabbitmq:latest
    container_name: rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: ${RABBITMQ_DEFAULT_USER}
      RABBITMQ_DEFAULT_PASS: ${RABBITMQ_DEFAULT_PASS}
    ports:
      - ${RABBITMQ_DEFAULT_PORT}:5672
    networks:
      - tanx
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 10s
      retries: 5
      start_period: 30s

  # Services
  email:
    build:
      context: ../../
      dockerfile: ./build/docker/email/Dockerfile
    container_name: email
    networks:
      - tanx
    environment:
      - RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER}
      - RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS}
      - RABBITMQ_DEFAULT_PORT=${RABBITMQ_DEFAULT_PORT}
      - RABBITMQ_DEFAULT_HOST=${RABBITMQ_DEFAULT_HOST}
      - SMTP_SERVER=${SMTP_SERVER}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_LOGIN=${SMTP_LOGIN}
      - SMTP_PASSWORD=${SMTP_PASSWORD}
    command: ["./app"]
    depends_on:
      rabbitmq:
        condition: service_healthy

  alert:
    build:
      context: ../../
      dockerfile: ./build/docker/alert/Dockerfile
    container_name: alert
    networks:
      - tanx
    environment:
      - RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER}
      - RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS}
      - RABBITMQ_DEFAULT_PORT=${RABBITMQ_DEFAULT_PORT}
      - RABBITMQ_DEFAULT_HOST=${RABBITMQ_DEFAULT_HOST}
      - MYSQL_PORT=${MYSQL_PORT}
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_HOST=${MYSQL_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_HOST=${REDIS_HOST}
    command: ["./app"]
    depends_on:
      rabbitmq:
        condition: service_healthy

  tanx:
    build:
      context: ../../
      dockerfile: ./build/docker/tanx/Dockerfile
    container_name: tanx
    networks:
      - tanx
    environment:
      - REST_API_PORT=${REST_API_PORT}
      - MYSQL_PORT=${MYSQL_PORT}
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_HOST=${MYSQL_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_HOST=${REDIS_HOST}
      - RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER}
      - RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS}
      - RABBITMQ_DEFAULT_PORT=${RABBITMQ_DEFAULT_PORT}
      - RABBITMQ_DEFAULT_HOST=${RABBITMQ_DEFAULT_HOST}
      - JWT_SECRET=${JWT_SECRET}
    ports:
      - "${REST_API_PORT}:${REST_API_PORT}"
    command: ["./app"]
    depends_on:
      rabbitmq:
        condition: service_healthy

networks:
  tanx:
    driver: bridge
```

### Solution for Sending Alerts

The application uses the following approach for sending alerts:

1. **Real-Time Price Updates**: Connects to Binance WebSocket for real-time price updates.
2. **Triggering Alerts**: Compares current prices with user-set alert prices and triggers alerts when prices match.
3. **Sending Emails**: Uses RabbitMQ for message brokering and an email service to send notification emails when an alert is triggered.

## Contact

For any questions or assistance, please reach out:

- **Twitter**: [@samarthasthan](https://twitter.com/samarthasthan)
- **Email**: [samarthasthan27@gmail.com](mailto:samarthasthan27@gmail.com)
