# PCF OpsMan Terraform Provider

This Terraform plugin exposes OpsMan primitives as Terraform resources.

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.9 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/sneal/terraform-provider-opsman`

```sh
$ mkdir -p $GOPATH/src/github.com/sneal; cd $GOPATH/src/github.com/sneal
$ git clone git@github.com:sneal/terraform-provider-opsman
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sneal/terraform-provider-sneal
$ ./build.sh
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `./build.sh`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ ./build.sh
...
$ $GOPATH/bin/terraform-provider-opsman
...
```

In order to run the full suite of Acceptance tests, run `./build_and_run.sh`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ ./build_and_run.sh
```

If you need to add a new package make sure you vendor it in the vendor directory using [govendor](https://github.com/kardianos/govendor).

If you need to add a new package in the vendor directory under `github.com/aws/aws-sdk-go`, create a separate PR handling _only_ the update of the vendor for your new requirement. Make sure to pin your dependency to a specific version, and that all versions of `github.com/aws/aws-sdk-go/*` are pinned to the same version.