package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

const (
	packetSize = 1024 // Bytes
	testDuration = 10 * time.Second
)

func benchmarkRawEBPF() (int, error) {
	// Load an eBPF program (e.g., XDP) to handle packets
	prog, err := ebpf.LoadProgram(&ebpf.ProgramSpec{
		Type: ebpf.XDP,
		License: "GPL",
		Instructions: ebpf.Instructions{
			// Insert minimal eBPF XDP code here (e.g., pass or drop packets).
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to load eBPF program: %w", err)
	}
	defer prog.Close()

	// Attach the program to a network interface (example: eth0)
	linkName := "eth0"
	link, err := ebpf.AttachLink(&ebpf.AttachLinkOptions{
		Target: linkName,
		Program: prog,
		AttachType: ebpf.AttachTypeXDP,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to attach eBPF program: %w", err)
	}
	defer link.Close()

	// Simulate traffic and measure throughput
	return simulateTraffic(linkName)
}

func benchmarkCilium() (int, error) {
	// Assume Cilium is already configured and running.
	// Simulate traffic on a network managed by Cilium.
	return simulateTraffic("eth0")
}

func simulateTraffic(interfaceName string) (int, error) {
	conn, err := net.ListenPacket("udp", "0.0.0.0:0")
	if err != nil {
		return 0, fmt.Errorf("failed to create UDP listener: %w", err)
	}
	defer conn.Close()

	destination := &net.UDPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}

	packet := make([]byte, packetSize)
	sentPackets := 0
	var mu sync.Mutex

	done := make(chan bool)
	go func() {
		time.Sleep(testDuration)
		done <- true
	}()

	for {
		select {
		case <-done:
			return sentPackets, nil
		default:
			n, err := conn.WriteTo(packet, destination)
			if err != nil {
				return 0, fmt.Errorf("error sending packet: %w", err)
			}
			if n == packetSize {
				mu.Lock()
				sentPackets++
				mu.Unlock()
			}
		}
	}
}

func main() {
	fmt.Println("Starting benchmark tests...")

	fmt.Println("Benchmarking raw eBPF...")
	ebpfPackets, err := benchmarkRawEBPF()
	if err != nil {
		log.Fatalf("Error benchmarking raw eBPF: %v", err)
	}

	fmt.Println("Benchmarking Cilium...")
	ciliumPackets, err := benchmarkCilium()
	if err != nil {
		log.Fatalf("Error benchmarking Cilium: %v", err)
	}

	fmt.Printf("Raw eBPF throughput: %d packets/s\n", ebpfPackets/int(testDuration.Seconds()))
	fmt.Printf("Cilium throughput: %d packets/s\n", ciliumPackets/int(testDuration.Seconds()))
}
