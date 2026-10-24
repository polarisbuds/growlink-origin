var re=a=>{throw TypeError(a)};var ye=(a,e,t)=>e.has(a)||re("Cannot "+t);var ee=(a,e,t)=>(ye(a,e,"read from private field"),t?t.call(a):e.get(a)),ne=(a,e,t)=>e.has(a)?re("Cannot add the same private member more than once"):e instanceof WeakSet?e.add(a):e.set(a,t);import"../chunks/Bzak7iHL.js";import{b as we,E as ke,n as R,q as ve,y as u,B as W,t as S,g as i,aN as Ne,ak as y,f as g,p as Q,A as c,C as d,k as b,F,o as G,l as U,i as X,j as m,m as be,x as ae,D as le,v as Se,G as Ce,$ as Me}from"../chunks/DCwSHdrR.js";import{B as ze,i as H}from"../chunks/DYKYtpev.js";import{h as Ae,e as se}from"../chunks/DXlHdHv8.js";import{b as Pe}from"../chunks/DgBdcBU0.js";import{p as oe}from"../chunks/DnXnZ1EV.js";import{e as he,i as me}from"../chunks/eicNEEtw.js";import{c as Ee}from"../chunks/D6lpQwG-.js";import{c as D,g as fe,s as te}from"../chunks/DLQUhaIf.js";import{p as Oe,l as C,s as M}from"../chunks/BoBBOnr4.js";import"../chunks/DDdDbyX6.js";import{I as z,s as A,C as Te}from"../chunks/1Cl8V32Q.js";function Le(a,e,...t){var n=new ze(a);we(()=>{const l=e()??null;n.ensure(l,l&&(o=>l(o,...t)))},ke)}const He=!1,je=!1,Ma=Object.freeze(Object.defineProperty({__proto__:null,prerender:je,ssr:He},Symbol.toStringTag,{value:"Module"}));function Be(a,e){const t=a.startsWith("/planning"),n=[{href:"/",label:"ワークスペース",icon:"home",active:a==="/"||a==="/workspace"},{href:"/planning/meetings",label:"プランニング",icon:"calendar-days",active:t}];return e.isSystemAdmin&&n.push({href:"/management",label:"管理",icon:"shield",active:a.startsWith("/management")}),(e.isSystemAdmin||e.isAnalyst)&&n.push({href:"/progress",label:"進捗",icon:"bar-chart-2",active:a.startsWith("/progress")}),n}function Ie(a){return{href:"/setting",label:"設定",icon:"settings",active:a.startsWith("/setting")}}const pe="growlink.sidebar_open";function We(){if(typeof localStorage>"u")return!0;const a=localStorage.getItem(pe);return a!=="false"}var K;class Fe{constructor(){ne(this,K,R(ve(We())))}get sidebarOpen(){return u(ee(this,K))}set sidebarOpen(e){W(ee(this,K),e,!0)}toggleSidebar(){this.sidebarOpen=!this.sidebarOpen,localStorage.setItem(pe,String(this.sidebarOpen))}}K=new WeakMap;const q=new Fe;var Ge=Ne('<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true"><path d="M12 2L4 6v6c0 5 3.5 9.5 8 10 4.5-.5 8-5 8-10V6l-8-4z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"></path><path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg>');function Ue(a,e){let t=Oe(e,"class",3,"h-5 w-5");var n=Ge();S(()=>D(n,0,fe(t()))),i(a,n)}function ie(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8"}],["path",{d:"M3 10a2 2 0 0 1 .709-1.528l7-5.999a2 2 0 0 1 2.582 0l7 5.999A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"}]];z(a,M({name:"house"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function Ve(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M8 2v4"}],["path",{d:"M16 2v4"}],["rect",{width:"18",height:"18",x:"3",y:"4",rx:"2"}],["path",{d:"M3 10h18"}],["path",{d:"M8 14h.01"}],["path",{d:"M12 14h.01"}],["path",{d:"M16 14h.01"}],["path",{d:"M8 18h.01"}],["path",{d:"M12 18h.01"}],["path",{d:"M16 18h.01"}]];z(a,M({name:"calendar-days"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function qe(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"}]];z(a,M({name:"shield"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function De(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["line",{x1:"18",x2:"18",y1:"20",y2:"10"}],["line",{x1:"12",x2:"12",y1:"20",y2:"4"}],["line",{x1:"6",x2:"6",y1:"20",y2:"14"}]];z(a,M({name:"chart-no-axes-column"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function de(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"}],["circle",{cx:"12",cy:"12",r:"3"}]];z(a,M({name:"settings"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}var Re=m('<div><div class="font-display text-sm font-semibold">Growlink</div> <div class="mt-0.5 text-[10px] text-base-content/50">1on1 ミーティング</div></div>'),ce=m("<span> </span>"),Ke=m("<a><!> <!></a>"),Ye=m('<aside><div class="flex h-16 shrink-0 items-center gap-3 border-b border-base-300 px-5"><div class="grid h-8 w-8 place-items-center rounded-xl bg-primary text-primary-content"><!></div> <!></div> <nav class="flex-1 min-h-0 space-y-1 overflow-y-auto p-3"></nav> <div class="mt-auto shrink-0 border-t border-base-300 p-3"><a><!> <!></a></div></aside>');function Je(a,e){Q(e,!0);const t={home:ie,"calendar-days":Ve,shield:qe,"bar-chart-2":De,settings:de};let n=F(()=>Be(e.pathname,e.me)),l=F(()=>Ie(e.pathname));function o(h){return h?"bg-primary/10 text-primary":"text-base-content/60 hover:bg-base-200 hover:text-base-content"}var r=Ye(),s=c(r),v=c(s),w=c(v);Ue(w,{}),d(v);var f=b(v,2);{var _=h=>{var x=Re();i(h,x)};H(f,h=>{q.sidebarOpen&&h(_)})}d(s);var k=b(s,2);he(k,21,()=>u(n),me,(h,x)=>{const B=F(()=>t[u(x).icon]??ie);var P=Ke(),V=c(P);Ee(V,()=>u(B),(N,I)=>{I(N,{class:"h-5 w-5 shrink-0"})});var Z=b(V,2);{var Y=N=>{var I=ce(),p=G(I,!0);S(()=>U(p,u(x).label)),i(N,I)};H(Z,N=>{q.sidebarOpen&&N(Y)})}d(P),S(N=>{te(P,"href",u(x).href),D(P,1,`flex h-11 items-center gap-3 rounded-xl px-3 text-sm font-medium transition ${N??""}`)},[()=>o(u(x).active)]),i(h,P)}),d(k);var E=b(k,2),$=c(E),O=c($);de(O,{class:"h-5 w-5 shrink-0"});var L=b(O,2);{var j=h=>{var x=ce(),B=G(x,!0);S(()=>U(B,u(l).label)),i(h,x)};H(L,h=>{q.sidebarOpen&&h(j)})}d($),d(E),d(r),S(h=>{D(r,1,`hidden shrink-0 border-r border-base-300 bg-base-100 transition-[width] duration-200 md:flex md:flex-col md:min-h-0 ${q.sidebarOpen?"w-64":"w-20"}`),te($,"href",u(l).href),D($,1,`flex h-11 items-center gap-3 rounded-xl px-3 text-sm font-medium transition ${h??""}`)},[()=>o(u(l).active)]),i(a,r),X()}const ge="growlink-theme",ue="growlink",Qe=[{value:"growlink",label:"Growlink"},{value:"wine",label:"Wine"},{value:"light",label:"Light"},{value:"dark",label:"Dark"},{value:"cupcake",label:"Cupcake"},{value:"bumblebee",label:"Bumblebee"},{value:"emerald",label:"Emerald"},{value:"corporate",label:"Corporate"},{value:"synthwave",label:"Synthwave"},{value:"retro",label:"Retro"},{value:"cyberpunk",label:"Cyberpunk"},{value:"valentine",label:"Valentine"},{value:"halloween",label:"Halloween"},{value:"garden",label:"Garden"},{value:"forest",label:"Forest"},{value:"aqua",label:"Aqua"},{value:"lofi",label:"Lofi"},{value:"pastel",label:"Pastel"},{value:"fantasy",label:"Fantasy"},{value:"wireframe",label:"Wireframe"},{value:"black",label:"Black"},{value:"luxury",label:"Luxury"},{value:"dracula",label:"Dracula"},{value:"cmyk",label:"Cmyk"},{value:"autumn",label:"Autumn"},{value:"business",label:"Business"},{value:"acid",label:"Acid"},{value:"lemonade",label:"Lemonade"},{value:"night",label:"Night"},{value:"coffee",label:"Coffee"},{value:"winter",label:"Winter"},{value:"dim",label:"Dim"},{value:"nord",label:"Nord"},{value:"sunset",label:"Sunset"},{value:"caramellatte",label:"Caramellatte"},{value:"abyss",label:"Abyss"},{value:"silk",label:"Silk"}];function Xe(a){document.documentElement.setAttribute("data-theme",a),localStorage.setItem(ge,a)}function Ze(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["circle",{cx:"13.5",cy:"6.5",r:".5",fill:"currentColor"}],["circle",{cx:"17.5",cy:"10.5",r:".5",fill:"currentColor"}],["circle",{cx:"8.5",cy:"7.5",r:".5",fill:"currentColor"}],["circle",{cx:"6.5",cy:"12.5",r:".5",fill:"currentColor"}],["path",{d:"M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z"}]];z(a,M({name:"palette"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function _e(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"m6 9 6 6 6-6"}]];z(a,M({name:"chevron-down"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}var ea=m('<li><button type="button"> </button></li>'),aa=m('<ul tabindex="0" class="menu dropdown-content z-50 mt-2 max-h-72 w-52 overflow-y-auto rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"></ul>'),ta=m('<div class="dropdown dropdown-end"><button type="button" tabindex="0" class="btn btn-ghost btn-sm gap-1" aria-label="テーマ"><!> <!></button> <!></div>');function ra(a,e){Q(e,!0);let t=R(!1),n=R(ve(typeof localStorage<"u"&&localStorage.getItem(ge)||ue));function l(_){W(n,_,!0),Xe(_),W(t,!1)}var o=ta(),r=c(o),s=c(r);Ze(s,{class:"h-4 w-4"});var v=b(s,2);_e(v,{class:"h-3 w-3 opacity-60"}),d(r);var w=b(r,2);{var f=_=>{var k=aa();he(k,21,()=>Qe,me,(E,$)=>{var O=ea(),L=c(O),j=G(L,!0);d(O),S(()=>{D(L,1,fe(u(n)===u($).value?"active":"")),U(j,u($).label)}),ae("click",L,()=>l(u($).value)),i(E,O)}),d(k),i(_,k)};H(w,_=>{u(t)&&_(f)})}d(o),ae("click",r,()=>W(t,!u(t))),i(a,o),X()}be(["click"]);function na(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["rect",{width:"18",height:"18",x:"3",y:"3",rx:"2"}],["path",{d:"M9 3v18"}]];z(a,M({name:"panel-left"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function la(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"}],["circle",{cx:"12",cy:"7",r:"4"}]];z(a,M({name:"user"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}function sa(a,e){const t=C(e,["children","$$slots","$$events","$$legacy"]);/**
 * @license lucide-svelte v0.487.0 - ISC
 *
 * ISC License
 *
 * Copyright (c) for portions of Lucide are held by Cole Bemis 2013-2022 as part of Feather (MIT). All other copyright (c) for Lucide are held by Lucide Contributors 2022.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */const n=[["path",{d:"M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"}],["polyline",{points:"16 17 21 12 16 7"}],["line",{x1:"21",x2:"9",y1:"12",y2:"12"}]];z(a,M({name:"log-out"},()=>t,{get iconNode(){return n},children:(l,o)=>{var r=y(),s=g(r);A(s,e,"default",{}),i(l,r)},$$slots:{default:!0}}))}var oa=m('<!> <span class="font-medium text-base-content"> </span>',1),ia=m('<img alt="" class="h-7 w-7 rounded-full object-cover"/>'),da=m('<span class="grid h-7 w-7 place-items-center rounded-full bg-primary text-xs font-semibold text-primary-content"> </span>'),ca=m('<header class="sticky top-0 z-30 flex h-16 items-center gap-3 border-b border-base-300 bg-base-100 px-4 md:px-6"><button type="button" class="btn btn-ghost btn-sm btn-square border border-base-300" aria-label="サイドバーを開閉"><!></button> <nav class="flex items-center gap-2 text-sm text-base-content/60"><a href="/" class="hover:text-base-content">ホーム</a> <!></nav> <div class="ml-auto flex items-center gap-2"><!> <div class="dropdown dropdown-end"><button type="button" tabindex="0" class="btn btn-ghost btn-sm gap-2"><!> <span class="hidden sm:inline"> </span> <!></button> <ul tabindex="0" class="menu dropdown-content z-50 mt-2 w-52 rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"><li><a href="/setting"><!> 設定</a></li> <li><a href="/logout" data-sveltekit-reload=""><!> ログアウト</a></li></ul></div></div></header>');function ua(a,e){Q(e,!0);let t=F(()=>e.me.imageUri?`/uploads/${e.me.imageUri}`:null),n=F(()=>(e.me.name||"?").slice(0,1));var l=ca(),o=c(l),r=c(o);na(r,{class:"h-5 w-5"}),d(o);var s=b(o,2),v=b(c(s),2);{var w=p=>{var T=oa(),J=g(T);Te(J,{class:"h-4 w-4"});var $e=b(J,2),xe=G($e,!0);S(()=>U(xe,e.breadcrumb)),i(p,T)};H(v,p=>{e.breadcrumb&&p(w)})}d(s);var f=b(s,2),_=c(f);ra(_,{});var k=b(_,2),E=c(k),$=c(E);{var O=p=>{var T=ia();S(()=>te(T,"src",u(t))),i(p,T)},L=p=>{var T=da(),J=G(T,!0);S(()=>U(J,u(n))),i(p,T)};H($,p=>{u(t)?p(O):p(L,-1)})}var j=b($,2),h=G(j,!0),x=b(j,2);_e(x,{class:"h-4 w-4 opacity-60"}),d(E);var B=b(E,2),P=c(B),V=c(P),Z=c(V);la(Z,{class:"h-4 w-4 inline"}),le(),d(V),d(P);var Y=b(P,2),N=c(Y),I=c(N);sa(I,{class:"h-4 w-4 inline"}),le(),d(N),d(Y),d(B),d(k),d(f),d(l),S(()=>U(h,e.me.name)),ae("click",o,()=>q.toggleSidebar()),i(a,l),X()}be(["click"]);var va=m('<div class="flex h-screen overflow-hidden"><!> <div class="flex min-w-0 flex-1 min-h-0 flex-col"><!> <main class="flex-1 min-h-0 overflow-y-auto p-4 md:p-8"><!></main></div></div>');function ba(a,e){var t=va(),n=c(t);Je(n,{get me(){return e.me},get pathname(){return e.pathname}});var l=b(n,2),o=c(l);ua(o,{get me(){return e.me},get breadcrumb(){return e.breadcrumb}});var r=b(o,2),s=c(r);Le(s,()=>e.children),d(r),d(l),d(t),i(a,t)}var ha=m('<div class="flex min-h-screen items-center justify-center bg-base-200 text-sm text-base-content/60">読み込み中…</div>');function za(a,e){Q(e,!0);let t=R(null),n=R(!0);Se(async()=>{W(t,await se(),!0),W(n,!1)}),Pe(async v=>{var f;if(((f=v.to)==null?void 0:f.url.pathname)==="/logout")return;await se()||v.cancel()});var l=y();Ae("12qhfyh",v=>{Ce(()=>{Me.title="Growlink"})});var o=g(l);{var r=v=>{var w=ha();i(v,w)},s=v=>{{let w=F(()=>{var f;return(f=oe.data)==null?void 0:f.breadcrumb});ba(v,{get me(){return u(t)},get pathname(){return oe.url.pathname},get breadcrumb(){return u(w)},get children(){return e.children}})}};H(o,v=>{u(n)?v(r):u(t)&&v(s,1)})}i(a,l),X()}export{za as component,Ma as universal};
