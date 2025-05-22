import { LightningElement, wire, api } from 'lwc';
import { setTabLabel, setTabIcon, IsConsoleNavigation, getFocusedTabInfo } from 'lightning/platformWorkspaceApi';

import { getPicklistValuesByRecordType } from './ui.js';

import Go from 'c/go';

// Apex proxy for REST calls
import callRest from '@salesforce/apex/GoBridge.callRest';

export default class Thunder extends LightningElement {
	@api set recordId(value) {
		// Expose recordId from Lightning record page to Go WASM environment
		globalThis.recordId = value;
	}
	get recordId() {
		return globalThis.recordId;
	}
	// URL of the WASM app to load
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
		// Expose REST methods to Go WASM
		// get, post, patch, delete should call Apex @AuraEnabled proxy
		globalThis.get = (url) => callRest({ method: 'GET', url, body: null });
		globalThis.post = (url, body) => callRest({ method: 'POST', url, body });
		globalThis.patch = (url, body) => callRest({ method: 'PATCH', url, body });
		globalThis.delete = (url) => callRest({ method: 'DELETE', url, body: null });

		globalThis.getPicklistValuesByRecordType = getPicklistValuesByRecordType;

		const resp = await fetch(this.app);
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
