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

// Concatenate ArrayBuffers into a single ArrayBuffer, preserving order. Used to
// reassemble a WASM bundle that was split across multiple static resources to
// stay under Salesforce's 5MB per-resource limit.
function concatBuffers(buffers) {
	if (buffers.length === 1) {
		return buffers[0];
	}
	const total = buffers.reduce((sum, buf) => sum + buf.byteLength, 0);
	const combined = new Uint8Array(total);
	let offset = 0;
	for (const buf of buffers) {
		combined.set(new Uint8Array(buf), offset);
		offset += buf.byteLength;
	}
	return combined.buffer;
}

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

		// A WASM bundle larger than Salesforce's 5MB static resource limit is
		// split across several resources. The base resource carries a parts.json
		// manifest recording the total count; the remaining chunks live in
		// sibling resources named <base>Part1, <base>Part2, ... Fetch the base
		// chunk and the manifest together, then any remaining chunks, and
		// concatenate them before instantiating. A base resource without a
		// manifest (apps deployed before this feature) is loaded as a single part.
		const baseResourceUrl = this.app.slice(0, this.app.lastIndexOf('/'));
		const [firstResp, manifestResp] = await Promise.all([
			fetch(this.app),
			fetch(baseResourceUrl + '/parts.json').catch(() => null)
		]);

		let partCount = 1;
		if (manifestResp && manifestResp.ok) {
			try {
				const manifest = await manifestResp.json();
				if (manifest && manifest.parts > 1) {
					partCount = manifest.parts;
				}
			} catch (e) {
				// Missing or invalid manifest: load as a single-part bundle.
			}
		}

		const responses = [firstResp];
		if (partCount > 1) {
			const restUrls = [];
			for (let i = 1; i < partCount; i++) {
				restUrls.push(baseResourceUrl + 'Part' + i + '/bundle.wasm');
			}
			const rest = await Promise.all(restUrls.map((url) => fetch(url)));
			responses.push(...rest);
		}

		const failed = responses.find((resp) => !resp.ok);
		if (failed) {
			this.isLoading = false;
			const pre = document.createElement('pre');
			pre.innerText = await failed.text();
			divElement.appendChild(pre);
		} else {
			const buffers = await Promise.all(responses.map((resp) => resp.arrayBuffer()));
			const src = concatBuffers(buffers);
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
