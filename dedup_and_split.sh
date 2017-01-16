#sort -n -t, -k4 < ecn-conspiracy-resolved-20170113.txt | ipdedup > ecn-conspiracy-dedup-20170113.txt
cut -d, -f1-3 ecn-conspiracy-dedup-20170113.txt | grep ,25, > public-targets-20170113-mail.csv
cut -d, -f1-3 ecn-conspiracy-dedup-20170113.txt | grep ,53, > public-targets-20170113-dns.csv
cut -d, -f1-3 ecn-conspiracy-dedup-20170113.txt | grep ,80, > public-targets-20170113-http.csv
