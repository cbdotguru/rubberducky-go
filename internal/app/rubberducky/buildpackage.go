package rubberducky

import (
	"encoding/json"
	"fmt"
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
	Bytecode         string                        `json:"bytecode"`
	DeployedBytecode string                        `json:"deployedBytecode"`
	ContractName     string                        `json:"contractName"`
	SourcePath       string                        `json:"sourcePath"`
	Compiler         *bytecode.CompilerInformation `json:"compiler"`
}

// BuildPackage TODO some comments
func BuildPackage(artifactname string, solidityname string) (string, error) {
	var dir string
	var err error
	var pmstring string
	fmt.Println()
	if os.Getenv("TEST") != TRUEENV {
		dir, err = os.Getwd()
		fmt.Println(dir)
		if err != nil {
			err = fmt.Errorf("Could not get the working directory: '%v'", err)
			return "", err
		}
	} else {
		dir = "../../../test/testdata"
	}
	artifactpath := dir + "/build/contracts/" + artifactname + ".json"
	var file *os.File
	var info os.FileInfo
	file, err = os.Open(filepath.FromSlash(artifactpath))
	if err != nil {
		err = fmt.Errorf("Could not open file '%v': '%v'", artifactpath, err)
		return "", err
	}
	info, _ = file.Stat()
	artifact := make([]byte, info.Size())
	_, err = file.Read(artifact)
	if err != nil {
		err = fmt.Errorf("Could not read file '%v': '%v'", artifactpath, err)
		return "", err
	}
	artifactObject := ArtifactInfo{}
	err = json.Unmarshal(artifact, &artifactObject)
	if err != nil {
		err = fmt.Errorf("Could not unpackage truffle artifact: '%v'", err)
		return "", err
	}
	nodepackagepath := dir + "/package.json"
	file, err = os.Open(filepath.FromSlash(nodepackagepath))
	if err != nil {
		err = fmt.Errorf("Could not open file '%v': '%v'", nodepackagepath, err)
		return "", err
	}
	info, _ = file.Stat()
	pkg := make([]byte, info.Size())
	_, err = file.Read(pkg)
	if err != nil {
		err = fmt.Errorf("Could not read file '%v': '%v'", artifactpath, err)
		return "", err
	}
	pkgObject := PackageInfo{}
	err = json.Unmarshal(pkg, &pkgObject)
	if err != nil {
		err = fmt.Errorf("Could not unpackage package.json: '%v'", err)
		return "", err
	}
	fmt.Printf("%v+\n", artifactObject)
	fmt.Printf("%v+\n", pkgObject)

	pm := ethpm.PackageManifest{}

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
		commithashpath = dir + "/.git/refs/heads/master"
	} else {
		commithashpath = dir + "/master"
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
	sources["./contracts/"+solidityname+".sol"] = pkgObject.Homepage + ".git@" + commit
	pm.Sources = sources
	contracttype := ethcontract.ContractType{}
	contracttype.ABI = artifactObject.ABI
	compiler := artifactObject.Compiler
	contracttype.Compiler = compiler
	dbc := bytecode.UnlinkedBytecode{}
	dbc.Bytecode = artifactObject.DeployedBytecode
	contracttype.RuntimeBytecode = &dbc
	rbc := bytecode.UnlinkedBytecode{}
	rbc.Bytecode = artifactObject.Bytecode
	contracttype.DeploymentBytecode = &rbc
	contracttypes := make(map[string]*ethcontract.ContractType)
	contracttypes[artifactname] = &contracttype
	pm.ContractTypes = contracttypes
	pmstring, err = pm.Write()
	if err != nil {
		err = fmt.Errorf("Could not write object: '%v'", err)
		fmt.Println(err)
		return "", err
	}
	return pmstring, nil
}

// WritePackage TODO some comments
func WritePackage(json string) error {
	var p *os.File
	var err error
	if os.Getenv("TEST") != "true" {
		p, err = os.Create("./ethpm.json")
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
