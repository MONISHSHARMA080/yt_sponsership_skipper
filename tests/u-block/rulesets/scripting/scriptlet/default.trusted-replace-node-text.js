/*******************************************************************************

    uBlock Origin Lite - a comprehensive, MV3-compliant content blocker
    Copyright (C) 2014-present Raymond Hill

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see {http://www.gnu.org/licenses/}.

    Home: https://github.com/gorhill/uBlock

*/

// ruleset: default

// Important!
// Isolate from global scope

// Start of local scope
(function uBOL_replaceNodeText() {

/******************************************************************************/

function replaceNodeText(
    nodeName,
    pattern,
    replacement,
    ...extraArgs
) {
    replaceNodeTextFn(nodeName, pattern, replacement, ...extraArgs);
}

function replaceNodeTextFn(
    nodeName = '',
    pattern = '',
    replacement = ''
) {
    const safe = safeSelf();
    const logPrefix = safe.makeLogPrefix('replace-node-text.fn', ...Array.from(arguments));
    const reNodeName = safe.patternToRegex(nodeName, 'i', true);
    const rePattern = safe.patternToRegex(pattern, 'gms');
    const extraArgs = safe.getExtraArgs(Array.from(arguments), 3);
    const reIncludes = extraArgs.includes || extraArgs.condition
        ? safe.patternToRegex(extraArgs.includes || extraArgs.condition, 'ms')
        : null;
    const reExcludes = extraArgs.excludes
        ? safe.patternToRegex(extraArgs.excludes, 'ms')
        : null;
    const stop = (takeRecord = true) => {
        if ( takeRecord ) {
            handleMutations(observer.takeRecords());
        }
        observer.disconnect();
        if ( safe.logLevel > 1 ) {
            safe.uboLog(logPrefix, 'Quitting');
        }
    };
    const textContentFactory = (( ) => {
        const out = { createScript: s => s };
        const { trustedTypes: tt } = self;
        if ( tt instanceof Object ) {
            if ( typeof tt.getPropertyType === 'function' ) {
                if ( tt.getPropertyType('script', 'textContent') === 'TrustedScript' ) {
                    return tt.createPolicy(getRandomTokenFn(), out);
                }
            }
        }
        return out;
    })();
    let sedCount = extraArgs.sedCount || 0;
    const handleNode = node => {
        const before = node.textContent;
        if ( reIncludes ) {
            reIncludes.lastIndex = 0;
            if ( safe.RegExp_test.call(reIncludes, before) === false ) { return true; }
        }
        if ( reExcludes ) {
            reExcludes.lastIndex = 0;
            if ( safe.RegExp_test.call(reExcludes, before) ) { return true; }
        }
        rePattern.lastIndex = 0;
        if ( safe.RegExp_test.call(rePattern, before) === false ) { return true; }
        rePattern.lastIndex = 0;
        const after = pattern !== ''
            ? before.replace(rePattern, replacement)
            : replacement;
        node.textContent = node.nodeName === 'SCRIPT'
            ? textContentFactory.createScript(after)
            : after;
        if ( safe.logLevel > 1 ) {
            safe.uboLog(logPrefix, `Text before:\n${before.trim()}`);
        }
        safe.uboLog(logPrefix, `Text after:\n${after.trim()}`);
        return sedCount === 0 || (sedCount -= 1) !== 0;
    };
    const handleMutations = mutations => {
        for ( const mutation of mutations ) {
            for ( const node of mutation.addedNodes ) {
                if ( reNodeName.test(node.nodeName) === false ) { continue; }
                if ( handleNode(node) ) { continue; }
                stop(false); return;
            }
        }
    };
    const observer = new MutationObserver(handleMutations);
    observer.observe(document, { childList: true, subtree: true });
    if ( document.documentElement ) {
        const treeWalker = document.createTreeWalker(
            document.documentElement,
            NodeFilter.SHOW_ELEMENT | NodeFilter.SHOW_TEXT
        );
        let count = 0;
        for (;;) {
            const node = treeWalker.nextNode();
            count += 1;
            if ( node === null ) { break; }
            if ( reNodeName.test(node.nodeName) === false ) { continue; }
            if ( node === document.currentScript ) { continue; }
            if ( handleNode(node) ) { continue; }
            stop(); break;
        }
        safe.uboLog(logPrefix, `${count} nodes present before installing mutation observer`);
    }
    if ( extraArgs.stay ) { return; }
    runAt(( ) => {
        const quitAfter = extraArgs.quitAfter || 0;
        if ( quitAfter !== 0 ) {
            setTimeout(( ) => { stop(); }, quitAfter);
        } else {
            stop();
        }
    }, 'interactive');
}

function getRandomTokenFn() {
    const safe = safeSelf();
    return safe.String_fromCharCode(Date.now() % 26 + 97) +
        safe.Math_floor(safe.Math_random() * 982451653 + 982451653).toString(36);
}

function runAt(fn, when) {
    const intFromReadyState = state => {
        const targets = {
            'loading': 1, 'asap': 1,
            'interactive': 2, 'end': 2, '2': 2,
            'complete': 3, 'idle': 3, '3': 3,
        };
        const tokens = Array.isArray(state) ? state : [ state ];
        for ( const token of tokens ) {
            const prop = `${token}`;
            if ( Object.hasOwn(targets, prop) === false ) { continue; }
            return targets[prop];
        }
        return 0;
    };
    const runAt = intFromReadyState(when);
    if ( intFromReadyState(document.readyState) >= runAt ) {
        fn(); return;
    }
    const onStateChange = ( ) => {
        if ( intFromReadyState(document.readyState) < runAt ) { return; }
        fn();
        safe.removeEventListener.apply(document, args);
    };
    const safe = safeSelf();
    const args = [ 'readystatechange', onStateChange, { capture: true } ];
    safe.addEventListener.apply(document, args);
}

function safeSelf() {
    if ( scriptletGlobals.safeSelf ) {
        return scriptletGlobals.safeSelf;
    }
    const self = globalThis;
    const safe = {
        'Array_from': Array.from,
        'Error': self.Error,
        'Function_toStringFn': self.Function.prototype.toString,
        'Function_toString': thisArg => safe.Function_toStringFn.call(thisArg),
        'Math_floor': Math.floor,
        'Math_max': Math.max,
        'Math_min': Math.min,
        'Math_random': Math.random,
        'Object': Object,
        'Object_defineProperty': Object.defineProperty.bind(Object),
        'Object_defineProperties': Object.defineProperties.bind(Object),
        'Object_fromEntries': Object.fromEntries.bind(Object),
        'Object_getOwnPropertyDescriptor': Object.getOwnPropertyDescriptor.bind(Object),
        'Object_hasOwn': Object.hasOwn.bind(Object),
        'RegExp': self.RegExp,
        'RegExp_test': self.RegExp.prototype.test,
        'RegExp_exec': self.RegExp.prototype.exec,
        'Request_clone': self.Request.prototype.clone,
        'String': self.String,
        'String_fromCharCode': String.fromCharCode,
        'String_split': String.prototype.split,
        'XMLHttpRequest': self.XMLHttpRequest,
        'addEventListener': self.EventTarget.prototype.addEventListener,
        'removeEventListener': self.EventTarget.prototype.removeEventListener,
        'fetch': self.fetch,
        'JSON': self.JSON,
        'JSON_parseFn': self.JSON.parse,
        'JSON_stringifyFn': self.JSON.stringify,
        'JSON_parse': (...args) => safe.JSON_parseFn.call(safe.JSON, ...args),
        'JSON_stringify': (...args) => safe.JSON_stringifyFn.call(safe.JSON, ...args),
        'log': console.log.bind(console),
        // Properties
        logLevel: 0,
        // Methods
        makeLogPrefix(...args) {
            return this.sendToLogger && `[${args.join(' \u205D ')}]` || '';
        },
        uboLog(...args) {
            if ( this.sendToLogger === undefined ) { return; }
            if ( args === undefined || args[0] === '' ) { return; }
            return this.sendToLogger('info', ...args);
            
        },
        uboErr(...args) {
            if ( this.sendToLogger === undefined ) { return; }
            if ( args === undefined || args[0] === '' ) { return; }
            return this.sendToLogger('error', ...args);
        },
        escapeRegexChars(s) {
            return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        },
        initPattern(pattern, options = {}) {
            if ( pattern === '' ) {
                return { matchAll: true, expect: true };
            }
            const expect = (options.canNegate !== true || pattern.startsWith('!') === false);
            if ( expect === false ) {
                pattern = pattern.slice(1);
            }
            const match = /^\/(.+)\/([gimsu]*)$/.exec(pattern);
            if ( match !== null ) {
                return {
                    re: new this.RegExp(
                        match[1],
                        match[2] || options.flags
                    ),
                    expect,
                };
            }
            if ( options.flags !== undefined ) {
                return {
                    re: new this.RegExp(this.escapeRegexChars(pattern),
                        options.flags
                    ),
                    expect,
                };
            }
            return { pattern, expect };
        },
        testPattern(details, haystack) {
            if ( details.matchAll ) { return true; }
            if ( details.re ) {
                return this.RegExp_test.call(details.re, haystack) === details.expect;
            }
            return haystack.includes(details.pattern) === details.expect;
        },
        patternToRegex(pattern, flags = undefined, verbatim = false) {
            if ( pattern === '' ) { return /^/; }
            const match = /^\/(.+)\/([gimsu]*)$/.exec(pattern);
            if ( match === null ) {
                const reStr = this.escapeRegexChars(pattern);
                return new RegExp(verbatim ? `^${reStr}$` : reStr, flags);
            }
            try {
                return new RegExp(match[1], match[2] || undefined);
            }
            catch {
            }
            return /^/;
        },
        getExtraArgs(args, offset = 0) {
            const entries = args.slice(offset).reduce((out, v, i, a) => {
                if ( (i & 1) === 0 ) {
                    const rawValue = a[i+1];
                    const value = /^\d+$/.test(rawValue)
                        ? parseInt(rawValue, 10)
                        : rawValue;
                    out.push([ a[i], value ]);
                }
                return out;
            }, []);
            return this.Object_fromEntries(entries);
        },
        onIdle(fn, options) {
            if ( self.requestIdleCallback ) {
                return self.requestIdleCallback(fn, options);
            }
            return self.requestAnimationFrame(fn);
        },
        offIdle(id) {
            if ( self.requestIdleCallback ) {
                return self.cancelIdleCallback(id);
            }
            return self.cancelAnimationFrame(id);
        }
    };
    scriptletGlobals.safeSelf = safe;
    if ( scriptletGlobals.bcSecret === undefined ) { return safe; }
    // This is executed only when the logger is opened
    safe.logLevel = scriptletGlobals.logLevel || 1;
    let lastLogType = '';
    let lastLogText = '';
    let lastLogTime = 0;
    safe.toLogText = (type, ...args) => {
        if ( args.length === 0 ) { return; }
        const text = `[${document.location.hostname || document.location.href}]${args.join(' ')}`;
        if ( text === lastLogText && type === lastLogType ) {
            if ( (Date.now() - lastLogTime) < 5000 ) { return; }
        }
        lastLogType = type;
        lastLogText = text;
        lastLogTime = Date.now();
        return text;
    };
    try {
        const bc = new self.BroadcastChannel(scriptletGlobals.bcSecret);
        let bcBuffer = [];
        safe.sendToLogger = (type, ...args) => {
            const text = safe.toLogText(type, ...args);
            if ( text === undefined ) { return; }
            if ( bcBuffer === undefined ) {
                return bc.postMessage({ what: 'messageToLogger', type, text });
            }
            bcBuffer.push({ type, text });
        };
        bc.onmessage = ev => {
            const msg = ev.data;
            switch ( msg ) {
            case 'iamready!':
                if ( bcBuffer === undefined ) { break; }
                bcBuffer.forEach(({ type, text }) =>
                    bc.postMessage({ what: 'messageToLogger', type, text })
                );
                bcBuffer = undefined;
                break;
            case 'setScriptletLogLevelToOne':
                safe.logLevel = 1;
                break;
            case 'setScriptletLogLevelToTwo':
                safe.logLevel = 2;
                break;
            }
        };
        bc.postMessage('areyouready?');
    } catch {
        safe.sendToLogger = (type, ...args) => {
            const text = safe.toLogText(type, ...args);
            if ( text === undefined ) { return; }
            safe.log(`uBO ${text}`);
        };
    }
    return safe;
}

/******************************************************************************/

const scriptletGlobals = {}; // eslint-disable-line
const argsList = [["script","\"adBlockWallEnabled\":true","\"adBlockWallEnabled\":false"],["script","/vastURL:.*?',/","vastURL: '',"],["script","/url:.*?',/","url: '',"],["script","a.onerror","xxx"],["script","popunder","p"],["script","/\\$.*embed.*.appendTo.*;/","","condition","appendTo"],["script","setInterval"],["script","/'globalConfig':.*?\",\\s};var exportz/s","};var exportz"],["script","/\\\"homad\\\",/"],["script","/\\\"homad\\\":\\{\\\"state\\\":\\\"enabled\\\"\\}/","\"homad\":{\"state\":\"disabled\"}"],["script","useAdBlockDefend: true","useAdBlockDefend: false"],["script","/5000|100/gm","10","condition","set"],["script","/^([^{])/","document.addEventListener('DOMContentLoaded',()=>{const i=document.createElement('iframe');i.style='height:0;width:0;border:0';i.id='aswift_0';document.body.appendChild(i);i.focus();const f=document.createElement('div');f.id='9JJFp';document.body.appendChild(f);});$1","sedCount","2"],["script","/false;/gm","true;","condition","isSubscribed"],["script","('t_modal')","('go_d2')"],["script","/9000|30000|1000|40000/gm","10"],["script","/1000|100|6|30/gm","1","condition","link"],["script","/^window\\.location\\.href.*\\'$/gms"],["script","getlink.removeClass('hidden');","gotolink.removeClass('hidden');"],["script","/initialTimeSeconds = \\d+/","initialTimeSeconds = 7"],["script","= false;","= true;","condition","innerHTML"],["script","'IFRAME'","'BODY'"],["script","/window\\.location.*?;/"],["script","typeof cdo == 'undefined' || document.querySelector('div.textads.banner-ads.banner_ads.ad-unit.ad-zone.ad-space.adsbox') == undefined","false"],["script","/window\\.location\\.href='.*';/","","condition","openLink"],["script","var seconde = 10;","var seconde = 0;"],["script","timeLeft = duration","timeLeft = 1"],["script","return!![]","return![]"],["script",";return;","","condition","_0x"],["script","/return Array[^;]+/","return true"],["script","function sutounlock_0x","//function sutounlock_0x"],["script","/\\d{4}/gm","10","condition","count"],["script","/getElementById\\('.*'\\).*'block';/gm","getElementById('btn6').style.display = 'block';","condition","count"],["script","3000)","10)"],["script","isadblock = 1;","isadblock = 0;"],["script","\"#sdl\"","\"#dln\""],["script","event.message);","event.message); /*start*/ !function(){\"use strict\";let t={log:window.console.log.bind(console),getPropertyValue:CSSStyleDeclaration.prototype.getPropertyValue,setAttribute:Element.prototype.setAttribute,getAttribute:Element.prototype.getAttribute,appendChild:Element.prototype.appendChild,remove:Element.prototype.remove,cloneNode:Element.prototype.cloneNode,Element_attributes:Object.getOwnPropertyDescriptor(Element.prototype,\"attributes\").get,Array_splice:Array.prototype.splice,Array_join:Array.prototype.join,createElement:document.createElement,getComputedStyle:window.getComputedStyle,Reflect:Reflect,Proxy:Proxy,crypto:window.crypto,Uint8Array:Uint8Array,Object_defineProperty:Object.defineProperty.bind(Object),Object_getOwnPropertyDescriptor:Object.getOwnPropertyDescriptor.bind(Object),String_replace:String.prototype.replace},e=t.crypto.getRandomValues.bind(t.crypto),r=function(e,r,l){return\"toString\"===r?e.toString.bind(e):t.Reflect.get(e,r,l)},l=function(r){let o=function(t){return t.toString(16).padStart(2,\"0\")},p=new t.Uint8Array((r||40)/2);e(p);let n=t.String_replace.call(t.Array_join.call(Array.from(p,o),\"\"),/^\\d+/g,\"\");return n.length<3?l(r):n},o=l(15);window.MutationObserver=new t.Proxy(window.MutationObserver,{construct:function(e,r){let l=r[0],p=function(e,r){for(let p=e.length,n=p-1;n>=0;--n){let c=e[n];if(\"childList\"===c.type&&c.addedNodes.length>0){let i=c.addedNodes;for(let a=0,y=i.length;a<y;++a){let u=i[a];if(u.localName===o){t.Array_splice.call(e,n,1);break}}}}0!==e.length&&l(e,r)};r[0]=p;let n=t.Reflect.construct(e,r);return n},get:r}),window.getComputedStyle=new t.Proxy(window.getComputedStyle,{apply(e,l,p){let n=t.Reflect.apply(e,l,p);if(\"none\"===t.getPropertyValue.call(n,\"clip-path\"))return n;let c=p[0],i=t.createElement.call(document,o);t.setAttribute.call(i,\"class\",t.getAttribute.call(c,\"class\")),t.setAttribute.call(i,\"id\",t.getAttribute.call(c,\"id\")),t.setAttribute.call(i,\"style\",t.getAttribute.call(c,\"style\")),t.appendChild.call(document.body,i);let a=t.getPropertyValue.call(t.getComputedStyle.call(window,i),\"clip-path\");return t.remove.call(i),t.Object_defineProperty(n,\"clipPath\",{get:(()=>a).bind(null)}),n.getPropertyValue=new t.Proxy(n.getPropertyValue,{apply:(e,r,l)=>\"clip-path\"!==l[0]?t.Reflect.apply(e,r,l):a,get:r}),n},get:r})}(); document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/"],["script","/\\.cloudfront\\.net|window\\.open/g","false"],["script","/element\\.contains\\(document\\.activeElement\\)|document\\.hidden && !timeCounted/g","true"],["script","!seen && ad","false"],["script","window.location.href","temp","includes","startDownloads"],["script","window.location.href","temp","includes","linkToOpen"],["script","(isAdsenseBlocked)","(false)"],["script","/.*adConfig.*frequency_period.*/","(async () => {const a=location.href;if(!a.includes(\"/download?link=\"))return;const b=new URL(a),c=b.searchParams.get(\"link\");try{location.assign(`${location.protocol}//${c}`)}catch(a){}} )();"],["script","break;case $."],["script","window.dataLayer =","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/common/css/etoday.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/window.dataLayer ="],["script","window.dataLayer =","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/css_renew/pc/common.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/window.dataLayer ="],["script","_paq.push","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/css/pc/ecn_common.min.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/_paq.push"],["script","window.dataLayer =","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/wp-content/themes/hts_v2/style.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/window.dataLayer ="],["script","window.dataLayer =","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/_css/css.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/window.dataLayer ="],["script","document.getElementById(\"divRecentQuotes\")","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/css/so-clean.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/document.getElementById(\"divRecentQuotes\")"],["script","var _paq =","/*start*/(function(){let link=document.createElement(\"link\");link.rel=\"stylesheet\";link.href=\"/Content/css/style.css\";document.head.appendChild(link)})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/var _paq =","sedCount","1"],["script","var localize =","/*start*/(function(){document.querySelectorAll(\"script[wp-data]\").forEach(element=>{const html=new DOMParser().parseFromString(atob(element.getAttribute(\"wp-data\")),\"text/html\");html.querySelectorAll(\"link:not([as])\").forEach(linkEl=>{element.after(linkEl)});element.parentElement.removeChild(element);})})();document.currentScript.textContent=document.currentScript.textContent.replace(/\\/\\*start\\*\\/(.*)\\/\\*end\\*\\//g,\"\");/*end*/var localize =","sedCount","1"],["script","/^.+/gms","!function(){var e=Object.getOwnPropertyDescriptor(Element.prototype,\"innerHTML\").set;Object.defineProperty(Element.prototype,\"innerHTML\",{set:function(t){return t.includes(\"html-load.com\")?e.call(this,\"\"):e.call(this,t)}})}();","condition","html-load.com"],["script","vastTag","v"],["script","/protect_block.*?,/"],["script","/window.open.*/gms","window.open(url, \"_self\");}"],["script","/window\\.location\\.href.*?;/"],["script","/\\(window\\.show[^\\)]+\\)/","(true)","condition","classList.add"],["script","(isAdblock)","(false)"],["style","visibility: visible !important;","display: none !important;"],["script","currentTime = 1500 * 2","currentTime = 0"],["script","added=false;","added=false;if (typeof localStorage !== 'undefined' && typeof JSON.parse(localStorage._ppp)['0_uid'] === 'undefined' && window.visualViewport.width != document.documentElement.offsetWidth+320) {const originalvisualViewport=window.visualViewport; Object.defineProperty(window, 'visualViewport', {value: new Proxy(originalvisualViewport, {get(target,property) {if (property === 'width') {return document.documentElement.offsetWidth+320} return target[property]}}), configurable:true});}"],["script","/^.+/s","navigator.serviceWorker.getRegistrations().then((registrations=>{for(const registration of registrations){if(registration.scope.includes(\"streamingcommunity.computer\")){registration.unregister()}}}));","condition","swDidInit"],["script","\"adsDisabled\":false","\"adsDisabled\":true"],["style","::after{content:\" \";display:table;box-sizing:border-box}","{display: none !important;}","condition","text-decoration:none;vertical-align:middle"],["script","?30:0","?0:0"],["script","/openNewTab\\(\".*?\"\\)/g","null"],["script","window.dataLayer =","window.HTMLAnchorElement.prototype.click=new Proxy(window.HTMLAnchorElement.prototype.click,{apply:(target,thisArg,args)=>{if(thisArg&&!thisArg.href.includes(\"streamsilk.com\"))return;return Reflect.apply(target,thisArg,args)}});window.dataLayer =","sedCount","1"],["script","({});","({}); function showHideElements(t,e){$(t).hide(),$(e).show()}function disableBtnclc(){let t=document.querySelector(\".submit-captcha\");t.disabled=!0,t.innerHTML=\"Loading...\"}function refreshButton(){$(\".refresh-capthca-btn\").addClass(\"disabled\")}function copyInput(){let t=document.querySelectorAll(\".copy-input\");t.forEach(t=>{navigator.clipboard.writeText(t.value)}),Materialize.toast(\"Copied!\",2e3)}function imgOnError(){$(\".ua-check\").html(window.atob(\"PGRpdiBjbGFzcz0idGV4dC1kYW5nZXIgZm9udC13ZWlnaHQtYm9sZCBoNSBtdC0xIj5DYXB0Y2hhIGltYWdlIGZhaWxlZCB0byBsb2FkLjxicj48YSBvbmNsaWNrPSJsb2NhdGlvbi5yZWxvYWQoKSIgc3R5bGU9ImNvbG9yOiM2MjcwZGE7Y3Vyc29yOnBvaW50ZXIiIGNsYXNzPSJ0ZXh0LWRlY29yYXRpb25lLW5vbmUiPlBsZWFzZSByZWZyZXNoIHRoZSBwYWdlLiA8aSBjbGFzcz0iZmEgZmEtcmVmcmVzaCI+PC9pPjwvYT48L2Rpdj4=\"))}$(window).on(\"load\",function(){$(\"body\").addClass(\"loaded\")}),window.history.replaceState&&window.history.replaceState(null,null,window.location.href),$(\".remove-spaces\").on(\"input\",function(){this.value=this.value.replace(/\\s/g,\"\")}),$(document).on(\"click\",\"#toast-container .toast\",function(){$(this).fadeOut(function(){$(this).remove()})}),$(\".tktemizle\").on(\"input propertychange\",function(){let t=$(this).val().match(\"access_token=(.*?)&\");t&&$(\".tktemizle\").val(t[1])}),$(document).ready(function(){let t=[{button:$(\".t-followers-button\"),menu:$(\".t-followers-menu\")},{button:$(\".t-hearts-button\"),menu:$(\".t-hearts-menu\")},{button:$(\".t-chearts-button\"),menu:$(\".t-chearts-menu\")},{button:$(\".t-views-button\"),menu:$(\".t-views-menu\")},{button:$(\".t-shares-button\"),menu:$(\".t-shares-menu\")},{button:$(\".t-favorites-button\"),menu:$(\".t-favorites-menu\")},{button:$(\".t-livestream-button\"),menu:$(\".t-livestream-menu\")},{button:$(\".ig-followers-button\"),menu:$(\".ig-followers-menu\")},{button:$(\".ig-likes-button\"),menu:$(\".ig-likes-menu\")}];$.each(t,function(t,e){e.button.click(function(){$(\".colsmenu\").addClass(\"nonec\"),e.menu.removeClass(\"nonec\")})})});"],["script","/devtoolsDetector\\.launch\\(\\)\\;/"],["script","//$('#btn_download').click();","$('#btn_download').click();","sedCount","1"],["script","/reymit_ads_for_categories\\.length>0|reymit_ads_for_streams\\.length>0/g","false"],["script","/data: \\[.*\\],/","data: [],","condition","ads_num"],["script","/try.*finally.*?}/"],["script","/for\\s*\\(\\s*(const|let|var).*?\\)\\;return\\;\\}_/g","_","condition","attribute"],["script","adv_","","condition","flashvars"],["script","typeof window.loadtwo","typeof window.loadtwo === true"],["script","(window.loadtwo","(window.loadtwo === undefined || true || window.googlescriptloaded"],["script","/^/","(()=>{window.admiral=function(d,a,b){if(\"function\"==typeof b)try{b({})}catch(a){}}})();","sedCount","1"],["script","/if \\(api && url\\).+/s","window.location.href = url","condition","quick-link"],["script","= getSetTimeout()","= function newTimeout(func, timer) {func()}"],["script","IFRAME","BODY"],["script","(hasBlocker)","(false)"],["P","/\\.[^.]+(1Password password manager|download 1Password)[^.]+/"],["script","startTime: '5'","startTime: '0'"],["script","/(function downloadHD\\(obj\\) {)[\\s\\S]*?(datahref.*)[\\s\\S]*?(window.location.href = datahref;)[\\s\\S]*/","$1$2$3}"],["script","buton.setAttribute","location.href=urldes;buton.setAttribute"],["script","clickCount === numberOfAdsBeforeCopy","numberOfAdsBeforeCopy >= clickCount"],["script","video_urls.length != activeItem+1","video_urls.length === activeItem+1"],["script","/if\\(.&&.\\.target\\)/","if(false)"],["script","(document.hasFocus())","(false)"],["script","(isNoAds)","(true)"],["script","popunder","","condition","popunder","stay","1"],["script","\"}};","\"}}; jQuery(document).ready(function(t){let e=document.createElement(\"link\");e.setAttribute(\"rel\",\"stylesheet\"),e.setAttribute(\"media\",\"all\"),e.setAttribute(\"href\",\"https://dragontea.ink/wp-content/cache/autoptimize/css/autoptimize_5bd1c33b717b78702e18c3923e8fa4f0.css\"),document.head.appendChild(e),t(\".dmpvazRKNzBib1IxNjh0T0cwUUUxekEyY3F6Wm5QYzJDWGZqdXFnRzZ0TT0nuobc\").parent().prev().prev().prev();var a=1,n=16,r=11,i=\"08\",g=\"\",c=\"\",d=0,o=2,p=3,s=0,h=100;s++,s*=2,h/=2,h/=2;var $=3,u=20;function b(){let e=t(\".entry-header.header\"),a=parseInt(e.attr(\"data-id\"));return a}function m(t,e,a,n,r){return CryptoJSAesJson.decrypt(t,e+a+n+r)}function f(t,e){return CryptoJSAesJson.decrypt(t,e)}function l(t,e){return parseInt(t.toString()+e.toString())}function k(t,e,a){return t.toString()+e.toString()+a.toString()}$*=2,u=u-2-2,i=\"03\",o++,r++,n=n/4-2,a++,a*=4,n++,n++,n++,a-=5,r++,i=\"07\",t(\".reading-content .page-break img\").each(function(){var e,g=t(this),c=f(g.attr(\"id\").toString(),(e=parseInt((b()+l(r,i))*a-t(\".reading-content .page-break img\").length),e=l(2*n+1,e)).toString());g.attr(\"id\",c)}),r=0,n=0,a=0,i=0,t(\".reading-content .page-break img\").each(function(){var e=t(this),a=parseInt(e.attr(\"id\").replace(/image-(\\d+)[a-z]+/i,\"$1\"));t(\".reading-content .page-break\").eq(a).append(e)}),t(\".reading-content .page-break img\").each(function(){var e=t(this).attr(\"id\");g+=e.substr(-1),t(this).attr(\"id\",e.slice(0,-1))}),d++,$++,$++,u/=4,u*=2,o*=2,p-=3,p++,t(\".reading-content .page-break img\").each(function(){var e,a=t(this),n=f(a.attr(\"dta\").toString(),(e=parseInt((b()+l($,u))*(2*d)-t(\".reading-content .page-break img\").length-(4*d+1)),e=k(2*o+p+p+1,g,e)).toString());a.attr(\"dta\",n)}),d=0,$=0,u=0,o=0,p=0,t(\".reading-content .page-break img\").each(function(){var e=t(this).attr(\"dta\").substr(-2);c+=e,t(this).removeAttr(\"dta\")}),s*=s,s++,h-=25,h++,h++,t(\".reading-content .page-break img\").each(function(){var e=t(this),a=f(e.attr(\"data-src\").toString(),(b(),k(b()+4*s,c,t(\".reading-content .page-break img\").length*(2*h))).toString());e.attr(\"data-src\",a)}),s=0,h=0,t(\".reading-content .page-break img\").each(function(){t(this).addClass(\"wp-manga-chapter-img img-responsive lazyload effect-fade\")}),_0xabe6x4d=!0});"],["script","scri12pts && ifra2mes && coo1kies","true"],["script","(scri12pts && ifra2mes)","(true)"],["script","/catch[\\s\\S]*?}/","","condition","fetch"],["script","/(function playVideo\\(\\) \\{[\\s\\S]*?\\.remove\\(\\);[\\s\\S]*?\\})/","$1 playVideo();"],["script","video_urls.length != activeItem","!1"],["script","adblockDetected = true","adblockDetected = false"],["script","window._taboola =","(()=>{const e={apply:(e,o,l)=>o.closest(\"body > video[src^=\\\"blob:\\\"]\")===o?Promise.resolve(!0):Reflect.apply(e,o,l)};HTMLVideoElement.prototype.play=new Proxy(HTMLVideoElement.prototype.play,e)})();window._taboola ="],["script","/window.open.*;/"],["script","!seen","false"],["script","/if.*Disable.*?;/g","","condition","blocker"],["script","this.ads.length > this.ads_start","1==2"],["script","/\\$\\(['\"]\\.play-overlay['\"]\\)\\.click.+/s","document.getElementById(\"mainvideo\").src=srclink;player.currentTrack=0;})})","condition","srclink"],["script","const ad_slot_","(()=>{window.addEventListener(\"load\",(()=>{document.querySelectorAll(\"ins.adsbygoogle\").forEach((element=>element.dataset.adsbygoogleStatus=\"done\"))}))})();const ad_slot_","sedCount","1"],["script","window.dataLayer =","(()=>{const time=parseInt(document.querySelector(\"meta[http-equiv=\\\"refresh\\\"]\").content.split(\";\")[0])*1000+1000;setTimeout(()=>{document.body.innerHTML=document.body.innerHTML},time)})();window.dataLayer =","sedCount","1"],["script","(self.__next_f=","[\"timeupdate\",\"durationchange\",\"ended\",\"enterpictureinpicture\",\"leavepictureinpicture\",\"loadeddata\",\"loadedmetadata\",\"loadstart\",\"pause\",\"play\",\"playing\",\"ratechange\",\"resize\",\"seeked\",\"seeking\",\"suspend\",\"volumechange\",\"waiting\"].forEach((e=>{window.addEventListener(e,(()=>{const e=document.getElementById(\"player\"),t=document.querySelector(\".plyr__time\");e.src.startsWith(\"https://i.imgur.com\")&&\"none\"===window.getComputedStyle(t).display&&(e.src=\"https://cdn.plyr.io/static/blank.mp4\",e.paused&&e.plyr.play())}))}));(self.__next_f=","sedCount","1"],["script","vid.vast","//vid.vast"],["script","/(function reklamla\\([^)]+\\) {)/","$1rekgecyen(0);"],["script","'G-1B4LC0KT6C');","'G-1B4LC0KT6C'); window.setTimeout(function(){blockPing()},200);"],["script","/manageAds\\(video_urls\\[activeItem\\], video_seconds\\[activeItem\\], ad_urls\\[activeItem],true\\);/","playVideo();","condition","playAdd"],["script","preroll_timer_current == 0 && preroll_player_called == false","true"],["script","/^var \\w+=\\[.+/","(()=>{let e=[];document.addEventListener(\"DOMContentLoaded\",(()=>{const t=document.querySelector(\"body script\").textContent.match(/\"] = '(.*?)'/g);if(!t)return;t.forEach((t=>{const r=t.replace(/.*'(.*?)'/,\"$1\");e.push(r)}));const r=document.querySelector('.dl_button[href*=\"preview\"]').href.split(\"?\")[1];e.includes(r)&&(e=e.filter((e=>e!==r)));document.querySelectorAll(\".dl_button[href]\").forEach((t=>{t.target=\"_blank\";let r=t.cloneNode(!0);r.href=t.href.replace(/\\?.*/,`?${e[0]}`),t.after(r);let o=t.cloneNode(!0);o.href=t.href.replace(/\\?.*/,`?${e[1]}`),t.after(o)}))}))})();","sedCount","1"],["script","/window.location.href[^?]+this[^?]+;/"],["script","!isAdTriggered","false"],["script","/function _.*JSON.*}}/gms","function checkName(){const a = document.querySelector(\".monsters .button_wrapper .button\");const b = document.querySelector(\"#nick\");const c = \"/?from_land=1&nick=\";a.addEventListener(\"click\", function () {document.location.href = c + b.value;}); } checkName();"],["script","(function($)","(function(){const a=document.createElement(\"div\");document.documentElement.appendChild(a),setTimeout(()=>{a&&a.remove()},100)})(); (function($)"],["script","/window\\.dataLayer.+?(location\\.replace\\(\\S+?\\)).*/","$1"],["script","WB.defer","window.wbads={public:{getDailymotionAdsParamsForScript:function(a,b){b(\"\")},setTargetingOnPosition:function(a,b){return}}};WB.defer","condition","wbads.public.setTargetingOnPosition"],["script","var ISMLIB","!function(){const o={apply:(o,n,r)=>(new Error).stack.includes(\"refreshad\")?0:Reflect.apply(o,n,r)};window.Math.floor=new Proxy(window.Math.floor,o)}();var ISMLIB"],["script","(function serverContract()","/*START*/\"YOUTUBE_PREMIUM_LOGO\"!==ytInitialData?.topbar?.desktopTopbarRenderer?.logo?.topbarLogoRenderer?.iconImage?.iconType&&(location.href.startsWith(\"https://www.youtube.com/tv#/\")||location.href.startsWith(\"https://www.youtube.com/embed/\")||document.addEventListener(\"DOMContentLoaded\",(function(){const t=()=>{const t=document.getElementById(\"movie_player\");if(!t)return;if(!t.getStatsForNerds?.()?.debug_info?.startsWith?.(\"SSAP, AD\"))return;const e=t.getProgressState?.();e&&e.duration>0&&(e.loaded<e.duration||e.duration-e.current>1)&&t.seekTo?.(e.duration)};t(),new MutationObserver((()=>{t()})).observe(document,{childList:!0,subtree:!0})})));(function serverContract()","sedCount","1"],["script","/  function [a-zA-Z]{1,2}\\([a-zA-Z]{1,2},[a-zA-Z]{1,2}\\).*?\\(\\)\\{return [a-zA-Z]{1,2}\\;\\}\\;return [a-zA-Z]{1,2}\\(\\)\\;\\}/"],["script","/\\}\\)\\;\\s+\\(function\\(\\)\\{var .*?\\)\\;\\}\\)\\(\\)\\;\\s+\\$\\(\\\"\\#reportChapte/","}); $(\"#reportChapte"],["script","'G-1B4LC0KT6C');","'G-1B4LC0KT6C'); localStorage.setItem(\"tuna\", \"dW5kZWZpbmVk\"); localStorage.setItem(\"sausage\", \"ZmFsc2U=\"); window.setTimeout(function(){fuckYouUblockAndJobcenterTycoon(false)},200);"],["script","_w.keyMap=","(()=>{const e={apply:(e,t,n)=>{let o=Reflect.apply(e,t,n);return o instanceof HTMLIFrameElement&&!o.src&&o.contentWindow&&(o.contentWindow.document.body.getElementsByTagName=window.document.body.getElementsByTagName,o.contentWindow.MutationObserver=void 0),o}};HTMLBodyElement.prototype.appendChild=new Proxy(HTMLBodyElement.prototype.appendChild,e);const t={apply:(e,t,n)=>(t instanceof HTMLLIElement&&\"b_algo\"===t?.classList?.value&&\"a\"===n?.[0]&&setTimeout((()=>{t.style.display=\"none\"}),100),Reflect.apply(e,t,n))};HTMLBodyElement.prototype.getElementsByTagName=new Proxy(HTMLBodyElement.prototype.getElementsByTagName,t)})();_w.keyMap="],["script","var CBSNEWS =","navigator.serviceWorker.getRegistrations().then((registrations=>{for(const registration of registrations){if(registration.scope.includes(\"cbsnews.com\")){registration.unregister()}}}));var CBSNEWS ="]];
const hostnamesMap = new Map([["bild.de",0],["ashemaletube.*",[1,2]],["pornblade.com",3],["pornfelix.com",3],["nhentai.net",4],["jizzbunker.com",5],["multiup.io",6],["player.buffed.de",7],["player.gamezone.de",7],["player.gamesaktuell.de",7],["player.pcgames.de",7],["player.videogameszone.de",7],["player.pcgameshardware.de",7],["giga.de",8],["kino.de",8],["spieletipps.de",8],["desired.de",8],["t-online.de",9],["sport.de",10],["modijiurl.com",11],["starkroboticsfrc.com",12],["sinonimos.de",12],["antonimos.de",12],["quesignifi.ca",12],["tiktokrealtime.com",12],["tiktokcounter.net",12],["tpayr.xyz",12],["poqzn.xyz",12],["ashrfd.xyz",12],["rezsx.xyz",12],["tryzt.xyz",12],["ashrff.xyz",12],["rezst.xyz",12],["dawenet.com",12],["erzar.xyz",12],["waezm.xyz",12],["waezg.xyz",12],["blackwoodacademy.org",12],["cryptednews.space",12],["vivuq.com",12],["swgop.com",12],["vbnmll.com",12],["telcoinfo.online",12],["dshytb.com",12],["link.djbassking.live",[13,14]],["tech5s.co",[13,14]],["game5s.com",[13,14]],["itijobalert.in",15],["techkhulasha.com",15],["tech.pubghighdamage.com",[16,17,18]],["tech.techkhulasha.com",[16,17,18]],["hipsonyc.com",[16,17,18]],["labgame.io",19],["blog.cryptowidgets.net",[20,21]],["blog.insurancegold.in",[20,21]],["blog.wiki-topia.com",[20,21]],["blog.coinsvalue.net",[20,21]],["blog.cookinguide.net",[20,21]],["blog.freeoseocheck.com",[20,21]],["blog.coinsrise.net",21],["blog.makeupguide.net",21],["blog.carstopia.net",21],["blog.carsmania.net",21],["bitzite.com",21],["appsbull.com",21],["diudemy.com",21],["maqal360.com",21],["advertisingexcel.com",21],["allcryptoz.net",21],["batmanfactor.com",21],["beautifulfashionnailart.com",21],["crewbase.net",21],["crewus.net",21],["documentaryplanet.xyz",21],["gametechreviewer.com",21],["midebalonu.net",21],["misterio.ro",21],["phineypet.com",21],["seory.xyz",21],["shinbhu.net",21],["shinchu.net",21],["substitutefor.com",21],["talkforfitness.com",21],["thefitbrit.co.uk",21],["thumb8.net",21],["thumb9.net",21],["topcryptoz.net",21],["uniqueten.net",21],["ultraten.net",21],["exactpay.online",21],["aylink.co",22],["cpmlink.pro",22],["suaurl.com",23],["mamahawa.com",24],["dlink2.net",25],["drinkspartner.com",26],["apcvpc.com",26],["posterify.net",26],["manishclasses.in",26],["dl.apkmoddone.com",27],["phongroblox.com",27],["www.apkmoddone.com",[28,29]],["apkmodvn.com",30],["mahitimanch.in",[31,32]],["indiamaja.com",[31,32]],["5ghindi.in",[31,32]],["inshorturl.com",33],["rfiql.com",[34,35]],["gujjukhabar.in",[34,35]],["smartfeecalculator.com",[34,35]],["djxmaza.in",[34,35]],["thecubexguide.com",[34,35]],["jytechs.in",[34,35]],["dropgalaxy.*",[36,37]],["financemonk.net",[36,37]],["dailytech-news.eu",38],["fuckingfast.co",39],["driveup.sbs",40],["xdl.my.id",41],["in91vip.win",42],["jpvhub.com",43],["aagmaal.*",44],["etoday.co.kr",45],["isplus.com",46],["economist.co.kr",47],["hometownstation.com",48],["sportalkorea.com",49],["thestockmarketwatch.com",50],["m.edaily.co.kr",51],["honkailab.com",52],["lifehacker.jp",53],["eroasmr.com",54],["bussyhunter.com",55],["xcloud.*",56],["xfile.*",57],["kusonime.com",58],["codingnepalweb.com",59],["demonoid.is",60],["freebie-ac.jp",61],["photopea.com",62],["streamingcommunity.*",63],["derstandard.at",64],["derstandard.de",64],["www.chip.de",65],["redketchup.io",66],["emturbovid.com",67],["findjav.com",67],["javggvideo.xyz",67],["mmtv01.xyz",67],["stbturbo.xyz",67],["streamsilk.com",[67,68]],["zefoy.com",69],["idoitmyself.xyz",70],["uploadboy.com",71],["reymit.ir",72],["empire-anime.*",73],["empire-stream.*",73],["empire-streaming.*",73],["empire-streamz.*",73],["foodxor.com",74],["download.megaup.net",75],["adultdeepfakes.com",76],["client.falixnodes.net",[77,78]],["247sports.com",79],["linkshortify.com",80],["nexusmods.com",81],["comidacaseira.me",82],["tvbanywherena.com",83],["haveibeenpwned.com",84],["trainerscity.com",85],["tikmate.app",86],["sonixgvn.net",87],["paste-drop.com",88],["dizipal730.com",89],["filext.com",90],["kiddyearner.com",91],["playhydrax.com",92],["fullxh.com",93],["galleryxh.site",93],["megaxh.com",93],["movingxh.world",93],["seexh.com",93],["unlockxh4.com",93],["valuexh.life",93],["xhaccess.com",93],["xhadult2.com",93],["xhadult3.com",93],["xhadult4.com",93],["xhadult5.com",93],["xhamster.*",93],["xhamster1.*",93],["xhamster10.*",93],["xhamster11.*",93],["xhamster12.*",93],["xhamster13.*",93],["xhamster14.*",93],["xhamster15.*",93],["xhamster16.*",93],["xhamster17.*",93],["xhamster18.*",93],["xhamster19.*",93],["xhamster20.*",93],["xhamster2.*",93],["xhamster3.*",93],["xhamster4.*",93],["xhamster42.*",93],["xhamster46.com",93],["xhamster5.*",93],["xhamster7.*",93],["xhamster8.*",93],["xhamsterporno.mx",93],["xhbig.com",93],["xhbranch5.com",93],["xhchannel.com",93],["xhdate.world",93],["xhday.com",93],["xhday1.com",93],["xhlease.world",93],["xhmoon5.com",93],["xhofficial.com",93],["xhopen.com",93],["xhplanet1.com",93],["xhplanet2.com",93],["xhreal2.com",93],["xhreal3.com",93],["xhspot.com",93],["xhtotal.com",93],["xhtree.com",93],["xhvictory.com",93],["xhwebsite.com",93],["xhwebsite2.com",93],["xhwebsite5.com",93],["xhwide1.com",93],["xhwide2.com",93],["xhwide5.com",93],["dragontea.ink",94],["perchance.org",[95,96]],["cheater.ninja",97],["dizikral.com",[98,99]],["dizikral1.pro",98],["dizikral2.pro",98],["dizikral3.pro",98],["dizikral4.pro",98],["dizikral5.pro",98],["soccerinhd.com",100],["streamcaster.live",100],["barstoolsports.com",101],["filmi7.net",102],["hidan.co",103],["hidan.sh",103],["authenticateme.xyz",104],["filmizlehdizle.com",105],["fullfilmizlesene.net",105],["strcloud.in",106],["streamta.site",106],["streamtape.*",106],["watchadsontape.com",106],["tempmail.ninja",107],["filmy4wap.co.in",108],["filmy4waps.org",108],["filiser.eu",109],["wishflix.cc",109],["zalukaj.io",109],["cdn.tiesraides.lv",110],["filmizleplus.cc",111],["panel.freemcserver.net",[112,126]],["dizipal911.com",113],["dizipal912.com",113],["dizipal913.com",113],["dizipal914.com",113],["dizipal915.com",113],["dizipal916.com",113],["dizipal917.com",113],["dizipal918.com",113],["dizipal919.com",113],["dizipal920.com",113],["dizipal921.com",113],["dizipal922.com",113],["dizipal923.com",113],["dizipal924.com",113],["dizipal925.com",113],["dizipal926.com",113],["dizipal927.com",113],["dizipal928.com",113],["dizipal929.com",113],["dizipal930.com",113],["dizipal931.com",113],["dizipal932.com",113],["dizipal933.com",113],["dizipal934.com",113],["dizipal935.com",113],["dizipal936.com",113],["dizipal937.com",113],["dizipal938.com",113],["dizipal939.com",113],["dizipal940.com",113],["dizipal941.com",113],["dizipal942.com",113],["dizipal943.com",113],["dizipal944.com",113],["dizipal945.com",113],["dizipal946.com",113],["dizipal947.com",113],["dizipal948.com",113],["dizipal949.com",113],["dizipal950.com",113],["dizipal951.com",113],["dizipal952.com",113],["dizipal953.com",113],["dizipal954.com",113],["dizipal955.com",113],["dizipal956.com",113],["dizipal957.com",113],["dizipal958.com",113],["dizipal959.com",113],["dizipal960.com",113],["dizipal961.com",113],["dizipal962.com",113],["dizipal963.com",113],["dizipal964.com",113],["dizipal965.com",113],["dizipal966.com",113],["dizipal967.com",113],["dizipal968.com",113],["dizipal969.com",113],["dizipal970.com",113],["dizipal971.com",113],["dizipal972.com",113],["dizipal973.com",113],["dizipal974.com",113],["dizipal975.com",113],["dizipal976.com",113],["dizipal977.com",113],["dizipal978.com",113],["dizipal979.com",113],["dizipal980.com",113],["dizipal981.com",113],["dizipal982.com",113],["dizipal983.com",113],["dizipal984.com",113],["dizipal985.com",113],["dizipal986.com",113],["dizipal987.com",113],["dizipal988.com",113],["dizipal989.com",113],["dizipal990.com",113],["dizipal991.com",113],["dizipal992.com",113],["dizipal993.com",113],["dizipal994.com",113],["dizipal995.com",113],["dizipal996.com",113],["dizipal997.com",113],["dizipal998.com",113],["dizipal999.com",113],["indianporngirl10.com",114],["exploader.net",115],["flixhouse.com",116],["fastt.gg",117],["eatcells.com",118],["9to5google.com",119],["9to5mac.com",119],["api.dock.agacad.com",120],["ozap.com",121],["jprime.jp",122],["www.youtube.com",123],["infinityscans.xyz",[124,125]],["infinityscans.net",[124,125]],["infinityscans.org",[124,125]],["bing.com",127],["cbsnews.com",128]]);
const exceptionsMap = new Map([["xcloud.eu",[56]],["xcloud.host",[56]]]);
const hasEntities = true;
const hasAncestors = false;

const collectArgIndices = (hn, map, out) => {
    let argsIndices = map.get(hn);
    if ( argsIndices === undefined ) { return; }
    if ( typeof argsIndices !== 'number' ) {
        for ( const argsIndex of argsIndices ) {
            out.add(argsIndex);
        }
    } else {
        out.add(argsIndices);
    }
};

const indicesFromHostname = (hostname, suffix = '') => {
    const hnParts = hostname.split('.');
    const hnpartslen = hnParts.length;
    if ( hnpartslen === 0 ) { return; }
    for ( let i = 0; i < hnpartslen; i++ ) {
        const hn = `${hnParts.slice(i).join('.')}${suffix}`;
        collectArgIndices(hn, hostnamesMap, todoIndices);
        collectArgIndices(hn, exceptionsMap, tonotdoIndices);
    }
    if ( hasEntities ) {
        const n = hnpartslen - 1;
        for ( let i = 0; i < n; i++ ) {
            for ( let j = n; j > i; j-- ) {
                const en = `${hnParts.slice(i,j).join('.')}.*${suffix}`;
                collectArgIndices(en, hostnamesMap, todoIndices);
                collectArgIndices(en, exceptionsMap, tonotdoIndices);
            }
        }
    }
};

const entries = (( ) => {
    const docloc = document.location;
    const origins = [ docloc.origin ];
    if ( docloc.ancestorOrigins ) {
        origins.push(...docloc.ancestorOrigins);
    }
    return origins.map((origin, i) => {
        const beg = origin.lastIndexOf('://');
        if ( beg === -1 ) { return; }
        const hn = origin.slice(beg+3)
        const end = hn.indexOf(':');
        return { hn: end === -1 ? hn : hn.slice(0, end), i };
    }).filter(a => a !== undefined);
})();
if ( entries.length === 0 ) { return; }

const todoIndices = new Set();
const tonotdoIndices = new Set();

indicesFromHostname(entries[0].hn);
if ( hasAncestors ) {
    for ( const entry of entries ) {
        if ( entry.i === 0 ) { continue; }
        indicesFromHostname(entry.hn, '>>');
    }
}

// Apply scriplets
for ( const i of todoIndices ) {
    if ( tonotdoIndices.has(i) ) { continue; }
    try { replaceNodeText(...argsList[i]); }
    catch { }
}

/******************************************************************************/

// End of local scope
})();

void 0;
