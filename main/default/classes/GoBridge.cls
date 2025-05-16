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
		// Handle SOQL cursor pagination via REST-style query locator URLs
		if (method.equalsIgnoreCase('GET') &&
			url.startsWith('/services/data/') &&
			url.contains('/query')) {
			String base = getBaseURL(url);
			String qs = getQueryString(url);
			String pathOnly = getPathOnly(url);

			Integer start;
			DAtabase.Cursor cursor;

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

		throw new UnsupportedUrlException();
	}

	class UnsupportedUrlException extends Exception {}
}
