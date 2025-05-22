import { getRecord as getRecordAdapter } from 'lightning/uiRecordApi';
import { getPicklistValuesByRecordType as getPicklistValuesByRecordTypeAdapter, getObjectInfo as getObjectInfoAdapter } from "lightning/uiObjectInfoApi";


function getRecord(config, cb) {
	var gr = new getRecordAdapter(result => {
		const { data, error } = result;
		if (error) {
			cb({ error: error})
		}
		if (data) {
			cb({ data: data})
		}
	});
	gr.connect();
	gr.update(config);
}

function getPicklistValuesByRecordType(config, cb) {
	var gr = new getPicklistValuesByRecordTypeAdapter(result => {
		const { data, error } = result;
		if (error) {
			cb({ error: error})
		}
		if (data) {
			cb({ data: data})
		}
	});
	gr.connect();
	gr.update(config);
}

function getObjectInfo(config, cb) {
	var oi = new getObjectInfoAdapter(result => {
		const { data, error } = result;
		if (error) {
			cb({ error: error });
		}
		if (data) {
			cb({ data: data });
		}
	});
	oi.connect();
	oi.update(config);
}

export {
	getRecord,
	getPicklistValuesByRecordType,
	getObjectInfo
};
