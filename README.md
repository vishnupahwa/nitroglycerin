# eznft - Distributed NFT Wrapper for vegeta

## TODO
- Think of cool name - nitroglycerin
- Divide load correctly for a scenario - setup scenario definitions in a way this 
can be done depending on what is called
- Add correct pacer for ramp up + constant load stage due to time between switching stages
- && change 'definitions' to be more clear
- Dynamic memory and cpu limits of job pods
- Add progress bar for both orchestrator and normal run
- Collate NFT results
- Make sure runs can be graphed, plotted and recorded correctly
with all required information
- Allow tweaking different things through both code and yaml
- Add assertions on metric results
- Add testing
- Add e2e tests with fake app

- Add throttling/before/after stage functions that
 run parallel rather than in between stages 
 
 - Allow for recording all request data or just recording metrics from orchestrator
 as the former requires lots of memory for orchestrator <- Is this needed?