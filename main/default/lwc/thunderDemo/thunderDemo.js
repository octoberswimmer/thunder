import { LightningElement } from 'lwc';
// Load the compiled Go WASM MASC app as a static resource
import THUNDER_APP from '@salesforce/resourceUrl/thunderDemo';

export default class ThunderDemo extends LightningElement {
	// Pass the Thunder WASM app URL into the thunder component
	thunderApp = THUNDER_APP;
}
