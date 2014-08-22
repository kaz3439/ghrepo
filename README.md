ghrepo
======
Simple Command Line Interface for Github Repository

Description
======
ghrepo is command line interface to operate repositories on Github.
You can create or edit a repository, and open repository on Browser with it.

ghrepo uses Github API ver. 3, please see (this page)[https://developer.github.com/v3/].

Installtion
======
Currentry, you can install ghrepo in GO way. ghrepo is installed in $GOPATH/bin
```
go get github.com/kaz3439/ghrepo
```

Usage
=====
1.Create a repositry
```
> ghrepo create [NAME]
```
2. Edit reposiotry
```
> ghrepo create [OWNER] [NAME]
```
3. Open repository
```
> githubrepo open
> githubrepo open [REMOTE]
> githubrepo open [REMOTE] [BRANCH]
> githubrepo open [GITHUB_PATH]
```
