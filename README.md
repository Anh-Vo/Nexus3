#Getting Started with Sonatype's Nexus3 API
This is a rough and dirty guide on how to get started with Sonatype's Nexus3 API written so even a dummy (like me) can easily understand.

## The Challenge 
To create a maven proxy repository in Nexus3 utilizing without ever touching the GUI (because I couldn't)

## My Setup 
I am running Nexus `3.22.1-02` off of the Docker official image. You can replicate my environment by running the following:

``` 
$ docker run -d -p 8080:8081 --name nexus -v nexus-data:/nexus-data sonatype/nexus3
```
Note that I am exposing this service out of port 8080. The configurations/scripts that I have included will reflect this change.
Your setup may vary so adjust accordingly.

## Nexus3 Docker Image Password
When running the Docker setup, the default password is not `admin123`, you'll need to access `exec` into the container and retrieve the
admin password before proceeding.

```
$ docker exec -it nexus bash       // Where `nexus` is the name of my container
$ cat ~/nexus-data/admin.password  // Copy and save the output
$ exit                             // Exit the container
```

If you are able to access the Nexus through the browser, you can login as `admin` and the copied password from the steps above where
you will be prompted to change the default password to something nicer. If not, then you can simply use the copied password.

## IMPORTANT
If you are utilizing Nexus `3.21.1` or older, make sure you enable script creation which is turned off by default and will cause 
everything to fail with a rare, but tool `HTTP 410 Gone` error.

To enable it script creation:

```
$ docker exec -it nexus bash
$ vi nexus-data/etc/nexus.properties

//Add the following:  nexus.scripts.allowCreation=true
//Save and exit
```

If your container is already running, you may need to restart it

```
$ docker restart nexus
```

More details here: 
https://issues.sonatype.org/browse/NEXUS-23205
## The Tools
1. Groovy `Mac: brew install groovy | Ubuntu: apt get install groovy`
2. Scripts:https://issues.sonatype.org/browse/NEXUS-23205
    1. [addUpdateScript](https://github.com/sonatype-nexus-community/nexus-scripting-examples/blob/master/complex-script/addUpdateScript.groovy)
    2. [provisionScript](https://github.com/sonatype-nexus-community/nexus-scripting-examples/blob/master/complex-script/provision.sh)

## Official Guide
https://github.com/sonatype-nexus-community/nexus-scripting-examples

## Overview
At the time of writing, Nexus3 only supports `GET` with `/service/rest/v1/repositories`. Which means that you cannot use this endpoint to create
repositories.

In order to create a repository, you'll need to utilize the [scripts-api](https://help.sonatype.com/repomanager3/rest-and-integration-api/script-api).
This will allows you do is basically `POST` JSON payloads containing Groovy code that can then be executed on the server (i.e. Nexus).

The workflow is basically this:
1. You write some Groovy code with your desired configurations (I want to create a repository that acts as a Maven proxy)
2. You `POST` your code to the Nexus  
3. You execute the code of interest  

For basic one liner scripts, you could simply write the JSON file yourself and push it up the Nexus using curl, as demonstrated in the docs.
However, for more complex scripts you can simply utilize the two scripts I linked to above to help you convert your multiline Groovy script
into a json payload and execute it all in one go.

For example groovy script called `maven.groovy` that creates a repository that serves as a maven proxy could be: 

`repository.createMavenProxy("test", "https://repo1.maven.org/maven2/", 'default', false, org.sonatype.nexus.repository.maven.VersionPolicy.SNAPSHOT, org.sonatype.nexus.repository.maven.LayoutPolicy.STRICT)`

Now, to get everything working, we need to include our files of interest in the `provision.sh` script. There are some default entries there to serve
as examples. However, you can simply remove them and place your Groovy script(s) there instead as I did.

`addAndRunScript HelloMaven maven.groovy`

Once this is done, make sure your provision script has the right permissions and execute it:

```
$ sudo chmod +x provision.sh
$ ./provision.sh
```

At this point `provision.sh` will kick off Groovy and get everything going. If Groovy complains about not being able to pull some dependencies, a
quick fix that worked for me was to remove the m2 and grape folder:

```
$ rm -rf ~/.m2
$ rm -rf ~/.groovy/grapes
```

If everything works correctly you can do a `GET` to `/service/rest/v1/repositories` to check if your repo is there.

```
$ curl -v -u admin:admin123 -X GET http://localhost:8080/service/rest/v1/script
```

If you dig through the output, you should see a maven repo called "test" (if you used my maven.groovy) or if not then look
for whatever you named your repo as.

Hopefully it's there... haha





