# HTTP Traffic Monitor

[![HTTP Traffic Monitor](https://goreportcard.com/badge/github.com/mpmlj/http-traffic-monitor)](https://goreportcard.com/report/github.com/mpmlj/http-traffic-monitor) [![Build Status](https://travis-ci.org/mpmlj/http-traffic-monitor.svg)](https://travis-ci.org/mpmlj/http-traffic-monitor)

## Code Challenge

Create a simple console program that monitors HTTP traffic on your machine:

- Consume an actively written-to w3c-formatted HTTP access log (https://en.wikipedia.org/wiki/Common_Log_Format)

- Every 10s, display in the console the sections of the web site with the most hits (a section is defined as being what's before the second '/' in a URL. i.e. the section for "http://my.site.com/pages/create' is "http://my.site.com/pages"), as well as interesting summary statistics on the traffic as a whole.

- Make sure a user can keep the console app running and monitor traffic on their machine

- Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”

- Whenever the total traffic drops again below that value on average for the past 2 minutes, add another message detailing when the alert recovered

- Make sure all messages showing when alerting thresholds are crossed remain visible on the page for historical reasons.

- Write a test for the alerting logic

- Explain how you’d improve on this application design

## General notes

- The program streams information to stdout in chronological order.
- There are three types of messages:<br>
-- point, notification about data check at every log file poll<br>
-- report, traffic summary prepared at user-defined intervals<br>
-- alert, notification about traffic average exceeding crossing a threshold<br>
- Each of these message types can be silenced for testing or usability purposes.

## Script arguments and defaults

##### Available for configuration by user:

` --log-file` - log file location, required.

`--alert-threshold` - alert threshold, _hits/sec._, default 1000 hits, optional.

`--poll-interval` - polling interval, _sec._, default 1 sec., optional.

`--mtf` - monitoring time frame, _sec._, default 120 sec., optional.

`--top-n` - # most visited sections, _sec._, default 10, optional.
 
`--report-interval` - interval for showing traffic report, _sec., default 10, optional.

Below - configuration for the code challenge (should be run from the repository root location):
 
`./bin/http-traffic-monitor --log-file=log/server.log --alert-threshold=2 --poll-interval=1 --mtf=120 --top-n=5 --report-interval=10`

It will run the monitor to:
- show a summary report every 10 seconds with top 5 most visited sections
- trigger an alert escalation and de-escalation based on a 2-minute moving average

##### Unavailable for external configuration, used for testing and debugging:

See Config struct in `config.go`:

- `MaxPolls` (int) limits duration of the monitor via maximum number of ticks (= polls, 1 second minimum) it should make before exiting. By default it is set to maximum int32 number and I hope 68 years should be enough for majority of users... :) However, set to a lower number when configuring monitor tests. For example, to test escalation and deescalation I set it to 2-3 seconds as this is enough to accumulate data for these tests.  

- `SendAlerts`, `SendReports` and `SendTicks` (bool) flags are used to silence certain types of output messages. Required to remove noise when testing one specific behavior. For example, I silenced ticks and reports to be able to test alert messaging from the monitor. 


## Testing

##### Notes on manual testing

- `"log"` folder contains a source log files for manual testing. 

- `server.log` is an (originally) empty file you can use to add records.

- `src.log` is an example/source log file with 1000 entries of Common Log Format.

- To test script behavior, append to file with `echo >>`, do not change the log file using vim, etc. and this will trigger `rename` OS event and the script will loose file handler and stop monitoring this particular file. 

- Example of a script configuration for testing, with low levels, to more easily trigger state changes:<br>
`./bin/http-traffic-monitor --log-file=log/server.log --alert-threshold=2 --poll-interval=1 --mtf=10 --top-n=5 --report-interval=3`

##### Sample echo commands

6 log entries with different response codes:

````
echo '
165.13.14.55 - - [28/Jul/1995:13:17:00 -0400] "GET /history/apollo/apollo-17/apollo-17-info.html HTTP/1.0" 200 1457
198.155.12.13 - - [28/Jul/1995:13:17:07 -0400] "GET /shuttle/countdown/liftoff.html HTTP/1.0" 200 5220
182.200.120.1 - - [28/Jul/1995:13:17:08 -0400] "GET /shuttle/countdown/count.html HTTP/1.0" 500 65536
198.155.12.16 - - [28/Jul/1995:13:17:09 -0400] "GET /images/NASA-logosmall.gif HTTP/1.0" 400 786
182.200.120.2 - - [28/Jul/1995:13:17:08 -0400] "GET /shuttle/countdown/count.html HTTP/1.0" 300 65536
182.200.120.3 - - [28/Jul/1995:13:17:08 -0400] "GET /shuttle/countdown/count.html HTTP/1.0" 500 65536
' >> server.log
````

20 log entries to trigger an alert:

````
echo '
131.182.170.137 - - [28/Jul/1995:13:16:31 -0400] "GET /images/MOSAIC-logosmall.gif HTTP/1.0" 200 363
210.166.12.00 - - [28/Jul/1995:13:16:31 -0400] "GET /htbin/cdt_clock.pl HTTP/1.0" 200 503
128.203.26.245 - - [28/Jul/1995:13:16:32 -0400] "GET /software HTTP/1.0" 302 1
128.203.26.245 - - [28/Jul/1995:13:16:33 -0400] "GET /software/ HTTP/1.0" 200 816
198.155.12.13 - - [28/Jul/1995:13:16:34 -0400] "GET /history/apollo/apollo-13/images/ HTTP/1.0" 200 1851
128.203.26.245 - - [28/Jul/1995:13:16:35 -0400] "GET /icons/blank.xbm HTTP/1.0" 401 509
128.203.26.245 - - [28/Jul/1995:13:16:35 -0400] "GET /icons/menu.xbm HTTP/1.0" 400 527
131.182.170.137 - - [28/Jul/1995:13:16:36 -0400] "GET /images/USA-logosmall.gif HTTP/1.0" 200 234
198.240.108.240 - - [28/Jul/1995:13:16:37 -0400] "GET /icons/blank.xbm HTTP/1.0" 200 509
198.240.108.240 - - [28/Jul/1995:13:16:40 -0400] "GET /icons/menu.xbm HTTP/1.0" 401 527
131.182.170.137 - - [28/Jul/1995:13:16:40 -0400] "GET /images/WORLD-logosmall.gif HTTP/1.0" 200 669
198.240.108.240 - - [28/Jul/1995:13:16:43 -0400] "GET /icons/image.xbm HTTP/1.0" 200 509
128.203.26.245 - - [28/Jul/1995:13:16:43 -0400] "GET /software/winvn/ HTTP/1.0" 500 2244
128.203.26.245 - - [28/Jul/1995:13:16:45 -0400] "GET /icons/image.xbm HTTP/1.0" 201 509
128.203.26.245 - - [28/Jul/1995:13:16:45 -0400] "GET /icons/text.xbm HTTP/1.0" 200 527
198.240.108.240 - - [28/Jul/1995:13:16:46 -0400] "GET /icons/unknown.xbm HTTP/1.0" 302 515
182.198.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/technology/sts-newsref/srb.html HTTP/1.0" 200 49553
182.200.120.1 - - [28/Jul/1995:13:16:47 -0400] "GET /htbin/wais.pl?SAREX-II HTTP/1.0" 502 7111
210.166.12.92 - - [28/Jul/1995:13:16:47 -0400] "GET /shuttle/missions/sts-71/sts-71-press-kit.txt   HTTP/1.0" 502 78588
198.155.12.13 - - [28/Jul/1995:13:19:30 -0400] "GET /shuttle/missions/sts-65/mission-sts-65.html HTTP/1.0" 200 131165
' >> server.log
````
 
##### Notes on testing with "go test"

 To test alerting (de)escalation logic only, run `go test -v -run TestMonitor`.


## UI description

#### Points

` ·  hits avg:      4  /  2`<br>First number - average number of hits received since last tick.
Second - number alert threshold.
When average number of hits exceeds threshold, the number will marked in red.

#### Report
````
--------------------------------------------------------------------------------
REPORT: 2017-02-06T01:44:41-05:00

Top 5 sections
| sections                                                       | count
--------------------------------------------------------------------------------
| /shuttle                                                       | 24
| /images                                                        | 6
| /history                                                       | 6

Summary:
| hits total | hits/s      | 2xx        | 3xx        | 4xx        | 5xx
--------------------------------------------------------------------------------
| 36         | 12         | 12         | 6          | 6          | 12
````

#### Alerts

````
 · High traffic generated an alert - hits = 2, triggered at 2017-02-06T01:48:10-05:00
````


````
 · High traffic alert recovered. Current hits = 0. At 2017-02-06T01:48:20-05:00
````

## Improvement considerations

- Introduce a warning threshold level, that would signal approaching to an actual alert level.
- Implement other senders, ex. SNS, Pager Duty, Slack, email.
- Other log formats (combined, extended, custom via regex expressions, etc.).
- Implement better logic for reaction on logrotate - currently loses file handler due to a rename event from OS.
- Remove dependencies, implement custom parser.
- Implement spike detection via stand deviation calculation.
- I would also like to revise data types used throughout the script as I feel like there can be some optimizations required.

## Other...

- Add average request size.
- Add request methods to report.
- Better test coverage.
