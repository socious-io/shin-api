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
	StatusRequested VerificationStatusType = "REQUESTED"
	StatusVerfied   VerificationStatusType = "VERIFIED"
	StatusFailed    VerificationStatusType = "FAILED"
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
	StatusIssued   CredentialStatusType = "ISSUED"
	StatusClaimed  CredentialStatusType = "CLAIMED"
	StatusCanceled CredentialStatusType = "CANCELED"
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
