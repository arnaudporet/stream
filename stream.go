// Copyright 2018 Arnaud Poret
// This work is licensed under the BSD 2-Clause License.
package main
import (
    "encoding/csv"
    "fmt"
    "os"
    "strings"
)
func main() {
    var (
        roots []string
        ward [][]string
        nodeSP map[string][]string
        edgeNames map[string]map[string]string
        edgeSP map[string]map[string][][]string
    )
    if len(os.Args)==4 {
        if !strings.HasSuffix(os.Args[1],".sif") {
            fmt.Println("ERROR: "+os.Args[1]+" must have the .sif file extension")
        } else if !strings.HasSuffix(os.Args[2],".txt") {
            fmt.Println("ERROR: "+os.Args[2]+" must have the .txt file extension")
        } else if (os.Args[3]!="up") && (os.Args[3]!="down") {
            fmt.Println("ERROR: direction must be \"up\" or \"down\"")
        } else {
            nodeSP,edgeSP,edgeNames=ReadNetwork(os.Args[1],os.Args[3])
            if len(edgeNames)==0 {
                fmt.Println("WARNING: "+os.Args[1]+" is empty after reading")
            } else {
                roots=ReadNodes(os.Args[2],nodeSP)
                if len(roots)==0 {
                    fmt.Println("WARNING: "+os.Args[2]+" is empty after reading")
                } else {
                    if os.Args[3]=="down" {
                        ward=ForwardEdges(roots,nodeSP,edgeSP)
                    } else if os.Args[3]=="up" {
                        ward=BackwardEdges(roots,nodeSP,edgeSP)
                    }
                    if len(ward)==0 {
                        fmt.Println("WARNING: no "+os.Args[3]+"stream paths found")
                    } else {
                        WriteNetwork(strings.Replace(os.Args[2],".txt",".sif",-1),ward,edgeNames)
                    }
                }
            }
        }
    } else if (len(os.Args)==2) && (os.Args[1]=="help") {
        fmt.Println(strings.Join([]string{
            "",
            "stream is a tool for finding the upstream/downstream paths starting from a couple of root nodes in a network.",
            "",
            "stream handles networks encoded in the SIF file format.",
            "",
            "stream does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge).",
            "",
            "Note that if a network contains duplicated edges then it is a multi-graph.",
            "",
            "Usage: stream networkFile rootFile direction",
            "",
            "    * networkFile: the network encoded in a SIF file",
            "",
            "    * rootFile:    the root nodes listed in a file (one node per line)",
            "",
            "    * direction:   follows the up stream (\"up\") or the down stream (\"down\")",
            "",
            "The returned file is a SIF file encoding the upstream/downstream paths of the root nodes in the network.",
            "",
            "For more information see https://github.com/arnaudporet/stream",
            "",
        },"\n"))
    } else if (len(os.Args)==2) && (os.Args[1]=="license") {
        fmt.Println(strings.Join([]string{
            "",
            "Copyright 2017-2018 Arnaud Poret",
            "",
            "Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:",
            "",
            "1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.",
            "",
            "2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.",
            "",
            "THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
            "",
        },"\n"))
    } else if (len(os.Args)==2) && (os.Args[1]=="usage") {
        fmt.Println(strings.Join([]string{
            "",
            "stream networkFile rootFile direction",
            "",
        },"\n"))
    } else {
        fmt.Println(strings.Join([]string{
            "ERROR: wrong number of arguments",
            "",
            "To print help:    stream help",
            "To print license: stream license",
            "To print usage:   stream usage",
            "",
            "For more information see https://github.com/arnaudporet/stream",
            "",
        },"\n"))
    }
}
func BackwardEdges(nroots []string,nodePred map[string][]string,edgePred map[string]map[string][][]string) [][]string {
    var (
        nroot,npred string
        eroot,check,epred []string
        backward,newCheck,toCheck [][]string
    )
    for _,nroot=range nroots {
        fmt.Println("backwarding "+nroot)
        for _,npred=range nodePred[nroot] {
            eroot=[]string{npred,nroot}
            if !IsInList2(backward,eroot) {
                backward=append(backward,CopyList(eroot))
                newCheck=[][]string{CopyList(eroot)}
                for {
                    toCheck=CopyList2(newCheck)
                    newCheck=[][]string{}
                    for _,check=range toCheck {
                        for _,epred=range edgePred[check[0]][check[1]] {
                            if !IsInList2(backward,epred) {
                                backward=append(backward,CopyList(epred))
                                newCheck=append(newCheck,CopyList(epred))
                            }
                        }
                    }
                    if len(newCheck)==0 {
                        break
                    }
                }
            }
        }
    }
    return backward
}
func CopyList(list []string) []string {
    var y []string
    y=make([]string,len(list))
    copy(y,list)
    return y
}
func CopyList2(list2 [][]string) [][]string {
    var (
        i int
        y [][]string
    )
    y=make([][]string,len(list2))
    for i=range list2 {
        y[i]=make([]string,len(list2[i]))
        copy(y[i],list2[i])
    }
    return y
}
func ForwardEdges(nroots []string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string) [][]string {
    var (
        nroot,nsucc string
        eroot,check,esucc []string
        forward,newCheck,toCheck [][]string
    )
    for _,nroot=range nroots {
        fmt.Println("forwarding "+nroot)
        for _,nsucc=range nodeSucc[nroot] {
            eroot=[]string{nroot,nsucc}
            if !IsInList2(forward,eroot) {
                forward=append(forward,CopyList(eroot))
                newCheck=[][]string{CopyList(eroot)}
                for {
                    toCheck=CopyList2(newCheck)
                    newCheck=[][]string{}
                    for _,check=range toCheck {
                        for _,esucc=range edgeSucc[check[0]][check[1]] {
                            if !IsInList2(forward,esucc) {
                                forward=append(forward,CopyList(esucc))
                                newCheck=append(newCheck,CopyList(esucc))
                            }
                        }
                    }
                    if len(newCheck)==0 {
                        break
                    }
                }
            }
        }
    }
    return forward
}
func IsInList(list []string,thatElement string) bool {
    var element string
    for _,element=range list {
        if element==thatElement {
            return true
        }
    }
    return false
}
func IsInList2(list2 [][]string,thatList []string) bool {
    var (
        found bool
        i int
        list []string
    )
    for _,list=range list2 {
        if len(list)==len(thatList) {
            found=true
            for i=range list {
                if list[i]!=thatList[i] {
                    found=false
                    break
                }
            }
            if found {
                return true
            }
        }
    }
    return false
}
func IsInNetwork(nodeSP map[string][]string,thatNode string) bool {
    var node string
    for node=range nodeSP {
        if node==thatNode {
            return true
        }
    }
    return false
}
func ReadNetwork(networkFile,direction string) (map[string][]string,map[string]map[string][][]string,map[string]map[string]string) {
    var (
        err error
        node1,node2,node3 string
        line []string
        lines [][]string
        nodeSP map[string][]string
        edgeNames map[string]map[string]string
        edgeSP map[string]map[string][][]string
        file *os.File
        reader *csv.Reader
    )
    fmt.Println("reading "+networkFile)
    file,err=os.Open(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=3
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err!=nil {
            fmt.Println("ERROR: "+err.Error())
        } else {
            nodeSP=make(map[string][]string)
            edgeSP=make(map[string]map[string][][]string)
            edgeNames=make(map[string]map[string]string)
            for _,line=range lines {
                for _,node1=range []string{line[0],line[2]} {
                    nodeSP[node1]=[]string{}
                }
                edgeSP[line[0]]=make(map[string][][]string)
                edgeNames[line[0]]=make(map[string]string)
            }
            if direction=="down" {
                for _,line=range lines {
                    if IsInList(nodeSP[line[0]],line[2]) {
                        fmt.Println("ERROR: multi-edges (or duplicated edges)")
                        nodeSP=make(map[string][]string)
                        edgeSP=make(map[string]map[string][][]string)
                        edgeNames=make(map[string]map[string]string)
                        break
                    } else {
                        nodeSP[line[0]]=append(nodeSP[line[0]],line[2])
                        edgeSP[line[0]][line[2]]=[][]string{}
                        edgeNames[line[0]][line[2]]=line[1]
                    }
                }
                for node1=range nodeSP {
                    for _,node2=range nodeSP[node1] {
                        for _,node3=range nodeSP[node2] {
                            edgeSP[node1][node2]=append(edgeSP[node1][node2],[]string{node2,node3})
                        }
                    }
                }
            } else if direction=="up" {
                for _,line=range lines {
                    if IsInList(nodeSP[line[2]],line[0]) {
                        fmt.Println("ERROR: multi-edges (or duplicated edges)")
                        nodeSP=make(map[string][]string)
                        edgeSP=make(map[string]map[string][][]string)
                        edgeNames=make(map[string]map[string]string)
                        break
                    } else {
                        nodeSP[line[2]]=append(nodeSP[line[2]],line[0])
                        edgeSP[line[0]][line[2]]=[][]string{}
                        edgeNames[line[0]][line[2]]=line[1]
                    }
                }
                for node1=range nodeSP {
                    for _,node2=range nodeSP[node1] {
                        for _,node3=range nodeSP[node2] {
                            edgeSP[node2][node1]=append(edgeSP[node2][node1],[]string{node3,node2})
                        }
                    }
                }
            }
        }
    }
    return nodeSP,edgeSP,edgeNames
}
func ReadNodes(nodeFile string,nodeSP map[string][]string) []string {
    var (
        err error
        line,nodes []string
        lines [][]string
        file *os.File
        reader *csv.Reader
    )
    fmt.Println("reading "+nodeFile)
    file,err=os.Open(nodeFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=1
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err!=nil {
            fmt.Println("ERROR: "+err.Error())
        } else {
            for _,line=range lines {
                if !IsInNetwork(nodeSP,line[0]) {
                    fmt.Println("WARNING: "+line[0]+" not in network")
                } else if !IsInList(nodes,line[0]) {
                    nodes=append(nodes,line[0])
                }
            }
        }
    }
    return nodes
}
func WriteNetwork(networkFile string,edges [][]string,edgeNames map[string]map[string]string) {
    var (
        err error
        edge []string
        lines [][]string
        file *os.File
        writer *csv.Writer
    )
    fmt.Println("writing "+networkFile)
    file,err=os.Create(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        for _,edge=range edges {
            lines=append(lines,[]string{edge[0],edgeNames[edge[0]][edge[1]],edge[1]})
        }
        writer=csv.NewWriter(file)
        writer.Comma='\t'
        writer.UseCRLF=false
        writer.WriteAll(lines)
    }
}
