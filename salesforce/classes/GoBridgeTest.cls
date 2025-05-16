@isTest
private class GoBridgeTest {
	@isTest
	static void should_handle_query_urls() {
		List<Account> accts = new List<Account>();
		for (Integer i = 0; i < 3; i++) {
			accts.add(new Account(Name = 'TestAcct' + i, Type = 'VA'));
		}
		insert accts;

		String soql = 'SELECT Id, Name FROM Account ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');
		String url = '/services/data/v58.0/query?q=' + soqlEnc;
		String jsonResp = GoBridge.callRest('GET', url, null);

		Map<String, Object> result = (Map<String, Object>)JSON.deserializeUntyped(jsonResp);
		System.assertEquals(3, (Integer)result.get('totalSize'), 'totalSize should match inserted records');
		System.assertEquals(true, (Boolean)result.get('done'), 'done should be true');
		List<Object> recs = (List<Object>)result.get('records');
		System.assertEquals(3, recs.size(), 'records list should have 3 entries');
	}

	@isTest
	static void should_use_cursor_for_pagination() {
		// Override page size for pagination
		GoBridge.PAGE_SIZE = 2;
		// Create test data exceeding page size
		List<Account> accts = new List<Account>();
		for (Integer i = 0; i < 5; i++) {
			accts.add(new Account(Name = 'Paginate' + i, Type = 'VA'));
		}
		insert accts;
		// Build query URL (initial page)
		String soql = 'SELECT Id, Name FROM Account WHERE Name LIKE \'Paginate%\' ORDER BY Name';
		String soqlEnc = EncodingUtil.urlEncode(soql, 'UTF-8');
		String url = '/services/data/v58.0/query?q=' + soqlEnc;

		// First page
		String jsonResp1 = GoBridge.callRest('GET', url, null);
		Map<String, Object> result1 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp1);
		System.assertEquals(5, (Integer)result1.get('totalSize'));
		System.assertEquals(false, (Boolean)result1.get('done'));
		List<Object> recs1 = (List<Object>)result1.get('records');
		System.assertEquals(2, recs1.size());
		String nextUrl1 = (String)result1.get('nextRecordsUrl');
		System.assertNotEquals(null, nextUrl1);

		// Second page
		String jsonResp2 = GoBridge.callRest('GET', nextUrl1, null);
		Map<String, Object> result2 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp2);
		System.assertEquals(false, (Boolean)result2.get('done'));
		List<Object> recs2 = (List<Object>)result2.get('records');
		System.assertEquals(2, recs2.size());
		String nextUrl2 = (String)result2.get('nextRecordsUrl');
		System.assertNotEquals(null, nextUrl2);

		// Final page
		String jsonResp3 = GoBridge.callRest('GET', nextUrl2, null);
		Map<String, Object> result3 = (Map<String, Object>)JSON.deserializeUntyped(jsonResp3);
		System.assertEquals(true, (Boolean)result3.get('done'));
		List<Object> recs3 = (List<Object>)result3.get('records');
		System.assertEquals(1, recs3.size());
	}

	@isTest
	static void unsupported_url_should_throw_exception() {
		String url = '/services/data/v58.0/sobjects/Account';
		String body = '{"Name":"FooTest"}';
		Boolean exceptionThrown = false;
		try {
			String resp = GoBridge.callRest('POST', url, body);
		} catch (GoBridge.UnsupportedUrlException e) {
			exceptionThrown = true;
		}
		System.assertEquals(true, exceptionThrown, 'exception thrown');
	}
}
