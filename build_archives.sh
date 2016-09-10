$GOPATH/bin/gox github.com/lnsp/bench
for i in bench_*; do tar -czf $i.tar.gz $i; done
