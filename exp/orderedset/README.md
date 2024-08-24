This container is foundermenatlly a hybrid of a slice and a map. One of the biggest problems with it is to properly handle the deletions. For a map, a deletion takes immediate effect, the size of the map changes immediately. But for a slice, the size of the slice only changes when a new element is added. The size of the slice doesn't change when an element is deleted. To combine the two, there has to be a way to consolidate these two very different behaviors.

## With Shift
A shift will occur and affect all the indices after the deleted element. This makes the container function more like a standard ordered map, which is easier to understand and use. But it is also problematic for a concurrent container for two reasons:

- **Inefficiency:** The shift has to be in realtime, it is inefficient for the large containers, if there are multiple deletions, the shift will be very costly.
  
- **Lower Concurrency:** It prevents many potential concurrent operations. If one thread deletes an element while another thread accesses the container by index, and the access is to elements after the deleted one, a conflict will be detected.  This occurs because the shift involves rewriting all the elements after the deleted one with the element next to it. Concurrent writes to the same location are not allowed.

## No Shift
It will make if more slice-like. For a slice, the size of the slice only changes only when a new element is added. The shift only happens at the commit time. The size of the container doesn't change when an element is deleted at run time. It allows much higher concurrency, threads can update the container without conflicts as long as they don't access the same element. There is no "Side Effect" of the deletions.

These key properties don't come for free. Some question arises:

- With no shift, if the same thread accesses an element it has deleted, what should happen? should it return a default value? How does it differiate between a deleted element and a non-existing element?
  
- Because the container is a combination of a map and a slice, if the element is removed in the slice, Should the key be removed from the map as well?
  
- How to get the real size of the container? Should it be the size of the slice or the size of the map? If 
it is the size of the map, the key shouldn't be removed from the map either.

## Solution

No-Shift is better because it allows higher concurrency. To mitigate the problems of random access to deleted elements, the container provide a function to specifically skip the deleted elements.