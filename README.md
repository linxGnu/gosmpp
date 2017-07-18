# gosmpp
Golang Smpp (3.4) Client Library, porting from [Java OpenSMPP Library](https://github.com/OpenSmpp/opensmpp). 

This library is tested well with several SMSC simulators:
- [smpp-smsc-simulator](http://www.voldrich.net/2015/01/11/smpp-smsc-simulator/): simulates a SMSC server – server which accepts SMS messages and handles its delivery to the mobile phone.
- [SMPPSim](http://www.seleniumsoftware.com/downloads.html): a SMPP SMSC simulation tool, designed to help you test your SMPP based application. SMPPSim is free of charge and open source.

## Installation
```
go get -u github.com/linxGnu/gosmpp
```

## Usage
Please refer to [Communication Test Case](https://github.com/linxGnu/gosmpp/blob/master/test/Communication_test.go) for sample code. If you are familiar with [OpenSMPP](https://github.com/OpenSmpp/opensmpp), you would know how to implement it easily.

TODO: more sample codes and test cases added in next minor version

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

## License
MIT License

Copyright (c) 2016 linxGnu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

