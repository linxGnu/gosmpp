# gosmpp
Golang Smpp (3.4) Client Library, porting from [Java OpenSMPP Library](https://github.com/OpenSmpp/opensmpp). 

This library is tested well with several SMSC simulators:
- [smpp-smsc-simulator](http://www.voldrich.net/2015/01/11/smpp-smsc-simulator/): simulates a SMSC server – server which accepts SMS messages and handles its delivery to the mobile phone.
- [SMPPSim](http://www.seleniumsoftware.com/downloads.html): a SMPP SMSC simulation tool, designed to help you test your SMPP based application. SMPPSim is free of charge and open source.

gosmpp has run well in production now:
- [My friend](https://github.com/tanlinhnd) at [traithivang.vn](http://traithivang.vn/) has used gosmpp as client to SMSC of Vietnamobile, a telecommunications company in Vietnam, without any problems for months.

## Installation
```
go get -u github.com/linxGnu/gosmpp
```

## Usage
Please refer to [Communication Test Case](https://github.com/linxGnu/gosmpp/blob/master/test/Communication_test.go) for sample code. If you are familiar with [OpenSMPP](https://github.com/OpenSmpp/opensmpp), you would know how to implement it easily.

Full project of building SMPP Client could be found at: [Telcos](https://github.com/linxGnu/gosmpp/examples/telcos)

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
- [ ] alert_notification
- [x] generic_nack

## Contributing
Please issue me for things gone wrong or:

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D