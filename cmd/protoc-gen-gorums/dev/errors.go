package dev

import (
	"bytes"
	"fmt"
)

// A NodeNotFoundError reports that a specified node could not be found.
type NodeNotFoundError uint32

func (e NodeNotFoundError) Error() string {
	return fmt.Sprintf("node not found: %d", e)
}

// A ConfigNotFoundError reports that a specified configuration could not be
// found.
type ConfigNotFoundError uint32

func (e ConfigNotFoundError) Error() string {
	return fmt.Sprintf("configuration not found: %d", e)
}

// An IllegalConfigError reports that a specified configuration could not be
// created.
type IllegalConfigError string

func (e IllegalConfigError) Error() string {
	return "illegal configuration: " + string(e)
}

// ManagerCreationError returns an error reporting that a Manager could not be
// created due to err.
func ManagerCreationError(err error) error {
	return fmt.Errorf("could not create manager: %s", err.Error())
}

// A QuorumCallError is used to report that a quorum call failed.
type QuorumCallError struct {
	Reason     string
	ReplyCount int
	Errors     []GRPCError
}

func (e QuorumCallError) Error() string {
	var b bytes.Buffer
	b.WriteString("quorum call error: ")
	b.WriteString(e.Reason)
	b.WriteString(fmt.Sprintf(" (errors: %d, replies: %d)", len(e.Errors), e.ReplyCount))
	if len(e.Errors) == 0 {
		return b.String()
	}
	b.WriteString("\ngrpc errors:\n")
	for _, err := range e.Errors {
		b.WriteByte('\t')
		b.WriteString(fmt.Sprintf("node %d: %v", err.NodeID, err.Cause))
		b.WriteByte('\n')
	}
	return b.String()
}

// GRPCError is used to report that a single gRPC call failed.
type GRPCError struct {
	NodeID uint32
	Cause  error
}

func (e GRPCError) Error() string {
	return fmt.Sprintf("node %d: %v", e.NodeID, e.Cause.Error())
}
