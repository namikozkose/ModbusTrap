# ModbusTrap 🏭🛡️

A lightweight ICS/SCADA Honeypot and Active Defense system written in Go.

I developed ModbusTrap to explore how Operational Technology (OT) networks can be defended against targeted attacks. It acts as a decoy Programmable Logic Controller (PLC) or cooling pump listening on the standard Modbus TCP port (502).

Instead of just functioning as a passive logger, I built a Deep Packet Inspection (DPI) mechanism into it. It reads the raw hexadecimal payloads to figure out exactly what the attacker is trying to do. If it detects a sabotage attempt, it acts as an IPS (Intrusion Prevention System) and instantly drops the connection.

## ⚙️ How It Works (Architecture)

*   **Deep Packet Inspection (Modbus TCP):** The system analyzes raw byte payloads. Specifically, it inspects Byte 7 of the Modbus packet, which holds the Function Code.
*   **Dynamic Threat Classification:**
    *   If the payload contains Function Codes `0x01` to `0x04`, the system flags it as *SCADA Reconnaissance* (the attacker is polling/reading data).
    *   If it catches `0x05`, `0x06`, `0x0F`, or `0x10`, it triggers an *Industrial Sabotage* alert (the attacker is trying to turn a physical valve or coil on/off).
*   **Concurrency & Memory Safety:** I used Go's `goroutines` to handle multiple concurrent connection attempts (like Nmap scans or botnets). The IP blacklist is protected using `sync.Mutex` locks to prevent race conditions in memory.
*   **Zero-Trust Isolation:** Once an IP sends a malicious command, it gets added to the memory blacklist. Any future connection attempts from that IP are instantly dropped at the socket level before any data is processed.

## 🚀 Running the Project

Since it listens on port 502, you might need admin/root privileges depending on your OS firewall.

```bash
# Clone the repository
git clone [https://github.com/namikozkose/modbustrap.git](https://github.com/namikozkose/modbustrap.git)
cd modbustrap

# Build the executable
go build -o modbustrap main.go

# Run the honeypot
./modbustrap