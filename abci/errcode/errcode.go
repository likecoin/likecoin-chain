package errcode

// Error Code Definition (5 digits)
// 1 0 0 0 0
// | | | | |
// | | | | Type of transaction
// | | Case
// | | - Common case 10-59
// | | - Other case  60-99
// Type
// - Before parsing request 00-09
// - Transaction            10-59
// - Query                  60-99

func RegisterCheckTxInvalidFormat() (uint32, string) {
	return 10000, "Invalid RegisterTransaction format in CheckTx"
}

func RegisterDeliverTxInvalidFormat() (uint32, string) {
	return 10001, "Invalid RegisterTransaction format in DeliverTx"
}

func RegisterCheckTxInvalidSignature() (uint32, string) {
	return 10010, "Invalid RegisterTransaction signature in CheckTx"
}

func RegisterDeliverTxInvalidSignature() (uint32, string) {
	return 10011, "Invalid RegisterTransaction signature in DeliverTx"
}

func RegisterCheckTxDuplicated() (uint32, string) {
	return 10020, "Duplicated RegisterTransaction in CheckTx"
}

func RegisterDeliverTxDuplicated() (uint32, string) {
	return 10021, "Duplicated RegisterTransaction in DeliverTx"
}

func DepositCheckTxInvalidFormat() (uint32, string) {
	return 11000, "Invalid DepositTransaction format in CheckTx"
}

func DepositDeliverTxInvalidFormat() (uint32, string) {
	return 11001, "Invalid DepositTransaction format in DeliverTx"
}

func DepositCheckTxDuplicated() (uint32, string) {
	return 11020, "Duplicated DepositTransaction in CheckTx"
}

func DepositDeliverTxDuplicated() (uint32, string) {
	return 11021, "Duplicated DepositTransaction in DeliverTx"
}

func TransferCheckTxInvalidFormat() (uint32, string) {
	return 12000, "Invalid TransferTransaction format in CheckTx"
}

func TransferDeliverTxInvalidFormat() (uint32, string) {
	return 12001, "Invalid TransferTransaction format in DeliverTx"
}

func TransferCheckTxInvalidSignature() (uint32, string) {
	return 12010, "Invalid TransferTransaction signature in CheckTx"
}

func TransferDeliverTxInvalidSignature() (uint32, string) {
	return 12011, "Invalid TransferTransaction signature in DeliverTx"
}

func TransferCheckTxDuplicated() (uint32, string) {
	return 12020, "Duplicated TransferTransaction in CheckTx"
}

func TransferDeliverTxDuplicated() (uint32, string) {
	return 12021, "Duplicated TransferTransaction in DeliverTx"
}

func WithdrawCheckTxInvalidFormat() (uint32, string) {
	return 13000, "Invalid WithdrawTransaction format in CheckTx"
}

func WithdrawDeliverTxInvalidFormat() (uint32, string) {
	return 13001, "Invalid WithdrawTransaction format in DeliverTx"
}

func WithdrawCheckTxInvalidSignature() (uint32, string) {
	return 13010, "Invalid WithdrawTransaction signature in CheckTx"
}

func WithdrawDeliverTxInvalidSignature() (uint32, string) {
	return 13011, "Invalid WithdrawTransaction signature in DeliverTx"
}

func WithdrawCheckTxDuplicated() (uint32, string) {
	return 13020, "Duplicated WithdrawTransaction in CheckTx"
}

func WithdrawDeliverTxDuplicated() (uint32, string) {
	return 13021, "Duplicated WithdrawTransaction in DeliverTx"
}
