# GOQ (Goku)

```
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/
```

## Installation
To install goq, run:

```bash
go install github.com/marianozunino/goq@latest
```

Or download the binary from the [GitHub releases](https://github.com/marianozunino/goq/releases) (Note: This link is hypothetical and should be updated with the actual releases page)

## Usage
```
Usage:
  goq [flags]

Flags:
      --config string        config file (default is $XDG_CONFIG_HOME/goq/goq.yaml)
  -u, --url string           RabbitMQ URL (e.g., localhost:5672)
  -e, --exchange string      RabbitMQ exchange name
  -q, --queue string         RabbitMQ queue name
  -o, --output string        Output file name (default "messages.txt")
  -s, --amqps                Use AMQPS instead of AMQP
  -v, --virtualhost string   RabbitMQ virtual host
  -k, --skip-tls-verify      Skip TLS certificate verification (insecure)
  -a, --auto-ack             Automatically acknowledge messages
  -m, --file-mode string     File mode (append or overwrite) (default "overwrite")
  -c, --stop-after-consume   Stop consuming after getting all messages from the queue
  -h, --help                 help for goq
```

### Examples
1. Dump messages from a local RabbitMQ server:
   ```
   goq -u localhost:5672 -q my_queue
   ```
2. Use AMQPS with a specific virtual host:
   ```
   goq -u rabbitmq.example.com:5671 -q important_queue -s -v my_vhost
   ```
3. Specify an exchange and output file:
   ```
   goq -u localhost:5672 -e my_exchange -q my_queue -o output.txt
   ```
4. Automatically acknowledge messages and stop after consuming:
   ```
   goq -u localhost:5672 -q my_queue -a -c
   ```
5. Append to an existing output file:
   ```
   goq -u localhost:5672 -q my_queue -m append -o existing_output.txt
   ```

## Configuration
goq can be configured using a YAML file. By default, it looks for a configuration file at `$XDG_CONFIG_HOME/goq/goq.yaml`.

Example configuration file:
```yaml
url: localhost:5672
exchange: my_exchange
queue: my_queue
output: messages.txt
amqps: false
virtualhost: /
skip-tls-verify: false
auto-ack: false
file-mode: overwrite
stop-after-consume: false
```

When using a config file the flags will override the values in the config file. So it is usefull to use the config file in combination with the flags,

## How Does goq Work?
1. **Connection**: Establishes a connection to the specified RabbitMQ server using either AMQP or AMQPS.
2. **Queue Binding**: Binds to the specified queue and exchange (if provided).
3. **Message Consumption**: Begins consuming messages from the queue.
4. **File Writing**: Writes consumed messages to the specified output file.
5. **Acknowledgment**: Acknowledges messages based on the auto-ack setting.
6. **Monitoring**: Displays progress of consumed messages.
7. **Termination**: Stops consuming messages based on user input (CTRL+C) or the stop-after-consume flag.

## Contributing
Contributions to goq are welcome! Please feel free to submit pull requests, create issues for bugs and feature requests, or contribute to the documentation.

## License
goq is released under the MIT License. See the LICENSE file for more details.
