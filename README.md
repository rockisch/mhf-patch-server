# MHF Patch Server

Patch server for the [new Monster Hunter Frontier launcher](https://github.com/rockisch/mhf-launcher).

## About

The pupose of this patch server is to provide a way for MHF servers to distribute file changes to users while minimizing the amount of time spent checksumming and the amount of data transfered.

With this, it should be much easier to ensure users always get the latest translations, custom skins, game fixes, etc.

## Running

1. Put the `patch_config.json` and the executable in the same folder (or use `--config` to specify a custom config location).
2. Update the `GameFolder` config to point to a folder containing the files you would like clients to sync to. This will most likely be an entire retail version of MHF with patches applied, but it could also be only a few files.
3. Update the `SignV2.PatchServer` config on your [Erupe](https://github.com/ZeruLight/Erupe) server to point to the address the patch server will run on. For example: `https://10.20.30.40:8081`.
4. Run the executable.

When changing files inside `GameFolder`, the server needs to be restarted.

## How does it work

As mentioned, the way the client and the server communicate was designed in order to minimize the time clients spend checking for files before before starting the game, and the amount of data the server will need to transfer after updates.

When clients connect for the first time to the server, the exchange goes more or less like this:

1. Client connects to patch server for the first time
2. Patch server returns a list of macthing paths and versions the client must have and a 'folder version' that will be used in subsequent requests
3. Client matches the returned paths and versions against the local files, and **only** downloads files that do not match

When doing subsequent connections, we can use the 'folder version' to skip the client check:

1. Client connects to patch server, sending the last received 'folder version'
    - If the folder version matches the patch server's version, it returns no data and the client can log in directly
    - If the folder version does not match the patch server's version, it follows the same flow as when connecting for the first time

The main launcher implementation is also configured to only do 1 request at a time, in order to not overwhelm the server.

**Important**: Do notice that the patch server does not give any guarantees that the client will have the same version as the server. This means it **does not stop cheaters from running custom launchers that bypass the verification**. It does make it much more convenient for normal users, but that's the extent of it.

## Config

| Name | Default | Description |
| - | - | - |
| `Port` | `8081` | Port the server will run on |
| `GameFolder` | `./game` | Folder containing the files you want your clients to match |
| `Force` | `false` | Forces the client to check all files client-side. The client will still only request files that differ from the server. This is disabled by default so clients can log in faster, but it can be enabled to guarantee file consistency. |

