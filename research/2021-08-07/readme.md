# August 6 2021

## Mahalanobis distance

- <https://wikipedia.org/wiki/Mahalanobis_distance>
- https://github.com/bitterfly/emotions/blob/master/emotions/kmeans.go
- https://stats.stackexchange.com/questions/172564

## variance

Using the variance of the variables and assuming that queries are in the same
distributions would probably go a long way towards a reasonable answer.

Calculate the variance for each of the variable, then scale by this (the
variables and the query), then choose based on minimum Euclidean distance. This
is a reasonable, but naive implementation.

- https://pkg.go.dev/gonum.org/v1/gonum/stat#MeanVariance
- https://pkg.go.dev/gonum.org/v1/gonum/stat#Variance
- https://wikipedia.org/wiki/Variance
- https://github.com/golang/perf/blob/master/internal/stats/sample.go