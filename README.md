# Word Frequencies

## Function SortedValues

This functions constructs a companion slice for the map that it gets as an argument.

Go maps are unsorted, and there is no way to sort the entries in a Go map itself. So, stepping through a map using 
`range m` 
will produce the entries in an undetermined order. But sometimes we do need to be able to step through a map, or, 
for that matter, a subset of that map, in a particular order. The method, or trick if you want, to accomplish that is to
construct a companion slice to the map, which contains the keys of interest in that map in a particular order. To 
access the map in that order, step through the slice containing the keys, and use the keys to access the values in 
the map.

In our case, we want to be able to access ngrams in order of their number of occurrences. To accomplish that, we 
construct a slice `s` with the ngram keys of interest first, and then sort `s` based on the number of 
occurrences. So, suppose we want to sort all bigrams in the ngrams map. The first thing we do is to construct a slice 
that contains all bigrams in the ngrams map.

To be able to construct a slice, we must tell the `make` function its 
maximum length. However, we do not know how many bigrams there are in the ngrams map. But it is certainly not more 
than the total number of entries in the ngrams map, so we could take the length of that map as an estimate. If our 
program is very pressed for space, we could at first walk through the ngrams map, counting the number of bigrams, and 
then construct a slice with the exact number of bigrams. It means we have to step through the ngrams map one extra 
time. As we will walk through the map in the next step as well, our algorithm is at least O(2 x length(ngramps map)).
If our program is tight on runtime, we could count the number of monograms, bigrams, trigrams, and so on, while we 
read the ngrams from file to store them in the ngrams map. The disadvantage from a maintenance point of view is that 
the code counting the number of ngrams in this situation is textually removed from the usage of those counts.

For now, we'll take the length of the ngrams map as an estimate for the length of the companion slice and take the 
extra unused space in the slice for granted.

Next, we copy all bigrams to the companion slice. And then we are ready to sort the slice in some particular order. 
We will use the `sort.Slice` function to do that: that is no need to write our own sorting algorithm here. That 
function, `sort.Slice`, requires an ordering function, a function through which the `sort.Slice` can compare two 
entries in the slice. For that, we introduce an ordering function that takes integer arguments (here: the number of 
occurrences of two ngrams), and returns true if the arguments are in order and false if they are in reverse order. 
Now, we want to order the slice by the number of occurrences of each bigram. So, taking two 
bigrams from the slice at positions `i` and `j`, we use the ngram map to find their respective number of 
occurrences: we compare `ngrams[s[i]]` with `ngrams[s[j]]`.

After `sort.Slice` has done its job, we will have a slice `s` where `ngrams[s[i]] >= ngrams[s[j]]` for all i>j. In 
other words, we now have an indirectly sorted map of bigrams.
