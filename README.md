# comuniGO

**comuniGO** is a [gRPC](https://grpc.io/docs/languages/go/) based P2P application developed in [GO](https://golang.org/) that supports both totally and causally ordered message delivery.

In the figure below **L. Lamport**, one of the main actors in _Distribuited Systems_ field, **S. Hykes**, a co-founder of Docker.inc, and **K. Thompson**, a co-designer of GO programming language, communicate with each other using **comuniGO**

![Demo of comuniGO](documentation/readme-gif/comunigo.gif)

## Building of components
The application is based on three different components:
1. The **peer** with whom you can interact via the web and internally uses gRPC to exchange messages with others
2. The **registration** service with which you can sign yourself and receive the peer list
3. The **sequencer** service necessary to implement the centralized totally ordered algorithm

So, thanks to the [builder script](build/builder.bash) it is possible to discriminate in fine grain what to build, in particular:
- You can use `-p` flag to build **peer** image
- You can use `-r` flag to build **registration service** image
- You can use `-s` flag to build **sequencer** image

as you can see running `bash builder.bash -h`

## Installation and configuration
Before building and running the application, you can edit the [configuration file](build/comunigo.cfg) to specify some environmental configuration settings.

After that you can use [startup script](build/startup.bash) to effectively startup peers and other components of the architecture.

As you can see running `bash startup -h` there are a few options that permits you to further configure your experience with comuniGO, for example:
- You have to use `-n` flag to set the number of peers to spawn
- You have to use `-t` flag to set the algorithm to use (`sequencer`, `scalar` or `vectorial`)
- You can use `-v` flag to enable verbose output on javascript console of your browser to see details of sent and/or received messages
- You can use `-a` flag to enable attach mode in `docker-compose up` 

## Testing
To test appplication simply:
1. Run the [startup script](build/startup.bash) as described above specifing which algorithm do you want to test
2. Run the [testing script](build/test.bash) specifing:
    - The type of service to test (`sequencer`, `scalar` or `vectorial`), using `-t` flag
    - The modality of the test (`single` or `multiple`), Using `-m` flag 
3. Finally, once collected results, you should shutdown components running `sh shutdown.sh`

## Usage
To use application simply follow one of this approach:
- Run [discovery script](build/discover.sh) to get links for peers' frontend
- Clone this [repository](https://gitlab.com/tibwere/comunigo-peer-discovery) and follow instructions inside the README file to achieve the same result

Enjoy comunication using comuniGO app and finally, you should shutdown components running `sh shutdown.sh`

## Compile documentation
To compile the documentation, first you need to `cd` into documentation folder and then run whatever tool you want for compiling `.tex` files (e.g. `pdflatex`, `latexmk`, ...)

For instance:
- Using `latexmk` you can run `latexmk -pdf -silent comunigo.tex`
- Using `pdflatex` you can run `pdflatex comunigo.tex && pdflatex comunigo.tex` (twice for the resolution of references in the document)
