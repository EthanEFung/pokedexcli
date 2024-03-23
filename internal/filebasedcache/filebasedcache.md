# So what should happen here exactly?

what we would like is to have a folder for all the incoming requests, and each request can have a seperate file that we can read from. conceptually it would look like:

```
	|- cached-files
		|- query/hello
```

the main issue that we have is that we cannot have filenames with slashes so we need to have a ledger file (something simple like a csv that keeps track of the request full uri and the filename of the stored response)

```
	|- ledger.csv
	|- cached-files
		|- ab2js
		|- c9exj
```

this ledger file should be an ever growing file of all the transaction history the program has undergone.

each ledger entry should have the following:

	1. timestamp of when the action was invoked
	2. an action this should be an enum of
		- ADDED: created a new file and was stored in the files directory
		- EVICTED: deleting a file that was stored in the files directory
		- READ: searched the ledger for a file using a key and was returned

## Add

this function should add a new file to the directory of cached files. when the add is called, we also want to check the ledger to see how many files there are stored in the cached directory. We loop over the ledger and the first N ADDED and READ and _not_ EVICTED files to a whitelist, and loop over the existing files removing all that did not make the whitelist.

## Get

After the routine validates that the file exists in the directory, we write a ledger entry stating
that the file has been READ. This file will be returned to the user.
this function should read the ledger from bottom to top searching for the last N added files
