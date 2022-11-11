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

### Highlight

- From `v0.1.4`, gosmpp is written in event-based style and fully-manage your smpp session, connection, error, rebinding, etc. You only need to implement some hooks:

```go
		trans, err := gosmpp.NewSession(
		gosmpp.TRXConnector(gosmpp.NonTLSDialer, auth),
		gosmpp.Settings{
			EnquireLink: 5 * time.Second,

			ReadTimeout: 10 * time.Second,

			OnSubmitError: func(_ pdu.PDU, err error) {
				log.Fatal("SubmitPDU error:", err)
			},

			OnReceivingError: func(err error) {
				fmt.Println("Receiving PDU/Network error:", err)
			},

			OnRebindingError: func(err error) {
				fmt.Println("Rebinding but error:", err)
			},

			OnPDU: handlePDU(),

			OnClosed: func(state gosmpp.State) {
				fmt.Println(state)
			},
		}, 5*time.Second)
		if err != nil {
		  log.Println(err)
		}
		defer func() {
		  _ = trans.Close()
		}()
```

### Version (0.1.4.RC+)

- Full example could be found: [here](https://github.com/linxGnu/gosmpp/blob/master/example)
  - In this example, you should run smsc first:
    - Build and run SMSC Simulator:
	```bash
	g++ example/smsc_simulator/smsc.cpp -o smsc
	./smsc &
	```
    - Run smpp client in the example: https://github.com/linxGnu/gosmpp/blob/master/example/main.go
    ```bash
	go run example/main.go
	```

### Old version (0.1.3 and previous)
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
