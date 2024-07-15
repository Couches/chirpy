package ChirpyDatabase

type Result struct {
	Code int
  Body *any
	Error  error
}

func GetErrorResult(code int, err error) Result {
  errMsg := any(err.Error())
	return Result {
		Code: code,
    Body: &errMsg,
		Error:  err,
	}
}

func GetOKResult(code int, body any) Result {
  return Result {
    Code: code,
    Body: &body,
    Error: nil,
  }
}
