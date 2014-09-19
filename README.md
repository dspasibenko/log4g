# Log4g Logging Library
Log4g is an open source project which intends to bring fast, flexible and scalable logging solution to the Go world. It allows the developer to control which log statements are output with arbitrary granularity. It is fully configurable at runtime using external configuration files. Among other logging approaches log4g borrows some cool stuff from log4j library, but it doesn't try to repeat the library architecture, even many things look pretty similar (including the name of the project with one letter different only)

## Architecture
log4g operates with the following terms and components: _log level_,  _logger_, _log event_, _logger name_, _logger context_, _appender_, and _log4g configuration_. 
`log4g.go` defines types, interfaces and public functions that allow to configure, use and expand the library functionalily. 

### Log Level
_Log level_ is an integer value which lies in [0..70] range. Every log message is associated with a certain _log level_ when the message is emitted. Depending on log4g configuration some messages can be filtered out and skipped from the processing flow because of their level. A level with lowest value has higher priority than the level with highest value. This means that if level X is set as maximum allowed, only messages with levels X1 <= X will be processed.

There are 7 _log level_ constants are pre-defined:
 * FATAL = 10
 * ERROR = 20
 * WARN = 30
 * INFO = 40
 * DEBUG = 50
 * TRACE = 60
 * ALL = 70

Any _log level_ value is associated with its name. Log message with certain _log level_ can be eventually formatted so the _log level_ name will be placed in the logging statement for the message. The pre-defined _log level_ values have the similar as the constat names associations (FATAL, ERROR etc.). Users can define their own _log levels_ names or change the pre-defined ones. To do this, users should make appropriate settings in _log4g configuration_ configuration file or do it programmatically, for example:

```
    ok := log4g.SetLogLevelName(23, "SEVERE")
```

In the example above the _log level_ with value 23 is named like "SEVERE". All messages emitted with the log level 23 will be named as "SEVERE" in the final log text, if the message formatting suppose to show messages log levels.
 
### Logger Name
The first and foremost advantage of log4g resides in its ability to disable certain log statements while allowing others to print unhindered. This capability assumes that the logging space, that is, the space of all possible logging statements, is categorized according to some developer-chosen criteria.

_Logger name_ is a string which should start from a letter, can contain letters `[A-Za-z]`,
digits `[0-9]` and dots `.`. The name cannot have `.` as a last symbol. The _root logger name_ is an empty string `""`. 

_Logger names_ are case-sensitive and they follow the hierarchical naming rule: A _logger name_ is said to be an ancestor of another _logger name_ if its name followed by a dot is a prefix of the descendant _logger name_. A _logger name_ is said to be a parent of a child _logger name_ if there are no ancestors between itself and the descendant _logger name_.

For example, the _Logger names_ `FileSystem` is a parent of the `FileSystem.ntfs`. Similarly, `a.b` is a parent of `a.b.c` and an ancestor of `a.b.c.d`. This naming scheme should be familiar to most developers and especially log4j users.

### Log Level Settings
log4g log level filtering configuration consists of list of pairs *<logger name : maximum allowed log level>*. Hierarchical relations of logger names allow to build flexible and advanced filtering configurations.

Adding or changing value in the list can be done by calling `SetLogLevel()` function declared in log4g package:

```
    func SetLogLevel(loggerName string, level Level)
```

Every logging message level is checked against nearest ancestor of the message logger name from the list. If the logger name level is low than the message level, the message will be skipped and not processed further. 

For example, lets suppose that log levels were set for 2 log names:

```
    log4g.SetLogLevel("FileSystem", INFO)
    log4g.SetLogLevel("FileSystem.ntfs", DEBUG)
```

which produce 2 pairs in the log level settings list. Logging messages made for `FileSystem.ext2` logger name will be filtered by `INFO` level, because the nearest ancestor for the name is `FileSystem`. But logging messages made for `FileSystem.ntfs` of `FileSystem.ntfs.hidden` will be filtered by `DEBUG` level, because the nearest ancestor for the namea is `FileSystem.ntfs`. 

log4g always has _root log level setting_ configured for _root logger name_

### Logger
`Logger` is an interface which allows to post logging messages to log4g for further processing. The instance of the interface can be retrieve by the function:

```
    func GetLogger(loggerName string) Logger
```

The function is idempotent, so it always returns same object for the same _logger name_ regardless of the log4g configuration and settings were made between different function calls. 

The `Logger` interface is only one "front end" element of the logging message processing, all logging messages in log4g are submitted through instances of the interface. 

Any logging message, submitted to log4g via `Logger`, is checked against _Log Level Settings_ and if the message should NOT be filtered because of its level, it is transformed to `LogEvent` object which is passed to _Logger Context_ for further processing. 

### Logger Context
_Logger Context_ is an internal component which allows to aggregate logging messages from different _loggers_ and distribute them between _Appenders_ associated with the _Logger Context_. 

There is no special functions to configure _Logger Contexts_, so users can specify their configurations via log4g configuration calls. 

Every _Logger_ is associated with one _Logger Context_, but one _Logger Context_ can be associated with multiple _Loggers_. This association is done by same manner how _Log Level Settings_ are applied to _Loggers_: every _Logger Context_ has a "Logger Name" associated with it. So as every _Logger_ is always associated with its "Logger Name" the _Logger Context_ is associated with the _Logger_ if its logger name is closed ancestor for the logger name. 

log4g always has _Logger Context_ associated with _root logger name_, so every _Logger_ will always be associated with at least this _root Logger Context_.

### Appender

### Log4g Configuration




