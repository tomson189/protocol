#!/bin/bash

#
# Test creating a single send transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TESTS=$GOPATH/src/github.com/Oneledger/protocol/node/tests

$CMD/resetOneLedger
$TESTS/register.sh

addrAlice=`$OLSCRIPT/lookup Alice RPCAddress tcp://127.0.0.1:`
addrBob=`$OLSCRIPT/lookup Bob RPCAddress tcp://127.0.0.1:`
addrCarol=`$OLSCRIPT/lookup Carol RPCAddress tcp://127.0.0.1:`
addrDavid=`$OLSCRIPT/lookup David RPCAddress tcp://127.0.0.1:`

# Put some money in the user accounts
SEQ=`$CMD/nextSeq`
olclient testmint -s $SEQ -a $addrAlice --party Alice --amount 123456 --currency OLT

sleep 8

# assumes fullnode is in the PATH
olclient send -s $SEQ -a $addrAlice --party Alice --counterparty Bob --amount 10 --currency OLT

sleep 8

olclient identity -a $addrAlice
olclient account -a $addrAlice
olclient account -a $addrBob

$CMD/stopOneLedger

