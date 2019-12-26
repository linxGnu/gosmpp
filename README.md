# gosmpp

[![](https://github.com/linxGnu/gosmpp/workflows/Build/badge.svg)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/linxGnu/gosmpp)](https://goreportcard.com/report/github.com/linxGnu/gosmpp)
[![Coverage Status](https://coveralls.io/repos/github/linxGnu/gosmpp/badge.svg?branch=master)](https://coveralls.io/github/linxGnu/gosmpp?branch=master)
[![godoc](https://img.shields.io/badge/docs-GoDoc-green.svg)](https://godoc.org/github.com/linxGnu/gosmpp)

SMPP (3.4) Client Library in pure Go.

This library is well tested with SMSC simulators:
- [Melroselabs SMSC](https://melroselabs.com/services/smsc-simulator/#smsc-simulator-try)

## Installation
```
go get -u github.com/linxGnu/gosmpp
```

## Usage

## Version (0.1.4.RC+)

Please refer to [Test Case And Sample Code](https://github.com/linxGnu/gosmpp/blob/master/transceiver_test.go).

## Old version (0.1.3 and previous)
Full example could be found: [gist](https://gist.github.com/linxGnu/b488997a0e62b3f6a7060ba2af6391ea)

## Supported PDUs

- [x] bind_transmitter
- [x] bind_transmitter_resp
- [x] bind_receiver
- [x] bind_receiver_resp
- [x] bind_transceiver
- [x] bind_transceiver_resp
- [x] outbind
- [x] unbind
- [x] unbind_resp
- [x] submit_sm
- [x] submit_sm_resp
- [x] submit_sm_multi
- [x] submit_sm_multi_resp
- [x] data_sm
- [x] data_sm_resp
- [x] deliver_sm
- [x] deliver_sm_resp
- [x] query_sm
- [x] query_sm_resp
- [x] cancel_sm
- [x] cancel_sm_resp
- [x] replace_sm
- [x] replace_sm_resp
- [x] enquire_link
- [x] enquire_link_resp
- [x] alert_notification
- [x] generic_nack
