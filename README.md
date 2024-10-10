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
  completion  Generate the autocompletion script for the specified shell
  configure   Create a sample configuration file
  dump        Dump RabbitMQ messages to a file
  help        Help about any command
  monitor     Create a temporary queue and consume messages from it having the specified routing keys
  update      Update GOQ to the latest version
  version     Print the version number of GOQ

Flags:
  -s, --amqps                Use AMQPS instead of AMQP
      --config string        Config file (default is $XDG_CONFIG_HOME/goq/goq.yaml)
  -e, --exchange string      RabbitMQ exchange name
  -m, --file-mode string     File mode (append or overwrite) (default "overwrite")
  -h, --help                 Help for GOQ
  -o, --output string        Output file name (default "messages.txt")
  -k, --skip-tls-verify      Skip TLS certificate verification (insecure)
  -u, --url string           RabbitMQ URL (e.g., localhost:5672)
  -v, --virtualhost string   RabbitMQ virtual host
  -p, --pretty-print         Pretty print JSON messages

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

