#!/bin/bash
sudo pfctl -e
sudo pfctl -f /etc/pf.conf
sudo pfctl -t blocked -T replace -f /etc/pf.blocked

