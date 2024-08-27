package main

import (
	"fmt"
	"net"
)

// Packet represents a network packet.
type Packet struct {
	SrcIP   net.IP
	DstIP   net.IP
	Payload string
}

// Handler is the interface that each handler in the chain will implement.
type Handler interface {
	SetNext(handler Handler)
	Handle(packet *Packet)
}

// BaseHandler provides a default implementation of the SetNext method.
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(packet *Packet) {
	if h.next != nil {
		h.next.Handle(packet)
	}
}

// IPFilterHandler filters packets based on their source IP address.
type IPFilterHandler struct {
	BaseHandler
	BlockedIPs []net.IP
}

func (h *IPFilterHandler) Handle(packet *Packet) {
	for _, blockedIP := range h.BlockedIPs {
		if packet.SrcIP.Equal(blockedIP) {
			fmt.Println("Packet blocked due to IP filter:", packet.SrcIP)
			return
		}
	}
	fmt.Println("IPFilterHandler passed packet.")
	h.BaseHandler.Handle(packet)
}

// PayloadInspectionHandler inspects the payload of a packet.
type PayloadInspectionHandler struct {
	BaseHandler
	ThreatSignatures []string
}

func (h *PayloadInspectionHandler) Handle(packet *Packet) {
	for _, signature := range h.ThreatSignatures {
		if signature == packet.Payload {
			fmt.Println("Threat detected in payload:", packet.Payload)
			return
		}
	}
	fmt.Println("PayloadInspectionHandler passed packet.")
	h.BaseHandler.Handle(packet)
}

// LoggerHandler logs packet details.
type LoggerHandler struct {
	BaseHandler
}

func (h *LoggerHandler) Handle(packet *Packet) {
	fmt.Printf("Logging packet: SrcIP=%s, DstIP=%s, Payload=%s\n", packet.SrcIP, packet.DstIP, packet.Payload)
	h.BaseHandler.Handle(packet)
}

func main() {
	// Create the handlers.
	ipFilter := &IPFilterHandler{
		BlockedIPs: []net.IP{
			net.IPv4(192, 168, 1, 2),
		},
	}
	payloadInspector := &PayloadInspectionHandler{
		ThreatSignatures: []string{"malicious_payload"},
	}
	logger := &LoggerHandler{}

	// Set up the chain of responsibility.
	ipFilter.SetNext(payloadInspector)
	payloadInspector.SetNext(logger)

	// Create some packets for testing.
	packets := []*Packet{
		{SrcIP: net.IPv4(192, 168, 1, 2), DstIP: net.IPv4(10, 0, 0, 1), Payload: "normal_payload"},
		{SrcIP: net.IPv4(192, 168, 1, 3), DstIP: net.IPv4(10, 0, 0, 1), Payload: "malicious_payload"},
		{SrcIP: net.IPv4(192, 168, 1, 4), DstIP: net.IPv4(10, 0, 0, 1), Payload: "normal_payload"},
	}

	// Process each packet through the chain.
	for _, packet := range packets {
		ipFilter.Handle(packet)
	}
}
