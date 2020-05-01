//go:generate stringer -type=CommandStatusType,CommandIDType -output header_data_string.go
package data

type CommandStatusType int32

type CommandIDType int32

const (
	// SMPP Command ID Set
	GENERIC_NACK          = CommandIDType(-2147483648)
	BIND_RECEIVER         = CommandIDType(0x00000001)
	BIND_RECEIVER_RESP    = CommandIDType(-2147483647)
	BIND_TRANSMITTER      = CommandIDType(0x00000002)
	BIND_TRANSMITTER_RESP = CommandIDType(-2147483646)
	QUERY_SM              = CommandIDType(0x00000003)
	QUERY_SM_RESP         = CommandIDType(-2147483645)
	SUBMIT_SM             = CommandIDType(0x00000004)
	SUBMIT_SM_RESP        = CommandIDType(-2147483644)
	DELIVER_SM            = CommandIDType(0x00000005)
	DELIVER_SM_RESP       = CommandIDType(-2147483643)
	UNBIND                = CommandIDType(0x00000006)
	UNBIND_RESP           = CommandIDType(-2147483642)
	REPLACE_SM            = CommandIDType(0x00000007)
	REPLACE_SM_RESP       = CommandIDType(-2147483641)
	CANCEL_SM             = CommandIDType(0x00000008)
	CANCEL_SM_RESP        = CommandIDType(-2147483640)
	BIND_TRANSCEIVER      = CommandIDType(0x00000009)
	BIND_TRANSCEIVER_RESP = CommandIDType(-2147483639)
	OUTBIND               = CommandIDType(0x0000000B)
	ENQUIRE_LINK          = CommandIDType(0x00000015)
	ENQUIRE_LINK_RESP     = CommandIDType(-2147483627)
	SUBMIT_MULTI          = CommandIDType(0x00000021)
	SUBMIT_MULTI_RESP     = CommandIDType(-2147483615)
	ALERT_NOTIFICATION    = CommandIDType(0x00000102)
	DATA_SM               = CommandIDType(0x00000103)
	DATA_SM_RESP          = CommandIDType(-2147483389)
)

