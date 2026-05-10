//
//  smscsimulator.cpp
//  SMPPLib
//
//  Created by Mark Hay on 12/05/2019.
//  Copyright © 2019 Melrose Labs. All rights reserved.
//

// Build: g++ -std=c++11 smscsimulator.cpp -o MLSMSCSimulator && ./MLSMSCSimulator

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <iostream>
#include <string>
#include <exception>
#include <map>
#include <list>

using namespace std;

#include <unistd.h>

#include <sys/ioctl.h>
#include <time.h>
#include <sys/time.h>

#include <netdb.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#define TRUE  (1==1)
#define FALSE (!TRUE)

//

uint64_t session_id_next = 0; // ID for each session

// General

uint64_t currentUSecsSinceEpoch( void )
{
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (uint64_t)tv.tv_sec * 1000000 + (uint64_t)tv.tv_usec;
}

// Client config

class ClientConfig {
private:
    ClientConfig() {}

    ~ClientConfig() {}

public:

    static ClientConfig& instance() {
        static ClientConfig c;
        return c;
    }

    typedef enum {
        BINDRESP_WITH_SC_INTERACE_VERSION_TLV
    } ConfigEntry;

    bool is(string esme_id,ConfigEntry entry) {
        return false;
    }
};

// Message deliverer

class MessageDeliverer {

public:
    class Message
    {
    public:
        uint8_t source_addr_ton;
        uint8_t source_addr_npi;
        char source_addr[32];
        uint8_t dest_addr_ton;
        uint8_t dest_addr_npi;
        char destination_addr[32];
        uint8_t esm_class;
        uint8_t protocol_id;
        uint8_t priority_flag;
        char schedule_delivery_time[32];
        char validity_period[32];
        uint8_t registered_delivery;
        uint8_t replace_if_present_flag;
        uint8_t data_coding;
        uint8_t sm_default_msg_id;
        uint8_t sm_length;
        uint8_t short_message[160];

        string smscMessageID;

    private:
        string content;

    public:
        Message() {}
        Message(char*in) { content = in; }
        ~Message() {}

        string getContent(void) { return content; }

        void setSource( uint8_t ton, uint8_t npi, char* addr ) { source_addr_ton = ton; source_addr_npi = npi; strcpy(source_addr,addr); }
        void setDestination( uint8_t ton, uint8_t npi, char* addr ) { dest_addr_ton = ton; dest_addr_npi = npi; strcpy(destination_addr,addr); }
        void setRegisteredDelivery( uint8_t val ) { registered_delivery = val; }
        void setShortMessage( uint8_t* sm_in, uint8_t sm_len_in ) { memcpy(short_message, sm_in, sm_len_in); sm_length = sm_len_in; }
        void setSMSCMessageID( char* id ) { smscMessageID = id; }
    };

private:
    // time-ordered list [using STL map] of lists containing messages
    typedef map<uint64_t/*time*/,list<Message>> MessageQueue;
    MessageQueue mq;

    MessageDeliverer() {}

    typedef map<string,MessageDeliverer*> MDMap;
    static MDMap mapInstance;

    bool empty(void) { return (mq.size()==0); }

    bool get( Message& msgout )
    {
        if (mq.size()==0) return false; // no time-lists

        uint64_t firstListTime = (*mq.begin()).first;
        uint64_t now = currentUSecsSinceEpoch();
        if ( firstListTime > now ) return false; // nothing due to be deliverered

        // we have a message (or messages) that are due to be delivered

        msgout = (*mq.begin()).second.front();

        (*mq.begin()).second.pop_front(); // remove from list

        if ((*mq.begin()).second.size()==0) // now empty - remove list
            mq.erase(mq.begin());

        return true;
    }

public:
    ~MessageDeliverer() {}

    static MessageDeliverer* getInstance(string& systemID)
    {
        MDMap::iterator it = mapInstance.find(systemID);
        if ( it == mapInstance.end() ) {
            std::pair<MDMap::iterator, bool> res = mapInstance.insert(MDMap::value_type(systemID,new MessageDeliverer));
            return (res.first)->second;
        }
        else return (*it).second;
    }

    static bool getInstance_get(string& systemID, Message& msgout)
    {
        MDMap::iterator it = mapInstance.find(systemID);
        if ( it == mapInstance.end() ) {
            return false;
        }

        bool found = (*it).second->get(msgout);

        if ((*it).second->empty()) {
            delete (*it).second;
            mapInstance.erase(systemID);
        }

        return found;
    }

    void add(uint64_t timeDeliver, Message msg)
    {
        MessageQueue::iterator it = mq.find(timeDeliver);
        if ( it == mq.end() ) {
            // no time-list exists - create
            list<Message> msgList; // empty message list
            std::pair<MessageQueue::iterator,bool> ret;
            ret = mq.insert({timeDeliver,msgList});
            if ( ret.second == false ) return; // failed to insert
            it = ret.first; // iterator to newly inserted time-list
        }

        (*it).second.push_back(msg); // add message to time-list
    }
};

MessageDeliverer::MDMap MessageDeliverer::mapInstance;


// Sub-SMPP layer

class SMPPSocket {
protected:
    int socket;

public:
    SMPPSocket() {}
    virtual ~SMPPSocket() {}

    virtual bool recvA( uint8_t& ) = 0;
    virtual void recv( void ) = 0;
    virtual bool send( uint8_t*, int len ) = 0;

    virtual long bytes_to_read( void ) = 0;
};

class SMPPSocketUnencrypted : public SMPPSocket {
private:

public:
    SMPPSocketUnencrypted() {

    }
    SMPPSocketUnencrypted(int socket_in) {
        socket = socket_in;
    }
    ~SMPPSocketUnencrypted() {}

    long bytes_to_read( void )
    {
        long count = 0;
        ioctl((int)socket, FIONREAD, &count);
        return count;
    }

    bool recvA( uint8_t& oct ) {
        int n = (int)::recv(socket,(void*)&oct,1,0);
        return (n>0);
    }
    void recv( void ) {}
    bool send( uint8_t* buf, int len )
    {
        return ::send(socket,(void*)buf,len,0);

    }
};

// SMPP

class SMPP {
public:

