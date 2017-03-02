# Law

Law perform continuous archiving of PostgreSQL WAL files along with taking
and restoring base backups.

## Installation

Download and install :

```
$ go get github.com/cyberdelia/law
```

## Usage

Law make use of environment-variable, the base one are:

 - ``STORAGE_URL``: URL indicating where files are stored.
   * For file storage: ``file:///tmp/``
   * For S3 storage: ``s3://s3.amazonaws.com/bucket_name``
 - ``DATABASE_URL``: URL to database.

S3 storage might requires one or more of theses variables:

 - ``AWS_ACCESS_KEY_ID``: An AWS access key.
 - ``AWS_SECRET_ACCESS_KEY``: An AWS secret key.
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


## PostgreSQL configuration

In order for law to work you'll need to setup PostgreSQL like so:

```
wal_level = archive
archive_mode = on
archive_command = 'law -storage <ssn> wal-push -segment %p'
archive_timeout = 60
restore_command = 'law -storage <ssn> wal-fetch -destination "%f" -segment "%p"'
```

## Limitations

This is alpha software, don't use it for anything serious.
