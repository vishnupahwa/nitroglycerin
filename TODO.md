# TODO

## TODO
- Think of cool name - nitroglycerin
- Add progress bar for both orchestrator and normal run
- Make sure runs can be graphed, plotted and recorded correctly
with all required information
- Allow tweaking different things through both code and yaml
- Add assertions on metric results
- Update code and tests to use both random or configured ports for easier testing
- Add e2e tests with fake app
  - replace bash/binary commands with kubernetes real API library:
    - use kustomize to build the resource but parse it into Go structs
    - Apply all structs & watch job.
  - test --args are passed from orchestration to pods
  - test distributed rate is correct
  - test NFT results are concatenated correctly
  
- Use gobbing instead of JSON for grpc
- Add versioning information through code generation, a
real container image default for jobs and 
versioning scripts in ./hack folder

- Add throttling/before/after stage functions that
 run parallel rather than in between stages 
 
- Allow for recording all request data or just recording metrics from orchestrator
 as the former requires lots of memory for orchestrator <- Is this needed?
 
- Reduce memory allocations for increasing efficiency. (See fasthttp) 
- see profiling - https://blog.golang.org/pprof
