package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"main.go/common"
)

type RepoErr struct {
	Message   string           `json:"message"`
	ErrorCode common.ErrorCode `json:"error_code"`
	Metadata  string           `json:"metadata"`
}

func (r *RepoErr) Error() string {
	return fmt.Sprintf("%s:%s", r.ErrorCode, r.Message)
}

func HandleDBError(err error) *RepoErr {
	if err == nil {
		return nil
	}

	// Handle GORM ErrRecordNotFound
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &RepoErr{
			Message:   "record not found",
			ErrorCode: common.ErrCodeNotFound,
		}
	}

	// Handle GORM Validation Errors
	if errors.Is(err, gorm.ErrInvalidData) {
		return &RepoErr{
			Message:   "invalid data",
			ErrorCode: common.ErrCodeInvalidData,
		}
	}

	// Handle SQL No Rows Error
	if errors.Is(err, sql.ErrNoRows) {
		return &RepoErr{
			Message:   "record not found",
			ErrorCode: common.ErrCodeNotFound,
		}
	}

	// Handle GORM Errors related to Database Constraints
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &RepoErr{
			Message:   "record not found",
			ErrorCode: common.ErrCodeNotFound,
		}
	}

	// Handle PostgreSQL Specific Errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // Unique constraint violation
			return &RepoErr{
				Message:   "duplicate entry",
				ErrorCode: common.ErrCodeDuplicateEntry,
				Metadata:  pqErr.Detail,
			}
		case "23503": // Foreign key violation
			return &RepoErr{
				Message:   "foreign key violation",
				ErrorCode: common.ErrCodeForeignKey,
				Metadata:  pqErr.Detail,
			}
		case "22P02": // Invalid input syntax
			return &RepoErr{
				Message:   "invalid input",
				ErrorCode: common.ErrCodeInvalidInput,
				Metadata:  pqErr.Detail,
			}
		case "08001", "08006": // Connection failure
			return &RepoErr{
				Message:   "database connection failed",
				ErrorCode: common.ErrCodeConnectionFailed,
				Metadata:  pqErr.Detail,
			}
		case "57014": // Query timeout
			return &RepoErr{
				Message:   "query execution timed out",
				ErrorCode: common.ErrCodeTimeout,
				Metadata:  pqErr.Detail,
			}
		case "40001": // Transaction failed (deadlock)
			return &RepoErr{
				Message:   "transaction failure due to deadlock",
				ErrorCode: common.ErrCodeTransactionFailed,
				Metadata:  pqErr.Detail,
			}
		case "42601": // Syntax error in SQL query
			return &RepoErr{
				Message:   "syntax error in query",
				ErrorCode: common.ErrCodeSyntaxError,
				Metadata:  pqErr.Detail,
			}
		case "42501": // Permission denied
			return &RepoErr{
				Message:   "insufficient privileges",
				ErrorCode: common.ErrCodePermissionDenied,
				Metadata:  pqErr.Detail,
			}
		case "53300": // Too many connections
			return &RepoErr{
				Message:   "too many connections to database",
				ErrorCode: common.ErrCodeTooManyConnections,
				Metadata:  pqErr.Detail,
			}
		case "40003": // Serialization failure (deadlocks)
			return &RepoErr{
				Message:   "serialization failure occurred",
				ErrorCode: common.ErrCodeSerializationFailure,
				Metadata:  pqErr.Detail,
			}
		case "42703": // Invalid column reference
			return &RepoErr{
				Message:   "invalid column reference",
				ErrorCode: common.ErrCodeInvalidColumnReference,
				Metadata:  pqErr.Detail,
			}
		default:
			return &RepoErr{
				Message:   "internal database error",
				ErrorCode: common.ErrCodeInternal,
				Metadata:  pqErr.Detail,
			}
		}
	}

	// Default case for unknown errors
	return &RepoErr{
		Message:   err.Error(),
		ErrorCode: common.ErrCodeInternal,
	}
}

func NewRepoErr(message string, errorcode common.ErrorCode) *RepoErr {
	return &RepoErr{
		Message:   message,
		ErrorCode: errorcode,
	}
}
