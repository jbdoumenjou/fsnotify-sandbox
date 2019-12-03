There is several point to synchronize the configuration with a mounted volume in docker.

It depends on the volume mounted (file or directory), on the item listen by inotify (file or directory) and the type of the action on the file (modify or rename).
With the delete action, the link with the binded file seems broken and the docker container doesn't see the file update anymore.

I used (on container and host):

 * Golang as editor to generate a write
 * Vim to generate the Moved + Chmod from host


| Volume mounted | Listen on     | File action   | From Host | From Container | Comment                                                                         |
|----------------|---------------|---------------|-----------|----------------|---------------------------------------------------------------------------------|
| file           | file          | Write         | OK        | OK             | Live reload on save. Strange behavior on delete, it may be the editor           |
| file           | file          | Rename        | KO        | OK             | First event receive then nothing. All local changes (in container) are detected |
| file           | directory     | Write         | KO        | OK             | Edition works only from container to host                                       |
| file           | directory     | Rename        | KO        | OK             | Edition works only from container to host                                       |
| file           | both          | Write         | OK        | OK             |                                                                                 |
| file           | both          | Rename        | KO        | OK             | Have the first event then is broken                                             |
| file           | both + reload | Rename        | KO        | OK             |                                                                                 |
| directory      | file          | Write         | OK        | OK             |                                                                                 |
| directory      | file          | Rename        | KO        | KO             | broke the link in the two directions                                            |
| directory      | directory     | Write         | OK        | OK             |                                                                                 |
| directory      | directory     | Rename        | OK        | OK             |                                                                                 |
| directory      | both          | Write         | OK        | OK             |                                                                                 |
| directory      | both          | Rename        | OK        | OK             |                                                                                 |