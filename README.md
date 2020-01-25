[error] some error happened in ... // / / .
[info] info messages
[warning] ...
[warning] ...
[info] info messages
[info] info messages
[info] info messages
[info] info messages
[error] some error happened in ... // / / .
[error] some error happened in ... // / / .
[error] some error happened in ... // / / .
[error] some error happened in ... // / / .


main.log =>
  info.log
  warning.log
  error.log

main thread -> main.log

3 workers -> line -> type -> file


main thread -> processing which level you have (error, warn, info) ->
3 threads (1 thread per log type)
  errors
  warning
  info

main thread - for reading
processing goroutine - parse the string (error, warning, info) -> errorChan, warningChan, infoChan
3 output threads
