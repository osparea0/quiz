## Quiz CLI 

Quiz a CLI based quiz application that asks the users multiple choice questions and allows them to select only one answer. It's written in Go and uses the [Cobra]https://github.com/spf13/cobra library as well as [promptui]https://github.com/manifoldco/promptui . 

### How To

### Install requirements

You will need to have Go installed in order to install the CLI app. You can follow the instructions for installing it here https://go.dev/ . The git CLI is also required https://git-scm.com/downloads . Once Go and Git are installed, clone the repository into a local directory using `git clone git@github.com:osparea0/quiz.git` . Once the clone is completed, `cd ./quiz` and run `go install` .

### Running the app

To launch the server portion of the app run the command `quiz start &` . You can also just run `quiz start` and then open another terminal session while keeping the previous session open for logging separation for each terminal. To "play" a quiz run the command `quiz play` . You will be prompted for a username, or you can just accept the default. Currently, a user may only take a single quiz. Use your arrow down and up keys to navigate the choices for each question. When you've highlighted your answer hit enter to accept it and move to the next question. A percentile is calculated along with the score for each completion of a quiz.
