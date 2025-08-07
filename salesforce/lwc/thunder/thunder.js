import { LightningElement, wire, api } from 'lwc';
import { setTabLabel, setTabIcon, IsConsoleNavigation, getFocusedTabInfo } from 'lightning/platformWorkspaceApi';
import { NavigationMixin } from 'lightning/navigation';
import { CloseActionScreenEvent } from 'lightning/actions';

import { getPicklistValuesByRecordType } from './ui.js';

import Go from 'c/go';

// Apex proxy for REST calls
import callRest from '@salesforce/apex/GoBridge.callRest';

// WeakMap to store instance-specific recordId associated with div elements
const divRecordIdMap = new WeakMap();

export default class Thunder extends NavigationMixin(LightningElement) {
	_recordId;

	@api set recordId(value) {
		// Store recordId as instance property
		this._recordId = value;
		// Update the WeakMap if we have a div element
		const divElement = this.template.querySelector('div');
		if (divElement && value) {
			divRecordIdMap.set(divElement, value);
		}
	}
	get recordId() {
		return this._recordId;
	}
	// URL of the WASM app to load
	@api app;
	// Label to display on the console tab when navigation is enabled
	@api appName = 'Thunder App';
	@wire(IsConsoleNavigation) isConsoleNavigation;

	renderMode = "shadow";
	initialized = false;
	isLoading = true;

	get appContainerClass() {
		return this.isLoading ? 'slds-hide' : '';
	}

	renderedCallback() {
		if (this.isConsoleNavigation && this.appName && !this.isQuickAction()) {
			getFocusedTabInfo().then((tabInfo) => {
				setTabLabel(tabInfo.tabId, this.appName);
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
		this.isLoading = true;
		var divElement = this.template.querySelector('div');

		// Store recordId in WeakMap if we have one (in case setter was called before div existed)
		if (this._recordId && divElement) {
			divRecordIdMap.set(divElement, this._recordId);
		}
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

		// Expose function to get recordId for a specific div
		globalThis.getRecordIdForDiv = (div) => divRecordIdMap.get(div);

		const resp = await fetch(this.app);
		if (!resp.ok) {
			this.isLoading = false;
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
			this.isLoading = false;
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