    class CmdStatus {
    public:
        static const uint64_t ESME_ROK = 0x00000000;
        static const uint64_t ESME_RINVBNDSTS = 0x00000004;
        static const uint64_t ESME_INVDSTADR = 0x0000000B;
        static const uint64_t ESME_RBINDFAIL = 0x0000000D;
        static const uint64_t ESME_RSUBMITFAIL = 0x00000045;
        static const uint64_t ESME_RINVSCHED = 0x00000061;
    };

    class CmdID {
    public:
        static const uint64_t BindReceiver = 0x00000001;
        static const uint64_t BindTransmitter = 0x00000002;
        static const uint64_t QuerySM = 0x00000003;
        static const uint64_t SubmitSM = 0x00000004;
        static const uint64_t DeliverSM = 0x00000005;
        static const uint64_t Unbind = 0x00000006;
        static const uint64_t ReplaceSM = 0x00000007;
        static const uint64_t CancelSM = 0x00000008;
        static const uint64_t BindTransceiver = 0x00000009;
        static const uint64_t Outbind = 0x0000000B;
        static const uint64_t EnquireLink = 0x00000015;
        static const uint64_t SubmitMulti = 0x00000021;
        static const uint64_t AlertNotification = 0x00000102;
        static const uint64_t DataSM = 0x00000103;
        static const uint64_t BroadcastSM = 0x00000111;
        static const uint64_t QueryBroadcastSM = 0x00000112;
        static const uint64_t CancelBroadcastSM = 0x00000113;
        static const uint64_t GenericNack = 0x80000000;
        static const uint64_t BindReceiverResp = 0x80000001;
        static const uint64_t BindTransmitterResp = 0x80000002;
        static const uint64_t QuerySMResp = 0x80000003;
        static const uint64_t SubmitSMResp = 0x80000004;
        static const uint64_t DeliverSMResp = 0x80000005;
        static const uint64_t UnbindResp = 0x80000006;
        static const uint64_t ReplaceSMResp = 0x80000007;
        static const uint64_t CancelSMResp = 0x80000008;
        static const uint64_t BindTransceiverResp = 0x80000009;
        static const uint64_t EnquireLinkResp = 0x80000015;
        static const uint64_t SubmitMultiResp = 0x80000021;
        static const uint64_t DataSMResp = 0x80000103;
        static const uint64_t BroadcastSMResp = 0x80000111;
        static const uint64_t QueryBroadcastSMResp = 0x80000112;
        static const uint64_t CancelBroadcastSMResp = 0x80000113;
    };

    typedef struct { uint64_t cmdid; const char* name; } CmdStrings;

    #define CMDEXP(A) CmdID::A, #A

    CmdStrings cmdStrings[34] = {
        { CMDEXP(BindReceiver) },
        { CMDEXP(BindTransmitter) },
        { CMDEXP(QuerySM) },
        { CMDEXP(SubmitSM) },
        { CMDEXP(DeliverSM) },
        { CMDEXP(Unbind) },
        { CMDEXP(ReplaceSM) },
        { CMDEXP(CancelSM) },
        { CMDEXP(BindTransceiver) },
        { CMDEXP(Outbind) },
        { CMDEXP(EnquireLink) },
        { CMDEXP(SubmitMulti) },
        { CMDEXP(AlertNotification) },
        { CMDEXP(DataSM) },
        { CMDEXP(BroadcastSM) },
        { CMDEXP(QueryBroadcastSM) },
        { CMDEXP(CancelBroadcastSM) },
        { CMDEXP(GenericNack) },
        { CMDEXP(BindReceiverResp) },
        { CMDEXP(BindTransmitterResp) },
        { CMDEXP(QuerySMResp) },
        { CMDEXP(SubmitSMResp) },
        { CMDEXP(DeliverSMResp) },
        { CMDEXP(UnbindResp) },
        { CMDEXP(ReplaceSMResp) },
        { CMDEXP(CancelSMResp) },
        { CMDEXP(BindTransceiverResp) },
        { CMDEXP(EnquireLinkResp) },
        { CMDEXP(SubmitMultiResp) },
        { CMDEXP(DataSMResp) },
        { CMDEXP(BroadcastSMResp) },
        { CMDEXP(QueryBroadcastSMResp) },
        { CMDEXP(CancelBroadcastSMResp) }
    };

    const char* cmdString(uint64_t id) {
        for(int i=0;i<34;i++) {
            if ( cmdStrings[i].cmdid == id ) return cmdStrings[i].name;
        }
        return "";
    }

    static void GSMTimeStringShort( time_t& t, char* szTimestamp, int nLen )
    {
        if ( t < 0 ) return;

        szTimestamp[0] = 0x00;

        tm* ptm = gmtime( &t );

        if ( ptm == NULL ) szTimestamp[0] = 0x00;
        else strftime( szTimestamp, nLen, "%y%m%d%H%M", ptm );
    }

    static time_t GSMStringTime( const char* szBuf )
    {
        tm atm;
        int tdif = 0;
        char chDir = 0;

        memset( &atm, 0, sizeof( tm ) );

        int nRead = sscanf( szBuf, "%02d%02d%02d%02d%02d%02d%*1d%02d%c",
            &atm.tm_year, &atm.tm_mon, &atm.tm_mday, &atm.tm_hour, &atm.tm_min, &atm.tm_sec,
            &tdif, &chDir);

        time_t t = 0;

        if ( nRead == 8 )
        {
            if ( atm.tm_year >= 70 ) atm.tm_year += 1900; else atm.tm_year += 2000;

            t = timegm( &atm );

            if ( t == -1 ) {  // invalid time
                t = 0;
            }

            if ( chDir == '+' ) { t = t - (tdif*15*60); }
            if ( chDir == '-' ) { t = t + (tdif*15*60); }
        }

        return t;
    }

    static time_t GSMRelativeTime( char* szBuf )
    {
        time_t t = 0;
        int mday, hour, min, sec;
        int nRead = sscanf( szBuf, "%*02d%*02d%02d%02d%02d%02d%*1d%*02d%*c", &mday, &hour, &min, &sec );
        if ( nRead == 4 ) t = ( ( mday * 24 + hour ) * 60 + min ) * 60 + sec;
        return t;
    }
};

class SMPPConnection {
public:
    class SMPPException {

    };

private:
    const int max_command_body_length = 1024;

