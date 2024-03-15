<p align="center">
    <h1 align="center">Log Parser</h1>
<p align="center">
   <a href='https://github.com/diegoclair/log-parser/commits/main'>
	<img src="https://img.shields.io/github/last-commit/diegoclair/log-parser?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
   </a>
   <a href="https://github.com/diegoclair/log-parser/actions">
     <img src="https://github.com/diegoclair/log-parser/actions/workflows/ci.yaml/badge.svg" alt="build status">
   </a>
  <a href='https://goreportcard.com/badge/github.com/diegoclair/log-parser'>
     <img src='https://goreportcard.com/badge/github.com/diegoclair/log-parser' alt='Go Report'/>
    </a>
<p>
  
## Description 
This project is a boilerplate for creating API projects in Go, incorporating key principles of Domain-Driven Design and Clean Architecture. It reflects my expertise in Golang, drawn from previous projects, and is structured to facilitate maintainability and scalability.

All layers of the codebase are tested to ensure reliability and robustness. The project is open to contributions and improvements. Feel free to fork the repository, submit pull requests, or open issues to discuss enhancements or report bugs.
  
### Project architecture:
<div align="center">
    <img src='./.github/assets/architecture.png' />
</div>

For the presentation layer, which I refer to as the transport layer, its purpose is to handle data transportation. It is responsible for receiving and responding to various types of requests, including API calls, gRPC, and messaging via RMQ, among others.

### Observations:
If player kill it self, will not be counted.  
`22:18 Kill: 2 2 7: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH`  
<br>
The player can change your username during the game, so we consider as player name, only the last used username for that player

### Tests:
For unit tests with MySQL and Redis, we are using real dependencies with [testcontainers](https://testcontainers.com/). 
It’s like putting our functions through a real-world.  
And we are also using mocks to test the errors scenarios, this way we can achieve 100% of cover. 💪 

## 💻 Getting Started 

### Prerequisites ❗
* Ensure Docker is installed on your machine.
* An installation of Go 1.18 or later. For installation instructions, see [Installing Go](https://go.dev/doc/install).

### ▶️ Launching the Application 
To start the application, execute the command: 
```bash
docker-compose up
```
Once you see the message `your server started on [::]:5000`, the application is up and running!

## 📝 API Documentation:
For detailed API endpoint documentation, navigate to the `/docs` directory where you will find the Swagger documentation.  
These swagger docs was generated using [goswag](https://github.com/diegoclair/goswag)

## Running tests
```bash
make tests
```
## Generating docs
```bash
make docs
```

##  Contributing

Contributions are welcome! Here are several ways you can contribute:

- **Submit Pull Requests**: Review open PRs, and submit your own PRs.
- **[Join the Discussions](https://github.com/diegoclair/log-parser/discussions)**: Share your insights, provide feedback, or ask questions.
- **[Report Issues](https://github.com/diegoclair/log-parser/issues)**: Submit bugs found or log feature requests for log-parser.

<details closed>
    <summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your GitHub account.
2. **Clone Locally**: Clone the forked repository to your local machine using a Git client.
   ```sh
   git clone https://github.com/<your_username>/log-parser
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to GitHub**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.

Once your PR is reviewed and approved, it will be merged into the main branch.

</details>  

---

##  License

This project is protected under the [MIT License](https://choosealicense.com/licenses/mit/) License. For more details, refer to the [LICENSE](./LICENSE) file.

<br>
:bowtie:
