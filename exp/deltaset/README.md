
## Delta Set
The ordered set uses a delta set to keep track of the element changes in the container. The delta set contains:

- **Updated set:** The set of elements that have been updated or added. If a newly added element is later deleted, it willn't be in either the updated set or the deleted set.
 
- **Deleted set:** Only elements previously committed can be included in this set; removal of a newly added element is not recorded in the deleted set; it is simply removed from the update set. If a deleted element is re-added, it’s removed from the deleted set.

