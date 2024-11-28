## goq dump

Dump messages from a RabbitMQ queue

### Synopsis

Dump messages from a specified RabbitMQ queue with flexible filtering and output options.

```
goq dump [flags]
```

### Examples

```
goq dump -q my_queue -o output.txt
```

### Options

```
  -a, --auto-ack             Automatically acknowledge messages
  -h, --help                 help for dump
  -q, --queue string         RabbitMQ queue name (required)
  -c, --stop-after-consume   Stop after consuming messages
```

### Options inherited from parent commands

```
  -s, --amqps                      Use AMQPS instead of AMQP
      --config string              Config file path (default "/home/forbi/.config/goq/goq.yaml")
  -e, --exchange string            RabbitMQ exchange name
  -x, --exclude-patterns strings   Exclude messages containing these patterns
  -m, --file-mode string           File mode (append or overwrite) (default "overwrite")
  -f, --full-message               Print complete message details
  -i, --include-patterns strings   Include messages containing these patterns
  -j, --json-filter string         JSON filter expression
  -z, --max-message-size int       Maximum message size in bytes (-1 for unlimited) (default -1)
  -o, --output string              Output file name
  -p, --pretty-print               Pretty print JSON messages
  -R, --regex-filter string        Regex pattern to filter messages
  -k, --tls-skip                   Skip TLS certificate verification
  -u, --url string                 RabbitMQ server URL (default "localhost:5672")
  -v, --virtualhost string         RabbitMQ virtual host (default "/")
  -w, --writer string              Output writer type (file or console) (default "file")
```

### SEE ALSO

* [goq](goq.md)	 - A tool to dump RabbitMQ messages to a file

