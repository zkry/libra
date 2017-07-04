# libra ♎️
![A sample example of the program](docs/media/example.png?raw=true "An example of the program running")
Libra is a command line utility that shows various statistics related with a directory. The output is similar to what you would find at the top of a GitHub repo page. **This app is currently in development.** This code has only been tested on OSX.

## Installation and Usage
If you have Go setup you can run  `go get github.com/zkry/libra`. If you have your Go bin setup you can just run  `libra [directory]`. Omitting directory performs the command in the current directory.

You can also analize a github project by providing a parameter in the form 'github.com/username/project'. 

I made this project in order to get an insight on the various projects I am going to work on or am thinging about working on. With this app you can get a good idea about the contents of a project. The bar section shows a percentage of relative size that a particular file size takes up (the largest filesize will take up the full bar and the other bars are filled in proportion to that). The number in parenthesis is the number of lines of the program. Below are hand coded stats for a particular filetype. For example, if the projecte directory contains go files then a report about the number of functions, and interfaces is shown.