    bool debug = false;

    SMPPSocket* socket = NULL;

    char ip[64];

    typedef enum { PV_INVALID, PV_LENGTH, PV_ALL } PDUValidity;

    typedef struct {
        PDUValidity valid; // indicate if this header structure has valid contents
        uint64_t command_length;
        uint64_t command_id;
        uint64_t command_status;
        uint64_t sequence_number;
    } PDUHeader;

    typedef struct {
        bool valid;
        uint8_t* body;
    } PDUBody;

    PDUHeader command_received_header = {PV_INVALID,0,0,0,0};
    PDUBody command_received_body = {false,NULL};

    // configuration
    int enquire_link_period = 0; // seconds

    // methods

    uint64_t getInteger( void )
    {
        if ( socket->bytes_to_read()<4 ) throw SMPPException();
        uint8_t v0,v1,v2,v3;
        socket->recvA(v0);
        socket->recvA(v1);
        socket->recvA(v2);
        socket->recvA(v3);
        return (uint64_t)(((v0<<8) | v1) <<8 | v2) <<8 | v3;
    }

    uint64_t getBytes( uint64_t len, uint8_t* mem )
    {
        if ( socket->bytes_to_read() < len ) throw SMPPException();
        for(int i=0;i<len;i++) {
            if (socket->recvA(mem[i])==0) throw SMPPException();
        }
        return 0;
    }

public:
    SMPPConnection()
    {
        socket = NULL;
        command_received_header.valid = PV_INVALID;
        command_received_body.body = new uint8_t[max_command_body_length];
        ip[0] = '\0';
    }

    ~SMPPConnection()
    {
        if ( command_received_body.body ) delete []command_received_body.body;

        if ( socket ) delete socket;
    }

    void setDebug(bool val) { debug = val; }

    void allocateSocket( void )
    {
        socket = new SMPPSocketUnencrypted;
    }

    void allocateSocket( int fdsocket )
    {
        socket = new SMPPSocketUnencrypted(fdsocket);
    }

    void setIP( char* ip_in )
    {
        if (ip_in != NULL) strcpy(ip,ip_in);
        else strcpy(ip,"");
    }

    char* getIP() { return ip; }

    uint64_t endian( uint64_t a) {
        uint64_t b;
        ((uint8_t*)&b)[0] = ((uint8_t*)&a)[3];
        ((uint8_t*)&b)[1] = ((uint8_t*)&a)[2];
        ((uint8_t*)&b)[2] = ((uint8_t*)&a)[1];
        ((uint8_t*)&b)[3] = ((uint8_t*)&a)[0];
        return b;
    }

    bool put(uint64_t sequence_number, uint64_t cmdID, uint64_t status, uint8_t* param, int len)
    {
        int pdu_len = 16+len;
        uint8_t buf[pdu_len];

        *((uint64_t*) (buf+0)) = endian(pdu_len); // command length
        *((uint64_t*) (buf+4)) = endian(cmdID); // command ID
        *((uint64_t*) (buf+8)) = endian(status); // command status
        *((uint64_t*) (buf+12)) = endian(sequence_number); // sequence number
        if ((param!=NULL)&&(len!=0)) memcpy((char*)(buf+16),param,len);

        //for(int i=0;i<pdu_len;i++) printf("%02x ",buf[i]);
        //printf("\n");

        return socket->send(buf,pdu_len);
    }

    uint8_t* getBodyPointer(int& len)
    {
        len = (int)(command_received_header.command_length-16);
        return command_received_body.body;
    }

    bool get()
    {
        // require "socket"

        if ( socket == NULL ) return false;

        // get header

        if ( command_received_header.valid != PV_ALL )
        {
            // don't have any of the header yet

            try {
                command_received_header.command_length = getInteger();
                command_received_header.command_id = getInteger();
                command_received_header.command_status = getInteger();
                command_received_header.sequence_number = getInteger();
                command_received_header.valid = PV_ALL; // all header values are now valid

            } catch (SMPPException e) {
                // no data
                if (debug) std::cout << "HEADER No data - ABORT!" << std::endl;
                return false;
            }
        }

        if (command_received_header.valid != PV_ALL) return false; // don't let past here until we have the header

        // get body (if present)

        command_received_body.valid = false;

        uint64_t body_length = command_received_header.command_length-16;

        if ( body_length > 0 )
        {
            try {
                getBytes(body_length,command_received_body.body);
            }
            catch(SMPPException e) {
                if (debug) std::cout << "BODY No data" << std::endl;
                return false;
            }
            command_received_body.valid = true;
            command_received_header.valid = PV_INVALID;

            return true;
        }
        else {
            // empty body
            command_received_body.valid = true; // empty body (valid)
            //command_received_body.body = NULL;
            command_received_header.valid = PV_INVALID;

            return true;
        }

        return false;
    }

    uint64_t pduCommandID( void ) { return command_received_header.command_id; }
    uint64_t pduSequenceNo( void ) { return command_received_header.sequence_number; }
};

class Session
{
protected:
    int sessionType;

public:
    Session() {}
    virtual ~Session() {}

    virtual bool timedCheck( void ) = 0;
    virtual bool run( void ) = 0;
};

class AdminSession : public Session
{
public:
    AdminSession() { sessionType = 0; }
    AdminSession(int fdsocket)
    {
        //conn.allocateSocket(fdsocket);
    }
    ~AdminSession() {}

    bool timedCheck( void )
    {
        // return true if to close
        return true;
    }

    bool run( void )
    {
        return true;
    }
};

class SMPPSession : public Session
{
private:
    SMPPConnection conn;

    uint64_t session_id;

    // session state
    string system_id;
    string system_type;
    long sequence_number_in = -1;
    long sequence_number_out = 1;
    uint64_t lastPDUFromESMETime = 0;
    uint64_t enquireLinkRespPending = 0;
    uint64_t closingTime = 0;

public:
    typedef enum { BS_NONE, BS_TRX, BS_TX, BS_RX } BindState;
    uint8_t version;

    BindState bindState;

public:
    SMPPSession()
    {
        sessionType = 1;

        bindState = BS_NONE;
        version = 0x00;

        session_id = session_id_next++;
    }

