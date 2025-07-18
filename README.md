# gator

Welcome to my implementation of Gator, an RSS blog aggrigator written in Go for the boot.dev course [Build a blog aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang). This project took me about 20h to complete and I'm quite happy with how it's progressed. 

## Setup

You'll need to make sure you have Golang installed on your machine. Instructions to do so will not be provided here. 

### Install Posgres

First thing you'll need to do to set up gator is to download and enable postgres. This proccess is explained quite well in the course linked above which is free to read, but I'll reitterate what they have to say here.


#### macOS with [brew](https://brew.sh/)

Install the package
```sh
brew install postgresql@15
```

Ensure the installation via the psql command line
```sh
psql --version
```

Start the postgres service
```sh
brew services start postgresql@15
```

open the psql shell
```sh
psql postgres
```

Finally, create the database that will store all the information required by Gator
```sql
CREATE DATABASE gator;
```

you can confirm that you've successfully created the database by connecting to it with 
```sh
\c gator
```

#### Ubuntu

Install postgres via your package manager
```sh
sudo apt install postgresql postgresql-contrib
```

Ensure the installation via the psql command line
```sh
psql --version
```

update the postgres password
```sh
sudo passwd postgres
```

Start the postgres service
```sh
sudo service postgresql start
```

open the psql shell
```sh
sudo -u postgres psql
```

Finally, create the database that will store all the information required by Gator
```sql
CREATE DATABASE gator;
```

you can confirm that you've successfully created the database by connecting to it with 
```sh
\c gator
```

### Use [Goose](https://github.com/pressly/goose) to migrate the database to a correct state

We'll need to set up the internals of the database to get this going. thankfully it's only one command. 


## WIP