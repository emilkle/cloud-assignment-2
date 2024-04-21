# Countries Dashboard Service

<!---   

## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

Already a pro? Just edit this README.md and make it your own. Want to make it easy? [Use the template at the bottom](#editing-this-readme)!

## Add your files

- [ ] [Create](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#create-a-file) or [upload](https://docs.gitlab.com/ee/user/project/repository/web_editor.html#upload-a-file) files
- [ ] [Add files using the command line](https://docs.gitlab.com/ee/gitlab-basics/add-file.html#add-a-file-using-the-command-line) or push an existing Git repository with the following command:

```
cd existing_repo
git remote add origin https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2024-workspace/mariusrd/countries-dashboard-service.git
git branch -M main
git push -uf origin main
```

## Integrate with your tools

- [ ] [Set up project integrations](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2024-workspace/mariusrd/countries-dashboard-service/-/settings/integrations)

## Collaborate with your team

- [ ] [Invite team members and collaborators](https://docs.gitlab.com/ee/user/project/members/)
- [ ] [Create a new merge request](https://docs.gitlab.com/ee/user/project/merge_requests/creating_merge_requests.html)
- [ ] [Automatically close issues from merge requests](https://docs.gitlab.com/ee/user/project/issues/managing_issues.html#closing-issues-automatically)
- [ ] [Enable merge request approvals](https://docs.gitlab.com/ee/user/project/merge_requests/approvals/)
- [ ] [Set auto-merge](https://docs.gitlab.com/ee/user/project/merge_requests/merge_when_pipeline_succeeds.html)

## Test and Deploy

Use the built-in continuous integration in GitLab.

- [ ] [Get started with GitLab CI/CD](https://docs.gitlab.com/ee/ci/quick_start/index.html)
- [ ] [Analyze your code for known vulnerabilities with Static Application Security Testing (SAST)](https://docs.gitlab.com/ee/user/application_security/sast/)
- [ ] [Deploy to Kubernetes, Amazon EC2, or Amazon ECS using Auto Deploy](https://docs.gitlab.com/ee/topics/autodevops/requirements.html)
- [ ] [Use pull-based deployments for improved Kubernetes management](https://docs.gitlab.com/ee/user/clusters/agent/)
- [ ] [Set up protected environments](https://docs.gitlab.com/ee/ci/environments/protected_environments.html)

***

# Editing this README

When you're ready to make this README your own, just edit this file and use the handy template below (or feel free to structure it however you want - this is just a starting point!). Thanks to [makeareadme.com](https://www.makeareadme.com/) for this template.

## Suggestions for a good README

Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Name
Choose a self-explaining name for your project.

## Description
Let people know what your project can do specifically. Provide context and add a link to any reference visitors might be unfamiliar with. A list of Features or a Background subsection can also be added here. If there are alternatives to your project, this is a good place to list differentiating factors.

## Badges
On some READMEs, you may see small images that convey metadata, such as whether or not all the tests are passing for the project. You can use Shields to add some to your README. Many services also have instructions for adding a badge.

## Visuals
Depending on what you are making, it can be a good idea to include screenshots or even a video (you'll frequently see GIFs rather than actual videos). Tools like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation
Within a particular ecosystem, there may be a common way of installing things, such as using Yarn, NuGet, or Homebrew.
However, consider the possibility that whoever is reading your README is a novice and would like more guidance. 
Listing specific steps helps remove ambiguity and gets people to using your project as quickly as possible. 
If it only runs in a specific context like a particular programming language version or operating system or has 
dependencies that have to be installed manually, also add a Requirements subsection. -->

## Assignment 2 in PROG2005 Cloud Technologies

### A fully functional GoLang REST API with configurable and dynamically populated information dashboards.

## Endpoints
<!--- Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README. -->

### Registrations
```
/dashboard/v1/registrations/
```
#### - HTTP request GET
Example GET requests:
* Get a configuration with a specific "id":
  * /dashboard/v1/registrations/1

* Example response:
    ```
    {
        "id": 1,
        "country": "Norway",
        "isoCode": "NO",
        "features": {
            "temperature": true,
            "precipitation": true,
            "capital": true,
            "coordinates": true,
            "population": true,
            "area": false,
            "targetCurrencies": [
                "EUR",
                "USD",
                "SEK"
            ]
        },
        "lastChange": "20240229 14:07"
    }
    ```

* Get multiple chosen configurations with a specific "id" for each document:
    * /dashboard/v1/registrations/1,2

* Example response:

    ```
    [
        {
            "id": 1,
            "country": "Norway",
            "isoCode": "NO",
            "features": {
                "temperature": true,
                "precipitation": true,
                "capital": true,
                "coordinates": true,
                "population": true,
                "area": false,
                "targetCurrencies": [
                    "EUR",
                    "USD",
                    "SEK"
                ]
            },
            "lastChange": "20240229 14:07"
        },
        {
            "id": 2,
            "country": "Sweden",
            "isoCode": "SE",
            "features": {
                "temperature": true,
                "precipitation": true,
                "capital": true,
                "coordinates": false,
                "population": false,
                "area": false,
                "targetCurrencies": [
                    "NOK",
                    "SEK",
                    "USD",
                    "DKK"
                ]
            },
            "lastChange": "20240324 10:57"
        }
    ]    
    ```

* Get all registered configurations:
    * /dashboard/v1/registrations/


* Example response:
    ```
    [
        {
            "id": 1,
            "country": "Norway",
            "isoCode": "NO",
            "features": {
                "temperature": true,
                "precipitation": true,
                "capital": true,
                "coordinates": true,
                "population": true,
                "area": false,
                "targetCurrencies": [
                    "EUR",
                    "USD",
                    "SEK"
                ]
            },
            "lastChange": "20240229 14:07"
        },
        {
            "id": 2,
            "country": "Sweden",
            "isoCode": "SE",
            "features": {
                "temperature": true,
                "precipitation": true,
                "capital": true,
                "coordinates": false,
                "population": false,
                "area": false,
                "targetCurrencies": [
                    "NOK",
                    "SEK",
                    "USD",
                    "DKK"
                ]
            },
            "lastChange": "20240324 10:57"
        },
        {
            "id": 3,
            "country": "Denmark",
            "isoCode": "DK",
            "features": {
                "temperature": false,
                "precipitation": true,
                "capital": true,
                "coordinates": true,
                "population": false,
                "area": true,
                "targetCurrencies": [
                    "NOK",
                    "MYR",
                    "JPY",
                    "EUR"
                ]
            },
            "lastChange": "20240324 16:19"
        },
    ]    
    ```

#### - HTTP request POST
Example POST request - Store a new registration on the server and return the associated "id":
* Path: /dashboard/v1/registrations/
* Request body:
    ```
    {
       "country": "France",
       "isoCode": "FR",
       "features": {
          "temperature": true,
          "precipitation": true,
          "capital": true,
          "coordinates": true,
          "population": true,
          "area": true,
          "targetCurrencies": [
                "EUR",
                "USD"
          ]
       }
    }
    ```
Example response:
```
{
    "id": 4
    "lastChange": "20240329 15:41"
}
```

#### - HTTP request PUT
Example PUT requests - Change a registered configuration's fields and return if the operation was successful:
* Path: /dashboard/v1/registrations/4
* Request body:
    ```
    {
       "country": "France",
       "isoCode": "FR",
       "features": {
          "temperature": false,   //The tempersture field wil be updated to false.
          "precipitation": true,
          "capital": true,
          "coordinates": true,
          "population": true,
          "area": true,
          "targetCurrencies": [
                "EUR",
                "CHF"    //The Currency is changed from USD to CHF.
          ]
       }
    }
    ```


Example response:
* Returns HTTP status code 200 if the operation was successful. 

#### - HTTP request DELETE
Example DELETE requests:
* Delete one registered configuration with a specific "id" field:
    * /dashboard/v1/registrations/4

Example response:
* Returns HTTP status code 204 if the operation was successful.

* Delete multiple registered configurations with specific "id" fields:
    * /dashboard/v1/registrations/3,4

Example response:
* Returns HTTP status code 204 if the operation was successful.


### Dashboards
```
/dashboard/v1/dashboards/
```
### HTTP request Dashboard
Example requests:
* Return a populated Dashboard based on a specific registration configuration with a specific "id" field:
  * /dashboard/v1/dashboards/1

Example response:
```
{
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": -2.4,
        "precipitation": 0,
        "capital": "Oslo",
        "coordinates": {
            "latitude": 62,
            "longitude": 10
        },
        "population": 5379475,
        "area": 0,
        "target_currencies": {
            "EUR": 0.085719,
            "SEK": 0.99713,
            "USD": 0.091073
        }
    },
    "last_retrieval": "20240418 10:04"
}
```

* No "id" specified:
 * /dashboard/v1/dashboards/

Example response:
```
Cannot retrieve dashboard because no ID was specified.
```

* More than one "id" specified:
 * /dashboard/v1/dashboards/1,2

Example response:
```
Cannot retrieve more than one dashboard, too many IDs specified.
```

* Non existing "id" specified:
 * /dashboard/v1/dashboards/999
 * /dashboard/v1/dashboards/abc

Example response:
```
Dashboard not found.
```


### Notifications
```
/dashboard/v1/notifications/
```
### HTTP request POST
Example requests:

To add a new webhook, send a POST request to the notifications endpoint detailing URL, Country and Event 
as described in example request body. A unique ID will be added upon storing in database and returned as a response.

- http://localhost:8080/dashboard/v1/notifications/

Request example body:
```
{
    "URL": "http://localhost:8080/dashboard/v1/dashboards/1",
    "Country": "NO",
    "Event": "Invoke"
}
```
Example response:
```
{
    "ID": "02595fec-75ce-4dda-ac8d-2450ddcc6d1e"
}
```
### HTTP request GET
Example requests:

To get all webhooks returned send GET request to URL: 

- http://localhost:8080/dashboard/v1/notifications/

Example response:
```
[
    {
        "ID": "9189d595-b00e-49a5-b5f4-90a322ae68c9",
        "URL": "https://webhook.site/20d8180f-b4d4-479e-9aa6-32d970dd21ae",
        "Country": "SE",
        "Event": "POST"
    },
    {
        "ID": "02595fec-75ce-4dda-ac8d-2450ddcc6d1e",
        "URL": "http://localhost:8080/dashboard/v1/dashboards/1",
        "Country": "NO",
        "Event": "Invoke"
    }
]
```
Example requests:

To get a single webhook identified by a unique ID specify the ID sending a GET request to URL:

- http://localhost:8080/dashboard/v1/notifications/9189d595-b00e-49a5-b5f4-90a322ae68c9

Example response:
```
[
    {
        "ID": "9189d595-b00e-49a5-b5f4-90a322ae68c9",
        "URL": "https://webhook.site/20d8180f-b4d4-479e-9aa6-32d970dd21ae",
        "Country": "SE",
        "Event": "POST"
    }
]
```
### HTTP request DELETE
Example requests:

To delete a webhook enter the webhook id of the corresponding webhook to be deleted as shown in the url below.
Id of deleted webhook is returned as a response.

- http://localhost:8080/dashboard/v1/notifications/02595fec-75ce-4dda-ac8d-2450ddcc6d1e

Example response:
```
{
    "ID": "02595fec-75ce-4dda-ac8d-2450ddcc6d1e"
}
The requested webhook were successfully deleted from the database.
```
### Invoke HTTP request POST

Registered webhooks are invoked if they have the same url as the methods invoking the following events:


REGISTER - webhook is invoked if a new configuration is registered

CHANGE - webhook is invoked if configuration is modified

DELETE - webhook is invoked if configuration is deleted

INVOKE - webhook is invoked if dashboard is retrieved (i.e., populated with values)

Example response:
```
{
  "Country": "EN",
  "Event": "GET",
  "ID": "b3fadc89-90e1-49c3-8098-99e67fd04d50",
  "time": "2024-04-16T22:06:35+02:00"
}

```

### Status
```
/dashboard/v1/status/
```
### HTTP request Dashboard
Example requests:
* Return HTTP status codes of Rest APIs used by the service, notoifications database, version and uptime of the service:
  * /dashboard/v1/status
Example response:
```
{
    "countries_api": 200,
    "meteo_api": 200,
    "currency_api": 200,
    "notification_db": 200,
    "webhooks": 7,
    "version": "V1",
    "uptime": 686
}
```

## Testing of the Countries Dashboard Service

In order to run tests for the Countries Dashboard Service, a locally emulated Firestore database needs to be configured and running on your computer.

### Setting up emulated Firestore database (Windows)

1. **Install Node.js**:
   Node.js is a runtime environment that lets you run JavaScript outside the browser. It includes npm (Node Package Manager) for managing dependencies.
   - Download Node.js (LTS) from [Node.js official site](https://nodejs.org/en).

2. **Install Java**:
   Java is required to run the Firestore Emulator.
   - Download JAVA JDK (LTS) from [Oracle Java download page](https://www.oracle.com/java/technologies/downloads/#java11).

3. **Configure Java Environment Variables**:
   - Open System Properties (System > Advanced System Settings > Environment Variables).
   - Under "System variables", click "New":
     - Variable name: `JAVA_HOME`
     - Variable value: path to your Java installation directory, e.g., `C:\Program Files\Java\jdk-<version>`
   - Update the System Path:
     - In "Environment Variables", under "System variables", find and edit "Path".
     - Add a new entry: `%JAVA_HOME%\bin`.
   - Restart your computer and verify the Java installation in the Terminal:
     - Run `java -version` and `javac -version` to confirm the installation and version.

4. **Install Firebase CLI**:
   - Open a terminal and install Firebase CLI using npm:
     ```bash
     npm install -g firebase-tools
     ```

5. **Set Firestore Emulator Environment Variable**:
   - In PowerShell, set the environment variable to "8081":
     ```powershell
     $env:FIRESTORE_EMULATOR_HOST="8081"
     ```
   - Verify the setting:
     ```powershell
     echo $env:FIRESTORE_EMULATOR_HOST
     ```
     The response should be "8081".

6. **Start the Emulated Firestore**:
   - Navigate to your project directory where the emulator setup files are located:
     ```bash
     cd path\to\your\project\countries-dashboard-service\firestoreEmulator\emulatorFiles
     ```
   - Start the emulator:
     ```bash
     firebase emulators:start --only firestore --project countries-dashboard-service
     ```
   - Access the Web UI via the link provided in the terminal output: `http://127.0.0.1:4000/firestore/data`
   - Link: [Emulator UI](http://127.0.0.1:4000/firestore/data)

7. **Testing with Emulated Firestore**:
    - With the emulator running, you can now run tests to verify functionality of functions without affecting the production database.

## Support
If you have any questions or issues, please contact us at:
- Marius Robsahm Dahl: <mariusrd@stud.ntnu.no>
- Kristian Welle Glomset: <kglomset@gmail.com>
- Emil Klevgård-Slåttsveen: <emilkle@stud.ntnu.no>