    SMPPSession(int fdsocket,char* ip)
    {
        bindState = BS_NONE;
        version = 0x00;
        conn.allocateSocket(fdsocket);
        conn.setIP(ip);

        session_id = session_id_next++;
    }

    ~SMPPSession() {
        if (bindState != BS_NONE) {
            logCommand("session aborted", "--");
        }
    }

    void setDebug( bool val )
    {
        conn.setDebug(val);
    }

    bool getCOctetString(uint8_t* buf_ptr,int buf_len,int& idx,string& strout,int max_param_len)
    {
        char str[max_param_len];
        int i= 0;
        while (idx!=buf_len)
        {
            str[i++] = buf_ptr[idx++];
            if (str[i-1] == '\0') {
                strout.assign(str);
                return true;
            }
            if (i==max_param_len) return false;
        }
        return false;
    }

    bool timedCheck( void )
    {
        // return true if session to close

        uint64_t now = currentUSecsSinceEpoch();

        // session management
        //

        if ((closingTime!=0)&&( now>=closingTime )) return true; // session was due to close now

        if (( lastPDUFromESMETime!= 0 ) && ( (now-lastPDUFromESMETime) > (120*1000000L) )) // wait up to 2 mins of inactivity before enquire link
        {
            if (enquireLinkRespPending==0)
            {
                // issue enquire link
                send(sequence_number_out++, SMPP::CmdID::EnquireLink, SMPP::CmdStatus::ESME_ROK, NULL, 0);
                enquireLinkRespPending = now;
                closingTime = 0;
            }
            else
            {
                if ( closingTime==0 )
                {
                    // waiting on enquire link
                    uint64_t tsinceenquirelinksent = now-enquireLinkRespPending;
                    if ( tsinceenquirelinksent > (60*1000000L) ) // wait up to 1 min on enquire_link response
                    {
                        send(sequence_number_out++, SMPP::CmdID::Unbind, SMPP::CmdStatus::ESME_ROK, NULL, 0);

                        closingTime = now + 5*1000000L; // close session in 5 seconds time
                    }
                }
            }
        }

        // simulate message delivery + delivery receipt generation
        //

        if ((system_id.length()>0) && ((bindState == SMPPSession::BindState::BS_TRX)||(bindState == SMPPSession::BindState::BS_RX)))
        {
            MessageDeliverer::Message msg;
            while (MessageDeliverer::getInstance_get(system_id,msg))
            {
                // indicate message delivered

                if ( msg.registered_delivery != 0 ) // ESME requested receipt so generate one
                {
                    uint8_t source_addr_ton = msg.source_addr_ton, source_addr_npi = msg.source_addr_npi;
                    string destination_addr = msg.destination_addr;
                    uint8_t dest_addr_ton = msg.dest_addr_ton, dest_addr_npi = msg.dest_addr_npi;
                    string source_addr = msg.source_addr;
                    string msgid = msg.smscMessageID;

                    generateReceipt(source_addr_ton,source_addr_npi,source_addr,dest_addr_ton,dest_addr_npi,destination_addr,msgid);
                }

                if ( strstr(msg.destination_addr,system_id.c_str()) != NULL ) { // system_id is in destination address - assume MO for testing purposes
                    generateMO(msg.source_addr_ton,msg.source_addr_npi,msg.source_addr,
                               msg.dest_addr_ton,msg.dest_addr_npi,msg.destination_addr,
                               msg.short_message,msg.sm_length);
                }
            }
        }

        return false;
    }

    void generateMO(uint8_t source_addr_ton,uint8_t source_addr_npi,string source_addr,
                    uint8_t dest_addr_ton,uint8_t dest_addr_npi,string destination_addr,
                    uint8_t* short_message, uint8_t sm_length)
    {
        uint8_t sbuf[1024];

        int sidx=0;

        sbuf[sidx++] = 0x00; // service type

        sbuf[sidx++] = source_addr_ton; //
        sbuf[sidx++] = source_addr_npi; //
        memcpy(sbuf+sidx, source_addr.c_str(), source_addr.length()+1); // destination_addr
        sidx += source_addr.length()+1;

        sbuf[sidx++] = dest_addr_ton; //
        sbuf[sidx++] = dest_addr_npi; //
        memcpy(sbuf+sidx, destination_addr.c_str(), destination_addr.length()+1); // source_addr
        sidx += destination_addr.length()+1;

        sbuf[sidx++] = 0x00; // esm_class
        sbuf[sidx++] = 0x00; // protocol_id
        sbuf[sidx++] = 0x00; // priority_flag
        sbuf[sidx++] = 0x00; // schedule_delivery_time
        sbuf[sidx++] = 0x00; // validity_period
        sbuf[sidx++] = 0x00; // registered_delivery
        sbuf[sidx++] = 0x00; // replace_if_present_flag
        sbuf[sidx++] = 0x00; // data_coding
        sbuf[sidx++] = 0x00; // sm_default_msg_id

        sbuf[sidx++] = sm_length; // sm_length

        memcpy((char*)(sbuf+sidx),short_message,sm_length);

        sidx += sm_length;

        send(sequence_number_out++,SMPP::CmdID::DeliverSM,SMPP::CmdStatus::ESME_ROK,sbuf,sidx);
    }

