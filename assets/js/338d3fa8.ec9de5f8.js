"use strict";(self.webpackChunkraccoon=self.webpackChunkraccoon||[]).push([[693],{5680:(e,t,r)=>{r.d(t,{xA:()=>p,yg:()=>y});var a=r(6540);function n(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,a)}return r}function l(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){n(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,a,n=function(e,t){if(null==e)return{};var r,a,n={},o=Object.keys(e);for(a=0;a<o.length;a++)r=o[a],t.indexOf(r)>=0||(n[r]=e[r]);return n}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)r=o[a],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(n[r]=e[r])}return n}var s=a.createContext({}),c=function(e){var t=a.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):l(l({},t),e)),r},p=function(e){var t=c(e.components);return a.createElement(s.Provider,{value:t},e.children)},u="mdxType",g={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},m=a.forwardRef((function(e,t){var r=e.components,n=e.mdxType,o=e.originalType,s=e.parentName,p=i(e,["components","mdxType","originalType","parentName"]),u=c(r),m=n,y=u["".concat(s,".").concat(m)]||u[m]||g[m]||o;return r?a.createElement(y,l(l({ref:t},p),{},{components:r})):a.createElement(y,l({ref:t},p))}));function y(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var o=r.length,l=new Array(o);l[0]=m;var i={};for(var s in t)hasOwnProperty.call(t,s)&&(i[s]=t[s]);i.originalType=e,i[u]="string"==typeof e?e:n,l[1]=i;for(var c=2;c<o;c++)l[c]=r[c];return a.createElement.apply(null,l)}return a.createElement.apply(null,r)}m.displayName="MDXCreateElement"},3597:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>s,contentTitle:()=>l,default:()=>g,frontMatter:()=>o,metadata:()=>i,toc:()=>c});var a=r(8168),n=(r(6540),r(5680));const o={},l="Release Process",i={unversionedId:"contribute/release",id:"contribute/release",title:"Release Process",description:"For maintainers, please read the sections below as a guide to create a new release.",source:"@site/docs/contribute/release.md",sourceDirName:"contribute",slug:"/contribute/release",permalink:"/raccoon/contribute/release",draft:!1,editUrl:"https://github.com/raystack/raccoon/edit/master/docs/docs/contribute/release.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Development Guide",permalink:"/raccoon/contribute/development"}},s={},c=[{value:"Create A New Release",id:"create-a-new-release",level:2},{value:"Important Notes",id:"important-notes",level:2}],p={toc:c},u="wrapper";function g(e){let{components:t,...r}=e;return(0,n.yg)(u,(0,a.A)({},p,r,{components:t,mdxType:"MDXLayout"}),(0,n.yg)("h1",{id:"release-process"},"Release Process"),(0,n.yg)("p",null,"For maintainers, please read the sections below as a guide to create a new release."),(0,n.yg)("h2",{id:"create-a-new-release"},"Create A New Release"),(0,n.yg)("p",null,"Please follow these steps to create a new release:"),(0,n.yg)("ul",null,(0,n.yg)("li",{parentName:"ul"},"create a new tag of the form ",(0,n.yg)("inlineCode",{parentName:"li"},"vM.m.p"),", where:",(0,n.yg)("ul",{parentName:"li"},(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"M")," = Major version, indicates there are breaking changes from the last Major version."),(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"m")," = Minor version, indicates there are backward-compatible changes."),(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"p")," = Patch version, indicates there are backward-compatible bug-fixes.")))),(0,n.yg)("p",null,"For example:"),(0,n.yg)("pre",null,(0,n.yg)("code",{parentName:"pre",className:"language-bash"},"$ git tag v1.2.0\n")),(0,n.yg)("ul",null,(0,n.yg)("li",{parentName:"ul"},"push the tags to trigger a release.")),(0,n.yg)("pre",null,(0,n.yg)("code",{parentName:"pre",className:"language-bash"},"$ git push --tags\n")),(0,n.yg)("p",null," Raccoon uses Goreleaser under the hood for release management. Each release pushes:"),(0,n.yg)("ul",null,(0,n.yg)("li",{parentName:"ul"},"A ",(0,n.yg)("a",{parentName:"li",href:"https://github.com/raystack/raccoon/releases/"},"github release")),(0,n.yg)("li",{parentName:"ul"},"A docker image to ",(0,n.yg)("a",{parentName:"li",href:"https://hub.docker.com/r/raystack/raccoon"},"raystack/raccoon")),(0,n.yg)("li",{parentName:"ul"},"Updates raystack's ",(0,n.yg)("a",{parentName:"li",href:"https://github.com/raystack/homebrew-tap"},"homebrew-tap")),(0,n.yg)("li",{parentName:"ul"},"Updates raystack's ",(0,n.yg)("a",{parentName:"li",href:"https://github.com/raystack/scoop-bucket"},"scoop-bucket"))),(0,n.yg)("p",null,"Additionally, the Github release will also contain with pre-built binaries for:"),(0,n.yg)("ul",null,(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"linux")),(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"darwin")," (macOS)"),(0,n.yg)("li",{parentName:"ul"},(0,n.yg)("inlineCode",{parentName:"li"},"windows"))),(0,n.yg)("h2",{id:"important-notes"},"Important Notes"),(0,n.yg)("ul",null,(0,n.yg)("li",{parentName:"ul"},"Raccoon release tags follow ",(0,n.yg)("a",{parentName:"li",href:"https://semver.org/"},"SEMVER")," convention."),(0,n.yg)("li",{parentName:"ul"},"Github workflow is used to build and push the built docker image to Docker hub."),(0,n.yg)("li",{parentName:"ul"},"A release is triggered when a github tag of format ",(0,n.yg)("inlineCode",{parentName:"li"},"vM.m.p")," is pushed."),(0,n.yg)("li",{parentName:"ul"},"Release tags should only point to main branch")))}g.isMDXComponent=!0}}]);