const (

	// Command_Status Error Codes
	ESME_ROK           = CommandStatusType(0x00000000)
	ESME_RINVMSGLEN    = CommandStatusType(0x00000001)
	ESME_RINVCMDLEN    = CommandStatusType(0x00000002)
	ESME_RINVCMDID     = CommandStatusType(0x00000003)
	ESME_RINVBNDSTS    = CommandStatusType(0x00000004)
	ESME_RALYBND       = CommandStatusType(0x00000005)
	ESME_RINVPRTFLG    = CommandStatusType(0x00000006)
	ESME_RINVREGDLVFLG = CommandStatusType(0x00000007)
	ESME_RSYSERR       = CommandStatusType(0x00000008)
	ESME_RINVSRCADR    = CommandStatusType(0x0000000A)
	ESME_RINVDSTADR    = CommandStatusType(0x0000000B)
	ESME_RINVMSGID     = CommandStatusType(0x0000000C)
	ESME_RBINDFAIL     = CommandStatusType(0x0000000D)
	ESME_RINVPASWD     = CommandStatusType(0x0000000E)
	ESME_RINVSYSID     = CommandStatusType(0x0000000F)
	ESME_RCANCELFAIL   = CommandStatusType(0x00000011)
	ESME_RREPLACEFAIL  = CommandStatusType(0x00000013)
	ESME_RMSGQFUL      = CommandStatusType(0x00000014)
	ESME_RINVSERTYP    = CommandStatusType(0x00000015)

	ESME_RADDCUSTFAIL  = CommandStatusType(0x00000019) // Failed to Add Customer
	ESME_RDELCUSTFAIL  = CommandStatusType(0x0000001A) // Failed to delete Customer
	ESME_RMODCUSTFAIL  = CommandStatusType(0x0000001B) // Failed to modify customer
	ESME_RENQCUSTFAIL  = CommandStatusType(0x0000001C) // Failed to Enquire Customer
	ESME_RINVCUSTID    = CommandStatusType(0x0000001D) // Invalid Customer ID
	ESME_RINVCUSTNAME  = CommandStatusType(0x0000001F) // Invalid Customer Name
	ESME_RINVCUSTADR   = CommandStatusType(0x00000021) // Invalid Customer Address
	ESME_RINVADR       = CommandStatusType(0x00000022) // Invalid Address
	ESME_RCUSTEXIST    = CommandStatusType(0x00000023) // Customer Exists
	ESME_RCUSTNOTEXIST = CommandStatusType(0x00000024) // Customer does not exist
	ESME_RADDDLFAIL    = CommandStatusType(0x00000026) // Failed to Add DL
	ESME_RMODDLFAIL    = CommandStatusType(0x00000027) // Failed to modify DL
	ESME_RDELDLFAIL    = CommandStatusType(0x00000028) // Failed to Delete DL
	ESME_RVIEWDLFAIL   = CommandStatusType(0x00000029) // Failed to View DL
	ESME_RLISTDLSFAIL  = CommandStatusType(0x00000030) // Failed to list DLs
	ESME_RPARAMRETFAIL = CommandStatusType(0x00000031) // Param Retrieve Failed
	ESME_RINVPARAM     = CommandStatusType(0x00000032) // Invalid Param

	ESME_RINVNUMDESTS = CommandStatusType(0x00000033)
	ESME_RINVDLNAME   = CommandStatusType(0x00000034)

	ESME_RINVDLMEMBDESC = CommandStatusType(0x00000035) // Invalid DL Member Description
	ESME_RINVDLMEMBTYP  = CommandStatusType(0x00000038) // Invalid DL Member Type
	ESME_RINVDLMODOPT   = CommandStatusType(0x00000039) // Invalid DL Modify Option

	ESME_RINVDESTFLAG = CommandStatusType(0x00000040)
	ESME_RINVSUBREP   = CommandStatusType(0x00000042)
	ESME_RINVESMCLASS = CommandStatusType(0x00000043)
	ESME_RCNTSUBDL    = CommandStatusType(0x00000044)
	ESME_RSUBMITFAIL  = CommandStatusType(0x00000045)
	ESME_RINVSRCTON   = CommandStatusType(0x00000048)
	ESME_RINVSRCNPI   = CommandStatusType(0x00000049)
	ESME_RINVDSTTON   = CommandStatusType(0x00000050)
	ESME_RINVDSTNPI   = CommandStatusType(0x00000051)
	ESME_RINVSYSTYP   = CommandStatusType(0x00000053)
	ESME_RINVREPFLAG  = CommandStatusType(0x00000054)
	ESME_RINVNUMMSGS  = CommandStatusType(0x00000055)
	ESME_RTHROTTLED   = CommandStatusType(0x00000058)

	ESME_RPROVNOTALLWD = CommandStatusType(0x00000059) // Provisioning Not Allowed

	ESME_RINVSCHED    = CommandStatusType(0x00000061)
	ESME_RINVEXPIRY   = CommandStatusType(0x00000062)
	ESME_RINVDFTMSGID = CommandStatusType(0x00000063)
	ESME_RX_T_APPN    = CommandStatusType(0x00000064)
	ESME_RX_P_APPN    = CommandStatusType(0x00000065)
	ESME_RX_R_APPN    = CommandStatusType(0x00000066)
	ESME_RQUERYFAIL   = CommandStatusType(0x00000067)

	ESME_RINVPGCUSTID      = CommandStatusType(0x00000080) // Paging Customer ID Invalid No such subscriber
	ESME_RINVPGCUSTIDLEN   = CommandStatusType(0x00000081) // Paging Customer ID length Invalid
	ESME_RINVCITYLEN       = CommandStatusType(0x00000082) // City Length Invalid
	ESME_RINVSTATELEN      = CommandStatusType(0x00000083) // State Length Invalid
	ESME_RINVZIPPREFIXLEN  = CommandStatusType(0x00000084) // Zip Prefix Length Invalid
	ESME_RINVZIPPOSTFIXLEN = CommandStatusType(0x00000085) // Zip Postfix Length Invalid
	ESME_RINVMINLEN        = CommandStatusType(0x00000086) // MIN Length Invalid
	ESME_RINVMIN           = CommandStatusType(0x00000087) // MIN Invalid (i.e. No such MIN)
	ESME_RINVPINLEN        = CommandStatusType(0x00000088) // PIN Length Invalid
	ESME_RINVTERMCODELEN   = CommandStatusType(0x00000089) // Terminal Code Length Invalid
	ESME_RINVCHANNELLEN    = CommandStatusType(0x0000008A) // Channel Length Invalid
	ESME_RINVCOVREGIONLEN  = CommandStatusType(0x0000008B) // Coverage Region Length Invalid
	ESME_RINVCAPCODELEN    = CommandStatusType(0x0000008C) // Cap Code Length Invalid
	ESME_RINVMDTLEN        = CommandStatusType(0x0000008D) // Message delivery time Length Invalid
	ESME_RINVPRIORMSGLEN   = CommandStatusType(0x0000008E) // Priority Message Length Invalid
	ESME_RINVPERMSGLEN     = CommandStatusType(0x0000008F) // Periodic Messages Length Invalid
	ESME_RINVPGALERTLEN    = CommandStatusType(0x00000090) // Paging Alerts Length Invalid
	ESME_RINVSMUSERLEN     = CommandStatusType(0x00000091) // int16 Message User Group Length Invalid
	ESME_RINVRTDBLEN       = CommandStatusType(0x00000092) // Real Time Data broadcasts Length Invalid
	ESME_RINVREGDELLEN     = CommandStatusType(0x00000093) // Registered Delivery Length Invalid
	ESME_RINVMSGDISTLEN    = CommandStatusType(0x00000094) // Message Distribution Length Invalid
	ESME_RINVPRIORMSG      = CommandStatusType(0x00000095) // Priority Message Length Invalid
	ESME_RINVMDT           = CommandStatusType(0x00000096) // Message delivery time Invalid
	ESME_RINVPERMSG        = CommandStatusType(0x00000097) // Periodic Messages Invalid
	ESME_RINVMSGDIST       = CommandStatusType(0x00000098) // Message Distribution Invalid
	ESME_RINVPGALERT       = CommandStatusType(0x00000099) // Paging Alerts Invalid
	ESME_RINVSMUSER        = CommandStatusType(0x0000009A) // int16 Message User Group Invalid
	ESME_RINVRTDB          = CommandStatusType(0x0000009B) // Real Time Data broadcasts Invalid
	ESME_RINVREGDEL        = CommandStatusType(0x0000009C) // Registered Delivery Invalid
	// public static final  ESME_RINVOPTPARSTREAM = CommandStatusType(0x0000009D)   // KIF IW Field out of data
	// public static final  ESME_ROPTPARNOTALLWD = CommandStatusType(0x0000009E)   // Optional Parameter not allowed
	ESME_RINVOPTPARLEN = CommandStatusType(0x0000009F) // Invalid Optional Parameter Length

	ESME_RINVOPTPARSTREAM = CommandStatusType(0x000000C0)
	ESME_ROPTPARNOTALLWD  = CommandStatusType(0x000000C1)
	ESME_RINVPARLEN       = CommandStatusType(0x000000C2)
	ESME_RMISSINGOPTPARAM = CommandStatusType(0x000000C3)
	ESME_RINVOPTPARAMVAL  = CommandStatusType(0x000000C4)
	ESME_RDELIVERYFAILURE = CommandStatusType(0x000000FE)
	ESME_RUNKNOWNERR      = CommandStatusType(0x000000FF)

	ESME_LAST_ERROR = CommandStatusType(0x0000012C) // THE VALUE OF THE LAST ERROR CODE
)
