# Price Alert Application

Welcome to the Price Alert Application project! This application allows users to set price alerts for cryptocurrencies. When the price of a cryptocurrency reaches the user's target price, an alert is triggered and an email is sent to the user.

## Technologies Used

- **Backend**: Go
- **Database**: MySQL
- **Caching**: Redis
- **Messaging**: RabbitMQ
- **Real-Time Data**: Binance WebSocket for real-time price updates
- **Email Service**: Brevo Emails
- **Reverse Proxy**: NGINX
- **Containerization**: Docker and Docker Compose

## Live URL

You can access the live version of the Price Alert Application APIs at [https://task.samarthasthan.com](https://task.samarthasthan.com).

**Note**: Please be aware that due to the server's location in London, you may experience some latency.

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


## Tolerance for Price Alerts

To demonstrate the tolerance functionality, the application uses an absolute tolerance of 90 units. This means that if the real-time price of a cryptocurrency is within 90 units of the target price set in the alert, the alert will be triggered. For example, if an alert is set for a price of $50,000, the alert will trigger if the price is between $49,910 and $50,090.

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
   docker-compose -f ./build/compose/docker-compose.yaml up -d
   ```

   or use Makefile

   ```bash
   make up
   ```

5. **Access the API**:
   - **Create Alert**: `POST /alerts/create/`
   - **Delete Alert**: `DELETE /alerts/delete/`
   - **Fetch Alerts**: `GET /alerts`

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

- **Endpoint**: `GET /alerts`
- **Description**: Fetch all alerts

## Docker Configuration

The project uses Docker Compose for containerization. Below is the `docker-compose.yml` configuration for the services:

```yaml
services:
  # Databases
  mysql:
    image: mysql:latest
    container_name: mysql
    hostname: mysql
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
    hostname: redis
    networks:
      - tanx
    ports:
      - "${REDIS_PORT}:6379"

  # Message Brokers
  rabbitmq:
    image: rabbitmq:latest
    container_name: rabbitmq
    hostname: rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: ${RABBITMQ_DEFAULT_USER}
      RABBITMQ_DEFAULT_PASS: ${RABBITMQ_DEFAULT_PASS}
    ports:
      - "${RABBITMQ_DEFAULT_PORT}:5672"
    networks:
      - tanx
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 30s
      timeout: 10s
      retries: 5

  # Services
  email:
    build:
      context: ../../
      dockerfile: ./build/docker/email/Dockerfile
    container_name: email
    hostname: email
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

  tanx:
    build:
      context: ../../
      dockerfile: ./build/docker/tanx/Dockerfile
    container_name: tanx
    hostname: tanx
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
    command: ["sh", "-c", "make migrate-up && ./app"]
    depends_on:
      rabbitmq:
        condition: service_healthy

  alert:
    build:
      context: ../../
      dockerfile: ./build/docker/alert/Dockerfile
    container_name: alert
    hostname: alert
    networks:
      - tanx
    environment:
      - MYSQL_PORT=${MYSQL_PORT}
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_HOST=${MYSQL_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_HOST=${REDIS_HOST}
      - RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER}
      - RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS}
      - RABBITMQ_DEFAULT_PORT=${RABBITMQ_DEFAULT_PORT}
      - RABBITMQ_DEFAULT_HOST=${RABBITMQ_DEFAULT_HOST}
    command: ["./app"]
    depends_on:
      rabbitmq:
        condition: service_healthy

networks:
  tanx:
    driver: bridge
```

## Solution for Sending Alerts

The application uses the following approach for sending alerts:

1. **Real-Time Price Updates**: Connects to Binance WebSocket for real-time price updates.
2. **Triggering Alerts**: Compares current prices with user-set alert prices and triggers alerts when prices match.
3. **Sending Emails**: Uses RabbitMQ for message brokering and an email service to send notification emails when an alert is triggered.

## Possible Improvements

While this application provides core functionalities, several enhancements could be considered for future development:

- **Centralized Logging**: Implement centralized logging to monitor and debug issues more effectively.
- **gRPC**: Integrate gRPC for efficient inter-service communication.
- **Metrics and Monitoring**: Add metrics and monitoring for better visibility into system performance and health.
- **Scalability Enhancements**: Consider architectural changes for improved scalability.

Due to time constraints, these improvements have been skipped in the current implementation.

## Contact

For any questions or assistance, please reach out:

- **Twitter**: [@samarthasthan](https://twitter.com/samarthasthan)
- **Email**: [samarthasthan27@gmail.com](mailto:samarthasthan27@gmail.com)
