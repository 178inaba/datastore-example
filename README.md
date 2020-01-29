# datastore-example

## Local Setup

### Datastore Emulator

https://cloud.google.com/datastore/docs/tools/datastore-emulator

Install datastore emulator.

```console
$ gcloud components install cloud-datastore-emulator
```

Run datastore emulator.

```console
$ gcloud beta emulators datastore start
```

Set environment variables.

```console
$ $(gcloud beta emulators datastore env-init)
```

## License

[MIT](LICENSE)

## Author

Masahiro Furudate (a.k.a. [178inaba](https://github.com/178inaba))  
<178inaba.git@gmail.com>