    void generateReceipt(uint8_t source_addr_ton,uint8_t source_addr_npi,string source_addr,
                         uint8_t dest_addr_ton,uint8_t dest_addr_npi,string destination_addr,
                         string msgid)
    {
        // receipt

        if ((bindState == SMPPSession::BindState::BS_TRX)||(bindState == SMPPSession::BindState::BS_RX))
        {
            uint8_t sbuf[1024];

            int sidx=0;
            sbuf[sidx++] = 0x00; // service type
            sbuf[sidx++] = dest_addr_ton; // source_addr_ton
            sbuf[sidx++] = dest_addr_npi; // source_addr_npi
            memcpy(sbuf+sidx, destination_addr.c_str(), destination_addr.length()+1); // source_addr
            sidx += destination_addr.length()+1;
            sbuf[sidx++] = source_addr_ton; // dest_addr_ton
            sbuf[sidx++] = source_addr_npi; // dest_addr_npi
            memcpy(sbuf+sidx, source_addr.c_str(), source_addr.length()+1); // destination_addr
            sidx += source_addr.length()+1;
            sbuf[sidx++] = 0x04; // esm_class [0x04 Short Message contains MC delivery report]
            sbuf[sidx++] = 0x00; // protocol_id
            sbuf[sidx++] = 0x00; // priority_flag
            sbuf[sidx++] = 0x00; // schedule_delivery_time
            sbuf[sidx++] = 0x00; // validity_period
            sbuf[sidx++] = 0x00; // registered_delivery
            sbuf[sidx++] = 0x00; // replace_if_present_flag
            sbuf[sidx++] = 0x00; // data_coding
            sbuf[sidx++] = 0x00; // sm_default_msg_id

            // - short_message containing receipt in text format
            time_t tSubmitStamp = time(NULL);
            time_t tDoneStamp = time(NULL);
            char szSubmitStamp[32];
            char szDoneStamp[32];
            SMPP::GSMTimeStringShort( tSubmitStamp, szSubmitStamp, sizeof( szSubmitStamp ) );
            SMPP::GSMTimeStringShort( tDoneStamp, szDoneStamp, sizeof( szDoneStamp) );

            char short_message[160];
            sprintf( short_message, "id:%s sub:000 dlvrd:%03d submit date:%s done date:%s stat:%s err:%03d text:",
                msgid.c_str(),
                1 /* 1 message delivered */,
                szSubmitStamp,
                szDoneStamp,
                "DELIVRD",
                0 /*error*/ );

            sbuf[sidx++] = strlen(short_message)+1; // sm_length

            memcpy((char*)(sbuf+sidx),short_message,(int)(strlen(short_message)+1));

            sidx += strlen(short_message)+1;

            // TLVs

            if (bindState == SMPPSession::BindState::BS_TRX)
            {
                uint8_t params1[] = {
                    // message_state TLV
                    0x04, 0x27, // tag
                    0x00, 0x01, // length
                    0x02, // value - delivered

                    // network_error TLV
                    0x04, 0x23, // tag
                    0x00, 0x03, // length
                    0x03, 0x00, 0x00, // value - GSM, 0, 0

                    // receipted_message_id TLV
                    0x00, 0x1e, // tag
                    0x00, 0x41, // length
                    // .. message_id (to be appended)
                };

                memcpy(sbuf+sidx,params1,sizeof(params1));
                sidx += sizeof(params1);
                memcpy(sbuf+sidx,msgid.c_str(),msgid.length()+1);
                sidx += msgid.length()+1;
            }

            send(sequence_number_out++,SMPP::CmdID::DeliverSM,SMPP::CmdStatus::ESME_ROK,sbuf,sidx);
        }
    }

