# Invertible Bloom Lookup Table

## Overview

Invertible Bloom Lookup Table (IBLT), is a probabilistic data structure for set-reconciliation problem. The general idea of IBLT was derived from Bloom filter and counting Bloom filter. By replacing the bits in standard Bloom filter with more bits, counting Bloom filter could count a certain location has been added how many times. Then it becomes possible to delete elements from it by simply reducing the count field. Invertible Bloom Lookup Table, pretty much like what we do in a hash table, hashes the inserted item and thus determines the location of insertions. Instead of only one hash function in the case of hash table, similar to Bloom filters, IBLT uses more than one hash functions and inserts items to multiple distinct locations. And the fields in Bloom filter now change to what we called `bucket`s in IBLT. 

## Principles

Each `bucket` contains `count`, `dataSum`, and `hashSum` fields.  
`count` field is analogy to counting Bloom filter. It keeps track of how many times this locations has been added or subtracted.  
`dataSum` is where the actually byte data stored, in the original design [?], this field is a key-value pair.  
`hashSum` is a checksum. It is useful if the IBLT's subtraction is introduced.

We allow different items to be inserted into same locations, unlike a hash table, where hash collision is not favored. Note that, an item will be inserted into multiple different locations. And these locations are unlike to be all the same if they are, then we have a hash collision in our context and IBLT fails. If the number of inserted items is low to a certain threshold, it's very likely that we would have lots of buckets with count field equals to one. In this case, the bucket is consider "pure", since it's been added only once. And the dataSum is the byte itself. We then could subtract the pure bucket from IBLT reversely as we did when insert them, specifically decrement count field, subtract from dataSum and hashSum. Subtraction will be very likely to creat more pure buckets. And thus, the IBLT is unravel layer by layer, all the bytes could be extracted successfully. 

Now, let's bring subtraction between IBLTs to our play. The subtraction itself is pretty much like subtract items one by one. Different from what we discussed above, the previous subtraction does not allow an items to be subtracted if it was not inserted before. To let it happen, we need to do some modification to IBLT. `count` field now has to be a signed integer, since subtraction is likely to introduce negative numbers. The criteria of a pure bucket now has changed. Because of subtraction, a bucket with count positive one and negative one could be potentially pure. If a bucket's absolute value of count is one, we could not guarantee it has been added only once. It may go through a number of addition and subtraction. `hashsum` must be added to the bucket's data structure. Every time an item is inserted or deleted, we would calculated its hash and add or subtract correspondingly to the `hashSum`. Now only the buckets with count equal to one as well as `dataSum`'s hash is exactly equal to `hashSum` could be considered pure. 

Following the same principle above, we could decode all the bytes, as long as the number of items is under IBLT's capacity. And we could distinguish whether the item was added to or deleted(subtracted) from the table. Now if we view it as a operation between two parties. Alice generates its IBLT by inserting all the elements in her set, and sends it the Bob. Upon receiving the IBLT, Bob also generates one IBLT by inserting all the elements in his set, and he perform a subtraction between Alice's and his IBLT. If decode is successful, Bob could be aware of what Alice has and Bob does not hold in his set, and know what items in his set are unique that Alice would not have. More application would be discussed under the application section.  

## Implementation

Following the design in [?], addition and subtraction were implemented as XOR (exclusive or) operation between bytes. Since it has good properties, for example, byte length does not grow, easy to implement.  
Unlike what IBLT was original designed in [?], key field and value field are separate. KV could actually be combined to one data field. All the operation defined could be supported as long as KV are provided at the same time, which is the case in most of our applications.  
One minor improvement in this implementation is that, an extra pure bucket condition was added to further reduces the length of hashSum, which is the storage overhead. This is a very simple and straightforward idea. If a bucket luckily satisfies `abs(count) == 1 && hash() == hashSum`, it would be falsely considered as pure. In [?] the author suggests to extends `hashSum` length to minimize the probability to be negligible. However, we could simply check whether the index of current bucket is in `index(dataSum)`. With this simple modification, the storage overhead of hash checksum could be further reduced.

## Example

Download the package by simple run  
`go get github.com/SheldonZhong/go-IBLT`

import the package
```go
import "github.com/SheldonZhong/go-IBLT"
````

without addition 
```go
    // number of bucket is 80
    // an item is a byte slice of 16 bytes
    // use 4 hash function, i.e., four different locations
    table := iblt.NewTable(80, 16, 4)
    
    for _, b := range bytes {
    	table.Insert(b)
    }
    
    diff, err := table.Decode()
```

for set reconciliation
```go
    tableAlice := iblt.NewTable(1024, 16, 4)
    
    for _, b := range bytesAlice {
    	tableAlice.Insert(b)
    }
    
    bytes := tableAlice.Serialized() // WIP: not implemented yet
```
and sends the serialized bytes to the other side
```go
    // parameters should be the same
    talbeAlice := iblt.DeSerialized() // WIP: not implemented yet
    tableBob := iblt.NewTable(1024, 16, 4)
    
    for _, b := range bytesAlice {
        	tableAlice.Insert(b)
    }
    
    diff, err := tableBob.Subtract(tableAlice)
    fmt.Println("Bob has and Alice does not have:")
    for _, b := range diff.Alpha { // WIP: better API
    	fmt.Println(b)
    }
    
    fmt.Println("Alice has and Bob does not have:")
    for _, b := range diff.Beta { // WIP: better API
        fmt.Println(b)
    }
```
## Applications



## References

[1] Invertible Bloom Filter was first proposed in [Straggler Identification in Round-Trip Data Streams via Newton's Identities and Invertible Bloom Filters](https://arxiv.org/abs/0704.3313).  
[2] ...
