// Copyright 2018 Arnaud Poret
// This work is licensed under the BSD 2-Clause License.
package main
import (
    "encoding/csv"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)
func main() {
    var (
        maxStep,maxWalk int64
        sources []string
        allPaths [][]string
        preds map[string][]string
        edges map[string]map[string]string
    )
    if (len(os.Args)==2) && (os.Args[1]=="help") {
        fmt.Println(strings.Join([]string{
            "",
            "stream is a tool for network traversal starting from source nodes along the up or down stream (the down stream is not yet implemented).",
            "",
            "Typical usage consists in extracting, from the network, the up/down-stream paths of a couple of nodes.",
            "",
            "stream handles networks encoded in the sif file format.",
            "",
            "stream does not handle multi-graphs (i.e. networks where nodes can be connected by more than one edge).",
            "",
            "Usage: stream networkFile sourceFile maxStep maxWalk",
            "",
            "    * networkFile: the network encoded in a sif file",
            "",
            "    * sourceFile: the source nodes listed in a file (one node per line)",
            "",
            "    * maxStep (>0): the maximum number of steps performed during a random walk when up/down streaming from a source node",
            "",
            "    * maxWalk (>0): the maximum number of random walks performed in the network when up/down streaming from a source node",
            "",
            "The returned file is a sif file encoding a subnetwork (of the provided network) containing the up/down-stream paths of the source nodes.",
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
            "stream networkFile sourceFile maxStep maxWalk",
            "",
        },"\n"))
    } else if len(os.Args)==5 {
        maxStep,_=strconv.ParseInt(os.Args[3],10,0)
        maxWalk,_=strconv.ParseInt(os.Args[4],10,0)
        if !strings.HasSuffix(os.Args[1],".sif") {
            fmt.Println("ERROR: "+os.Args[1]+" must have the .sif file extension")
        } else if !strings.HasSuffix(os.Args[2],".txt") {
            fmt.Println("ERROR: "+os.Args[2]+" must have the .txt file extension")
        } else if int(maxStep)<1 {
            fmt.Println("ERROR: maxStep must 1 or more")
        } else if int(maxWalk)<1 {
            fmt.Println("ERROR: maxWalk must 1 or more")
        } else {
            preds,edges=ReadNetwork(os.Args[1])
            if len(edges)==0 {
                fmt.Println("ERROR: "+os.Args[1]+" is empty after reading")
            } else {
                sources=ReadNodes(os.Args[2],preds)
                if len(sources)==0 {
                    fmt.Println("ERROR: "+os.Args[2]+" is empty after reading")
                } else {
                    rand.Seed(int64(time.Now().Nanosecond()))
                    allPaths=FindAllPaths(sources,int(maxStep),int(maxWalk),preds)
                    if len(allPaths)==0 {
                        fmt.Println("WARNING: no paths found")
                    } else {
                        WriteNetwork(strings.Replace(os.Args[2],".txt",".sif",-1),allPaths,edges)
                    }
                }
            }
        }
    } else {
        fmt.Println(strings.Join([]string{
            "ERROR: wrong number of arguments",
            "",
            "To print help:    stream help",
            "To print usage:   stream usage",
            "To print license: stream license",
            "",
            "For more information see https://github.com/arnaudporet/stream",
            "",
        },"\n"))
    }
}
func CopyPath(path []string) []string {
    var y []string
    y=make([]string,len(path))
    copy(y,path)
    return y
}
func FindAllPaths(sources []string,maxStep,maxWalk int,preds map[string][]string) [][]string {
    var (
        i int
        tail string
        path []string
        paths,allPaths [][]string
    )
    tail="/"+strconv.FormatInt(int64(len(sources)),10)+")"
    for i=range sources {
        fmt.Println("streaming "+sources[i]+" ("+strconv.FormatInt(int64(i+1),10)+tail)
        paths=FindPaths(sources[i],maxStep,maxWalk,preds)
        for _,path=range paths {
            if !IsInPaths(allPaths,path) {
                allPaths=append(allPaths,CopyPath(path))
            }
        }
    }
    return allPaths
}
func FindPaths(source string,maxStep,maxWalk int,preds map[string][]string) [][]string {
    var (
        i int
        path []string
        paths [][]string
    )
    for i=0;i<maxWalk;i++ {
        path=RandomWalk(source,maxStep,preds)
        if (len(path)!=0) && !IsInPaths(paths,path) {
            paths=append(paths,CopyPath(path))
        }
    }
    return paths
}
func IsInPath(path []string,thatNode string) bool {
    var node string
    for _,node=range path {
        if node==thatNode {
            return true
        }
    }
    return false
}
func IsInPaths(paths [][]string,thatPath []string) bool {
    var path []string
    for _,path=range paths {
        if PathEq(path,thatPath) {
            return true
        }
    }
    return false
}
func IsInPreds(preds map[string][]string,thatNode string) bool {
    var node string
    for node=range preds {
        if node==thatNode {
            return true
        }
    }
    return false
}
func PathEq(path1,path2 []string) bool {
    var i int
    if len(path1)!=len(path2) {
        return false
    } else {
        for i=range path1 {
            if path1[i]!=path2[i] {
                return false
            }
        }
        return true
    }
}
func RandomWalk(source string,maxStep int,preds map[string][]string) []string {
    var (
        i int
        current,next string
        path []string
    )
    if len(preds[source])!=0 {
        current=source
        path=append(path,source)
        for i=0;i<maxStep;i++ {
            next=preds[current][rand.Intn(len(preds[current]))]
            if (len(preds[next])==0) || IsInPath(path,next) {
                path=append(path,next)
                return path
            } else {
                path=append(path,next)
                current=next
            }
        }
    }
    return []string{}
}
func ReadNetwork(networkFile string) (map[string][]string,map[string]map[string]string) {
    var (
        err error
        node string
        line []string
        lines [][]string
        preds map[string][]string
        edges map[string]map[string]string
        reader *csv.Reader
        file *os.File
    )
    fmt.Println("reading "+networkFile)
    file,_=os.Open(networkFile)
    reader=csv.NewReader(file)
    reader.Comma='\t'
    reader.Comment=0
    reader.FieldsPerRecord=3
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    reader.ReuseRecord=true
    lines,err=reader.ReadAll()
    file.Close()
    preds=make(map[string][]string)
    edges=make(map[string]map[string]string)
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        for _,line=range lines {
            for _,node=range []string{line[0],line[2]} {
                preds[node]=[]string{}
            }
            edges[line[0]]=make(map[string]string)
        }
        for _,line=range lines {
            if IsInPath(preds[line[2]],line[0]) {
                fmt.Println("ERROR: contains multi-edges")
                preds=make(map[string][]string)
                edges=make(map[string]map[string]string)
                break
            } else {
                preds[line[2]]=append(preds[line[2]],line[0])
                edges[line[0]][line[2]]=line[1]
            }
        }
    }
    return preds,edges
}
func ReadNodes(nodeFile string,preds map[string][]string) []string {
    var (
        err error
        nodes,line []string
        lines [][]string
        reader *csv.Reader
        file *os.File
    )
    fmt.Println("reading "+nodeFile)
    file,_=os.Open(nodeFile)
    reader=csv.NewReader(file)
    reader.Comma='\t'
    reader.Comment=0
    reader.FieldsPerRecord=1
    reader.LazyQuotes=false
    reader.TrimLeadingSpace=true
    reader.ReuseRecord=true
    lines,err=reader.ReadAll()
    file.Close()
    if err!=nil {
        fmt.Println("ERROR: "+err.Error())
    } else {
        for _,line=range lines {
            if !IsInPreds(preds,line[0]) {
                fmt.Println("WARNING: "+line[0]+" not in network")
            } else if !IsInPath(nodes,line[0]) {
                nodes=append(nodes,line[0])
            }
        }
    }
    return nodes
}
func WriteNetwork(networkFile string,paths [][]string,edges map[string]map[string]string) {
    var (
        i int
        path,line []string
        lines [][]string
        writer *csv.Writer
        file *os.File
    )
    fmt.Println("writing "+networkFile)
    for _,path=range paths {
        for i=0;i<len(path)-1;i++ {
            line=[]string{path[i+1],edges[path[i+1]][path[i]],path[i]}
            if !IsInPaths(lines,line) {
                lines=append(lines,CopyPath(line))
            }
        }
    }
    file,_=os.Create(networkFile)
    writer=csv.NewWriter(file)
    writer.Comma='\t'
    writer.UseCRLF=false
    writer.WriteAll(lines)
    file.Close()
}
