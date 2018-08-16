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
        ok bool
        wardFile string
        nodes,roots []string
        edges,ward [][]string
        edgeNames map[string]map[string]string
    )
    if len(os.Args)==4 {
        ok=true
        if !strings.HasSuffix(os.Args[1],".sif") {
            fmt.Println("ERROR: "+os.Args[1]+" must have the \".sif\" file extension")
            ok=false
        }
        if !strings.HasSuffix(os.Args[2],".txt") {
            fmt.Println("ERROR: "+os.Args[2]+" must have the \".txt\" file extension")
            ok=false
        }
        if (os.Args[3]!="up") && (os.Args[3]!="down") {
            fmt.Println("ERROR: direction must be \"up\" or \"down\"")
            ok=false
        }
        if ok {
            fmt.Println("INFO: reading "+os.Args[1])
            nodes,edges,edgeNames=ReadNetwork(os.Args[1])
            if len(edges)==0 {
                fmt.Println("WARNING: "+os.Args[1]+" is empty after reading")
            } else {
                fmt.Println("INFO: reading "+os.Args[2])
                roots=ReadNodes(os.Args[2],nodes)
                if len(roots)==0 {
                    fmt.Println("WARNING: "+os.Args[2]+" is empty after reading")
                } else {
                    fmt.Println("INFO: "+os.Args[3]+"streaming "+os.Args[2])
                    if os.Args[3]=="down" {
                        ward=ForwardEdges(roots,edges)
                    } else if os.Args[3]=="up" {
                        ward=BackwardEdges(roots,edges)
                    }
                    if len(ward)==0 {
                        fmt.Println("WARNING: no "+os.Args[3]+"stream paths found")
                    } else {
                        wardFile=strings.Replace(os.Args[2],".txt","-"+os.Args[3]+"stream.sif",-1)
                        fmt.Println("INFO: writing "+wardFile)
                        WriteNetwork(wardFile,ward,edgeNames)
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
            "Copyright 2018 Arnaud Poret",
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
func BackwardEdges(roots []string,edges [][]string) [][]string {
    var (
        root,npred string
        edge,epred []string
        toCheck,newCheck,backward [][]string
        nodePred map[string][]string
        edgePred map[string]map[string][][]string
    )
    nodePred,edgePred=GetPredecessors(edges)
    for _,root=range roots {
        for _,npred=range nodePred[root] {
            backward=append(backward,[]string{npred,root})
            newCheck=append(newCheck,[]string{npred,root})
        }
    }
    for {
        toCheck=CopyList2(newCheck)
        newCheck=[][]string{}
        for _,edge=range toCheck {
            for _,epred=range edgePred[edge[0]][edge[1]] {
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
func ForwardEdges(roots []string,edges [][]string) [][]string {
    var (
        root,nsucc string
        edge,esucc []string
        toCheck,newCheck,forward [][]string
        nodeSucc map[string][]string
        edgeSucc map[string]map[string][][]string
    )
    nodeSucc,edgeSucc=GetSuccessors(edges)
    for _,root=range roots {
        for _,nsucc=range nodeSucc[root] {
            forward=append(forward,[]string{root,nsucc})
            newCheck=append(newCheck,[]string{root,nsucc})
        }
    }
    for {
        toCheck=CopyList2(newCheck)
        newCheck=[][]string{}
        for _,edge=range toCheck {
            for _,esucc=range edgeSucc[edge[0]][edge[1]] {
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
    return forward
}
func GetPredecessors(edges [][]string) (map[string][]string,map[string]map[string][][]string) {
    var (
        node,node2,node3 string
        edge []string
        nodePred map[string][]string
        edgePred map[string]map[string][][]string
    )
    nodePred=make(map[string][]string)
    edgePred=make(map[string]map[string][][]string)
    for _,edge=range edges {
        for _,node=range edge {
            nodePred[node]=[]string{}
        }
        edgePred[edge[0]]=make(map[string][][]string)
    }
    for _,edge=range edges {
        nodePred[edge[1]]=append(nodePred[edge[1]],edge[0])
        edgePred[edge[0]][edge[1]]=[][]string{}
    }
    for node=range nodePred {
        for _,node2=range nodePred[node] {
            for _,node3=range nodePred[node2] {
                edgePred[node2][node]=append(edgePred[node2][node],[]string{node3,node2})
            }
        }
    }
    return nodePred,edgePred
}
func GetSuccessors(edges [][]string) (map[string][]string,map[string]map[string][][]string) {
    var (
        node,node2,node3 string
        edge []string
        nodeSucc map[string][]string
        edgeSucc map[string]map[string][][]string
    )
    nodeSucc=make(map[string][]string)
    edgeSucc=make(map[string]map[string][][]string)
    for _,edge=range edges {
        for _,node=range edge {
            nodeSucc[node]=[]string{}
        }
        edgeSucc[edge[0]]=make(map[string][][]string)
    }
    for _,edge=range edges {
        nodeSucc[edge[0]]=append(nodeSucc[edge[0]],edge[1])
        edgeSucc[edge[0]][edge[1]]=[][]string{}
    }
    for node=range nodeSucc {
        for _,node2=range nodeSucc[node] {
            for _,node3=range nodeSucc[node2] {
                edgeSucc[node][node2]=append(edgeSucc[node][node2],[]string{node2,node3})
            }
        }
    }
    return nodeSucc,edgeSucc
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
func ReadNetwork(networkFile string) ([]string,[][]string,map[string]map[string]string) {
    var (
        err error
        node string
        nodes,edge,line []string
        edges,lines [][]string
        edgeNames map[string]map[string]string
        file *os.File
        reader *csv.Reader
    )
    file,err=os.Open(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+networkFile+" "+err.Error())
        return []string{},[][]string{},map[string]map[string]string{}
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
            fmt.Println("ERROR: "+networkFile+" "+err.Error())
            return []string{},[][]string{},map[string]map[string]string{}
        } else {
            edgeNames=make(map[string]map[string]string)
            for _,line=range lines {
                edge=[]string{line[0],line[2]}
                if IsInList2(edges,edge) {
                    fmt.Println("ERROR: "+networkFile+" contains multi-edges (or duplicated edges)")
                    return []string{},[][]string{},map[string]map[string]string{}
                } else {
                    edges=append(edges,CopyList(edge))
                    for _,node=range edge {
                        if !IsInList(nodes,node) {
                            nodes=append(nodes,node)
                        }
                    }
                    edgeNames[line[0]]=make(map[string]string)
                }
            }
            for _,line=range lines {
                edgeNames[line[0]][line[2]]=line[1]
            }
        }
    }
    return nodes,edges,edgeNames
}
func ReadNodes(nodeFile string,networkNodes []string) []string {
    var (
        err error
        line,nodes []string
        lines [][]string
        file *os.File
        reader *csv.Reader
    )
    file,err=os.Open(nodeFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+nodeFile+" "+err.Error())
        return []string{}
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
            fmt.Println("ERROR: "+nodeFile+" "+err.Error())
            return []string{}
        } else {
            for _,line=range lines {
                if !IsInList(networkNodes,line[0]) {
                    fmt.Println("WARNING: "+nodeFile+"/"+line[0]+" not in network")
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
    file,err=os.Create(networkFile)
    defer file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+networkFile+" "+err.Error())
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
