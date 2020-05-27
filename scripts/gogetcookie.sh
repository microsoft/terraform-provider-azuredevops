#!/bin/bash

touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-8787xu.gmail.com=1//0gMSiXWNyBFAsCgYIARAAGBASNwF-L9IrlaDV1-0F-IeVEw1aPlpBD3-DVq2rCVx3JR2sxcXxFgkq4oWXJLpaKfT5dmnUgasVCHc
go-review.googlesource.com,FALSE,/,TRUE,2147483647,o,git-8787xu.gmail.com=1//0gMSiXWNyBFAsCgYIARAAGBASNwF-L9IrlaDV1-0F-IeVEw1aPlpBD3-DVq2rCVx3JR2sxcXxFgkq4oWXJLpaKfT5dmnUgasVCHc
__END__