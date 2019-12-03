# Goal

Analyse the behavior of [fsnotify](https://github.com/fsnotify/fsnotify) with a docker mounted volume.

# Process

This project contains a simple Go app that show the events caught by fsnotify.
The docker-compose file contains different configuration of the Go app to check the combination of mounted volumes and watched target. 
The file modifications are done on both the container and the host.

I used:

 * Vi/Golang as editor to generate a `Write` notification.
 * Vim to generate the `Rename` notification.


| Volume mounted | Listen on     | File action   | From Host | From Container | Comment                                                                         |
|----------------|---------------|---------------|-----------|----------------|---------------------------------------------------------------------------------|
| file           | file          | Write         | OK        | OK             | Live reload on save.                                                            |
| file           | file          | Rename        | KO        | OK             | First event receive then nothing. All local changes (in container) are detected |
| file           | directory     | Write         | KO        | OK             | Edition works only from container to host                                       |
| file           | directory     | Rename        | KO        | OK             | Edition works only from container to host                                       |
| file           | both          | Write         | OK        | OK             |                                                                                 |
| file           | both          | Rename        | KO        | OK             | Have the first event then the link is broken                                    |
| file           | both + reload | Rename        | KO        | OK             |                                                                                 |
| directory      | file          | Write         | OK        | OK             |                                                                                 |
| directory      | file          | Rename        | KO        | KO             | The link is broken in the two directions                                        |
| directory      | directory     | Write         | OK        | OK             |                                                                                 |
| directory      | directory     | Rename        | OK        | OK             |                                                                                 |
| directory      | both          | Write         | OK        | OK             |                                                                                 |
| directory      | both          | Rename        | OK        | OK             |                                                                                 |

# Conclusion

The main issue occurred with the `Rename` event, the link between the host and the container file is broken and no more event are shared.
The other issue occurred with the `Write` event only when we watch for directory event on a mounted file.

# App Usage

The simple app using `fsnotify` has two options that are available via flags.

## `--target`

The watcher target, could be one of the following elements :

* `file` watch the file (default value)
* `dir` watch the file parent directory
* `both` watch the file and its parent directory

## `--reload`

If `true`, try to create a new watcher on the file when a `rename` event occurred.
The default value is `false`.
