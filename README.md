# cover_two
defensive security tool for IP address blocking. 

cover_two/
├── server.go               # go server that stores the block list
├── cover_two.go            # go program that calls the pf block list
├── conf/
│   ├── pf.conf.example     # example pf.conf configuration file
│   └── pf.blocked.example  # example pf.blocked file with sample IPs
├── scripts/
│   └── update_cover_two.sh # script to fetch and update the block list
