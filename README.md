# PRICE INDEX

[Test task](Index_price_test_task.pdf) can be described like: Need to create the "fair" price index (avg. for all datasets, in given range) for limited amount of sources that send data with random delays without guaranties.

## Principle of application

Application based on SOLID principles to achieve easy way extension and re-use of components.

## Implementation for the simple scenario

At the first glance the task is simple. Exchanges are mostly sending ticker data not often than 1s. That leads us to mind that channels is ok, data synchronization based on mutex is more that enough for the limited amount of data sources. The example of index with mutex can be found in [simple application](cmd/simple/main.go).


## "Fairness" challenge

It's good to understand how good component in real-time. I Decide to add fairness property which represents that metric for each type of indexes. Meanwhile, the delays and errors the biggest obstacle to achieve higher fairness value. But we can assume that our system has multiple indexes which can be used to make more "fair" data by cross conversion. The new [index](internal/indexes/cross-rate.go) that implements that behaviour, used in resulted [cross-rate application](cmd/cross-rate/main.go).

## Further improvements (optional)

In opposite of the simple case we could imagine that data streams can provide data much master, f.ex. getting all trades or/and capture order book and get highly changeable ask/bid prices. In that case data transfers will be intensive, especially in high-frequency trading. Channels and tons of goroutines are not a good choice to achieve index accuracy, because of context switching, channels performance. More than that it will be nice to have optimization of the index itself. I have made optimized version of [mutex-based](internal/indexes/mutex-based.go) index [mb-optimized](internal/indexes/mb-optimized.go) version with improvements of allocations and mutexes. But in ideal scenario this is not enough since we are still using mutexes which can lock our writing/reading with context switches. That's why I added the [lock-free](internal/indexes/lock-free-optimized.go) version of the index which almost always ready for reading and writing.

*p.s. I know that it's looks like pre-mature optimization. It's not battle tested and most of the time the performance will be blocked by the network latency, etc. But it was interesting to dig in that for me. So you can freely skip that part ;)*