    bool run( void )
    {
        // return true if closed

        uint64_t now = currentUSecsSinceEpoch();

        if ((closingTime!=0)&&( now>=closingTime )) return true; // session was due to close now

        // handle PDUs from ESME
        //

        bool allowBind = true;

        uint64_t cmdid,seqno;
        uint8_t sbuf[1024];

        if (recv(cmdid,seqno))
        {
            //SMPP b;
            //printf("<< 0x%08llx %s\n",cmdid,b.cmdString(cmdid));

            lastPDUFromESMETime = now;

            if (( cmdid == SMPP::CmdID::BindTransceiver )||( cmdid == SMPP::CmdID::BindReceiver )||( cmdid == SMPP::CmdID::BindTransmitter ))
            {
                if ( bindState != SMPPSession::BindState::BS_NONE ) // PROTOCOL must be in unbound state for ESME to send bind command
                {
                    send(seqno,cmdid+0x80000000, SMPP::CmdStatus::ESME_RINVBNDSTS, NULL, 0);
                }
                else {
                    string esme_system_id = "";
                    string esme_password = "";
                    string esme_system_type ="";
                    uint8_t esme_smpp_ver = 0x00;

                    int ptr_max = 0;
                    uint8_t* ptr = conn.getBodyPointer(ptr_max);

                    int idx = 0;
                    if (!getCOctetString(ptr,ptr_max,idx,esme_system_id,16)) { allowBind = false; goto parse_complete; }
                    if (!getCOctetString(ptr,ptr_max,idx,esme_password,9)) { allowBind = false; goto parse_complete; }
                    if (!getCOctetString(ptr,ptr_max,idx,esme_system_type,13)) { allowBind = false; goto parse_complete; }

                    if (idx<ptr_max) esme_smpp_ver = ptr[idx++];
                    else { allowBind = false; goto parse_complete; }

                    if (( esme_smpp_ver < 0x34 ) && (cmdid == SMPP::CmdID::BindTransceiver )) allowBind = false; // PROTOCOL transceiver bind only allowed for v3.4+

                    parse_complete:  // parse complete label

                    if (allowBind)
                    {
                        system_id = esme_system_id;

                        char smsc_system_id[] = "MelroseLabsSMSC";
                        memcpy(sbuf,smsc_system_id,strlen(smsc_system_id)+1);

                        if (ClientConfig::instance().is(esme_system_id,ClientConfig::ConfigEntry::BINDRESP_WITH_SC_INTERACE_VERSION_TLV))
                        {
                            uint8_t params[] = {
                                // sc_interface_version TLV
                                0x02, 0x10, // tag
                                0x00, 0x01, // length
                                0x34, // value [v3.4 supported]
                            };

                            memcpy(sbuf+strlen(smsc_system_id)+1,params,sizeof(params));
                            send(seqno,cmdid+0x80000000, SMPP::CmdStatus::ESME_ROK, sbuf, strlen(smsc_system_id)+1+sizeof(params));
                        }
                        else send(seqno,cmdid+0x80000000, SMPP::CmdStatus::ESME_ROK, sbuf, strlen(smsc_system_id)+1);

                        if ( cmdid == SMPP::CmdID::BindTransceiver ) bindState = SMPPSession::BindState::BS_TRX;
                        else if ( cmdid == SMPP::CmdID::BindTransmitter ) bindState = SMPPSession::BindState::BS_TX;
                        else if ( cmdid == SMPP::CmdID::BindReceiver ) bindState = SMPPSession::BindState::BS_RX;
                        else bindState = SMPPSession::BindState::BS_NONE;

                        version = esme_smpp_ver;
                    }
                    else
                    {
                        send(seqno,SMPP::CmdID::BindTransceiverResp, SMPP::CmdStatus::ESME_RBINDFAIL, NULL, 0); // responding indicating bind unsuccessful
                        bindState = SMPPSession::BindState::BS_NONE;
                    }
                }
            }
            else if ( cmdid == SMPP::CmdID::Unbind )
            {
                if (bindState != SMPPSession::BindState::BS_NONE)
                {
                    send(seqno,SMPP::CmdID::UnbindResp, SMPP::CmdStatus::ESME_ROK, NULL, 0);
                    bindState = SMPPSession::BindState::BS_NONE;
                }
                else
                {
                    send(seqno,SMPP::CmdID::UnbindResp, SMPP::CmdStatus::ESME_RINVBNDSTS, NULL, 0);
                }
            }
            else if ( cmdid == SMPP::CmdID::EnquireLinkResp )
            {
                enquireLinkRespPending = 0; // received pending enquire link
            }
            else if ( cmdid == SMPP::CmdID::EnquireLink )
            {
                send(seqno,SMPP::CmdID::EnquireLinkResp, SMPP::CmdStatus::ESME_ROK, NULL, 0);
            }
            else if ( cmdid == SMPP::CmdID::GenericNack )
            {
                // no response to GenericNack
            }
            else if ( cmdid == SMPP::CmdID::SubmitSM )
            {
                if ((bindState == SMPPSession::BindState::BS_NONE)||(bindState == SMPPSession::BindState::BS_RX))
                {
                    send(seqno,SMPP::CmdID::SubmitSMResp,SMPP::CmdStatus::ESME_RINVBNDSTS,NULL,0);
                }
                else
                {
                    // process requested

                    uint64_t tNow = currentUSecsSinceEpoch();

                    uint64_t tESMEDeliveryTimeRequired = tNow; // default to immediate delivery

                    bool allowSubmit = true;
                    uint64_t smppSubmitError = SMPP::CmdStatus::ESME_RSUBMITFAIL;

                    string service_type = "";
                    uint8_t source_addr_ton = 0;
                    uint8_t source_addr_npi = 0;
                    string source_addr = "";
                    uint8_t dest_addr_ton = 0;
                    uint8_t dest_addr_npi = 0;
                    string destination_addr = "";
                    string schedule_delivery_time = "";
                    string validity_period = "";
                    uint8_t registered_delivery = 0;
                    uint8_t sm_length;
                    uint8_t short_message[160];

                    int ptr_max = 0;
                    uint8_t* ptr = conn.getBodyPointer(ptr_max);

                    int idx = 0;
                    if (!getCOctetString(ptr,ptr_max,idx,service_type,6)) { allowSubmit = false; goto parse_complete_submit; }
                    if (idx<ptr_max) {source_addr_ton = ptr[idx++];} else { allowSubmit = false; goto parse_complete_submit; }
                    if (idx<ptr_max) {source_addr_npi = ptr[idx++];} else { allowSubmit = false; goto parse_complete_submit; }
                    if (!getCOctetString(ptr,ptr_max,idx,source_addr,21)) { allowSubmit = false; goto parse_complete_submit; }
                    if (idx<ptr_max) {dest_addr_ton = ptr[idx++];} else { allowSubmit = false; goto parse_complete_submit; }
                    if (idx<ptr_max) {dest_addr_npi = ptr[idx++];} else { allowSubmit = false; goto parse_complete_submit; }
                    if (!getCOctetString(ptr,ptr_max,idx,destination_addr,21)) { allowSubmit = false; goto parse_complete_submit; }

                    idx++; // esm_class
                    idx++; // protocol_id
                    idx++; // priority_flag
                    if (!getCOctetString(ptr,ptr_max,idx,schedule_delivery_time,17)) { allowSubmit = false; goto parse_complete_submit; }
                    if (!getCOctetString(ptr,ptr_max,idx,validity_period,17)) { allowSubmit = false; goto parse_complete_submit; }
                    registered_delivery = ptr[idx++]; // registered_delivery
                    idx++; // replace_if_present_flag
                    idx++; // data_coding
                    idx++; // sm_default_msg_id
                    sm_length = ptr[idx++];
                    memcpy( short_message, ptr+idx, sm_length );
                    idx += sm_length;

                    if (schedule_delivery_time.length()>0) // scheduled delivery
                    {
                        if (schedule_delivery_time[schedule_delivery_time.length()-1] == 'R') // relative
                            tESMEDeliveryTimeRequired = tNow + SMPP::GSMRelativeTime((char*)schedule_delivery_time.c_str())*1000000L;
                        else
                            tESMEDeliveryTimeRequired = (SMPP::GSMStringTime(schedule_delivery_time.c_str()) + 3 + (rand()%10))*1000000L;

                        if (tESMEDeliveryTimeRequired < tNow) // scheduled delivery time is before now
                        {
                            allowSubmit = false;
                            smppSubmitError = SMPP::CmdStatus::ESME_RINVSCHED;
                            goto parse_complete_submit;
                        }
                    }

                    // ...

                    if (destination_addr == "333") { allowSubmit = false; goto parse_complete_submit; } // force ESME_RSUBMITFAIL

                    if (destination_addr.length()<8) { allowSubmit = false; smppSubmitError = SMPP::CmdStatus::ESME_INVDSTADR; goto parse_complete_submit; }

                    parse_complete_submit:

                    // send response

                    if (allowSubmit)
                    {
                        // accept message and send message ID back to ESME

                        int msgid_len = 64; // PROTOCOL v3.4 message_id field is COctet String of up to 65 chars (inc NULL terminator)
                        if ( version < 0x34 ) msgid_len = 8; // PROTOCOL v3.3 message_id field is COctet String (hex) of up to 9 characters (inc NULL terminator)

                        char msgid[msgid_len+1];
                        for(int i=0;i<msgid_len;i++) msgid[i] = (rand()%16)<10?('0'+rand()%10):('a'+rand()%6);
                        msgid[msgid_len] = 0;
                        send(seqno,SMPP::CmdID::SubmitSMResp,SMPP::CmdStatus::ESME_ROK,(uint8_t*)msgid,msgid_len+1);

                        // add message to deliverer so that receipt will then be sent to ESME

                        uint64_t timeDeliver = tESMEDeliveryTimeRequired + (3+(rand()%10))*1000000L;



                        MessageDeliverer::Message msg;
                        msg.setSource(source_addr_ton,source_addr_npi,(char*)source_addr.c_str());
                        msg.setDestination(dest_addr_ton,dest_addr_npi,(char*)destination_addr.c_str());
                        msg.setRegisteredDelivery(registered_delivery);
                        msg.setShortMessage(short_message, sm_length);
                        msg.setSMSCMessageID(msgid);
                        MessageDeliverer::getInstance(system_id)->add(timeDeliver,msg);
                    }
                    else send(seqno,SMPP::CmdID::SubmitSMResp, smppSubmitError, NULL, 0); // don't accept message and send NAK back to ESME
                }
            }
        }
        else {
            // error or closed
            return true;
        }
        return false;
    }

