# MHF Patch Server

Patch server for the [new Monster Hunter Frontier launcher](https://github.com/rockisch/mhf-launcher).

## About

The pupose of this patch server is to provide a way for MHF servers to distribute file changes to users while minimizing the amount of time spent checking for files and the amount of data transfered.

With this, it should be much easier to ensure users always get the latest translations, custom skins, game fixes, etc.

## Running

1. Put the `patcher_config.json` and the executable in the same folder (or use `--config` to specify a custom config location).
2. Configure the `GameFolder` config to point to a folder containing the files you would like clients to sync to. This will most likely be the entire retail version of MHF with patches applied.
3. Update your Erupe server's `SignV2.PatchServer` config to point to the address the patch server will run on. For example: `https://10.20.30.40:8081`.
4. Run the executable.

The server needs to be restarted after changing files inside `GameFolder`.

## How does it work

As mentioned, the way the client and the server communicate was designed in order to minimize the time clients spend checking for files before before starting the game, and the amount of data transfer the server will need to do.

When clients connect for the first time to the server, the exchange goes more or less like this:

- Client connects to patch server for the first time
- Patch server returns a list of paths and versions the client must match and a 'folder version' that will be used in subsequent requests
- Client checks its local files, and **only** downloads files that do not match their local version

When doing subsequent connections, we can use the 'folder version' to skip the client check:

- Client connects to patch server, sending the last received 'folder version'
    - If the received folder version matches the patch server version, it returns no data and the client can log in without verifying data
    - If the received folder version does not match the patch server version, it follows the same flow as when connecting for the first time

The main launcher implementation is also configured to only do 1 request at a time, in order to not overwhelm the server.

**Important**: Do notice that the patch server does not give any guarantees that the client will have the same version as the server. This means it **does not stop cheaters from running custom launchers that bypass the verification**. It does make it much more convenient for normal users, but that's the extent of it.

## Config

| Name | Default | Description |
| - | - | - |
| `Port` | `8081` | Port the server will run on
| `GameFolder` | `./game` | Folder containing the files you want your clients to match

