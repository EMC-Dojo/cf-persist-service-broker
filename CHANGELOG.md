# CloudFoundry Persistence Service Broker Version Log
## v1
vccfe2a3622428bd93521f0b6ce8fbb0116b3413f Initial Commit With CI
## v2
v8cee9155abeb71d0f38359f2c834240f21461ea1 Resolve depencies problem on libstorage by deleting api/server, import/routers, and api/tests folders
## v3
v39a8ed044f1e16e21c78182f36b711ca2163b1a2 Fix integration test script
## v4
vbdfe5338002907e8638c1aa94ddce01257000efd Add bind, unbind, deprovision, provision into service broker
## v5
vbad235205fce7458ec536602ceb2bc624342c4a8 Modified catalog and binding for service broker to match that of cf-scaleio-broker
## v6
v5b20ade2d236ecea4b628aad0a8efefcc6a4e9a4 Change service broker to use volumeName for Binding
## v8
v88dc765d825417ce657ada3be5c43330f351d2a1 Reuse service broker instance from lifecycle test
## v9
va5ed98d54ea2e083e4cfc1bdfd5f95746d54e827 Change to non-interactive for delete app
## v10
v2d5fe573ecbf0eb974b9167f4cc4a13db54e6a87 Change delete app to restage app
## v11
v48f8eb128e036f37b15c2a489abab31101dfbe35 Combine acceptance tests to lifecycle
## v12
vde0ffcd82e6956aca93b5422e5db52ed3d0cb60e Changed EMC-CMD to EMC-Dojo
## v13
v710575ae06bd28fcfa28699d705e3999cda85376 Add configure persistence docs
## v14
vd043756be46d96a5d9f8d1dd038206cc5b0a5c46 copy README.md inside docs folder
## v15
v49339eb278b1e73c2bf389896ac8764d3743683b Format configure persistence document
## v16
v94e735020a0f7cdefed34de1eea487604ce266d3 Error handling for delete service instance when service is still attached to App [#121134559](https://www.pivotaltracker.com/story/show/121134559)
## v17
v1fdbc21b77cc0064eba43f82bc970731f4472e05 Adding unbind service to lifecycle step of Pipeline
## v18
v7ff5dc1dd8cdb8544f6300be65f4cb4fa6ea295e Created new block in pipeline for releasing new candidate
## v19
v8293a3dbe354b1ce1e5b567c10ed2d0591d81cc5 Setup Lifecycle to check all Diego cells when running multimount storage
