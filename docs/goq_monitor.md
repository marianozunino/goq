## goq monitor

Monitor RabbitMQ messages using routing keys

### Synopsis

Monitor RabbitMQ messages by consuming from a temporary queue with specified routing keys.

```
goq monitor [flags]
```

### Examples

```
  # Monitor all messages from an exchange
  goq monitor -K "#" -e "my_exchange" -w console -p

  # Monitor specific routing keys with filtering
  goq monitor -K "user.created,user.updated" -e "events" -i "admin" -o users.log

  # Monitor with secure connection
  goq monitor -K "order.*" -e "orders" -s -k -u "rabbitmq.example.com:5671"
```

### Options

```
  -a, --auto-ack               Automatically acknowledge messages
  -h, --help                   help for monitor
  -K, --routing-keys strings   List of routing keys to monitor (required)
```

### Options inherited from parent commands

```
      --config string              Config file path (default "/home/forbi/.config/goq/goq.yaml")
  -e, --exchange string            RabbitMQ exchange name
  -x, --exclude-patterns strings   Exclude messages containing these patterns
  -m, --file-mode string           File mode (append or overwrite) (default "overwrite")
  -i, --include-patterns strings   Include messages containing these patterns
  -k, --insecure                   Skip TLS certificate verification
  -j, --json-filter string         JSON filter expression
  -z, --max-message-size int       Maximum message size in bytes (-1 for unlimited) (default -1)
  -o, --output string              Output file name
  -p, --pretty-print               Pretty print JSON messages
  -r, --regex-filter string        Regex pattern to filter messages
  -s, --secure                     Use AMQPS (secure) instead of AMQP
  -u, --url string                 RabbitMQ server URL (default "localhost:5672")
  -v, --virtualhost string         RabbitMQ virtual host (default "/")
  -w, --writer string              Output writer type (file or console) (default "file")
```

### SEE ALSO

* [goq](goq.md)	 - A tool to dump RabbitMQ messages to a file

