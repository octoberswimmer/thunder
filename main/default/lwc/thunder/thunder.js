import { LightningElement, wire, api } from 'lwc';
import { setTabLabel, setTabIcon, IsConsoleNavigation, getFocusedTabInfo } from 'lightning/platformWorkspaceApi';

import Go from 'c/go';

// Default static resource for Thunder WASM app; can be overridden via @api app
import DEFAULT_THUNDER_APP from "@salesforce/resourceUrl/hello";
// Apex proxy for REST calls
import callRest from '@salesforce/apex/GoBridge.callRest';

// An example Thunder application
export default class Thunder extends LightningElement {
	// URL of the WASM app to load; uses DEFAULT_THUNDER_APP if not provided
	@api app;
	// Label to display on the console tab when navigation is enabled
	@api appName;
	@wire(IsConsoleNavigation) isConsoleNavigation;

	renderMode = "shadow";
	initialized = false;

	renderedCallback() {
		if (this.isConsoleNavigation) {
			getFocusedTabInfo().then((tabInfo) => {
				// Use provided appName or fallback to default
				const label = this.appName || 'Thunder App';
				setTabLabel(tabInfo.tabId, label);
				setTabIcon(tabInfo.tabId, 'apex');
			});
		}
		if (this.initialized) {
			return;
		}
		this.init();
	}

	async init() {
		this.initialized = true;
		var divElement = this.template.querySelector('div');
		const appUrl = this.app || DEFAULT_THUNDER_APP;
		// Expose REST methods to Go WASM
		// get, post, put, delete should call Apex @AuraEnabled proxy
		globalThis.get = (url) => callRest({ method: 'GET', url, body: null });
		globalThis.post = (url, body) => callRest({ method: 'POST', url, body });
		globalThis.put = (url, body) => callRest({ method: 'PUT', url, body });
		globalThis.delete = (url) => callRest({ method: 'DELETE', url, body: null });

		const resp = await fetch(appUrl);
		if (!resp.ok) {
			const pre = document.createElement('pre');
			pre.innerText = await resp.text();
			divElement.appendChild(pre);
		} else {
			const src = await resp.arrayBuffer();
			const go = new Go();
			const result = await WebAssembly.instantiate(src, go.importObject);
			go.run(result.instance);
			await new Promise(resolve => setTimeout(resolve, 1000));
			startWithDiv(divElement);
		}
	}
}