    uint8_t getVersion(void) { return version; }
    void setVersion(uint8_t version_in) { version = version_in; }

    void logCommand(uint64_t cmdID,const char* direction)
    {
        char buf[640];
        SMPP b;
        sprintf(buf,"0x%08llx %s",cmdID,b.cmdString(cmdID));
        logCommand(buf,direction);
    }

    void logCommand(char* logline,const char* direction)
    {
        time_t rawtime;
        struct tm * timeinfo;
        char buffer [80];

        time (&rawtime);
        timeinfo = localtime (&rawtime);

        strftime (buffer,80,"%F %X ",timeinfo);

        printf("%s %s [s%06llx:%-15s] %s\n",buffer,direction,session_id,conn.getIP(),logline);
    }

    bool recv( uint64_t& cmdID, uint64_t& seqNo )
    {
        bool ret = conn.get();
        if (ret) { cmdID = conn.pduCommandID(); seqNo  = conn.pduSequenceNo();
            logCommand(cmdID,">>");
        }
        return ret;
    }

    bool send( uint64_t seqNo, uint64_t cmdID, uint64_t cmdStatus, uint8_t* param, int len )
    {
        //SMPP b;
        //printf("<< 0x%08llx %s\n",cmdID,b.cmdString(cmdID));
        logCommand(cmdID,"<<");

        return conn.put(seqNo, cmdID, cmdStatus, param, len);
    }
};

// Sockets reference: https://www.ibm.com/support/knowledgecenter/ssw_ibm_i_74/rzab6/xnonblock.htm

