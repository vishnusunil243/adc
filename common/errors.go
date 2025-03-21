package common

// Define Error Codes
const (
	ErrCodeNotFound               ErrorCode = "NOT_FOUND"
	ErrCodeDuplicateEntry         ErrorCode = "DUPLICATE_ENTRY"
	ErrCodeForeignKey             ErrorCode = "FOREIGN_KEY_VIOLATION"
	ErrCodeInvalidInput           ErrorCode = "INVALID_INPUT"
	ErrCodeInternal               ErrorCode = "INTERNAL_ERROR"
	ErrCodeConnectionFailed       ErrorCode = "CONNECTION_FAILED"
	ErrCodeTimeout                ErrorCode = "TIMEOUT"
	ErrCodeTransactionFailed      ErrorCode = "TRANSACTION_FAILED"
	ErrCodeSyntaxError            ErrorCode = "SYNTAX_ERROR"
	ErrCodePermissionDenied       ErrorCode = "PERMISSION_DENIED"
	ErrCodeTooManyConnections     ErrorCode = "TOO_MANY_CONNECTIONS"
	ErrCodeSerializationFailure   ErrorCode = "SERIALIZATION_FAILURE"
	ErrCodeInvalidColumnReference ErrorCode = "INVALID_COLUMN_REFERENCE"
	ErrCodeInvalidData            ErrorCode = "INVALID_DATA"

	ErrCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
)

// RepoErr represents a structured error
type ErrorCode string
