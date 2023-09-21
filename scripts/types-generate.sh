#!/usr/bin/env bash

# get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
# get up level directory
TOP_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"

GENERATOR=$TOP_DIR/tools/gen-ast-types.go

# delete generated contents
awk '{print} /\/\/ 下面的代码是自动生成的/ {exit}' $GENERATOR > $GENERATOR.tmp
mv $GENERATOR.tmp $GENERATOR

# exec scripts/pkgreflect.go
go run $SCRIPT_DIR/pkgreflect.go -nofuncs -novars -norecurs -noconsts -stdout $TOP_DIR/pkg/semantic  >> $GENERATOR

# exec tools/gen-ast-types.go
go run $GENERATOR -dir $TOP_DIR/pkg/semantic/internal

# exec go generate
go generate $TOP_DIR/pkg/semantic/ast.go

