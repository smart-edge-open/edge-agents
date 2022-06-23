# INTEL CONFIDENTIAL
#
# Copyright 2021-2021 Intel Corporation.
#
# This software and the related documents are Intel copyrighted materials, and your use of
# them is governed by the express license under which they were provided to you ("License").
# Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
# disclose or transmit this software or the related documents without Intel's prior written permission.
#
# This software and the related documents are provided as is, with no express or implied warranties,
# other than those that are expressly stated in the License.

FROM alpine:3.12.0 AS edgedns-deps-image

RUN printf "http://nl.alpinelinux.org/alpine/v3.12/main\nhttp://nl.alpinelinux.org/alpine/v3.12/community" >> /etc/apk/repositories
RUN apk add --no-cache sudo

FROM edgedns-deps-image

ENV DNS_FWDR=""

ARG username=edgednssvr
ARG user_dir=/home/$username

RUN addgroup -S sudo && adduser -S $username -G sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
# Workaround for sudo error: https://gitlab.alpinelinux.org/alpine/aports/issues/11122
RUN echo 'Set disable_coredump false' >> /etc/sudo.conf

USER $username
WORKDIR $user_dir

COPY ./edgednssvr ./
ENTRYPOINT ["sudo", "./edgednssvr"]
CMD ["-port=53"]
