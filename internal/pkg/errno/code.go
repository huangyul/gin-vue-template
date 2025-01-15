package errno

var (
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	BadRequest          = &Errno{Code: 10002, Message: "Bad request"}
	UsernameConflict    = &Errno{Code: 20001, Message: "username is already taken"}
	UserNotFound        = &Errno{Code: 20001, Message: "user not found"}
)
