# GOQ (Goku)

```
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/
```

GOQ is a command-line tool for working with RabbitMQ queues, allowing you to dump, monitor, and manage messages efficiently.

## Installation
To install GOQ, run:

```bash
go install github.com/marianozunino/goq@latest
```

Or download the binary from the [GitHub releases](https://github.com/marianozunino/goq/releases)

## Documentation
For detailed usage instructions, configuration options, and examples, please refer to the [GOQ Documentation](docs/goq.md).

## Quick Start
1. **Dump messages from a local RabbitMQ server**:
   ```bash
   goq dump -u localhost:5672 -q my_queue
   ```

2. **Monitor messages with specific routing keys**:
   ```bash
   goq monitor -u localhost:5672 -r routing_key1,routing_key2,routing_key_pattern.#
   ```

## Configuration
GOQ can be configured using a YAML file. By default, it looks for a configuration file at `$XDG_CONFIG_HOME/goq/goq.yaml`. For configuration options and examples, see the [documentation](docs/goq.md#configuration).

## Contributing
Contributions to GOQ are welcome! Please feel free to submit pull requests, create issues for bugs and feature requests, or contribute to the documentation.

## License
GOQ is released under the MIT License. See the LICENSE file for more details.
