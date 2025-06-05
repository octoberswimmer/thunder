public with sharing class GoBridge {
	/**
	 * Generic REST proxy for client-side Go WASM calls.
	 * Mirrors Salesforce REST API endpoints via Apex @AuraEnabled method.
	 *
	 * @param method HTTP method (GET, POST, PUT, DELETE)
	 * @param url    Relative Salesforce REST API URL (e.g., '/services/data/v58.0/query?q=...')
	 * @param body   Request body for POST/PUT (JSON string), or null for GET/DELETE
	 * @return Response body as JSON string
	 */
	// Page size for cursor-based pagination; override in tests if needed
	public static Integer PAGE_SIZE = 2000;

	private static String getBaseURL(String url) {
		String pathOnly = getPathOnly(url);
		String qs = getQueryString(url);
		// Derive base path
		Integer vIdx = pathOnly.indexOf('/services/data/') +
			'/services/data/'.length();
		Integer vEnd = pathOnly.indexOf('/', vIdx);
		String version = pathOnly.substring(vIdx, vEnd);
		String base = '/services/data/' + version + '/query';
		return base;
	}

	private static String getQueryString(String url) {
		// Separate path and query string
		Integer qi = url.indexOf('?');
		String qs = qi > 0 ? url.substring(qi + 1) : '';
		return qs;
	}

	private static String getPathOnly(String url) {
		// Separate path and query string
		Integer qi = url.indexOf('?');
		String pathOnly = qi > 0 ? url.substring(0, qi) : url;
		return pathOnly;
	}

	private static Database.Cursor getCursor(String url) {
		String pathOnly = getPathOnly(url);
		Integer slash = pathOnly.lastIndexOf('/');
		String raw = pathOnly.substring(slash + 1);
		Integer dashPos = raw.lastIndexOf('-');
		String locatorId = dashPos >= 0 ? raw.substring(0, dashPos) : raw;
		String locJson = '{"queryId":"' + locatorId + '"}';
		Database.Cursor cursor = (Database.Cursor)JSON.deserialize(locJson, Database.Cursor.class);
		return cursor;
	}

	private static Integer getStart(String url) {
		String pathOnly = getPathOnly(url);
		Integer slash = pathOnly.lastIndexOf('/');
		String raw = pathOnly.substring(slash + 1);
		Integer dashPos = raw.lastIndexOf('-');
		Integer start = dashPos >= 0 ? Integer.valueOf(raw.substring(dashPos + 1)) : 0;
		return start;
	}

	@AuraEnabled
	public static String callRest(String method, String url, String body) {
		// Dispatch to appropriate handler based on request type
		if (url.startsWith('/services/apexrest/GoBridge/getThunderSettings')) {
			return handleThunderSettingsRequest();
		} else if (isCompositeRequest(method, url)) {
			return handleCompositeRequest(body);
		} else if (isQueryRequest(method, url)) {
			return handleQueryRequest(url);
		} else if (isSObjectRequest(method, url)) {
			return handleSObjectRequest(method, url, body);
		} else {
			throw new UnsupportedUrlException();
		}
	}

	/**
	 * Check if this is a composite request
	 */
	@TestVisible
	private static Boolean isCompositeRequest(String method, String url) {
		return method.equalsIgnoreCase('POST') &&
			url.startsWith('/services/data/') &&
			url.contains('/composite');
	}

	/**
	 * Check if this is a query request
	 */
	@TestVisible
	private static Boolean isQueryRequest(String method, String url) {
		return method.equalsIgnoreCase('GET') &&
			url.startsWith('/services/data/') &&
			url.contains('/query');
	}

	/**
	 * Check if this is an sObject request
	 */
	@TestVisible
	private static Boolean isSObjectRequest(String method, String url) {
		String pathOnly = getPathOnly(url);
		return url.startsWith('/services/data/') && pathOnly.indexOf('/sobjects/') >= 0;
	}

	/**
	 * Handle Thunder Settings requests
	 */
	private static String handleThunderSettingsRequest() {
		return getThunderSettings();
	}

	/**
	 * Handle SOQL query requests with cursor pagination
	 */
	private static String handleQueryRequest(String url) {
		String base = getBaseURL(url);
		String qs = getQueryString(url);
		String pathOnly = getPathOnly(url);

		Integer start;
		Database.Cursor cursor;

		if (qs.startsWith('q=')) {
			// Initial page with q=<SOQL>
			String soqlEnc = qs.substring(2);
			String soql = EncodingUtil.urlDecode(soqlEnc, 'UTF-8');

			start = 0;
			cursor = Database.getCursor(soql);
		} else {
			// Subsequent page via locator-offset
			start = getStart(url);
			cursor = getCursor(url);
		}

		Integer totalSize = cursor.getNumRecords();
		Integer remaining = totalSize - start;
		Integer count = remaining < PAGE_SIZE ? remaining : PAGE_SIZE;
		List<sObject> records = cursor.fetch(start, count);
		Boolean done = (start + records.size() >= totalSize);

		Map<String, Object> result = new Map<String, Object>{
			'totalSize' => totalSize,
			'done' => done,
			'records' => records
		};

		if (!done) {
			String locJson = JSON.serialize(cursor);
			Map<String, Object> locMap = (Map<String, Object>)JSON.deserializeUntyped(locJson);
			String locatorId = (String)locMap.get('queryId');
			Integer nextStart = start + records.size();
			String nextUrl = base + '/' +
				EncodingUtil.urlEncode(locatorId + '-' + nextStart, 'UTF-8');
			result.put('nextRecordsUrl', nextUrl);
		}

		return JSON.serialize(result);
	}

	/**
	 * Handle sObject create, update (PATCH), and delete requests
	 */
	private static String handleSObjectRequest(String method, String url, String body) {
		String pathOnly = getPathOnly(url);
		Integer soIdx = pathOnly.indexOf('/sobjects/');
		String rest = pathOnly.substring(soIdx + '/sobjects/'.length());
		List<String> parts = rest.split('/');
		String objectName = parts.size() > 0 ? parts[0] : null;
		Schema.SObjectType sobt = objectName != null ?
			Schema.getGlobalDescribe().get(objectName) :
			null;
		if (sobt == null) {
			throw new UnsupportedUrlException();
		}

		// Create new sObject
		if (method.equalsIgnoreCase('POST') &&
			parts.size() == 1) {
			return handleSObjectCreate(body, sobt);
		}
		// Update existing sObject via PATCH
		else if (method.equalsIgnoreCase('PATCH') &&
			parts.size() == 2) {
			String recId = parts[1];
			return handleSObjectUpdate(body, sobt, recId);
		}
		// Delete existing sObject
		else if (method.equalsIgnoreCase('DELETE') &&
			parts.size() == 2) {
			String recId = parts[1];
			return handleSObjectDelete(sobt, recId);
		} else {
			throw new UnsupportedUrlException();
		}
	}

	/**
	 * Handle sObject creation
	 */
	private static String handleSObjectCreate(String body, Schema.SObjectType sobt) {
		Map<String, Object> data = (Map<String, Object>)JSON.deserializeUntyped(body);
		SObject sobj = sobt.newSObject();
		setFieldsWithTypeConversion(sobj, data, sobt);
		insert sobj;
		Map<String, Object> result = new Map<String, Object>{
			'id' => (String)sobj.get('Id'),
			'success' => true,
			'errors' => new List<String>()
		};
		return JSON.serialize(result);
	}

	/**
	 * Handle sObject updates
	 */
	private static String handleSObjectUpdate(
		String body,
		Schema.SObjectType sobt,
		String recId
	) {
		Map<String, Object> data = (Map<String, Object>)JSON.deserializeUntyped(body);
		SObject sobj = sobt.newSObject();
		sobj.put('Id', recId);
		setFieldsWithTypeConversion(sobj, data, sobt);
		update sobj;
		Map<String, Object> result = new Map<String, Object>{
			'success' => true,
			'errors' => new List<String>()
		};
		return JSON.serialize(result);
	}

	/**
	 * Handle sObject deletion
	 */
	private static String handleSObjectDelete(Schema.SObjectType sobt, String recId) {
		SObject sobj = sobt.newSObject();
		sobj.put('Id', recId);
		delete sobj;
		Map<String, Object> result = new Map<String, Object>{
			'success' => true,
			'errors' => new List<String>()
		};
		return JSON.serialize(result);
	}

	private static String handleCompositeRequest(String body) {
		try {
			Map<String, Object> compositeRequest = (Map<String, Object>)JSON.deserializeUntyped(body);
			List<Object> subRequests = (List<Object>)compositeRequest.get('compositeRequest');
			Boolean allOrNone = compositeRequest.get('allOrNone') != null ?
				(Boolean)compositeRequest.get('allOrNone') :
				false;
			List<Map<String, Object>> results = new List<Map<String, Object>>();

			if (allOrNone) {
				return handleAllOrNoneCompositeRequest(subRequests);
			} else {
				return handleStandardCompositeRequest(subRequests);
			}
		} catch (Exception e) {
			Map<String, Object> errorResponse = new Map<String, Object>{
				'message' => 'Error processing composite request: ' + e.getMessage(),
				'errorCode' => 'COMPOSITE_REQUEST_ERROR'
			};
			return JSON.serialize(errorResponse);
		}
	}

	private static String handleStandardCompositeRequest(List<Object> subRequests) {
		List<Map<String, Object>> results = new List<Map<String, Object>>();
		Map<String, Map<String, Object>> referenceMap = new Map<String, Map<String, Object>>();

		for (Object subRequestObj : subRequests) {
			Map<String, Object> subRequest = (Map<String, Object>)subRequestObj;
			String method = (String)subRequest.get('method');
			String url = (String)subRequest.get('url');
			String subBody = subRequest.get('body') != null ?
				JSON.serialize(subRequest.get('body')) :
				null;
			String referenceId = (String)subRequest.get('referenceId');

			// Resolve references in URL and body
			url = resolveReferences(url, referenceMap);
			if (subBody != null) {
				subBody = resolveReferences(subBody, referenceMap);
			}

			Map<String, Object> subResult = new Map<String, Object>();
			subResult.put('referenceId', referenceId);

			try {
				String response = processSubRequest(method, url, subBody);
				subResult.put('httpStatusCode', 200);
				Map<String, Object> responseBody = (Map<String, Object>)JSON.deserializeUntyped(response);
				subResult.put('body', responseBody);

				// Store successful response for future references
				referenceMap.put(referenceId, responseBody);
			} catch (Exception e) {
				subResult.put('httpStatusCode', 400);
				subResult.put('body', new Map<String, Object>{
					'message' => e.getMessage(),
					'errorCode' => 'COMPOSITE_SUB_REQUEST_ERROR'
				});
			}

			results.add(subResult);
		}

		Map<String, Object> compositeResponse = new Map<String, Object>{ 'compositeResponse' => results };
		return JSON.serialize(compositeResponse);
	}

	private static String handleAllOrNoneCompositeRequest(List<Object> subRequests) {
		Savepoint sp = Database.setSavepoint();
		List<Map<String, Object>> results = new List<Map<String, Object>>();
		Map<String, Map<String, Object>> referenceMap = new Map<String, Map<String, Object>>();
		Boolean hasError = false;

		try {
			for (Object subRequestObj : subRequests) {
				Map<String, Object> subRequest = (Map<String, Object>)subRequestObj;
				String method = (String)subRequest.get('method');
				String url = (String)subRequest.get('url');
				String subBody = subRequest.get('body') != null ?
					JSON.serialize(subRequest.get('body')) :
					null;
				String referenceId = (String)subRequest.get('referenceId');

				// Resolve references in URL and body
				url = resolveReferences(url, referenceMap);
				if (subBody != null) {
					subBody = resolveReferences(subBody, referenceMap);
				}

				Map<String, Object> subResult = new Map<String, Object>();
				subResult.put('referenceId', referenceId);

				try {
					String response = processSubRequest(method, url, subBody);
					subResult.put('httpStatusCode', 200);
					Map<String, Object> responseBody = (Map<String, Object>)JSON.deserializeUntyped(response);
					subResult.put('body', responseBody);

					// Store successful response for future references
					referenceMap.put(referenceId, responseBody);
				} catch (Exception e) {
					hasError = true;
					subResult.put('httpStatusCode', 400);
					subResult.put('body', new Map<String, Object>{
						'message' => e.getMessage(),
						'errorCode' => 'COMPOSITE_SUB_REQUEST_ERROR'
					});
				}

				results.add(subResult);

				if (hasError) {
					break;
				}
			}

			if (hasError) {
				Database.rollback(sp);
				for (Map<String, Object> result : results) {
					if ((Integer)result.get('httpStatusCode') == 200) {
						result.put('httpStatusCode', 400);
						result.put('body', new Map<String, Object>{
							'message' => 'Transaction rolled back due to allOrNone=true',
							'errorCode' => 'ALL_OR_NONE_OPERATION_ROLLED_BACK'
						});
					}
				}
			}

			Map<String, Object> compositeResponse = new Map<String, Object>{ 'compositeResponse' => results };
			return JSON.serialize(compositeResponse);
		} catch (Exception e) {
			Database.rollback(sp);
			throw e;
		}
	}

	private static String processSubRequest(String method, String url, String body) {
		// Use the new specialized handlers for consistency
		if (isQueryRequest(method, url)) {
			return handleQueryRequest(url);
		} else if (isSObjectRequest(method, url)) {
			return handleSObjectRequest(method, url, body);
		} else {
			throw new UnsupportedUrlException();
		}
	}

	private static void setFieldsWithTypeConversion(
		SObject sobj,
		Map<String, Object> data,
		Schema.SObjectType sobjType
	) {
		// Get field describe info once for the entire object
		Map<String, Schema.SObjectField> fieldMap = sobjType.getDescribe().fields.getMap();

		for (String key : data.keySet()) {
			setFieldWithTypeConversion(sobj, key, data.get(key), fieldMap);
		}
	}

	private static void setFieldWithTypeConversion(
		SObject sobj,
		String fieldName,
		Object value,
		Map<String, Schema.SObjectField> fieldMap
	) {
		if (value == null) {
			sobj.put(fieldName, null);
			return;
		}

		// Get field describe info from the pre-fetched field map
		Schema.SObjectField field = fieldMap.get(fieldName);

		if (field == null) {
			// Field doesn't exist, let Salesforce handle the error
			sobj.put(fieldName, value);
			return;
		}

		Schema.DisplayType fieldType = field.getDescribe().getType();

		// Convert string dates to proper Date/DateTime objects
		if (value instanceof String) {
			String stringValue = (String)value;

			try {
				if (fieldType == Schema.DisplayType.Date) {
					// Parse ISO date format: YYYY-MM-DD
					if (Pattern.matches('\\d{4}-\\d{2}-\\d{2}', stringValue)) {
						List<String> dateParts = stringValue.split('-');
						Integer year = Integer.valueOf(dateParts[0]);
						Integer month = Integer.valueOf(dateParts[1]);
						Integer day = Integer.valueOf(dateParts[2]);
						sobj.put(fieldName, Date.newInstance(year, month, day));
						return;
					}
				} else if (fieldType == Schema.DisplayType.DateTime) {
					// Parse ISO datetime format: YYYY-MM-DDTHH:MM:SS or YYYY-MM-DD HH:MM:SS
					if (Pattern.matches('\\d{4}-\\d{2}-\\d{2}[T ]\\d{2}:\\d{2}:\\d{2}.*', stringValue)) {
						// Use JSON deserialization for datetime parsing as it handles ISO format well
						String jsonWrapper = '{"dt":"' + stringValue + '"}';
						Map<String, Object> wrapper = (Map<String, Object>)JSON.deserializeUntyped(jsonWrapper);
						Object deserializedDt = wrapper.get('dt');
						if (deserializedDt instanceof DateTime) {
							sobj.put(fieldName, deserializedDt);
							return;
						}
					}
				}
			} catch (Exception e) {
			// If conversion fails, fall through to default assignment
			}
		}

		// Default assignment for all other cases
		sobj.put(fieldName, value);
	}

	@TestVisible
	private static String resolveReferences(String input, Map<String, Map<String, Object>> referenceMap) {
		if (input == null || !input.contains('@{')) {
			return input;
		}

		String result = input;
		Pattern refPattern = Pattern.compile('@\\{([^}]+)\\}');
		Matcher matcher = refPattern.matcher(input);

		while (matcher.find()) {
			String fullMatch = matcher.group(0);
			String referenceExpression = matcher.group(1);

			Object resolvedValue = resolveReferenceExpression(referenceExpression, referenceMap);
			String replacement = resolvedValue != null ? String.valueOf(resolvedValue) : 'null';
			result = result.replace(fullMatch, replacement);
		}

		return result;
	}

	private static Object resolveReferenceExpression(
		String expression,
		Map<String, Map<String, Object>> referenceMap
	) {
		// Split by first dot to get referenceId and remaining path
		Integer firstDot = expression.indexOf('.');
		if (firstDot == -1) {
			// No dot, just return the whole reference object
			return referenceMap.get(expression);
		}

		String referenceId = expression.substring(0, firstDot);
		String remainingPath = expression.substring(firstDot + 1);

		if (!referenceMap.containsKey(referenceId)) {
			return null;
		}

		Object currentObject = referenceMap.get(referenceId);
		return resolveFieldPath(currentObject, remainingPath);
	}

	private static Object resolveFieldPath(Object currentObject, String fieldPath) {
		if (currentObject == null || String.isEmpty(fieldPath)) {
			return currentObject;
		}

		// Check for array indexing first
		if (fieldPath.contains('[')) {
			Integer bracketStart = fieldPath.indexOf('[');
			Integer bracketEnd = fieldPath.indexOf(']');

			if (bracketStart > 0 && bracketEnd > bracketStart) {
				// Get field name before bracket
				String fieldName = fieldPath.substring(0, bracketStart);
				String indexStr = fieldPath.substring(bracketStart + 1, bracketEnd);
				String remainingPath = bracketEnd + 1 < fieldPath.length() ?
					fieldPath.substring(bracketEnd + 1) :
					'';

				// Remove leading dot from remaining path
				if (remainingPath.startsWith('.')) {
					remainingPath = remainingPath.substring(1);
				}

				// Navigate to the field first
				Object fieldValue = getFieldValue(currentObject, fieldName);
				if (!(fieldValue instanceof List<Object>)) {
					return null;
				}

				// Get array index
				List<Object> arrayValue = (List<Object>)fieldValue;
				Integer index;
				try {
					index = Integer.valueOf(indexStr);
				} catch (TypeException e) {
					return null;
				}

				if (index < 0 || index >= arrayValue.size()) {
					return null;
				}

				Object indexedValue = arrayValue[index];

				// Continue with remaining path if any
				if (String.isEmpty(remainingPath)) {
					return indexedValue;
				} else {
					return resolveFieldPath(indexedValue, remainingPath);
				}
			}
		}

		// Handle simple field access
		Integer nextDot = fieldPath.indexOf('.');
		if (nextDot == -1) {
			// Last field in path
			return getFieldValue(currentObject, fieldPath);
		} else {
			// More fields to navigate
			String fieldName = fieldPath.substring(0, nextDot);
			String remainingPath = fieldPath.substring(nextDot + 1);
			Object fieldValue = getFieldValue(currentObject, fieldName);
			return resolveFieldPath(fieldValue, remainingPath);
		}
	}

	private static Object getFieldValue(Object obj, String fieldName) {
		if (obj instanceof Map<String, Object>) {
			Map<String, Object> objMap = (Map<String, Object>)obj;
			return objMap.get(fieldName);
		}
		return null;
	}

	/**
	 * Get Thunder Settings for the current user/org
	 * Returns custom settings data as JSON string
	 */
	@AuraEnabled
	public static String getThunderSettings() {
		try {
			Thunder_Settings__c settings = Thunder_Settings__c.getInstance();

			Map<String, Object> result = new Map<String, Object>();

			if (settings != null) {
				result.put('Google_Maps_API_Key__c', settings.Google_Maps_API_Key__c);
			} else {
				// Return empty object if no settings configured
				result.put('Google_Maps_API_Key__c', null);
			}

			return JSON.serialize(result);
		} catch (Exception e) {
			// Return error information
			Map<String, Object> errorResult = new Map<String, Object>{
				'error' => true,
				'message' => 'Failed to retrieve Thunder Settings: ' + e.getMessage()
			};
			return JSON.serialize(errorResult);
		}
	}

	class UnsupportedUrlException extends Exception {}
}
