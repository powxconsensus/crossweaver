package types

import "fmt"

type CustomError struct {
	ErrorID int64
	Err     error
	Retry   uint8
}

func NewCustomError(errorID int64, Err error, retry uint8) *CustomError {
	return &CustomError{
		ErrorID: errorID,
		Err:     Err,
		Retry:   retry,
	}
}
func (e *CustomError) Error() error {
	return fmt.Errorf("id %d: err %v", e.ErrorID, e.Err)
}

func RelayerOutOfFundsError() *CustomError {
	return NewCustomError(700, fmt.Errorf("relayer is out of funds"), 3)
}

func InternalError() *CustomError {
	return NewCustomError(701, fmt.Errorf("internal error occured"), 3)
}
func RPCFailureError() *CustomError {
	return NewCustomError(702, fmt.Errorf("rpc failure error occured"), 3)
}
func ContractError() *CustomError {
	return NewCustomError(703, fmt.Errorf("unexpected contract error occured"), 3)
}
func DbError() *CustomError {
	return NewCustomError(704, fmt.Errorf("Error while interacting with DB"), 3)
}

func OutOfGasError() *CustomError {
	return NewCustomError(705, fmt.Errorf("GasLimit exhausted"), 0)
}
