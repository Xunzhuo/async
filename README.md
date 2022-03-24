<div id="top"></div>


[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
<!-- [![Apache License][license-shield]][license-url] -->



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/Xunzhuo/async">
    <img src="docs/assets/images/logo.png" alt="Logo" width="80" height="80">
  </a>

<h1 align="center"> Async </h3>
  <p align="center">
    Async is a lightwight, easy-to-use, high performance, more human-being Asynchronous Engine
  </p>
</div>

[![GoDoc](https://godoc.org/github.com/Xunzhuo/async?status.svg)](https://godoc.org/github.com/Xunzhuo/async)
[![Build Status](https://travis-ci.org/KXunzhuo/async.svg?branch=master)](https://travis-ci.org/Xunzhuo/async)
[![Go Report Card](https://goreportcard.com/badge/github.com/Xunzhuo/async)](https://goreportcard.com/report/github.com/Xunzhuo/async)
[![Coverage Status](https://coveralls.io/repos/github/Xunzhuo/async/badge.svg?branch=master)](https://coveralls.io/github/Xunzhuo/async?branch=master)


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#architecture">Architecture</a></li>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


<!-- ABOUT THE PROJECT -->
## About this Project

Async is a lightwight, easy-to-use, high performance, more human-being Asynchronous Engine

### Spotlights

+ Async is very easy to use, creating the jobs just by a couple of lines.
+ Support **master/slave** job mode or **standalone** job mode 
+ Provide powerful options to control the jobs like the MaxNumber of WorkQueue.
+ Provide inner cache to speed up to get the cached jobs data.
+ Help you easily manage your jobs into asynchronous way like:
    + reducing the time like the long time of http response in large number of requests
    + reducing the time like when interating with DataBase to query SQL
    + reducing ....

<!-- GETTING STARTED -->
## Getting Started

<div align="center">
<img src="docs/assets/images/code.png" alt="Code" width="75%">
</div>

### Installation

``` go
  go get github.com/Xunzhuo/async
```

#### Built with

* [Golang](https://go.dev/)
* [Docker](https://www.docker.com/)
* [Docker-Compose](https://docs.docker.com/compose/)

<!-- USAGE EXAMPLES -->
## Concepts

Async supports two Running Mode:
+ The Standalone Job mode
  In this mode, each job has an unique JobID, you can create one job for one JobID
+ The Master/Slave Job mode
  In this mode, the job ID can be called as master job ID, the master job ID is unique as well
  one master job ID can contains a few slave jobs with subID, you can create one master job with many slave jobs

**Async takes JobID as the key to create/find/update/delete Job**

JobID in Async has two kinds:
+ jobID: 
  + the unique job id in standalone job mode
  + the master job id in master/slave job mode
+ subID: the slave id in master/slave job mode

### Quick Start

``` go
  // create a job
	job := async.NewJob("Unique JobID", JobFunc, JobFuncParams)
  // add job to default engine
	async.Engine.AddJobAndRun(job)
  // get job data by job id
	async.Engine.GetJobData("Unique JobID")
```

### Demo

+ [The Standalone Job mode](demos/standalone/standalone.go)
+ [The Master/Slave Job mode](demos/masterSlave/masterSlave.go)

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/Xunzhuo/async/issues) for a full list of proposed features (and known issues).

## Contributors

<a href="https://github.com/merbridge/merbridge/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Xunzhuo/async" />
</a>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

<!-- CONTACT -->
## Contact

Project Link: [https://github.com/Xunzhuo/async](https://github.com/Xunzhuo/async)


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/Xunzhuo/async.svg?style=for-the-badge
[contributors-url]: https://github.com/Xunzhuo/async/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Xunzhuo/async.svg?style=for-the-badge
[forks-url]: https://github.com/Xunzhuo/async/network/members
[stars-shield]: https://img.shields.io/github/stars/Xunzhuo/async.svg?style=for-the-badge
[stars-url]: https://github.com/Xunzhuo/async/stargazers
[issues-shield]: https://img.shields.io/github/issues/Xunzhuo/async.svg?style=for-the-badge
[issues-url]: https://github.com/Xunzhuo/async/issues
