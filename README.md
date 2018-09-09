rubberducky-go
=========================
A go tool tightly integrated with geth which packages your project in accordance with the [EthPM v2 package manifest](https://github.com/ethpm/ethpm-spec) and publishes to Ethereum.

# Layout
This repository abides by the standard layout [as defined here](https://github.com/golang-standards/project-layout)

# Tools
This repository uses:   
* [dep for dependency management](https://golang.github.io/dep/)   
* [gitflow for branch workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)    

## Usage
**Please read carefully because this is currently a hackathon build**   

RubberDucky is a cli tool which will currently package a truffle project into the ethpm v2 spec and post to the PackageDB contract on Rinkeby.   

*Note* rubberducky.test currently points to the PackageIndex associated with this db, however, this is not hooked up yet. Additionally, this has been built specifically with complete truffle projects in mind. Therefore, for now, you must have a truffle project direcotry that is:   
1. Complete   
2. You have run `truffle compile` and have a `build/contracts` folder with .json artifacts   
3. `git` is installed in the directory   
4. If you need one, [here is a simple repo to clone and use](https://github.com/Modular-Network/ethereum-libraries-basic-math).

### Steps

1. Start a geth node syncing to Rinkeby and ensure you have eth available on a local rinkeby private key. RubberDucky uses all defaults for now so any custom filepaths won't work.
2. View [instructions here to install golang on your machine](https://golang.org/doc/install).
3. Ensure your `$GOPATH` is set.
4. `go get github.com/Hackdom/rubberducky-go`
5. `cd $GOPATH/src/github.com/Hackdom/rubberducky-go`
	* You will notice the vendor folder was committed. There are several unstable dependencies that have been fixed which are not fully committed on github, therefore, keep the vendor folder to use those dependencies and don't `dep ensure`.
6. There are several possibilities for execution:
	* If your `$GOPATH/bin` is in your `$PATH` then you should do following:
		1. Change to the parent directory of your truffle projects
		2. Run `rubberducky` and follow the prompts. The first prompt asks for a relative path to a truffle project, so if you followed this, you should type `./<yourTruffleProject>`
		3. Enjoy
	* If you don't have your `$PATH` set, then you can move the `rubberducky` executable in `$GOPATH/bin` to the parent directory of your truffle projects, or if you're confident with typing relative paths, just leave it where it is. Then:
		1. Change to the directory where the executable resides
		2. Run `./rubberducky` and follow the prompts. It will ask for a relative path to a truffle project so you just type `../../../../<howevermanytimesyouneed>/then/finda/<truffleProject>`
		3. Enjoy

### ETHPM and Smart Contracts

This tool relies on the core [build of ethpm-go found here](https://github.com/ethpm/ethpm-go). There is still a lot of work to do on the core tool as well. Additionally, I was able to make extensive use of Piper's excellent work on the current build of the [ethpm smart contract repository found here](https://github.com/ethpm/escape-truffle).

## Feedback

Any feedback would be appreciated here. Feel free to use issues and PR's or email me at contact@modulartech.io. I'm interested to here whether or not this was a worthwhile endeavor.
