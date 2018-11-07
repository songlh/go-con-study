#### Critical Section Mining


##### mutex_channel_patch_finder

The script mutex_channel_patch_finder is used to find using channel to patch mutex and using mutex to patch channel.  
- Clone the target repo into your local path
- cd the target repo path and run the command in your terminal "git log --pretty="%H" > commit-hash-log.txt"
- python3 mutex_channel_patch_finder.py

** NOTE: The script is only consider that the buggy code and the patch code are in a contiguous block. It may be has a FN. **
