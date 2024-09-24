# Atlassian Admin API Event Log Exporter

This Go application fetches events from the Atlassian Admin API, processes them, and logs the results. It supports pagination, rate limiting, and state persistence.

## Features

- Fetches events from the Atlassian Admin API
- Supports custom date ranges for event retrieval
- Handles API rate limiting
- Logs events to console and optionally to a file
- Persists the last processed event date to resume from where it left off
- Configurable via command-line flags and environment variables

## Prerequisites

- Go 1.x or higher
- Atlassian Admin API Token
- Atlassian Organization ID

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/m1keru/atlassian_log_exporter.git
   ```

2. Navigate to the project directory:

   ```sh
   cd atlassian_log_exporter
   ```

3. Install dependencies:

   ```sh
   go mod tidy
   ```

4. Build:

   ```sh
   go build
   ```

## Usage

Run the application with the following command:

```sh
./atlassian_log_exporter --help
```

### Flags

- `-api_user_agent`: API User Agent (default "curl/7.54.0")
- `-api_token`: Atlassian Admin API Token (can also be set via ATLASSIAN_ADMIN_API_TOKEN environment variable)
- `-from`: (Optional) From date (RFC3339 format)
- `-org_id`: Organization ID (can also be set via ATLASSIAN_ORGID environment variable)
- `-log-to-file`: (Optional) Enable logging to file
- `-log-file`: (Optional) Path to log file (default "log.txt")
- `-debug`: Enable debug mode
- `-query`: Query to filter the events
- `-sleep`: Sleep time in milliseconds between requests (default 200)

### Environment Variables

- `ATLASSIAN_ADMIN_API_TOKEN`: Atlassian Admin API Token
- `ATLASSIAN_ORGID`: Atlassian Organization ID

## Example

```sh
./atlassian_log_exporter -api_token=your_api_token -org_id=your_org_id -from=2023-09-01T00:00:00Z -log-to-file -debug
```

or

```sh
ATLASSIAN_ADMIN_API_TOKEN=123 ATLASSIAN_ORGID=123-123-123 ./atlassian_log_exporter
```

This command will fetch events from September 1, 2023, log the results to both console and file, and enable debug mode.

## State Persistence

The application saves the timestamp of the last processed event in a file named `jira_state.json`. This allows the application to resume from where it left off in subsequent runs.

## Error Handling

The application handles various error scenarios, including API rate limiting. When the rate limit is exceeded, it will wait for the specified time before retrying.

## Logging

Logs are output to the console by default. If the `-log-to-file` flag is set, logs will also be written to the specified file (default: "log.txt").

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache License Version 2.0, January 2004
