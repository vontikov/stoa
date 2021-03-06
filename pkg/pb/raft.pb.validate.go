// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: raft.proto

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

// Validate checks the field values on Discovery with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Discovery) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for BindIp

	// no validation rules for BindPort

	// no validation rules for GrpcIp

	// no validation rules for GrpcPort

	// no validation rules for Leader

	return nil
}

// DiscoveryValidationError is the validation error returned by
// Discovery.Validate if the designated constraints aren't met.
type DiscoveryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DiscoveryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DiscoveryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DiscoveryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DiscoveryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DiscoveryValidationError) ErrorName() string { return "DiscoveryValidationError" }

// Error satisfies the builtin error interface
func (e DiscoveryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDiscovery.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DiscoveryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DiscoveryValidationError{}

// Validate checks the field values on ClusterCommand with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ClusterCommand) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Command

	// no validation rules for TtlEnabled

	// no validation rules for TtlMillis

	// no validation rules for BarrierEnabled

	// no validation rules for BarrierMillis

	switch m.Payload.(type) {

	case *ClusterCommand_Entity:

		if v, ok := interface{}(m.GetEntity()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClusterCommandValidationError{
					field:  "Entity",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ClusterCommand_Value:

		if v, ok := interface{}(m.GetValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClusterCommandValidationError{
					field:  "Value",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ClusterCommand_Key:

		if v, ok := interface{}(m.GetKey()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClusterCommandValidationError{
					field:  "Key",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ClusterCommand_KeyValue:

		if v, ok := interface{}(m.GetKeyValue()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClusterCommandValidationError{
					field:  "KeyValue",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *ClusterCommand_ClientId:

		if v, ok := interface{}(m.GetClientId()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ClusterCommandValidationError{
					field:  "ClientId",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ClusterCommandValidationError is the validation error returned by
// ClusterCommand.Validate if the designated constraints aren't met.
type ClusterCommandValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterCommandValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterCommandValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterCommandValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterCommandValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterCommandValidationError) ErrorName() string { return "ClusterCommandValidationError" }

// Error satisfies the builtin error interface
func (e ClusterCommandValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterCommand.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterCommandValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterCommandValidationError{}
