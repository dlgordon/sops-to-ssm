# sops-to-ssm
Simply put, this takes a sops file of a certain structure and sync the data intelligently with SSM.

The SOPS file should be a YAML file in 2 sections, passed in with the `--sops-file-path` flag

```
environment:
  NAME: VALUE
secrets:
  NAME: VALUE
```

Of course, the file should be encrypted - sops-to-ssm will decrypt it based on the `sops:` section in the file.

Specify the path prefix you want to upload/sync (ntoe that sops-to-ssm current can't support root syncing) with the `--path-prefix` flag and sops-to-ssm will match the key/values 
in the sops file with the appropriate SSM parameter and create it if it is missing or update it if it changed.