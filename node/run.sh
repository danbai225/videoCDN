sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
nohup ./node >/dev/null 2>&1 &
