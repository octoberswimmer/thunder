import { LightningElement, wire } from 'lwc';
import { setTabLabel, setTabIcon, IsConsoleNavigation, getFocusedTabInfo } from 'lightning/platformWorkspaceApi';

import Go from 'c/go';

import VECTY_APP from "@salesforce/resourceUrl/hello";

// An example Vecty application
export default class Vecty extends LightningElement {
	@wire(IsConsoleNavigation) isConsoleNavigation;

	renderMode = "shadow";
	initialized = false;

	renderedCallback() {
		if (this.isConsoleNavigation) {
			getFocusedTabInfo().then((tabInfo) => {
				setTabLabel(tabInfo.tabId, 'Masc Example');
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
		const resp = await fetch(VECTY_APP);
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
