# Streaming over networks

Copyright 2018 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [BSD 2-Clause License](https://raw.githubusercontent.com/arnaudporet/stream/master/license.txt).

## stream

[stream](https://github.com/arnaudporet/stream) is a tool implemented in [Go](https://golang.org) for finding the upstream/downstream paths starting from a couple of root nodes in a network.

stream handles networks encoded in the SIF file format (see below): the provided network must be encoded in the SIF file format.

Together with the network encoded in a SIF file, stream requires the root nodes to be listed in a file (see below).

## The SIF file format

In a SIF file encoding a network, each line encodes an edge of the network as follows:
* `source \t interaction \t target`

Note that the field separator is the tabulation `\t`: the SIF file format is the tab-separated values format (TSV) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a SIF file encoded as follows:
* `HRAS \t activation \t RAF1`

## Usage

In a terminal emulator:
1. `go build stream.go`
2. `./stream networkFile rootFile direction`

or simply
* `go run stream.go networkFile rootFile direction`

Note that `go run` builds stream each time before running it.

The Go package can have different names depending on your operating system. For example, with [Ubuntu](https://www.ubuntu.com), the Go package is named golang. Consequently, running a Go file with Ubuntu might be `golang-go run yourfile.go` instead of `go run yourfile.go` with [Arch Linux](https://www.archlinux.org).

Arguments:
* `networkFile`: the network encoded in a SIF file (see above)
* `rootFile`: the root nodes listed in a file (one node per line)
* `direction`: follows the up stream ("up") or the down stream ("down")

The returned file is a SIF file encoding the upstream/downstream paths of the root nodes in the network.

## Cautions

* stream does not handle multi-edges (i.e. two or more edges having the same source node and the same target node)
* note that duplicated edges are multi-edges
* the network must be provided as a SIF file (see above)
* in the file containing the root nodes (see above): one node per line

## Examples

All the networks used in these examples are adapted from pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html).

* ErbB signaling pathway
    * `stream ErbB_signaling_pathway.sif roots.txt up`
    * networkFile: the ErbB signaling pathway (239 edges)
    * rootFile: contains the nodes JUN and MYC
    * direction: up
    * result: roots-upstream.sif (133 edges), also in svg for visualization

The ErbB signaling pathway is a growth-promoting signaling pathway typically activated by the epidermal growth factor (EGF).

JUN and MYC are two transcription factors influencing the expression of target genes following EGF stimulation.

The resulting file `roots-upstream.sif` converted to svg shows the upstream paths (i.e. the regulating paths) of JUN (red) and MYC (green) in the ErbB signaling pathway. It highlights that JUN and MYC share common elements in there regulating paths (red and green) and also specific elements (red or green). Note that other regulating paths outside of the ErbB signaling pathway exist.

* Toll-like receptor signaling pathway
    * `stream Toll-like_receptor_signaling_pathway.sif roots.txt down`
    * networkFile: the Toll-like receptor signaling pathway (219 edges)
    * rootFile: contains the nodes TLR3 and TLR4
    * direction: down
    * result: roots-downstream.sif (152 edges), also in svg for visualization

The Toll-like receptors (TLRs) are cell surface receptors which can be activated by various pathogen associated molecular patterns (PMAPs).

PMAPs are molecules coming from microorganisms and the TLRs can detect them in order to signal the presence of potential infectious pathogens. There are several types of TLRs, each being able to detect specific PMAPs.

TLR3 can detect the presence of double-stranded RNA (dsRNA) coming from RNA viruses whereas TLR4 can detect lipopolysaccharide (LPS) coming from Gram-negative bacteria.

The resulting file `roots-downstream.sif` converted to svg shows the downstream paths (i.e. the effector paths) of TLR3 (red) and TLR4 (green). It highlights that TLR3 and TLR4 share common effectors (red and green) and also specific ones (red or green).

## Forthcoming

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:
* go (Arch Linux)
* golang (Ubuntu)

Otherwise see https://golang.org/doc/install
