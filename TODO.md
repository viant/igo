- [ ] add goroutine (`go` statement) support
- [ ] add channel send/receive support
- [ ] add `select` support
- [ ] add closure support

- [x] stabilize complex struct state-assignment paths seen in tests
- [x] add more filter test variants (case-insensitive, multi-condition, empty/nil handling)
- [x] add benchmarks for row filtering workloads (`contains` + append/filter path)

- [ ] optimize row-filter benchmark hot path (current Igo benchmark still allocates ~26 allocs/op)
