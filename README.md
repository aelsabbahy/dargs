# Dargs - Dynamic CLI arguments and completion
[![Build Status](https://circleci.com/gh/aelsabbahy/dargs.png)](https://circleci.com/gh/aelsabbahy/dargs)
**
[![Twitter Follow](https://img.shields.io/twitter/follow/aelsabbahy1.svg?style=social&label=Follow&maxAge=2592000)]()
[![Blog](https://img.shields.io/badge/follow-blog-brightgreen.svg)](https://medium.com/@aelsabbahy)

# Introduction

<a href="https://asciinema.org/a/AwMi8FkS33bQtna7yTlvAEYwf?autoplay=1" target="_blank"><img src="https://user-images.githubusercontent.com/6783261/30841588-3d75422a-a24b-11e7-8108-238d805dbde6.gif" alt="asciicast"></a>

## What is Dargs?

Dargs is a tool that allows you to define dynamic argument replacements and completions for any CLI command.

For example, using dargs one can define rules to:

* **ssh into an ec2 instance by instance id:**

  `ssh i-xxxxxxxx`

* **Run packer build using YAML file format instead of JSON:**

    `packer build demo.yml`

* **ssh completions for AWS ec2 instance names:**

    `ssh ec2:instance_name<tab>`

## Installation

If you have go
```bash
go get -u github.com/aelsabbahy/dargs
```

If not, download the stand-alone binary

Linux:
```bash
curl -L https://github.com/aelsabbahy/dargs/releases/download/v0.0.1/dargs_linux_amd64 -o /usr/local/bin/dargs
chmod +rx /usr/local/bin/dargs
```

OSX:
```bash
curl -L https://github.com/aelsabbahy/dargs/releases/download/v0.0.1/dargs_darwin_amd64 -o /usr/local/bin/dargs
chmod +rx /usr/local/bin/dargs
```


### Completions installation

**Bash:**

Add the following to your `~/.bashrc`
```bash
if [[ "$(ls -A ~/.dargs/completions/bash)" ]];then
  for f in ~/.dargs/completions/bash/*;do source "$f";done
fi
```
**Zsh:**

Add the following to your `~/.zshrc`
```bash
if [[ "$(ls -A ~/.dargs/completions/zsh)" ]];then
  for f in ~/.dargs/completions/zsh/*;do source "$f";done
fi
```

## Tutorial

**Note:** This section walks through the basics of dargs, for more detailed information, see [manual](#manual)

Now that you have dargs installed, lets use dargs to enhance ssh with ec2 abilities and give packer YAML support.

First, create the following `~/.dargs.yml` config file:
```yaml
imports:
  - https://raw.githubusercontent.com/aelsabbahy/dargs/master/examples/quick_start.yml

commands:
  - name: /usr/bin/ssh
    wrapper: ssh
    # Use fzf to fuzzy complete (must have fzf installed)
    #fzf-complete: true
    completers:
      - ec2_name
    transformers:
      - ec2_name # transforms ec2:instance_name -> instance-id (i-xxxxxxxx)
      - ec2_id # Transforms instance-id -> PrivateIpAddress

  - name: /usr/local/bin/packer
    wrapper: packer
    transformers:
      - yaml2json # Converts yaml to json
```

Lets test out the `ec2_id` and `ec2_name` filters using `dargs run`.
```
# -n = dry run, only print what would have been executed
$ dargs run -n -- ssh i-05d3e9e370805d0b2
ssh 10.0.0.9

# Same instance but now I'm passing in the instance name
# This was processed by both transformers, try running with `-d` to see what's happening
$ dargs run -n -- ssh ec2:test
ssh 10.0.0.9

# We can pass in the transformers on the command line and run any command with dargs.
# The ~/.dargs.yml config allows us to omit the --transformers, -t flag for commands defined in it
$ dargs run -t ec2_id -- echo i-05d3e9e370805d0b2
10.0.0.9

# Here's an example of the yaml2json filter using cat
$ cat /tmp/test.yml
---
foo: bar
moo: cow

# -v is verbose, prints out the command before executing it
$ dargs run -v -t yaml2json -- cat /tmp/test.yml
cat /tmp/test.json
{
  "foo": "bar",
  "moo": "cow"
}
```

Typing all that out wouldn't be very efficient. To make this simpler we can either:
```
alias ssh='dargs run -- ssh'
```
OR:
```
$ dargs generate-bins
INFO Wrote /home/***/bin/ssh
INFO Wrote /home/***/bin/packer

# and ensure ~/bin has high precidence in your PATH
$ export PATH=$HOME/bin:$PATH
```

Lets test out completions
```
$ dargs generate-completions -f
INFO Wrote /home/***/.dargs/completions/bash/zzdargs_ssh
INFO Wrote /home/***/.dargs/completions/bash/zzdargs_packer

$ source ~/.dargs/completions/bash/*

$ ssh ec2:te<tab>

# We can run the dargs completion command stand-alone to see the results
# syntax: dargs completions -- command prev_arg current_arg
$ dargs completions -- ssh "" "ec2:te"
ec2:test
```

## Manual

### Commands
Dargs commands:
* `dargs run [flags] -- command command_args..`
  * Transforms arguments and runs command.
* `dargs generate-bins`
  * Generate wrapper bin scripts, alternative to `alias cmd='dargs run -- cmd'`
* `dargs generate-completions [flags]`
  * Generate bash/zsh completion scripts
* `dargs complete -- cmd prev_arg cur_arg`
  * Used by bash/zsh completion scripts, useful for debugging

### The dargs config file: ~/.dargs.yml

Dargs configuration consists of four top level keys:
* **imports**    - imports dargs config from file glob or URLs
* **transformers** - Transform matching CLI args before executing the command
* **completers** - Complete matching CLI args
* **commands**   - Mapping transformers/completers to commands

```yaml
# Files to import transformers, completers, and (optionally) commands from
imports:
  # Can be url
  - https://raw.githubusercontent.com/aelsabbahy/dargs/master/examples/quick_start.yml
  # Local path or glob
  - ~/.dargs/config.d/*.yml

# These transform CLI arguments
transformers:
    # Name of transformer
  - name: ec2_id
    # Transformer only runs on CLI args matching this regex
    # Regex groups are available as $RE_0, $RE_1 or if named group $RE_GROUPNAME
    match: '^i-[0-9a-z]+$'
    # (optional) Only match if the previous argument also matches this regex
    # prev-match: '--some-flag'
    # Seconds to cache results (604800 seconds = 1 week)
    cache: 604800
    # Command to run, replacing argument with output of the command
    # To convert an argument into many arguments, use newline as a separator
    # Example:
    # echo -e '-i\n~/.ssh/some_key.pem\n10.11.12.13'
    # Will expand into -i ~/.ssh/some_key.pem 10.11.12.13
    command: |
      aws ec2 describe-instances \
        --filters 'Name=instance-state-name,Values=running' --instance-ids "$RE_0" \
        --query 'Reservations[*].Instances[*].[PrivateIpAddress]' --output text

# These complete CLI arguments
completers:
    # name, match, prev-match, cache have same meaning as transformers
    - name: ec2_name
      match: '^ec2:(?P<name_prefix>.+)$'
      cache: 120
      # Newline separated list of completions
      command: |
      aws ec2 describe-instances \
      --filters "Name=tag:Name,Values=${RE_name_prefix}*" 'Name=instance-state-name,Values=running' \
      --query "Reservations[*].Instances[*].[Tags[?Key=='Name'].Value]" --output text | sed -e "s/^/ec2:/"

commands:
    # Command to match
  - name: /usr/bin/ssh
    # Wrapper name when running `dargs generate-bins`
    wrapper: ssh
    # Use fzf to fuzzy complete (must have fzf installed)
    #fzf-complete: true
    # List of completers/transformers to associate with this command by default
    completers:
      - ec2_name
    # Transformers are run in order which allows them to be chained
    transformers:
      - ec2_id
```
