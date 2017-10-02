#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


info() { echo -e "$*"; }
fail() { echo -e "FAIL $*"; exit 1; }

run_test() {
    got=$("${cmd[@]}" || true)
    [[ "$got" == "$want" ]] || fail "cmd: ${cmd[*]}\ngot: $got\nwant: $want"
    info "pass\n"
}

dargs_bin=$1

cp "$dargs_bin" "$SCRIPT_DIR/dargs"

dargs="$SCRIPT_DIR/dargs"

rm -rf "$SCRIPT_DIR/.dargs"

info 'Runs transfomer'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -t test -- echo "test:hi")
want='{hi}'
run_test

info 'Runs command transformers'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -- echo "test:hi")
want='_hi_'
run_test

info 'Ignores non-matching'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -- echo "foobar")
want='foobar'
run_test

info 'Match honors prev-match'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -t prev -- echo 'prev' "foobar")
want='prev true'
run_test

info 'Match ignores miss on prev-match'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -t prev -- echo "foobar")
want='foobar'
run_test

info 'Multiple argument expansion'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -n -t multi -- echo "t")
want='echo arg1 arg2'
run_test

info 'Uses cache'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" run -t cache -- echo "cache")
"${cmd[@]}" > /dev/null
# Should not take longer than .1s
timeout 0.1s "${cmd[@]}" > /dev/null || fail 'Timed out, did not use cache'
info "pass\n"

info 'Completion'
cmd=("$dargs" --config "$SCRIPT_DIR/dargs.yml" completions -c test -- echo "echo" "cmp:foobar")
want='cmp:foobar1
cmp:foobar2'
run_test
