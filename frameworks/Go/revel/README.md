# [Revel](https://github.com/revel/revel) (Go) Benchmarking Test

This is the go [Revel](https://github.com/TechEmpower/FrameworkBenchmarks/tree/master/frameworks/Go/revel) portion of a benchmarking test suite comparing a variety of web development platforms.

"Revel - A high-productivity web framework for the Go language"

## Source

* [All test source](https://github.com/TechEmpower/FrameworkBenchmarks/tree/master/frameworks/Go/revel/src/benchmark/app/controllers/app.go)
* Revel Framework version: v0.13.1

## Test URLs
```
Test 1: JSON              => http://localhost:8080/json
Test 2: Single DB Query   => http://localhost:8080/db
Test 3: Multiple DB Query => http://localhost:8080/dbs?queries=[1-500]
Test 4: Fortunes          => http://localhost:8080/fortune
Test 5: Database Updates  => http://localhost:8080/update?queries=[1-500]
Test 6: Plaintext         => http://localhost:8080/plaintext
```
