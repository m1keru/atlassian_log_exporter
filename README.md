# Atlassian Log Exporter

This project is a Go application that exports logs from Atlassian Jira using the Jira API. It fetches audit records and saves the state to resume from the last fetched record.
It is designed to be runned from systemd timer so it will log all events directlry into journalctl.

## Prerequisites

- Go 1.16 or higher
- Jira API credentials (email, token, and endpoint)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/atlassian_log_exporter.git
    cd atlassian_log_exporter
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Build binary:

    ```sh
    go build .
    ```

## Usage

1. Set the required environment variables:

    ```sh
    export JIRA_API_EMAIL=your-email@example.com
    export JIRA_API_TOKEN=your-api-token
    export JIRA_API_ENDPOINT=https://your-domain.atlassian.net
    ```

2. Run the application:

    ```sh
    go run main.go
    ```

3. Optional flags:
    - `-jira_api_email`: Jira API Email (can be set with `JIRA_API_EMAIL` env variable)
    - `-jira_api_token`: Jira API Token (can be set with `JIRA_API_TOKEN` env variable)
    - `-jira_api_endpoint`: Jira API Endpoint (can be set with `JIRA_API_ENDPOINT` env variable)
    - `-debug`: Enable debug mode

### Example

```sh
./atlassian_log_exporter -debug
```

### Logging

The application uses `zap` for logging. Logs are printed to the console. In debug mode, more detailed logs are provided.

### State Management

The application saves its state in a JSON file (`jira_state.json`) to resume fetching logs from where it left off. If the state file is not found, it starts from the beginning.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
