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
as you can see running `builder.bash -h`

## Installation and configuration
Before building and running the application, you can edit the [configuration file](build/comunigo.cfg) to specify some environment configuration settings.

After that you can use [startup script](build/startup.bash) to effectively startup p2p nodes and other components of the net.

As you can see running `bash startup -h` there are a few options that permits you to further configure your experience with comuniGO, for example:
- You have to use `-n` flag to set the number of peers to spawn
- You have to use `-t` flag to set the algorithm to use (_sequencer_, _scalar_ or _vectorial_)
- You can use `-v` flag to enable verbose output on javascript console of your browser to see details of sent and/or received messages
- You can use `-a` flag to enable attach mode in `docker-compose up` 

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## Testing
