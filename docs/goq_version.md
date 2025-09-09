## goq version

Display the current version of the goq tool.

```
goq version [flags]
```

### Examples

```
goq version
```

### Options

```
  -h, --help   help for version
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

