# b7s-attributes

## Description

This code in this repo aims to provide a control panel that can be used to create, examine, validate and otherwise manage b7s attributes files.
Attributes file describes the parameters and capabilities of a node.

The executable used to manage attributes files can be found in `/cmd/b7s-attributes`.

### Attributes

Attributes are key value pairs, for example:
```
Arch: x86_64
CPU: 4
Comment: ThisIsAComment
Database: MySQL
Env: Production
ExpirationDate: 1683255281
Filesystem: ext4
OS: Linux
PrivateIP: 10.0.0.2
RAM: 16GB
Storage: 100GB
StorageType: SSD
```

Attributes are collected from the environment variables set on the node.
Attributes are typically set with the same prefix, which is `B7S_` by default.
In order to set the attribute for the `Storage` attribute, you can set the `B7S_Storage` environment variable.
An arbitrary value can be used for the prefix.

## CLI Tool Usage

This section describes usage and actions possible from the `cmd/b7s-attributes` CLI tool.

### Examine Attributes

It is possible to do a dry-run and simply see which attributes would be collected from the environment, without creating the attributes file.
This can be done using the following command:

```console
b7s-attributes show --prefix B7S_ --limit 50 --ignore B7S_INTEG_RUNTIME_DIR
```

Some of this values are also the default ones, so the same can be achieved by using:

```console
b7s-attributes show --ignore B7S_INTEG_RUNTIME_DIR
```

Tool will return the list of attributes set.

### Create Attributes File

When ready, you can write the current values of the attributes to an attribute file.
Usage is similar to the previous one:

```console
 ./b7s-attributes create --ignore B7S_INTEG_RUNTIME_DIR --output attributes.bin
```

This creates a new file called `attributes.bin`.

### Examine Attributes File

To print the content of an attributes file, you can use:

```console
b7s-attributes print ./attributes.bin
```

### Sign an Attributes File

To prove the origin of an attributes file, the operator running the node can sign it with the nodes private key.
This can be done using the following command:

```console
b7s-attributes update sign --signing-key /path/to/private-key.bin ./attributes.bin
```

This command creates a backup of the old file under the `old-filename.old` name.
For example, in the example above, the file created would be `attributes.bin.old`.

This command alters the content of the attributes file by adding information about the signer (the ID) and the signature.
The signature guarantees that the attributes written in the file have not been tampered with and come from the node with the recorded ID.

## Provide an Attestation for the Attributes FIle

If a third-party wants to provide an attestation for a given signed attributes file, it can do this by executing the following command:

```console
b7s-attributes update attest --signing-key /path/to/attestors/private-key.bin ./attributes.bin
```

The attestation serves as a sort of testimony for the correctness of the value of the attributes listed in the file, as well as the node from which they originate.
The attestation cannot be given to an *unsigned* attributes file.
This is because the attestation serves to confirm that the node in question has the specified attributes.

Since this action also modifies the provided file, it will also create a backup under the `filename.old` name.

### Validate Attributes File

To validate the content of the file, use:

```console
b7s-attributes validate ./attributes.bin
```

This action will verify:
    - attributes are listed in the correct format
    - signature (if signed) is valid and corresponds to the attributes listed
    - attestations (if any) are valid and correspond to the attributes and signature combination

### Upload Attributes File to IPFS

When satisfied with the configuration, the CLI tool can be used to upload the attributes file to IPFS.
Note that this tool uses the [web3storage](https://web3.storage/products/web3storage/) API to upload the file to IPFS.
To do this, it relies on having the `WEB3STORAGE_TOKEN` environment variable with the token value for the mentioned API.

This functionality aims to be a convenience to the user, and uploading the file to IPFS in any other method is also completely valid.

On successful upload, the tool returns the CID of the uploaded file:

```console
b7s-attributes upload ./attributes.bin
bafybeih62bldr42slbkegccznrdvmtinlxmvy2x77o2anenaehba76qy4y
```

The file can be retrieved using any public IPFS gateway, for example https://ipfs.io/ipfs/bafybeih62bldr42slbkegccznrdvmtinlxmvy2x77o2anenaehba76qy4y/attributes.bin.

### Create an IPNS Record

After uploading the file to IPNS, it's possible to create an IPNS record.
IPNS record will typically be created with the private key of the node that generated the attributes.

Input for this command is the CID of the attributes file uploaded to IPFS in the previous step.

```console
b7s-attributes add-name --key ./signer/priv.bin bafybeih62bldr42slbkegccznrdvmtinlxmvy2x77o2anenaehba76qy4y
k51qzi5uqu5dl3kqdoz5i2r1ijy7rm660rrynfvwlkdgpnxxtvyarlzrx8bami
```

Like in the previous step, the tool relies on the HTTP API from [web3name](https://web3.storage/products/w3name/), and having the `WEB3STORAGE_TOKEN` environment variable set.

The file can be retrieved using any public IPFS gateway that supports IPNS, for example https://ipfs.io/ipns/k51qzi5uqu5dl3kqdoz5i2r1ijy7rm660rrynfvwlkdgpnxxtvyarlzrx8bami/attributes.bin.
