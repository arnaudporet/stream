# Streaming over networks

Copyright 2018 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [BSD 2-Clause License](https://raw.githubusercontent.com/arnaudporet/stream/master/license.txt).

## stream

[stream](https://github.com/arnaudporet/stream) is a tool implemented in [Go](https://golang.org) for network traversal starting from source nodes along the up or down stream (the down stream is not yet implemented).

To do so, in the network, stream performs random walks starting from the nodes given as sources.

Typical usage consists in extracting, from the network, the up/down-stream paths of a couple of nodes.

stream handles networks encoded in the sif file format (see below): the provided network must be encoded in the sif file format.

Together with the network encoded in a sif file, stream requires the source nodes to be listed in a file (see below).

## The sif file format

In a sif file encoding a network, each line encodes an edge of the network as follows:
* `source \t interaction \t target`

Note that the field separator is the tabulation `\t`: the sif file format is the tab-separated values format (tsv) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a sif file encoded as follows:
* `HRAS \t activation \t RAF1`

## Usage

In a terminal emulator:
1. `go build stream.go`
2. `./stream networkFile sourceFile maxStep maxWalk`

or simply
* `go run stream.go networkFile sourceFile maxStep maxWalk`

Note that `go run` builds stream each time before running it.

The Go package can have different names depending on your operating system. For example, with [Ubuntu](https://www.ubuntu.com), the Go package is named golang. Consequently, running a Go file with Ubuntu might be `golang-go run yourfile.go` instead of `go run yourfile.go` with [Arch Linux](https://www.archlinux.org).

Arguments:
* `networkFile`: the network encoded in a sif file (see above)
* `sourceFile`: the source nodes listed in a file (one node per line)
* `maxStep` (`>0`): the maximum number of steps performed during a random walk when up/down streaming from a source node
* `maxWalk` (`>0`): the maximum number of random walks performed in the network when up/down streaming from a source node

The returned file is a sif file encoding a subnetwork (of the provided network) containing the up/down-stream paths of the source nodes.

## Cautions

* stream does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge)
* the network must be provided as a sif file (see above)
* in the file containing the source nodes (see above): one node per line
* since stream uses random walks:
    * the results can be subject to variability
    * returning all the possible up/down-stream paths is not guaranteed
* increasing `maxWalk`:
    * increases the robustness of the results
    * but also increases the calculation time

## Examples

All the used example networks are adapted from pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html).

* example 1: ErbB signaling pathway
    * `stream ErbB_signaling_pathway.sif sources.txt 1000 1000000`
    * networkFile: the ErbB signaling pathway (138 edges)
    * sourceFile: contains the nodes JUN and MYC
    * maxStep: 1000
    * maxWalk: 1000000
    * result: sources.sif (133 edges), also in svg for visualization

The ErbB signaling pathway is a growth-promoting signaling pathway typically activated by the epidermal growth factor (EGF).

JUN and MYC are two transcription factors influencing the expression of target genes following EGF stimulation.

The resulting file `sources.sif` converted in svg shows the upstream paths (i.e. the regulating paths) of JUN (red) and MYC (green) in the ErbB signaling pathway. It highlights that JUN and MYC share common elements in there regulating paths (red and green) and also specific elements (red or green). Note that other regulating paths outside of the ErbB signaling pathway exists.

## Forthcoming

* implementing the down-stream traversal
* adding more complex examples

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:
* go (Arch Linux)
* golang (Ubuntu)

Otherwise see https://golang.org/doc/install
