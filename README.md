# eznft - Distributed NFT Wrapper for vegeta

## TODO
- Think of cool name - nitroglycerin
- Add progress bar for both orchestrator and normal run
- Make sure runs can be graphed, plotted and recorded correctly
with all required information
- Allow tweaking different things through both code and yaml
- Add assertions on metric results

- Add missing tests
- Add e2e tests with fake app

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
