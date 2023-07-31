# Go Sorted Set

This project is a Go implementation of a sorted set (zset) data structure, similar to the one used in Redis. The sorted set is implemented using a combination of a skip list and a map for efficient operations.

The skip list allows for efficient sorted insertions, deletions, and searches, while the map provides quick access to elements by member.

## Features

- `ZAdd`: Add member with a specific score to the set.
- `ZRange`: Get a range of members in the set by their scores.
- `ZRem`: Remove a member from the set.
- `ZRemByScore`: Remove members with a specific score from the set.
- `ZRemRangeByScore`: Remove all members in a range of scores.
- `ZIncrBy`: Increment the score of a member in the set.

## Testing

This repository includes tests for the SortedSet object and its methods. The test code is located in the `sset_test.go` file.
