# gobci
A package intended to use with the [buffer_bci](https://github.com/jadref/buffer_bci) framework, to connect the hardware of your choice to your Go application. 

## Status
Currently it only *receives* data from the buffer server. There is no intention to build the *sending* part yet, 
because there is nothing useful to send (as we're not a hardware driver). 

There's still some testing to be done, and there are no guarantees about the backwards compatibility just yet. 

## Examples
There is currently one example, that gathers all samples from the buffer server, and saving a plot of three sample
channels to `output.jpg`. It is available in the [examples](https://github.com/EtienneBruines/gobci/tree/master/examples/) directory. 

There's a [blog post] that accompanies this repository. 

## Contributing
Contributions are very welcome. Just create an issue, work on it yourself (if you want), and create a PR. 

If there haven't been a lot of commits, it's because I ran out of ideas to implement. 