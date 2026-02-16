# BGPEFL – BGP Easy for Labs

Application written in **Go** to simplify the creation of **BGP** sessions and the injection of real ASN prefixes (via IRR) into lab environments such as **EVE-NG** and **PnetLab**.

Ideal for simulating realistic routing scenarios using real public prefixes without manually configuring hundreds of routes.

---

## 🚀 Demonstration

<p align="center">
  <img src="example.gif" alt="BGPEFL Demo" width="700">
</p>

## 📦 About the Project

**BGPEFL (BGP Easy for Labs)** allows you to:

* ✅ Quickly establish a BGP session
* ✅ Fetch real prefixes from an ASN via IRR (e.g., `whois.radb.net`)
* ✅ Automatically inject routes into GoBGP
* ✅ Easily clean the RIB
* ✅ Control the application lifecycle

It uses **GoBGP (`gobgpd`)** as the BGP daemon.

---

## 🔧 Requirements

* Linux
* Golang installed (1.23+)
* GoBGP installed (`gobgpd`)
* Root privileges (required for IP/interface manipulation)
* Connectivity to an IRR server (default: `whois.radb.net`)
* Lab environment such as:

  * EVE-NG
  * PnetLab

---

## 🚀 Installation

Start a Linux node (Debian/Ubuntu, for example) in EVE-NG or PnetLab with at least **two interfaces**:

* One connected to the Internet
* One dedicated to the BGP session

You can either compile from source, use a prebuilt binary, or deploy the BGPEFL Appliance.

---

## Build from Source

```bash
git clone https://github.com/jeffersonraimon/bgpefl.git
cd bgpefl
go build -o bgpefl
cp bgpefl /usr/bin/bgpefl
```

---

## Prebuilt Binary (GoBGP required)

```bash
wget https://github.com/jeffersonraimon/bgpefl/releases/download/v1.0.2/bgpefl
apt install gobgpd
chmod +x bgpefl
mv bgpefl /usr/bin/bgpefl
```

---

## BGPEFL Appliance

Lightweight QEMU image based on Alpine 3.20.3 with BGPEFL preinstalled and ready to use.

**Login:**
`root` / no password

### Installation Steps:

On EVE-NG/PNETLAB Host:
```bash
wget https://github.com/jeffersonraimon/bgpefl/releases/download/v1.0.2/bgpefl-appliance_v1.0.2.zip
unzip bgpefl-appliance_v1.0.2.zip
```

* Move the folder to:
  `/opt/unetlab/addons/qemu/`

* Place the `.yml` file in:
  `/opt/unetlab/html/templates/amd`
  or
  `/opt/unetlab/html/templates/intel`

* Fix permissions:

  ```bash
  /opt/unetlab/wrappers/unl_wrapper -a fixpermissions
  ```

---

## 📌 General Usage

```bash
bgpefl [command]
```

Available commands:

| Command  | Description                             |
| -------- | --------------------------------------- |
| init     | Initializes a BGP session for lab usage |
| gen      | Generates BGP routes based on IRR data  |
| clearrib | Clears the BGPEFL RIB                   |
| status   | Shows BGPEFL status                     |
| stop     | Stops BGPEFL                            |
| help     | Displays help information               |

---

## 🔹 init – Initialize BGP Session

Creates the BGP session, configures the interface IP address, and starts `gobgpd`.

### Usage:

```bash
bgpefl init --ip <IP> --cidr <CIDR> --int <INTERFACE> \
            --local-as <LOCAL_AS> \
            --neighbor <NEIGHBOR_IP> \
            --remote-as <REMOTE_AS> \
            [--router-id <ROUTER_ID>]
```

### Required Flags:

| Flag        | Description          |
| ----------- | -------------------- |
| --ip        | Interface IP address |
| --cidr      | CIDR subnet mask     |
| --int       | Network interface    |
| --local-as  | Local ASN            |
| --neighbor  | Neighbor IP address  |
| --remote-as | Remote ASN           |

### Example:

```bash
bgpefl init \
  --ip 192.168.0.2 \
  --cidr 30 \
  --int eth1 \
  --local-as 65001 \
  --neighbor 192.168.0.1 \
  --remote-as 65000
```

---

## 🔹 gen – Generate Prefixes via IRR

Fetches prefixes from an ASN via an IRR server and injects them into GoBGP.

### Usage:

```bash
bgpefl gen --as <ASN> [flags]
```

### Flags:

| Flag       | Description                              |
| ---------- | ---------------------------------------- |
| --as       | ASN to fetch prefixes from               |
| --irr      | IRR server (default: whois.radb.net)     |
| --dry-run  | Simulation mode (does not inject routes) |
| --limit    | Total prefix limit                       |
| --limit-v4 | IPv4 prefix limit                        |
| --limit-v6 | IPv6 prefix limit                        |
| --min-v4   | Minimum IPv4 prefix length               |
| --min-v6   | Minimum IPv6 prefix length               |
| --only-v4  | IPv4 only                                |
| --only-v6  | IPv6 only                                |

### Example:

```bash
bgpefl gen --as 15169 --limit 100
```

### Simulation Mode:

```bash
bgpefl gen --as 13335 --dry-run
```
---

## 🔹 clearrib – Clear Routes

Removes routes currently injected into GoBGP.

### Usage:

```bash
bgpefl clearrib [flags]
```

### Flags:

| Flag    | Description                         |
| ------- | ----------------------------------- |
| --soft  | Removes routes one by one (default) |
| --force | Uses direct `del all`               |
| --ipv4  | Removes IPv4 routes only            |
| --ipv6  | Removes IPv6 routes only            |

⚠️ If `gobgpd` is not running, the command will return an error.

---

## 🔹 status – Check Status

Displays the current BGPEFL status and GoBGP information.

```bash
bgpefl status
```

### Example output:

```bash
========== BGPEFL STATUS ==========
gobgpd: STOPPED
====================================
```

---

## 🔹 stop – Stop BGPEFL

Stops `gobgpd` and can optionally clear routes and remove the IP address from the interface.

### Usage:

```bash
bgpefl stop [flags]
```

### Flags:

| Flag          | Description                           |
| ------------- | ------------------------------------- |
| --clear-rib   | Removes all routes                    |
| --force       | Forcefully kills the process          |
| --remove-ip   | IP address to remove                  |
| --remove-cidr | CIDR of the IP                        |
| --remove-int  | Interface from which to remove the IP |

### Example:

```bash
bgpefl stop --clear-rib --force
```

### Removing IP from the interface:

```bash
bgpefl stop \
  --remove-ip 192.168.0.2 \
  --remove-cidr 30 \
  --remove-int eth1
```

---

## 🧠 Recommended Workflow

Initialize BGP session:

```bash
bgpefl init ...
```

Generate prefixes:

```bash
bgpefl gen --as <ASN> --limit 20
```

Check status:

```bash
bgpefl status
```

Clear routes (if needed):

```bash
bgpefl clearrib
```

Stop the environment:

```bash
bgpefl stop
```

---

## 🧪 Use Cases

* BGP labs in EVE-NG / PnetLab
* BGP policy testing
* Studying filters and route-maps
* RPKI and origin validation testing
* Certification training (e.g., CCNP / CCIE / JNCIP)

---

## ⚠️ Disclaimer

This project is intended **exclusively for lab environments**.

Do NOT use it for:

* Announcing real prefixes to the public Internet
* Production environment testing
* Uncontrolled routing scenarios

---

## 📄 License

MIT License

