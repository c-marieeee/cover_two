# Cover_Two

Defensive security tool for IP address blocking.

## File Structure

### Source

- **server.go**: Go server that stores the block list.
- **cover_two.go**: Go program that calls the pf block list.

### Configuration

- **conf/**
  - `pf.conf.example`: Example `pf.conf` configuration file.
  - `pf.blocked.example`: Example `pf.blocked` file with sample IPs.
