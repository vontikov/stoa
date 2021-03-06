// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: stoa.proto

package pb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on Entity with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Entity) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetEntityName()) < 1 {
		return EntityValidationError{
			field:  "EntityName",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// EntityValidationError is the validation error returned by Entity.Validate if
// the designated constraints aren't met.
type EntityValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EntityValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EntityValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EntityValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EntityValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EntityValidationError) ErrorName() string { return "EntityValidationError" }

// Error satisfies the builtin error interface
func (e EntityValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEntity.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EntityValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EntityValidationError{}

// Validate checks the field values on ClientId with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *ClientId) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for EntityName

	if len(m.GetId()) < 1 {
		return ClientIdValidationError{
			field:  "Id",
			reason: "value length must be at least 1 bytes",
		}
	}

	// no validation rules for Payload

	return nil
}

// ClientIdValidationError is the validation error returned by
// ClientId.Validate if the designated constraints aren't met.
type ClientIdValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClientIdValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClientIdValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClientIdValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClientIdValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClientIdValidationError) ErrorName() string { return "ClientIdValidationError" }

// Error satisfies the builtin error interface
func (e ClientIdValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClientId.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClientIdValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClientIdValidationError{}

// Validate checks the field values on Value with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Value) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetEntityName()) < 1 {
		return ValueValidationError{
			field:  "EntityName",
			reason: "value length must be at least 1 runes",
		}
	}

	if len(m.GetValue()) < 1 {
		return ValueValidationError{
			field:  "Value",
			reason: "value length must be at least 1 bytes",
		}
	}

	return nil
}

// ValueValidationError is the validation error returned by Value.Validate if
// the designated constraints aren't met.
type ValueValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValueValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValueValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValueValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValueValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValueValidationError) ErrorName() string { return "ValueValidationError" }

// Error satisfies the builtin error interface
func (e ValueValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValue.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValueValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValueValidationError{}

// Validate checks the field values on Key with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Key) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetEntityName()) < 1 {
		return KeyValidationError{
			field:  "EntityName",
			reason: "value length must be at least 1 runes",
		}
	}

	if len(m.GetKey()) < 1 {
		return KeyValidationError{
			field:  "Key",
			reason: "value length must be at least 1 bytes",
		}
	}

	return nil
}

// KeyValidationError is the validation error returned by Key.Validate if the
// designated constraints aren't met.
type KeyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e KeyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e KeyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e KeyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e KeyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e KeyValidationError) ErrorName() string { return "KeyValidationError" }

// Error satisfies the builtin error interface
func (e KeyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sKey.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = KeyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = KeyValidationError{}

// Validate checks the field values on KeyValue with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *KeyValue) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetEntityName()) < 1 {
		return KeyValueValidationError{
			field:  "EntityName",
			reason: "value length must be at least 1 runes",
		}
	}

	if len(m.GetKey()) < 1 {
		return KeyValueValidationError{
			field:  "Key",
			reason: "value length must be at least 1 bytes",
		}
	}

	if len(m.GetValue()) < 1 {
		return KeyValueValidationError{
			field:  "Value",
			reason: "value length must be at least 1 bytes",
		}
	}

	return nil
}

// KeyValueValidationError is the validation error returned by
// KeyValue.Validate if the designated constraints aren't met.
type KeyValueValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e KeyValueValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e KeyValueValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e KeyValueValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e KeyValueValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e KeyValueValidationError) ErrorName() string { return "KeyValueValidationError" }

// Error satisfies the builtin error interface
func (e KeyValueValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sKeyValue.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = KeyValueValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = KeyValueValidationError{}
