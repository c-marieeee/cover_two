# cover_two

Defensive security tool for IP address blocking.

## File Structure

### Main Files

- **server.go**: Go server that stores the block list.
- **cover_two.go**: Go program that calls the pf block list.

### Configuration Files

- **conf/**
  - `pf.conf.example`: Example `pf.conf` configuration file.
  - `pf.blocked.example`: Example `pf.blocked` file with sample IPs.

### Scripts

- **scripts/**
  - `update_cover_two.sh`: Script to fetch and update the block list.
