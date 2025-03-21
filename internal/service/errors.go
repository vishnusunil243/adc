package service

import (
	"main.go/common"
	"main.go/internal/repository"
)

type ServiceError struct {
	Message    string           `json:"message"`
	StatusCode int              `json:"status_code"`
	ErrorCode  common.ErrorCode `json:"error_code"`
	MetaData   string           `json:"metadata"`
}

func (s ServiceError) Error() string { return s.Message }

func HandleRepoErr(repoErr *repository.RepoErr, msg string) *ServiceError {
	return &ServiceError{
		Message:   msg,
		ErrorCode: repoErr.ErrorCode,
		MetaData:  repoErr.Error(),
	}
}

func NewServiceError(msg string, errorCode common.ErrorCode, statuscode int) *ServiceError {
	return &ServiceError{
		Message:    msg,
		ErrorCode:  errorCode,
		StatusCode: statuscode,
	}
}
