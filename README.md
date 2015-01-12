Logging Utility for Go
======================

An opinonated telemetry & logging utility for Go. 

+ one global logging stream per application
+ formatted as a stream of mostly unstructured JSON data events 
+ transport via UDP to an aggregator

Log Events
-----------

Log is a stream of data events that fall into 3 categories:

+ *Action*: an event that indicates a problem that needs attention to resolve. 
+ *Info*: an audit event adding context around any events requiring action. 
+ *Telemetry*: timing or count events that are to be used within visulation tools to better understand system dynamics.

Benchmark
---------

    go test -bench .