# Law

Law perform continuous archiving of PostgreSQL WAL files along with managing backups.

## Installation

Download and install :

```
$ go get github.com/cyberdelia/law
```

## Usage

Law make use of environment-variable, the base one are:

 - ``STORAGE_URL``: URL indicating where files are stored.
   * For file storage: ``file:///tmp/law``
   * For S3 storage: ``s3://law.s3.amazonaws.com/``
 - ``DATABASE_URL``: URL to database.

S3 storage might requires one or more of theses variables:

 - ``AWS_ACCESS_KEY``: An AWS access key.
 - ``AWS_SECRET_KEY``: An AWS secret key.
 - ``AWS_SECURITY_TOKEN``: An AWS STS Token.

Law has 4 subcommands :

 - ``wal-push``: Push wal archive to storage. 
  
  Example: ``law wal-push -segment %p``

 - ``wal-fetch``: Fetch wal archive from storage.
 
   Example: ``law wal-fetch -segment %p -destination %p``
   
 - ``backup-push``: Push a backup to storage.
  
   Example: ``law backup-push -cluster /var/lib/database``

 - ``backup-fetch``: Fetch a backup from storage.
   
   Example: ``law backup-fetch -cluster /var/lib/database``


## Limitations

This is alpha software, don't use it for anything serious.
Take also note of theses specifics problems:

 - Upload or download aren't made in parallel, this is highly inefficient.
