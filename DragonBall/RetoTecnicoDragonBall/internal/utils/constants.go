package utils

// Error Code
const (
	EC_CRITICAL_ERROR = "-1"
	EC_SIN_ERROR      = "0"
)

// Rest Service
const C_ContentText = "text"

// Formatos
const (
	FORMAT_YYYYMMDD_DASH   = "2006-01-02"
	FORMAT_DDMMYYYY_STRING = "02 %s 2006"
)

// Middleware
const (
	C_Request                   = "objRequest"
	C_Response                  = "objResponse"
	C_ContentType               = "Content-Type"
	C_ContentTypeValue          = "text/plain"
	C_X_Forwarded_Authorization = "X-Forwarded-Authorization"
	C_SubID                     = "SubID"
	C_TraceID                   = "trace-id"
	C_TokenApplication          = "tokenApplication"
	C_RequestLog                = "RequestIDLog"
	C_Counter                   = "counter"
)

// Security encrypt
const (
	URL_ENCRYPT = "/api/credit-loan/security/encrypt"
	URL_DECRYPT = "/api/credit-loan/security/decrypt"
)

// HTTP
const (
	CONTENT_TYPE     = "Content-Type"
	APPLICATION_JSON = "application/json"
	TEXT_PLAIN       = "text/plain"
)

const BEARER = "Bearer "

// Env
const (
	ENV_DIGITAL_BANKING_CONFIG_VERIFY_CERTIFICATE = "DIGITAL_BANKING_CONFIG_VERIFY_CERTIFICATE"
)

// TimeZone
const (
	TIME_ZONE_UTC_PERU = -5
)

// Frauds
const (
	CURRENCY_ID_PEN     = 1
	CURRENCY_ID_USD     = 2
	CURRENCY_SYMBOL_USD = "$"
	CURRENCY_SYMBOL_PEN = "S/"
)

const MOCK_USER = "Rai Delgado"
