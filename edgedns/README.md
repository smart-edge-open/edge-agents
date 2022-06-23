```text
INTEL CONFIDENTIAL

Copyright 2021-2021 Intel Corporation.

This software and the related documents are Intel copyrighted materials, and your use of
them is governed by the express license under which they were provided to you ("License").
Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
disclose or transmit this software or the related documents without Intel's prior written permission.

This software and the related documents are provided as is, with no express or implied warranties,
other than those that are expressly stated in the License.
```

# Edge DNS Responder

This project provides a standards compliant DNS server that exposes gRPC interfaces for the realtime creation of records.

Feature|Community Edition|Enterprise Edition|
|---|:---:|:---:|
|gRPC Control API|✅|✅|
|Embedded database|✅|✅|
|Embedded Forwarder Cache||✅|
|Nested dynamic Forwarder chains||✅
|IPv6 Listeners||✅|
|IPv6 Record Types||✅|
|Authoritative TXT Record||✅|
|Authoritative SRV Record||✅|
|Dynamic logging levels||✅|
|Logging to syslog||✅|

This Community Edition server implements:

* DNS Authoritative server
* Control via gRPC API

## Usage

All queries are processed in the following order:

1. Authoritative lookup (default TTL of 10 seconds)
2. Forwarder lookup

The Enterprise Edition allows the dynamic definition of forwarders on a per FQDN basis with hierarchical traversal of forwarders if a given forwarder does not return an answer for the query.

### API Client

See the test [API client](pkg/edgednssvr/test/control_client.go) for example usage of the control API.

### Logging

By default only major events related to the listeners or databases, as well as control API requests, are sent to `STDERR`.

### CLI

You can specify the following options:

|flag|required|default|description|
|---|---|---|---|
|4|NO|anyhost|IPv4 Listen address|
|port|NO|5053|UDP Listen port|
|sock|NO|`/run/edgedns.sock`|Filesystem path for the UNIX gRPC socket|
|address|NO|``|API IP address. If defined, socket parameter is not used|
|db|NO|`/var/lib/edgedns/rrsets.db`|Filesystem path for persistent database file|
|statsdip|NO|``|IP address of external statsd service. By default mock statsd service gets initiated|
|statsdport|NO|0|Port of external statsd service|
|hb|NO|60|Heartbeat interval. Heartbeats sent to the stats service at the specified interval in s|
|log|NO|`info`|Log level. Supported values: debug, info, notice, warning, error, critical, alert, emergency|
|syslog|NO|``|Syslog address|
|cert|NO|`certs/cert.pem`|PKI Cert Path|
|key|NO|`certs/key.pem`|PKI Key Path|
|ca|NO|`certs/root.pem`|PKI CA Path|

## Configuration

The following operations are available via the gRPC inteface:

* Add(Create/Update), Delete and Get operations for A records
* Add/Delete operations for forwarders

