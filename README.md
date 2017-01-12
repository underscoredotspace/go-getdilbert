# Go Get Dilbert
Downloads Dilbert strip for today's or specified date and saves it as`./images/[yyyy]/[mm]/[dd].gif`

### Usage
| Command | Result |
|:----------|:-------|
| ./go-getdilbert 1986-08-28 | downloads strip from 28th August 1986 |
| ./go-getdilbert | downloads today's strip |

### Build
````
# go test
# go build
````

### Installation
Ideally used with crontab (`crontab -e`). I run it every day at 0600 UK time, but I haven't checked exactly what time they are uploaded so YMMV:
````
0 6 * * *  $HOME/$PATHTOBINARY
````

### Web Strip Viewer
If you're feeling brave, long ago I created a PHP site to present the comics with nginx config. It's not quite set up for this getter, but it's easily modified to work. [underscoredotspace/d-downloader](https://github.com/underscoredotspace/d-downloader)