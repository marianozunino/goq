## goq

A tool to dump RabbitMQ messages to a file

### Synopsis


```
  ______    ______    ______
 /\  ___\  /\  __ \  /\  __ \
 \ \ \__ \ \ \ \/\ \ \ \ \/\_\
  \ \_____\ \ \_____\ \ \___\_\
   \/_____/  \/_____/  \/___/_/
```

This application connects to a RabbitMQ server and dumps queue messages to a file.

### Options

```
  -s, --amqps                      Use AMQPS instead of AMQP
      --config string              Config file path (default "/home/forbi/.config/goq/goq.yaml")
  -e, --exchange string            RabbitMQ exchange name
  -x, --exclude-patterns strings   Exclude messages containing these patterns
  -m, --file-mode string           File mode (append or overwrite) (default "overwrite")
  -f, --full-message               Print complete message details
  -h, --help                       help for goq
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

* [goq configure](goq_configure.md)	 - Generate a sample configuration file for goq.
* [goq dump](goq_dump.md)	 - Dump messages from a RabbitMQ queue
* [goq monitor](goq_monitor.md)	 - Monitor RabbitMQ messages using routing keys
* [goq update](goq_update.md)	 - Update the goq tool to the latest available version.
* [goq version](goq_version.md)	 - Display the current version of the goq tool.

