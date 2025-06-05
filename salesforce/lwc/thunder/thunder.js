import { LightningElement, wire, api } from 'lwc';
import { setTabLabel, setTabIcon, IsConsoleNavigation, getFocusedTabInfo } from 'lightning/platformWorkspaceApi';
import { NavigationMixin } from 'lightning/navigation';
import { CloseActionScreenEvent } from 'lightning/actions';

import { getPicklistValuesByRecordType } from './ui.js';

import Go from 'c/go';

// Apex proxy for REST calls
import callRest from '@salesforce/apex/GoBridge.callRest';

export default class Thunder extends NavigationMixin(LightningElement) {
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
		globalThis.get = (url) => {
			console.log('get called with ' + url);
			return callRest({ method: 'GET', url, body: null, cachebust: Date.now() })
				.then((res) => {
					console.log('got response for ' + url);
					console.log('type ' + typeof res);
					console.log(res);
					return res;
				})
				.catch((err) => {
					console.log('got error for ' + url);
					console.log('type ' + typeof err);
					console.log(JSON.stringify(err));
					throw err;
				});
		};
		globalThis.post = (url, body) => callRest({ method: 'POST', url, body });
		globalThis.patch = (url, body) => callRest({ method: 'PATCH', url, body });
		globalThis.delete = (url) => callRest({ method: 'DELETE', url, body: null });

		globalThis.getPicklistValuesByRecordType = getPicklistValuesByRecordType;

		// Expose Thunder exit functions to Go WASM
		globalThis.thunderExit = () => this.exitApp();
		globalThis.thunderExitToRecord = (recordId) => this.exitToRecord(recordId);
		globalThis.thunderCloseModal = () => this.closeModal();

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

	// Exit the Thunder app based on context
	exitApp() {
		if (this.isQuickAction()) {
			this.closeModal();
		} else {
			this.exitToRecord(this.recordId);
		}
	}

	// Navigate to a record's standard view page
	exitToRecord(recordId) {
		if (!recordId) {
			recordId = this.recordId;
		}

		if (recordId) {
			this[NavigationMixin.Navigate]({
				type: 'standard__recordPage',
				attributes: {
					recordId: recordId,
					actionName: 'view'
				}
			});
		}
	}

	// Close modal if running in quick action context
	closeModal() {
		// Use CloseActionScreenEvent for quick actions
		this.dispatchEvent(new CloseActionScreenEvent());
	}

	// Determine if running in a quick action context
	isQuickAction() {
		// Quick actions typically run in overlays or modals
		// This heuristic checks for common quick action contexts
		return window.location.pathname.includes('/one/one.app') &&
			   (window.parent !== window || document.body.classList.contains('forceOverlay'));
	}
}
