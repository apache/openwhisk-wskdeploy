# Creating Tagged Releases of ```wskdeploy```

The most convenient way to create a tagged release for wskdeploy is to build the binaries by adding tag to upstream master.


1. Add a tag to a commit id: ```git tag -a 0.8.9<tag> c08b0f<commit id>```
2. Push the tag upstream: ```git push -f upstream 0.8.9<tag>```

Travis will start the build of 0.8.9 automatically by the event of tag creation.

If the travis build passed, binaries will be pushed into releases page.

If we modify the tag by pointing to a different commit, use ```git push -f upstream 0.8.9<tag>``` to overwrite the old tag. New binaries from travis build will overwrite the old binaries as well.

You can download the binaries, and delete them from the releases page in GitHub if we do not want them to be public.
