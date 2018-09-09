package rubberducky

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethpm/ethpm-go/pkg/bytecode"
	"github.com/ethpm/ethpm-go/pkg/ethcontract"
	"github.com/ethpm/ethpm-go/pkg/ethpm"
)

const TRUEENV = "true"

type OriginStruct struct {
	Url string `toml:"url"`
}

type GitConfig struct {
	Remote OriginStruct `toml:"remote origin"`
}

type PackageInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Author      map[string]string `json:"author"`
	License     string            `json:"license"`
	Homepage    string            `json:"homepage"`
	Description string            `json:"description"`
}

type ArtifactInfo struct {
	ABI              []*ethcontract.ABIObject      `json:"abi"`
	Address          string                        `json:"address"`
	Bytecode         string                        `json:"bytecode"`
	DeployedBytecode string                        `json:"deployedBytecode"`
	ContractName     string                        `json:"contractName"`
	SourcePath       string                        `json:"sourcePath"`
	Compiler         *bytecode.CompilerInformation `json:"compiler"`
}

// BuildAndWritePackage TODO some comments
func BuildAndWritePackage(directoryname string) error {
	json, err := buildPackage(directoryname)
	if err != nil {
		return err
	}

	err = writePackage(directoryname, json)
	return err
}

func buildPackage(directoryname string) (string, error) {
	var dir string
	var err error
	var pmstring string
	var artifactpath string
	var nodepackagepath string

	pm := ethpm.PackageManifest{}
	contracttype := ethcontract.ContractType{}
	contracttypes := make(map[string]*ethcontract.ContractType)
	dbc := bytecode.UnlinkedBytecode{}
	rbc := bytecode.UnlinkedBytecode{}

	if os.Getenv("TEST") != TRUEENV {
		dir, _ = os.Getwd()
		if err != nil {
			err = fmt.Errorf("Could not access working directory: '%v'", err)
			return "", err
		}
		artifactpath = dir + "/" + directoryname + "/build/contracts/"
	} else {
		artifactpath = "../../../test/testdata/build/contracts/"
	}
	artifactpath = filepath.FromSlash(artifactpath)
	files, err := ioutil.ReadDir(artifactpath)
	if err != nil {
		err = fmt.Errorf("Could not access directory: '%v'", err)
		return "", err
	}

	var file *os.File
	var info os.FileInfo
	artifactObjects := []ArtifactInfo{}
	for i, f := range files {
		name := f.Name()

		file, err = os.Open(artifactpath + name)
		if err != nil {
			err = fmt.Errorf("First, tell Bryant I said hi, and second, this could not open file '%v': '%v'", artifactpath, err)
			return "", err
		}
		info, _ = file.Stat()
		artifact := make([]byte, info.Size())
		_, err = file.Read(artifact)
		if err != nil {
			err = fmt.Errorf("Could not read file '%v': '%v'", artifactpath, err)
			return "", err
		}
		tempao := ArtifactInfo{}
		err = json.Unmarshal(artifact, &tempao)
		if err != nil {
			err = fmt.Errorf("Could not unpackage truffle artifact: '%v'", err)
			return "", err
		}
		contracttype.ABI = tempao.ABI
		compiler := tempao.Compiler
		dbc.Bytecode = tempao.DeployedBytecode
		rbc.Bytecode = tempao.Bytecode

		contracttype.Compiler = compiler
		contracttype.RuntimeBytecode = &dbc
		contracttype.DeploymentBytecode = &rbc
		artifactname := name[:len(name)-5]
		contracttypes[artifactname] = &contracttype
		artifactObjects = append(artifactObjects, tempao)
		fmt.Printf("%v+\n", artifactObjects[i])
	}

	if os.Getenv("TEST") != TRUEENV {
		nodepackagepath = dir + "/" + directoryname + "/package.json"
	} else {
		nodepackagepath = "../../../test/testdata/package.json"
	}
	file, err = os.Open(filepath.FromSlash(nodepackagepath))
	if err != nil {
		err = fmt.Errorf("Could not open file '%v': '%v'", nodepackagepath, err)
		return "", err
	}
	info, _ = file.Stat()
	pkg := make([]byte, info.Size())
	_, err = file.Read(pkg)
	if err != nil {
		err = fmt.Errorf("Could not read file '%v': '%v'", nodepackagepath, err)
		return "", err
	}
	pkgObject := PackageInfo{}
	err = json.Unmarshal(pkg, &pkgObject)
	if err != nil {
		err = fmt.Errorf("Could not unpackage package.json: '%v'", err)
		return "", err
	}

	fmt.Printf("%v+\n", pkgObject)

	author := make([]string, 1)
	author[0] = pkgObject.Author["name"]
	theMeta := &ethpm.PackageMeta{}

	pm.PackageName = pkgObject.Name
	pm.Version = pkgObject.Version
	theMeta.Authors = author
	theMeta.License = pkgObject.License
	theMeta.Description = pkgObject.Description
	links := make(map[string]string)
	links["documentation"] = pkgObject.Homepage
	theMeta.Links = links
	pm.Meta = theMeta
	sources := make(map[string]string)
	var commithashpath string
	if os.Getenv("TEST") != "true" {
		commithashpath = dir + "/" + directoryname + "/.git/refs/heads/master"
	} else {
		commithashpath = "../../../test/testdata/master"
	}
	file, err = os.Open(filepath.FromSlash(commithashpath))
	if err != nil {
		err = fmt.Errorf("Could not open file '%v': '%v'", commithashpath, err)
		return "", err
	}
	info, _ = file.Stat()
	commitbytes := make([]byte, info.Size())
	_, err = file.Read(commitbytes)
	if err != nil {
		err = fmt.Errorf("Could not read file '%v': '%v'", commithashpath, err)
		return "", err
	}
	commitbytes = commitbytes[:len(commitbytes)-1]
	commit := string(commitbytes)

	var contractpath string
	if os.Getenv("TEST") != "true" {
		contractpath = dir + "/" + directoryname + "/contracts/"
	} else {
		contractpath = "../../../test/testdata/contracts/"
	}

	files, err = ioutil.ReadDir(contractpath)
	if err != nil {
		err = fmt.Errorf("Could not access directory: '%v'", err)
		return "", err
	}

	for _, f := range files {
		name := f.Name()
		sources["./contracts/"+name] = pkgObject.Homepage + ".git@" + commit
	}
	pm.Sources = sources

	pm.ContractTypes = contracttypes
	pmstring, err = pm.Write()
	if err != nil {
		err = fmt.Errorf("Could not write object: '%v'", err)
		fmt.Println(err)
		return "", err
	}
	return pmstring, nil
}

func writePackage(directoryname string, json string) error {
	var p *os.File
	var err error
	if os.Getenv("TEST") != "true" {
		p, err = os.Create(directoryname + "/ethpm.json")
	} else {
		p, err = os.Create("../../../test/testdata/ethpm.json")
	}
	if err != nil {
		err = fmt.Errorf("Could create file ethpm.json: '%v'", err)
		return err
	}
	defer p.Close()
	_, err = p.WriteString(json)
	if err != nil {
		err = fmt.Errorf("Could not write ethpm.json: '%v'", err)
		return err
	}
	p.Sync()
	return err
}
