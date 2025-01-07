package errno

var (
	// Common errors

	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	BadRequest          = &Errno{Code: 10002, Message: "Bad request"}

	// User errors
)
