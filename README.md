# cmd-rotate
Run command and log rotation.

# Installation

Download from  https://github.com/tkuchiki/cmd-rotate/releases

# Usage

```shell
$ ./cmd-rotate --help
usage: cmd-rotate [<flags>] [<args>...]

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --stdout-log="stdout.log"  stdout log name
  --stderr-log="stderr.log"  stderr log name
  --merge-log                stdout and stderr write to same file
  --logdir=$TMPDIR           log directory location
  --file-mode="0644"         file permission
  --file-size=10485760       rotate file size
  --file-num=20              number of files
  --version                  Show application version.

Args:
  [<args>]  command

```

# Examples

```shell
$ ./cmd-rotate --file-size=10000 "for i in \$(seq 1 10000); do echo \$i; echo \$((\$i + 100000)) 1>&2; done"

$ ll ${TMPDIR}/std*.log*
-rw-r--r--  1 tkuchiki  tkuchiki   9982  1 20 16:01 /tmp/stderr.log
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705740104248
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705764957069
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705790779494
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705815530451
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705841166791
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stderr.log-1484895705865680355
-rw-r--r--  1 tkuchiki  tkuchiki   8876  1 20 16:01 /tmp/stdout.log
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:01 /tmp/stdout.log-1484895705753610012
-rw-r--r--  1 tkuchiki  tkuchiki  10005  1 20 16:01 /tmp/stdout.log-1484895705789271570
-rw-r--r--  1 tkuchiki  tkuchiki  10005  1 20 16:01 /tmp/stdout.log-1484895705824162040

$ wc -l ${TMPDIR}/std*.log*
    1426 /tmp/stderr.log
    1429 /tmp/stderr.log-1484895705740104248
    1429 /tmp/stderr.log-1484895705764957069
    1429 /tmp/stderr.log-1484895705790779494
    1429 /tmp/stderr.log-1484895705815530451
    1429 /tmp/stderr.log-1484895705841166791
    1429 /tmp/stderr.log-1484895705865680355
    1775 /tmp/stdout.log
    2222 /tmp/stdout.log-1484895705753610012
    2001 /tmp/stdout.log-1484895705789271570
    2001 /tmp/stdout.log-1484895705824162040
    2001 /tmp/stdout.log-1484895705859756905
   20000 total
```

```shell
$ ./cmd-rotate --file-size=10000 --merge-log "for i in \$(seq 1 10000); do echo \$i; echo \$((\$i + 100000)) 1>&2; done"

$ ll ${TMPDIR}/std*.log*
-rw-r--r--  1 tkuchiki  tkuchiki   8869  1 20 16:04 /tmp/stdout.log
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:04 /tmp/stdout.log-1484895876068506745
-rw-r--r--  1 tkuchiki  tkuchiki  10005  1 20 16:04 /tmp/stdout.log-1484895876083855300
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:04 /tmp/stdout.log-1484895876099313607
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:04 /tmp/stdout.log-1484895876114493978
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:04 /tmp/stdout.log-1484895876130060418
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:04 /tmp/stdout.log-1484895876145418714
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:04 /tmp/stdout.log-1484895876163829302
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:04 /tmp/stdout.log-1484895876181443424
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:04 /tmp/stdout.log-1484895876196617106
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:04 /tmp/stdout.log-1484895876211929819
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:04 /tmp/stdout.log-1484895876229089344

$ wc -l ${TMPDIR}/std*.log*
    1478 /tmp/stdout.log
    1838 /tmp/stdout.log-1484895876068506745
    1681 /tmp/stdout.log-1484895876083855300
    1667 /tmp/stdout.log-1484895876099313607
    1667 /tmp/stdout.log-1484895876114493978
    1667 /tmp/stdout.log-1484895876130060418
    1667 /tmp/stdout.log-1484895876145418714
    1667 /tmp/stdout.log-1484895876163829302
    1667 /tmp/stdout.log-1484895876181443424
    1667 /tmp/stdout.log-1484895876196617106
    1667 /tmp/stdout.log-1484895876211929819
    1667 /tmp/stdout.log-1484895876229089344
   20000 total

```

```shell
$ ./cmd-rotate --file-size=10000 --merge-log --file-num=5 "for i in \$(seq 1 10000); do echo \$i; echo \$((\$i + 100000)) 1>&2; done"
removed /tmp/stdout.log-1484896835231790682
removed /tmp/stdout.log-1484896835246696757
removed /tmp/stdout.log-1484896835262294460
removed /tmp/stdout.log-1484896835276790627
removed /tmp/stdout.log-1484896835291380222
removed /tmp/stdout.log-1484896835305508999
removed /tmp/stdout.log-1484896835320462226

$ ll ${TMPDIR}/std*.log*
-rw-r--r--  1 tkuchiki  tkuchiki   8869  1 20 16:20 /tmp/stdout.log
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:20 /tmp/stdout.log-1484896835336282138
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:20 /tmp/stdout.log-1484896835351348796
-rw-r--r--  1 tkuchiki  tkuchiki  10001  1 20 16:20 /tmp/stdout.log-1484896835368019663
-rw-r--r--  1 tkuchiki  tkuchiki  10003  1 20 16:20 /tmp/stdout.log-1484896835382824530

$ wc -l ${TMPDIR}/std*.log*
    1478 /tmp/stdout.log
    1667 /tmp/stdout.log-1484896835336282138
    1667 /tmp/stdout.log-1484896835351348796
    1667 /tmp/stdout.log-1484896835368019663
    1667 /tmp/stdout.log-1484896835382824530
    8146 total
```
