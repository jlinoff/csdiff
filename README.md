# csdiff
side by side diff tool with colorization and regular expression support for filtering

_Watch this space - still under development_

Command line tool that does a side by side diff of two text files
with regular expression filtering and ANSI terminal colorization.

It is useful for analyzing text files that have patterns like
timestamps that can easily be filtered out.

Diffs.
```bash
$ bin/Darwin-x86_64/csdiff test/td03.txt test/td04.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766790/21bf7818-494d-11e7-88c2-84eea6022a0e.png" alt="example-1">

Filter out timestamps.
```bash
$ bin/Darwin-x86_64/csdiff -r '\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}' 'yyyy-mm-dd HH:MM:SS' test/td03.txt test/td04.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766793/2d0d2530-494d-11e7-849b-a03bec7a1a5c.png" alt="example-2">

Customize colors.

```bash
$ bin/Darwin-x86_64/csdiff -c cd=bold,fgRed \
                           -c cm=bold,fgBlue \
                           -c s=bold,fgMagenta \
                           -c lm=bold,fgGreen \
                           -c llo=bold,fgCyan \
                           -c rlo=bold,fgCyan \
                           test/td01.txt test/td02.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766795/32be864a-494d-11e7-9b37-1554c4821494.png" alt="example-3">

TODO
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766798/38e82800-494d-11e7-8e0e-e429322d993e.png" alt="example-4">
