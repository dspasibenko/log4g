# Log4g Logging Library
Log4g is an open source project which intends to bring fast, flexible and scalable logging solution to the Go world. It allows the developer to control which log statements are output with arbitrary granularity. It is fully configurable at runtime using external configuration files. Among other logging approaches log4g borrows some cool stuff from log4j library, but it doesn't try to repeat the library architecture, even many things look pretty similar (including the name of the project with one letter different only)

## Architecture
log4g operates with the following terms and components: _log level_, _logger_, _logger name_, _logger contex_t, _appender_, and _log4g configuration_. 
