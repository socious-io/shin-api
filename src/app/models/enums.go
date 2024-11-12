package models

import (
	"database/sql/driver"
	"fmt"
)

type AttributeType string

const (
	Text     AttributeType = "TEXT"
	Number   AttributeType = "NUMBER"
	Boolean  AttributeType = "BOOLEAN"
	Url      AttributeType = "URL"
	Datetime AttributeType = "DATETIME"
	Email    AttributeType = "EMAIL"
)

func (a *AttributeType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan attribute type: %v", value)
	}
	*a = AttributeType(strValue)
	return nil
}

func (a AttributeType) Value() (driver.Value, error) {
	return string(a), nil
}

type VerificationStatusType string

const (
	StatusVerifCreated   VerificationStatusType = "CREATED"
	StatusVerifRequested VerificationStatusType = "REQUESTED"
	StatusVerifVerfied   VerificationStatusType = "VERIFIED"
	StatusVerifFailed    VerificationStatusType = "FAILED"
)

func (c *VerificationStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = VerificationStatusType(string(v))
	case string:
		*c = VerificationStatusType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (c VerificationStatusType) Value() (driver.Value, error) {
	return string(c), nil
}

type CredentialStatusType string

const (
	StatusCreated  CredentialStatusType = "CREATED"
	StatusIssued   CredentialStatusType = "ISSUED"
	StatusClaimed  CredentialStatusType = "CLAIMED"
	StatusCanceled CredentialStatusType = "CANCELED"
	StatusRevoked  CredentialStatusType = "REVOKED"
)

func (c *CredentialStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = CredentialStatusType(string(v))
	case string:
		*c = CredentialStatusType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (c CredentialStatusType) Value() (driver.Value, error) {
	return string(c), nil
}

type VerificationOperatorType string

const (
	OperatorEqual   VerificationOperatorType = "EQUAL"
	OperatorNot     VerificationOperatorType = "NOT"
	OperatorBigger  VerificationOperatorType = "BIGGER"
	OperatorSmaller VerificationOperatorType = "SMALLER"
)

func (o *VerificationOperatorType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = VerificationOperatorType(string(v))
	case string:
		*o = VerificationOperatorType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o VerificationOperatorType) Value() (driver.Value, error) {
	return string(o), nil
}

type KybVerificationStatusType string

const (
	KYBStatusPending  KybVerificationStatusType = "PENDING"
	KYBStatusApproved KybVerificationStatusType = "APPROVED"
	KYBStatusRejected KybVerificationStatusType = "REJECTED"
)

func (o *KybVerificationStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = KybVerificationStatusType(string(v))
	case string:
		*o = KybVerificationStatusType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o KybVerificationStatusType) Value() (driver.Value, error) {
	return string(o), nil
}

type CSVImportDocType string

const (
	CSVDocTypeCredentials CSVImportDocType = "CREDENTIALS"
)

func (o *CSVImportDocType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = CSVImportDocType(string(v))
	case string:
		*o = CSVImportDocType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o CSVImportDocType) Value() (driver.Value, error) {
	return string(o), nil
}

type CSVImportStatus string

const (
	CSVStatusFailed CSVImportStatus = "FAILED"
	CSVStatusDone   CSVImportStatus = "DONE"
)

func (cis *CSVImportStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*cis = CSVImportStatus(string(v))
	case string:
		*cis = CSVImportStatus(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (cis CSVImportStatus) Value() (driver.Value, error) {
	return string(cis), nil
}

type VerificationType string

const (
	VerificationSingle VerificationType = "SINGLE"
	VerificationMulti  VerificationType = "MULTI"
)

func (o *VerificationType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = VerificationType(string(v))
	case string:
		*o = VerificationType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o VerificationType) Value() (driver.Value, error) {
	return string(o), nil
}

type OperationServiceTrigger string

const (
	OperationCredentialRevoke OperationServiceTrigger = "CREDENTIAL_REVOKE"
)

func (ost *OperationServiceTrigger) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ost = OperationServiceTrigger(string(v))
	case string:
		*ost = OperationServiceTrigger(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (ost OperationServiceTrigger) Value() (driver.Value, error) {
	return string(ost), nil
}
