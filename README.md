# GOQ (Goku)

```
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/
```

## Installation
To install GOQ, run:

```bash
go install github.com/marianozunino/goq@latest
```

Or download the binary from the [GitHub releases](https://github.com/marianozunino/goq/releases)

## Usage
```
Usage:
  goq [command]

Available Commands:
  configure   Generate a sample configuration file for goq.
  dump        Dump messages from a RabbitMQ queue to a file.
  monitor     Monitor RabbitMQ messages using routing keys and a temporary queue.
  update      Update the goq tool to the latest available version.
  version     Display the current version of the goq tool.

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -s, --amqps                      Use AMQPS instead of AMQP
      --config string              config file (default "/home/forbi/.config/goq/goq.yaml")
  -e, --exchange string            RabbitMQ exchange name
  -x, --exclude-patterns strings   Exclude messages containing these patterns
  -m, --file-mode string           File mode (append or overwrite, only valid for file writer) (default "overwrite")
  -f, --full-message               Print full message
  -h, --help                       help for goq
  -i, --include-patterns strings   Include messages containing these patterns
  -j, --json-filter string         JSON filter expression
  -z, --max-message-size int       Maximum message size in bytes (default -1)
  -o, --output string              Output file name (required when writer is 'file')
  -p, --pretty-print               Pretty print JSON messages
  -R, --regex-filter string        Regex pattern to filter messages
  -k, --skip-tls-verify            Skip TLS certificate verification (insecure)
  -u, --url string                 RabbitMQ URL (e.g., localhost:5672)
  -v, --virtualhost string         RabbitMQ virtual host
  -w, --writer string              Output writer type (console or file) (default "file")

Use "goq [command] --help" for more information about a command.
```

### Example Commands
1. **Dump messages from a local RabbitMQ server**:
   ```bash
   goq dump -u localhost:5672 -q my_queue
   ```

2. **Monitor messages with specific routing keys**:
   ```bash
   goq monitor -u localhost:5672 -r routing_key1,routing_key2,routing_key_pattern.#
   ```

3. **Use AMQPS with a specific virtual host**:
   ```bash
   goq dump -u rabbitmq.example.com:5671 -q important_queue -s -v my_vhost
   ```

4. **Automatically acknowledge messages and stop after consuming**:
   ```bash
   goq dump -u localhost:5672 -q my_queue -a -c
   ```

5. **Append to an existing output file**:
   ```bash
   goq dump -u localhost:5672 -q my_queue -m append -o existing_output.txt
   ```
6. **Print the full message including headers, exchange, and timestamp:**
   ```bash
   goq dump -u localhost:5672 -q my_queue -f
   ```

## Configuration
GOQ can be configured using a YAML file. By default, it looks for a configuration file at `$XDG_CONFIG_HOME/goq/goq.yaml`.

### Example Configuration File
```yaml
# Configuration for rabbitmq-dumper

# Use AMQPS instead of AMQP
amqps: true

# Skip TLS certificate verification (insecure)
skip-tls-verify: true

# RabbitMQ server URL
url: "127.0.0.1:5672"

# RabbitMQ exchange name, leave empty if not used
exchange: "some_exchange"

# RabbitMQ virtual host, leave empty if not used
virtualhost: "some-virtual-host"

# Output file name (relative to the current working directory)
output: "messages.txt"

# File mode (append or overwrite)
file-mode: "overwrite"

# Pretty print JSON messages
pretty-print: true

# Print full message including headers, exchange, timestamp, and body
full-message: false
```

When using a config file, the flags will override the values in the config file. It's useful to combine the config file with command-line flags for more flexibility.

## How Does GOQ Work?
1. **Connection**: Establishes a connection to the specified RabbitMQ server using either AMQP or AMQPS.
2. **Queue Binding**: Binds to the specified queue and exchange (if provided).
3. **Message Consumption**: Begins consuming messages from the queue.
4. **File Writing**: Writes consumed messages to the specified output file.
5. **Acknowledgment**: Acknowledges messages based on the auto-ack setting.
6. **Monitoring**: Displays the progress of consumed messages.
7. **Termination**: Stops consuming messages based on user input (CTRL+C) or the stop-after-consume flag.

## Contributing
Contributions to GOQ are welcome! Please feel free to submit pull requests, create issues for bugs and feature requests, or contribute to the documentation.

## License
GOQ is released under the MIT License. See the LICENSE file for more details.

