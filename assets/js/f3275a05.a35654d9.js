"use strict";(self.webpackChunkfirehose=self.webpackChunkfirehose||[]).push([[403],{3905:function(e,t,n){n.d(t,{Zo:function(){return l},kt:function(){return d}});var o=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);t&&(o=o.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,o)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,o,r=function(e,t){if(null==e)return{};var n,o,r={},i=Object.keys(e);for(o=0;o<i.length;o++)n=i[o],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(o=0;o<i.length;o++)n=i[o],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var s=o.createContext({}),u=function(e){var t=o.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},l=function(e){var t=u(e.components);return o.createElement(s.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return o.createElement(o.Fragment,{},t)}},p=o.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,s=e.parentName,l=c(e,["components","mdxType","originalType","parentName"]),p=u(n),d=r,f=p["".concat(s,".").concat(d)]||p[d]||m[d]||i;return n?o.createElement(f,a(a({ref:t},l),{},{components:n})):o.createElement(f,a({ref:t},l))}));function d(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,a=new Array(i);a[0]=p;var c={};for(var s in t)hasOwnProperty.call(t,s)&&(c[s]=t[s]);c.originalType=e,c.mdxType="string"==typeof e?e:r,a[1]=c;for(var u=2;u<i;u++)a[u]=n[u];return o.createElement.apply(null,a)}return o.createElement.apply(null,n)}p.displayName="MDXCreateElement"},5057:function(e,t,n){n.r(t),n.d(t,{assets:function(){return l},contentTitle:function(){return s},default:function(){return d},frontMatter:function(){return c},metadata:function(){return u},toc:function(){return m}});var o=n(7462),r=n(3366),i=(n(7294),n(3905)),a=["components"],c={},s="Contribution Process",u={unversionedId:"contribute/contribution",id:"contribute/contribution",title:"Contribution Process",description:"The following is a set of guidelines for contributing to Raccoon. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request. Here are some important resources:",source:"@site/docs/contribute/contribution.md",sourceDirName:"contribute",slug:"/contribute/contribution",permalink:"/raccoon/contribute/contribution",draft:!1,editUrl:"https://github.com/raystack/raccoon/edit/master/docs/docs/contribute/contribution.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Metrics",permalink:"/raccoon/reference/metrics"},next:{title:"Development Guide",permalink:"/raccoon/contribute/development"}},l={},m=[{value:"How can I contribute?",id:"how-can-i-contribute",level:2},{value:"Becoming a maintainer",id:"becoming-a-maintainer",level:2},{value:"Guidelines",id:"guidelines",level:2}],p={toc:m};function d(e){var t=e.components,n=(0,r.Z)(e,a);return(0,i.kt)("wrapper",(0,o.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"contribution-process"},"Contribution Process"),(0,i.kt)("p",null,"The following is a set of guidelines for contributing to Raccoon. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request. Here are some important resources:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/raccoon/contribute/contribution"},"Concepts")," section will explain you about Raccoon architecture,"),(0,i.kt)("li",{parentName:"ul"},"Our ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/raystack/raccoon/docs/roadmap.md"},"roadmap")," is the 10k foot view of how we envision Raccoon to evolve"),(0,i.kt)("li",{parentName:"ul"},"Github ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/raystack/raccoon/issues"},"issues")," track the ongoing and reported issues.")),(0,i.kt)("h2",{id:"how-can-i-contribute"},"How can I contribute?"),(0,i.kt)("p",null,"We use RFCS and GitHub issues to communicate ideas."),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"You can report a bug or suggest a feature enhancement or can just ask questions. Reach out on Github discussions for this purpose."),(0,i.kt)("li",{parentName:"ul"},"You are also welcome to add new features, improve monitoring,logging and code quality."),(0,i.kt)("li",{parentName:"ul"},"You can help with documenting new features or improve existing documentation."),(0,i.kt)("li",{parentName:"ul"},"You can also review and accept other contributions if you are a maintainer.")),(0,i.kt)("p",null,"Please submit a PR to the master branch of the Raccoon repository once you are ready to submit your contribution. Code submission to Raccoon, including submission from project maintainers, require review and approval from maintainers or code owners. PRs that are submitted by the general public need to pass the build. Once build is passed community members will help to review the pull request."),(0,i.kt)("h2",{id:"becoming-a-maintainer"},"Becoming a maintainer"),(0,i.kt)("p",null,"We are always interested in adding new maintainers. What we look for is a series of contributions, good taste, and an ongoing interest in the project."),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"maintainers will have write access to the Raccoon repositories."),(0,i.kt)("li",{parentName:"ul"},"There is no strict protocol for becoming a maintainer or PMC member. Candidates for new maintainers are typically people that are active contributors and community members."),(0,i.kt)("li",{parentName:"ul"},"Candidates for new maintainers can also be suggested by current maintainers or PMC members."),(0,i.kt)("li",{parentName:"ul"},"If you would like to become a maintainer, you should start contributing to Raccoon in any of the ways mentioned. You might also want to talk to other maintainers and ask for their advice and guidance.")),(0,i.kt)("h2",{id:"guidelines"},"Guidelines"),(0,i.kt)("p",null,"Please follow these practices for you change to get merged fast and smoothly:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Contributions can only be accepted if they contain appropriate testing ","(","Unit and Integration Tests.",")"),(0,i.kt)("li",{parentName:"ul"},"If you are introducing a completely new feature or making any major changes in an existing one, we recommend to start with an RFC and get consensus on the basic design first."),(0,i.kt)("li",{parentName:"ul"},"Make sure your local build is running with all the tests and ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/golang/lint"},"golint")," checks passing."),(0,i.kt)("li",{parentName:"ul"},"Docs live in the code repo under ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/raystack/raccoon/docs/README.md"},(0,i.kt)("inlineCode",{parentName:"a"},"docs"))," so that changes to that can be done in the same PR as changes to the code.")))}d.isMDXComponent=!0}}]);