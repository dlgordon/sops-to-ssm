# sops-to-ssm
Simply put, this takes a sops file of a certain structure and sync the data intelligently with SSM.

```
ssm:
  path-prefix: <path to prefix>
environment:
  NAME: VALUE
secrets:
  NAME: VALUE
```

Of course, the file should be encrypted - sops-to-ssm will decrypt it based on the `sops:` section in the file.

There are currently two commands - `check` and `push`

## check
This will just do a diff against your local file, specified with `--sops-file-path` and the stored data in SSM. For example:
```
./sops-to-ssm check  -sops-file-path muhenvironment.yaml
No New Parameters Found
~ changed parameter PARAM 1
~ changed parameter PARAM 2
- removed parameter PARAM 3
```

## push
Push your local file up to SSM and store the environment data under the prefix defined in the `--sops-file-path` file. For example:
```
./sops-to-ssm push -sops-file-path muhenvironment.yaml
```

By default, items that exist in SSM but not your local file will not be removed. You can chage this with the `-remove-missing` flag. For example:
```
./sops-to-ssm push -remove-missing -sops-file-path muhenvironment.yaml
```
