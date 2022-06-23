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

# Contribution Guide

Please consider the following criteria before suggesting or implementing any changes:

* This project's goal is to provide the absolute bare _minimum_ set of DNS features, it is expected that a forwarder will be used if more advanced features or controls are required
* Performance is extremely important because this service will impact almost all mobile user traffic
* Zone guards should be implemented by a calling services, this service only provides the responder features

Updating the gRPC inteface

`protoc -I pb --go_out=plugins=grpc:pb pb/resolver.proto`

