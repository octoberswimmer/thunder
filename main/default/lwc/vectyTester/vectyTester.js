import { LightningElement } from 'lwc';
// Load the compiled Go WASM MASC app as a static resource
import MASC_APP from '@salesforce/resourceUrl/masc';

/**
 * Simple test harness for the Vecty WASM component.
 */
export default class VectyTester extends LightningElement {
    // Pass the MASC WASM app URL into the vecty component
    mascApp = MASC_APP;
}
