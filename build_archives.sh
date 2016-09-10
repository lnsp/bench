$GOPATH/bin/gox github.com/lnsp/bench
for i in *.tar.gz; do tar -czf $i.tar.gz $i; done
