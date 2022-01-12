# rexport
Re-encrypt existing EBS snapshots with a new KMS customer managed key and share them to another AWS account

# Setup
Download a binary for your architecture from the releases tab.

Set up your config file (YAML) as seen in `sample-config.yml`

# Run
Run the binary with the argument for your config file

e.g.,
```
./rexport -config myconfig.yml
```

It will do the following:
1. Create a new volume from your existing snapshots, setting the new encryption
1. Create a new snapshot from the new volume
1. Dump the info about the volumes, snapshots, etc to a JSON file with a `(random session id).json`. The script refers to this export as a "folio"

Since snapshot creation can take a very long time the script was designed with a concept of resuming a session later.

Now, to actually share the snapshot you resume the session by running something like:

```
./rexport -config myconfig.yml -folios abcdefg.json -resume
```

From there it checks status of all the snapshots and shares the ones that are "available". You can run this "resume" mode as many times as you want as re-sharing a snapshot multiple times causes no harm.