int dolisten( int portno )
{
    struct sockaddr_in serv_addr;

    /* First call to socket() function */
    int listensockfd = socket(AF_INET, SOCK_STREAM, 0);

    if (listensockfd < 0) {
        perror("ERROR opening socket");
        return -1;
    }

    int rc, on = 1;

    /*************************************************************/
    /* Allow socket descriptor to be reuseable                   */
    /*************************************************************/
    rc = setsockopt(listensockfd, SOL_SOCKET,  SO_REUSEADDR,
                    (char *)&on, sizeof(on));
    if (rc < 0)
    {
        perror("setsockopt() failed");
        close(listensockfd);
        exit(-1);
    }

    /*************************************************************/
    /* Set socket to be nonblocking. All of the sockets for      */
    /* the incoming connections will also be nonblocking since   */
    /* they will inherit that state from the listening socket.   */
    /*************************************************************/
    rc = ioctl(listensockfd, FIONBIO, (char *)&on);
    if (rc < 0)
    {
        perror("ioctl() failed");
        close(listensockfd);
        exit(-1);
    }

    /* Initialize socket structure */
    bzero((char *) &serv_addr, sizeof(serv_addr));

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(portno);

    /* Now bind the host address using bind() call.*/
    if (::bind(listensockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0) {
        perror("ERROR on binding");
        return -1;
    }

    /* Now start listening for the clients, here
     * process will go in sleep mode and will wait
     * for the incoming connection
     */

    listen(listensockfd,32);

    return listensockfd;
}

Session* sessionSockMap[32000]; // array of SMPP session by socket

int main(int argc, const char * argv[])
{
    printf("%s build time: %s %s\n",argv[0],__DATE__,__TIME__);

    //

    int portSMPP = 2775;

    int listensockfdSMPP = dolisten(portSMPP);

    if (listensockfdSMPP==-1) {
        perror("Failed to listen on SMPP port");
        exit(1);
    }

    std::cout << "Listening for SMPP on port " << portSMPP << std::endl;

    //

    int portAdmin = 8775;

    int listensockfdAdmin = dolisten(portAdmin);

    if (listensockfdAdmin==-1) {
        perror("Failed to listen on admin port");
        exit(1);
    }

    std::cout << "Listening for admin on port " << portAdmin << std::endl;

    // Initialize the master fd_set
    fd_set master_set, working_set;
    FD_ZERO(&master_set);
    int max_sd = listensockfdSMPP;
    FD_SET(listensockfdSMPP, &master_set);
    if ( listensockfdAdmin > max_sd ) max_sd = listensockfdAdmin;
    FD_SET(listensockfdAdmin, &master_set);

    //

    int rc,i,desc_ready;
    int end_server = FALSE;
    int new_sd;
    int close_conn;
    struct timeval timeout;

    // Loop waiting for incoming connects or for incoming data
    do
    {
        // Copy the master fd_set over to the working fd_set.
        memcpy(&working_set, &master_set, sizeof(master_set));

        // Initialize the timeval struct to 1 second.  If no
        // activity after 1 second then wake-up and cycle.
        timeout.tv_sec  = 1;
        timeout.tv_usec = 0;

        // Call select() and wait for it to complete.
        rc = select(max_sd + 1, &working_set, NULL, NULL, &timeout);

        // Check to see if the select call failed.
        if (rc < 0)
        {
            perror("  select() failed");
            break;
        }

        // Check to see if select call timed out.
        if (rc == 0)
        {
            // perform periodic session tasks

            for (i=0; i <= max_sd; ++i)
            {
                if (sessionSockMap[i]!=NULL)
                    if (sessionSockMap[i]->timedCheck())
                    {
                        printf("  %ld Force close on connection - %d.\n",time(NULL),i);

                        if ( sessionSockMap[i] != NULL )
                        {
                            delete sessionSockMap[i];
                            sessionSockMap[i] = NULL;
                        }

                        /*
                         If the close_conn flag was turned on, we need
                         to clean up this active connection.  This
                         clean up process includes removing the
                         descriptor from the master set and
                         determining the new maximum descriptor value
                         based on the bits that are still turned on in
                         the master set.
                        */

                        close(i);
                        FD_CLR(i, &master_set);
                        if (i == max_sd)
                        {
                            while (FD_ISSET(max_sd, &master_set) == FALSE) { max_sd -= 1; }
                        }
                    }
            }

            continue;
        }

        /**********************************************************/
        /* One or more descriptors are readable.  Need to         */
        /* determine which ones they are.                         */
        /**********************************************************/
        desc_ready = rc;
        for (i=0; i <= max_sd  &&  desc_ready > 0; ++i)
        {
            /*******************************************************/
            /* Check to see if this descriptor is ready            */
            /*******************************************************/
            if (FD_ISSET(i, &working_set))
            {
                /****************************************************/
                /* A descriptor was found that was readable - one   */
                /* less has to be looked for.  This is being done   */
                /* so that we can stop looking at the working set   */
                /* once we have found all of the descriptors that   */
                /* were ready.                                      */
                /****************************************************/
                desc_ready -= 1;

                /***********************************************************/
                /* Check to see if this is the listening socket (SMPP)     */
                /***********************************************************/
                if (i == listensockfdSMPP)
                {
                    //printf("  Listening socket (SMPP) is readable\n");

                    /*************************************************/
                    /* Accept all incoming connections that are      */
                    /* queued up on the listening socket before we   */
                    /* loop back and call select again.              */
                    /*************************************************/
                    do
                    {
                        /**********************************************/
                        /* Accept each incoming connection.  If       */
                        /* accept fails with EWOULDBLOCK, then we     */
                        /* have accepted all of them.  Any other      */
                        /* failure on accept will cause us to end the */
                        /* server.                                    */
                        /**********************************************/
                        struct sockaddr_in client_addr;
                        int clen = sizeof(sockaddr_in);
                        new_sd = accept(listensockfdSMPP, (struct sockaddr *)&client_addr, (socklen_t*) &clen);
                        if (new_sd < 0)
                        {
                            if (errno != EWOULDBLOCK)
                            {
                                perror("  accept() failed (SMPP)");
                                if (errno != EMFILE) end_server = TRUE;
                            }
                            break;
                        }

                        char* ip = inet_ntoa(client_addr.sin_addr);

                        //

                        SMPPSession* newsession = new SMPPSession(new_sd,ip);

                        sessionSockMap[new_sd] = newsession;

                        /**********************************************/
                        /* Add the new incoming connection to the     */
                        /* master read set                            */
                        /**********************************************/
                        //printf("  New incoming connection - %d\n", new_sd);
                        FD_SET(new_sd, &master_set);
                        if (new_sd > max_sd)
                            max_sd = new_sd;

                        /**********************************************/
                        /* Loop back up and accept another incoming   */
                        /* connection                                 */
                        /**********************************************/
                    } while (new_sd != -1);
                }
                else
                /***********************************************************/
                /* Check to see if this is the listening socket (admin)    */
                /***********************************************************/
                if (i == listensockfdAdmin)
                {
                    //printf("  Listening socket (admin) is readable\n");
                    /*************************************************/
                    /* Accept all incoming connections that are      */
                    /* queued up on the listening socket before we   */
                    /* loop back and call select again.              */
                    /*************************************************/
                    do
                    {
                        /**********************************************/
                        /* Accept each incoming connection.  If       */
                        /* accept fails with EWOULDBLOCK, then we     */
                        /* have accepted all of them.  Any other      */
                        /* failure on accept will cause us to end the */
                        /* server.                                    */
                        /**********************************************/
                        new_sd = accept(listensockfdAdmin, NULL, NULL);
                        if (new_sd < 0)
                        {
                            if (errno != EWOULDBLOCK)
                            {
                                perror("  accept() failed (admin)");
                                if (errno != EMFILE) end_server = TRUE;
                            }
                            break;
                        }

                        //

                        AdminSession* newsession = new AdminSession(new_sd);

                        sessionSockMap[new_sd] = newsession;

                        /**********************************************/
                        /* Add the new incoming connection to the     */
                        /* master read set                            */
                        /**********************************************/
                        //printf("  New incoming connection (Admin) - %d\n", new_sd);
                        FD_SET(new_sd, &master_set);
                        if (new_sd > max_sd)
                            max_sd = new_sd;

                        /**********************************************/
                        /* Loop back up and accept another incoming   */
                        /* connection                                 */
                        /**********************************************/
                    } while (new_sd != -1);
                }

                /****************************************************/
                /* This is not the listening socket, therefore an   */
                /* existing connection must be readable             */
                /****************************************************/
                else
                {
                    //printf("  Descriptor %d is readable\n", i);
                    close_conn = FALSE;

                    bool closed = sessionSockMap[i]->run();

                    if (closed)
                    {
                        //printf("  Closing connection - %d\n", i);

                        if ( sessionSockMap[i] != NULL )
                        {
                            delete sessionSockMap[i];
                            sessionSockMap[i] = NULL;
                        }

                        close_conn = TRUE;
                    }
                    else close_conn = FALSE;

                    /*************************************************/
                    /* If the close_conn flag was turned on, we need */
                    /* to clean up this active connection.  This     */
                    /* clean up process includes removing the        */
                    /* descriptor from the master set and            */
                    /* determining the new maximum descriptor value  */
                    /* based on the bits that are still turned on in */
                    /* the master set.                               */
                    /*************************************************/
                    if (close_conn)
                    {
                        close(i);
                        FD_CLR(i, &master_set);
                        if (i == max_sd)
                        {
                            while (FD_ISSET(max_sd, &master_set) == FALSE)
                                max_sd -= 1;
                        }
                    }
                } /* End of existing connection is readable */
            } /* End of if (FD_ISSET(i, &working_set)) */
        } /* End of loop through selectable descriptors */

    } while (end_server == FALSE);

    //

    close(listensockfdSMPP);

    close(listensockfdAdmin);

    std::cout << "No longer listening on port " << portSMPP << std::endl;

    return 0;
}