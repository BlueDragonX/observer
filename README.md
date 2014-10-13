Observer
========
Observer is a modular system for monitoring a host system. With some minor
amount of duct tape you can monitor a host system from inside of a Docker
container. This is Observer's primary use case.

Config
------
Observer, being generic and modular, uses the concept of sources, sinks, and
pipes. A source is a source of metric data. A sink is someplace you can put
that metric data. And a pipe lets you connect one or more sources to one or
more sinks.

The best way to get started is to take a look at the two configuration
examples. The [config-print.yml][1] configures Observer to print system load,
memory usage, and disk space metrics. The [config-aws.yml][2] does the same
thing but puts them in AWS CloudWatch instead of printing them.

License
-------
Copyright (c) 2014 Ryan Bourgeois. Licensed under BSD-Modified. See the
[LICENSE][2] file for a copy of the license.

[1]: https://github.com/BlueDragonX/observer/blob/master/config-print.yml   Dummy Example
[1]: https://github.com/BlueDragonX/observer/blob/master/config-aws.yml   AWS Example
